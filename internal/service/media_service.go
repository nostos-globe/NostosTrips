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
    url, err := config.MinioClient.PresignedGetObject(context.TODO(), "media", media.FilePath, time.Minute*5, nil)
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
