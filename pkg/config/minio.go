// secrets.go
package config

import (
    "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/credentials"
    "log"
    "context"
)

var MinioClient *minio.Client

func InitMinIO() *minio.Client {
    endpoint := "minioapi.nostos-globe.me"
    client, err := minio.New(endpoint, &minio.Options{
        Creds:  credentials.NewStaticV4("nostos", "nostos123", ""),
        Secure: true,
    })
    if err != nil {
        log.Fatalf("Failed to init MinIO: %v", err)
    }

    // Check if bucket exists, if not create it
    bucketName := "nostos-media"
    exists, err := client.BucketExists(context.Background(), bucketName)
    if err != nil {
        log.Printf("Error checking bucket: %v", err)
    }

    if !exists {
        err = client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
        if err != nil {
            log.Printf("Error creating bucket: %v", err)
        } else {
            log.Printf("Successfully created bucket: %s", bucketName)
        }
    }

    MinioClient = client
    log.Printf("Successfully connected to MinIO")
    return client
}
