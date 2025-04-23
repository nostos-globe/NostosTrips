package models

type Trip struct {
	TripID      int    `gorm:"primaryKey;autoIncrement"`
	UserID      uint   `json:"user_id" db:"user_id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description,omitempty" db:"description"`
	Visibility  string `json:"visibility" db:"visibility" default:"PRIVATE"`
	StartDate   string `json:"start_date,omitempty" db:"start_date"`
	EndDate     string `json:"end_date,omitempty" db:"end_date"`
}

type TripRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Visibility  string `json:"visibility"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	AlbumID     string `json:"album_id"`
}


