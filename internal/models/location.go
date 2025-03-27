package models

type Location struct {
    LocationID int64  `gorm:"primaryKey;autoIncrement"`
    Name       string `json:"name"`
    Country    string `json:"country"`
    City       string `json:"city"`
}
