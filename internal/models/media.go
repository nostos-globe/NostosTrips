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
    MediaID      int64          `json:"media_id" gorm:"column:media_id"`
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