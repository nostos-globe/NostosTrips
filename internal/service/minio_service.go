package service

import (
	"context"
	"fmt"
	"main/pkg/config"
	"mime/multipart"
	"time"

	"github.com/minio/minio-go/v7"
)

type MinioService struct {
	BucketName string
}

func NewMinioService() *MinioService {
	return &MinioService{
		BucketName: "nostos-media",
	}
}

func (s *MinioService) UploadFile(file multipart.File, header *multipart.FileHeader) (string, error) {
	objectName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), header.Filename)
	fmt.Printf("Generated object name: %s\n", objectName)

	fmt.Printf("Uploading file to MinIO bucket '%s'. Size: %d, Content-Type: %s\n",
		s.BucketName,
		header.Size,
		header.Header.Get("Content-Type"))

	_, err := config.MinioClient.PutObject(
		context.Background(),
		s.BucketName,
		objectName,
		file,
		header.Size,
		minio.PutObjectOptions{
			ContentType: header.Header.Get("Content-Type"),
		},
	)

	if err != nil {
		fmt.Printf("Error uploading to MinIO: %v\n", err)
		return "", err
	}
	fmt.Println("Successfully uploaded file to MinIO")

	return objectName, nil
}

func (s *MinioService) GetPresignedURL(objectName string, duration time.Duration) (string, error) {
	url, err := config.MinioClient.PresignedGetObject(
		context.Background(),
		s.BucketName,
		objectName,
		duration,
		nil,
	)
	if err != nil {
		return "", err
	}

	return url.String(), nil
}

func (s *MinioService) DeleteObject(objectName string) error {
	err := config.MinioClient.RemoveObject(
		context.Background(),
		s.BucketName,
		objectName,
		minio.RemoveObjectOptions{},
	)
	if err != nil {
		return fmt.Errorf("failed to delete object from MinIO: %w", err)
	}
	return nil
}