package models

import (
    "time"
)

type VisibilityEnum string

const (
    Public  VisibilityEnum = "PUBLIC"
    Private VisibilityEnum = "PRIVATE"
	Friends VisibilityEnum = "FRIENDS"
)

type Media struct {
    MediaID      int64          `gorm:"primaryKey;autoIncrement"`
    TripID       int64          `json:"trip_id" gorm:"column:trip_id"`
    UserID       int64          `json:"user_id" gorm:"column:user_id"`
    LocationID   int64          `json:"location_id" gorm:"column:location_id"`
    Type         string         `json:"type"`
    FilePath     string         `json:"file_path"`
    Visibility   VisibilityEnum `json:"visibility"`
    UploadDate   time.Time      `json:"upload_date"`
    CaptureDate  time.Time      `json:"capture_date"`
    GpsLatitude  float64        `json:"gps_latitude"`
    GpsLongitude float64        `json:"gps_longitude"`
    GpsAltitude  float64        `json:"gps_altitude"`
}

type MediaMetadata struct {
    Type                 string
    LocationID           int64
    CaptureDate          time.Time
    Latitude             float64
    Longitude            float64
    Altitude             float64
    City                 string
    Country              string
    RequiresManualLocation bool // New field to indicate if manual location is needed
}