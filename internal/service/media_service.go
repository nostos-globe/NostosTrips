package service

import (
	"context"
	"fmt"
	"main/internal/db"
	"main/internal/models"
	"mime/multipart"
    "main/pkg/config"
	"strconv"
	"time"
	"github.com/minio/minio-go/v7"
    "github.com/rwcarlsen/goexif/exif"
    "path/filepath"
    "strings"
    "encoding/json"
    "net/http"
)

type MediaService struct {
	MediaRepo *db.MediaRepository
}

func (s *MediaService) UploadMedia(userID int64, file multipart.File, header *multipart.FileHeader, visibility models.VisibilityEnum) (string, error) {
    fmt.Printf("Starting upload for user %d with file %s\n", userID, header.Filename)
    
    defer file.Close()

    // Define bucket name as a constant
    const bucketName = "nostos-media"

    objectName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), header.Filename)
    fmt.Printf("Generated object name: %s\n", objectName)

    fmt.Printf("Uploading file to MinIO bucket '%s'. Size: %d, Content-Type: %s\n", 
        bucketName,
        header.Size, 
        header.Header.Get("Content-Type"))

    _, err := config.MinioClient.PutObject(context.Background(), bucketName, objectName, file, header.Size, minio.PutObjectOptions{
        ContentType: header.Header.Get("Content-Type"),
    })

    if err != nil {
        fmt.Printf("Error uploading to MinIO: %v\n", err)
        return "", err
    }
    fmt.Println("Successfully uploaded file to MinIO")

    return objectName, nil
}

func (s *MediaService) GetMediaURL(mediaID int64, userID int64) (string, error) {
    media, err := s.MediaRepo.GetMediaByID(mediaID)
    if err != nil {
        return "", err
    }

    // Permissions check
    switch media.Visibility {
    case "PRIVATE":
        if media.UserID != userID {
            return "", fmt.Errorf("not authorized")
        }
    case "FRIENDS":
        if !s.MediaRepo.AreFriends(userID, media.UserID) {
            return "", fmt.Errorf("not authorized")
        }
    }

    // Signed URL
    url, err := config.MinioClient.PresignedGetObject(context.TODO(), "nostos-media", media.FilePath, time.Minute*5, nil)
    if err != nil {
        return "", err
    }

    return url.String(), nil
}

func (s *MediaService) DeleteMediaByTripID(tripID string) error {
	id, err := strconv.Atoi(tripID)
	if err != nil {
		return err
	}

	result := s.MediaRepo.DeleteMediaByTripID(id)
	if result != nil {
		return result
	}

	return nil
}

func (s *MediaService) GetMediaByID(id int) (any, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *MediaService) SaveMedia(media *models.Media) error {
    return s.MediaRepo.SaveMedia(media)
}

func (s *MediaService) DeleteMedia(mediaID int64, tripID string) error {
    // First verify the media belongs to the trip
    media, err := s.MediaRepo.GetMediaByID(mediaID)
    if err != nil {
        return err
    }

    tripIDInt, err := strconv.ParseInt(tripID, 10, 64)
    if err != nil {
        return err
    }

    if media.TripID != tripIDInt {
        return fmt.Errorf("media does not belong to this trip")
    }

    return s.MediaRepo.DeleteMedia(mediaID)
}

func (s *MediaService) ExtractMetadata(file multipart.File, header *multipart.FileHeader) (*models.MediaMetadata, error) {
    metadata := &models.MediaMetadata{
        CaptureDate: time.Now(),
    }

    // Determine media type
    fileExt := strings.ToLower(filepath.Ext(header.Filename))
    switch fileExt {
    case ".jpg", ".jpeg", ".png", ".gif":
        metadata.Type = "photo"
    case ".mp4", ".mov", ".avi":
        metadata.Type = "video"
    default:
        metadata.Type = "unknown"
    }

    fmt.Printf("Processing file: %s (type: %s)\n", header.Filename, metadata.Type)

    if metadata.Type == "photo" {
        fmt.Println("Attempting to extract EXIF data...")
        exifData, err := exif.Decode(file)
        if err != nil {
            fmt.Printf("Failed to decode EXIF: %v\n", err)
        } else {
            // Get GPS coordinates
            if lat, long, err := exifData.LatLong(); err == nil {
                metadata.Latitude = lat
                metadata.Longitude = long
                fmt.Printf("Found coordinates: lat=%f, long=%f\n", lat, long)

                // Get altitude if available
                if altTag, err := exifData.Get(exif.GPSAltitude); err == nil {
                    ratValue, err := altTag.Rat(0)
                    if err == nil {
                        metadata.Altitude, _ = ratValue.Float64()
                        fmt.Printf("Found altitude: %f\n", metadata.Altitude)
                    }
                }

                // Get location info if coordinates are available
                if locationInfo, err := s.getLocationInfo(lat, long); err == nil {
                    locationFound, err := s.GetLocationByCountryAndCity(locationInfo)
                    fmt.Printf("Location found: %s, %s\n", locationInfo.City, locationInfo.Country)

                    if err == nil{
                        metadata.LocationID = locationFound.LocationID
                    	metadata.City = locationFound.City
                        metadata.Country = locationFound.Country
                    } else {
                        fmt.Printf("Location not found in database\n")
                        locationCreated, err := s.setLocationInfo(locationInfo)
                        if err == nil {
                            metadata.LocationID = locationCreated.LocationID
                            metadata.City = locationCreated.City
                            metadata.Country = locationCreated.Country                        
                        }

                    }
                    

                    
                    fmt.Printf("Location found: %s, %s\n", metadata.City, metadata.Country)
                } else {
                    fmt.Printf("Failed to get location info: %v\n", err)
                }
            } else {
                fmt.Printf("No GPS coordinates found: %v\n", err)
            }

            // Get capture date
            if dateTag, err := exifData.DateTime(); err == nil {
                metadata.CaptureDate = dateTag
                fmt.Printf("Found capture date: %v\n", metadata.CaptureDate)
            } else {
                fmt.Printf("No capture date found: %v\n", err)
            }
        }
        
        // Reset file pointer for future operations
        file.Seek(0, 0)
    }

    fmt.Println("Final metadata summary:")
    fmt.Printf("- Type: %s\n", metadata.Type)
    fmt.Printf("- Capture Date: %v\n", metadata.CaptureDate)
    fmt.Printf("- Coordinates: lat=%f, long=%f, alt=%f\n", metadata.Latitude, metadata.Longitude, metadata.Altitude)
    fmt.Printf("- Location: %s, %s\n", metadata.City, metadata.Country)

    return metadata, nil
}

func (s *MediaService) getLocationInfo(lat, long float64) (*models.Location, error) {
    // Using OpenStreetMap Nominatim API (free, no API key required)
    url := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?format=json&lat=%f&lon=%f", lat, long)
    
    client := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("User-Agent", "NostosTrips/1.0")
    
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result struct {
        Address struct {
            City    string `json:"city"`
            Country string `json:"country"`
        } `json:"address"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    return &models.Location{
        City:    result.Address.City,
        Country: result.Address.Country,
    }, nil
}

func (s *MediaService) GetLocationByCountryAndCity(location *models.Location) (*models.Location, error) {
    return s.MediaRepo.GetLocationByCountryAndCity(location)
}


func (s *MediaService) setLocationInfo(location *models.Location) (*models.Location, error) {
    if location == nil {
        return nil, fmt.Errorf("location cannot be nil")
    }
    
    if err := s.MediaRepo.SaveLocationInfo(location); err != nil {
        return nil, fmt.Errorf("failed to save location info: %w", err)
    }
    
    return location, nil
}