package service

import (
	"context"
	"encoding/json"
	"fmt"
	"main/internal/db"
	"main/internal/models"
	"main/pkg/config"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/rwcarlsen/goexif/exif"
)

type MediaService struct {
	MediaRepo    *db.MediaRepository
	MinioService *MinioService
}

func (s *MediaService) GetMediaByTripID(tripID int64, userID int64) ([]models.MediaByTrip, error) {
	// Get media list from repository
	mediaList, err := s.MediaRepo.GetMediaByTripID(tripID)
	if err != nil {
		return nil, fmt.Errorf("failed to get media: %w", err)
	}

	var response []models.MediaByTrip
	for _, media := range mediaList {
		// Check visibility permissions
		if media.Visibility == "PRIVATE" && media.UserID != userID {
			continue
		}
		if media.Visibility == "FRIENDS" && !s.MediaRepo.AreFriends(userID, media.UserID) {
			continue
		}

		// Get presigned URL for the media
		url, err := s.MinioService.GetPresignedURL(media.FilePath, time.Minute*5)
		if err != nil {
			continue
		}

		response = append(response, models.MediaByTrip{
			MediaID: media.MediaID,
			URL:     url,
		})
	}

	return response, nil
}

func (s *MediaService) ChangeMediaVisibility(mediaID int64, i int64, visibility models.VisibilityEnum) error {
	media, err := s.MediaRepo.GetMediaByID(mediaID)
	if media == nil || err != nil {
		return err
	}

	media.Visibility = visibility

	err = s.MediaRepo.UpdateMedia(mediaID, media)
	if err != nil {
		return fmt.Errorf("failed to update media: %w", err)
	}
	return nil
}

func (s *MediaService) UpdateMediaMetadata(mediaID int64, i int64, latitude float64, longitude float64, altitude float64) error {
	media, err := s.MediaRepo.GetMediaByID(mediaID)
	if media == nil || err != nil {
		return err
	}

	media.GpsLatitude = latitude
	media.GpsLongitude = longitude
	media.GpsAltitude = altitude

	// Get location info based on coordinates
	locationInfo, err := s.getLocationInfo(latitude, longitude)
	if err != nil {
		return err
	}

	// Try to find existing location in database
	locationFound, err := s.GetLocationByCountryAndCity(locationInfo)
	if err == nil {
		media.LocationID = locationFound.LocationID
	} else {
		// Create new location if not found
		locationCreated, err := s.setLocationInfo(locationInfo)
		if err != nil {
			return err
		}
		media.LocationID = locationCreated.LocationID
		locationInfo.City = locationCreated.City
		locationInfo.Country = locationCreated.Country
	}

	err = s.MediaRepo.UpdateMedia(media.MediaID, media)
	if err != nil {
		return fmt.Errorf("failed to update media: %w", err)
	}
	return nil
}

func (s *MediaService) UploadMedia(userID int64, file multipart.File, header *multipart.FileHeader, visibility models.VisibilityEnum) (string, error) {
	fmt.Printf("Starting upload for user %d with file %s\n", userID, header.Filename)
	defer file.Close()

	// Use the MinioService to upload the file
	objectName, err := s.MinioService.UploadFile(file, header)
	if err != nil {
		return "", err
	}

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

	// Get presigned URL using MinioService
	url, err := s.MinioService.GetPresignedURL(media.FilePath, time.Minute*5)
	if err != nil {
		return "", err
	}

	return url, nil
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
	// Get file extension and MIME type
	fileExt := strings.ToLower(filepath.Ext(header.Filename))
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return nil, err
	}
	file.Seek(0, 0) // Reset file pointer

	mimeType := http.DetectContentType(buffer)

	// Check both extension and MIME type
	switch {
	case fileExt == ".jpg" || fileExt == ".jpeg" || fileExt == ".png" || fileExt == ".gif" || fileExt == ".webp" || fileExt == ".tiff" ||
		strings.HasPrefix(mimeType, "image/"):
		metadata.Type = "photo"
	case fileExt == ".mp4" || fileExt == ".mov" || fileExt == ".avi" || fileExt == ".mkv" || fileExt == ".webm" || fileExt == ".flv" ||
		strings.HasPrefix(mimeType, "video/"):
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

					if err == nil {
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

	// Check if location data is missing or has default values
	if (metadata.Latitude == 0 && metadata.Longitude == 0) ||
		(metadata.City == "" && metadata.Country == "") {
		// We'll still return the metadata, but with a special error code
		// that indicates manual location input is required
		metadata.RequiresManualLocation = true
		return metadata, fmt.Errorf("MANUAL_LOCATION_REQUIRED")
	}

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
		Name:    result.Address.City + ", " + result.Address.Country,
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

func (s *MediaService) DeleteMediaCompletely(mediaID int64, userID int64) error {
	// First get the media to check permissions and get the file path
	media, err := s.MediaRepo.GetMediaByID(mediaID)
	if err != nil {
		return fmt.Errorf("failed to find media: %w", err)
	}

	// Check if the user owns this media
	if media.UserID != userID {
		return fmt.Errorf("not authorized to delete this media")
	}

	// Delete from MinIO using MinioService
	err = s.MinioService.DeleteObject(media.FilePath)
	if err != nil {
		return fmt.Errorf("failed to delete from storage: %w", err)
	}

	// Delete from database
	err = s.MediaRepo.DeleteMedia(mediaID)
	if err != nil {
		return fmt.Errorf("failed to delete from database: %w", err)
	}

	return nil
}

func (s *MediaService) deleteFromStorage(filePath string) error {
	// Delete the object from MinIO
	err := config.MinioClient.RemoveObject(context.Background(), "nostos-media", filePath, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}
