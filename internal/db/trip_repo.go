package db

import (
	"main/internal/models"

	"gorm.io/gorm"
)

type TripsRepository struct {
	DB *gorm.DB
}

func (repo *TripsRepository) CreateTrip(trip models.Trip) (any, error) {
	result := repo.DB.Table("trips.trips").Create(&trip)
	if result.Error != nil {
		return nil, result.Error
	}

	return trip, nil
}

func (repo *TripsRepository) GetTripByID(tripID int) (models.Trip, error) {
	var trip models.Trip
	result := repo.DB.Table("trips.trips").First(&trip, tripID)
	if result.Error != nil {
		return models.Trip{}, result.Error
	}

	return trip, nil
}
