package db

import (
	"fmt"
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

func (repo *TripsRepository) UpdateTrip(trip models.Trip) (any, error) {
	result := repo.DB.Table("trips.trips").Where("trip_id = ?", trip.TripID).Updates(&trip)
	if result.Error != nil {
		return nil, result.Error
	}

	return trip, nil
}

func (repo *TripsRepository) DeleteTrip(tripID int) error {
	result := repo.DB.Table("albums.album_trips").Where("trip_id =?", tripID).Delete("trips.trips")
	if result.Error != nil {
		return result.Error
	}

	result = repo.DB.Table("trips.trips").Delete(&models.Trip{}, tripID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *TripsRepository) GetAllTrips() ([]models.Trip, error) {
	var trips []models.Trip
	result := repo.DB.Table("trips.trips").Find(&trips)
	if result.Error != nil {
		return nil, result.Error
	}
	return trips, nil
}

func (repo *TripsRepository) GetMyTrips(userID uint) ([]models.Trip, error) {
	var trips []models.Trip
	result := repo.DB.Table("trips.trips").Where("user_id = ?", userID).Find(&trips)
	if result.Error != nil {
		return nil, result.Error
	}
	return trips, nil
}

func (repo *TripsRepository) SearchTrips(query string, userID uint) ([]models.Trip, error) {
	var trips []models.Trip
	searchQuery := "%" + query + "%"
	result := repo.DB.Table("trips.trips").
		Where("user_id != ? AND (name ILIKE ? OR description ILIKE ?)",
			userID, searchQuery, searchQuery).
		Find(&trips)
	if result.Error != nil {
		return nil, result.Error
	}
	fmt.Printf("Search results from DB: %+v\n", trips)
	return trips, nil
}

func (repo *TripsRepository) GetTripByID(tripID int) (models.Trip, error) {
	var trip models.Trip
	result := repo.DB.Table("trips.trips").First(&trip, tripID)
	if result.Error != nil {
		return models.Trip{}, result.Error
	}

	return trip, nil
}

func (r *TripsRepository) GetAllPublicTrips() ([]models.Trip, error) {
	var trips []models.Trip
	result := r.DB.Table("trips.trips").Where("visibility = ?", "PUBLIC").Find(&trips)
	if result.Error != nil {
		return nil, result.Error
	}
	return trips, nil
}

func (repo *TripsRepository) GetTripsByUserID(userID uint) ([]models.Trip, error) {
	var trips []models.Trip
	result := repo.DB.Table("trips.trips").Where("user_id =?", userID).Find(&trips)

	if result.Error != nil {
		return nil, result.Error
	}
	return trips, nil
}

func (repo *TripsRepository) GetPublicTripsForEveryone(userID uint) ([]models.Trip, error) {
	var trips []models.Trip
	result := repo.DB.Table("trips.trips").
		Where("user_id != ? AND visibility = ?", userID, "PUBLIC").
		Find(&trips)

	if result.Error != nil {
		return nil, result.Error
	}
	return trips, nil
}

func (repo *TripsRepository) GetPublicTripsForUser(userID uint) ([]models.Trip, error) {
	var trips []models.Trip
	result := repo.DB.Table("trips.trips").
		Where("user_id = ? AND visibility = ?", userID, "PUBLIC").
		Find(&trips)

	if result.Error != nil {
		return nil, result.Error
	}
	return trips, nil
}

func (repo *TripsRepository) GetPublicAndFriendsTripsForUser(userID uint) ([]models.Trip, error) {
	var trips []models.Trip
	result := repo.DB.Table("trips.trips").
		Where("user_id != ? AND (visibility = ? OR visibility = ?)", userID, "PUBLIC", "FRIENDS").
		Find(&trips)

	if result.Error != nil {
		return nil, result.Error
	}
	return trips, nil
}
