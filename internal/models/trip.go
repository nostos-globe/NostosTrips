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
