package models

import "time"

type Trip struct {
	TripID      uint      `json:"trip_id" db:"trip_id"`
	UserID      uint      `json:"user_id" db:"user_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description,omitempty" db:"description"`
	Visibility  string    `json:"visibility" db:"visibility" default:"PRIVATE"`
	StartDate   time.Time `json:"start_date,omitempty" db:"start_date"`
	EndDate     time.Time `json:"end_date,omitempty" db:"end_date"`
}
