package service

import (
	"main/internal/db"
	"main/internal/models"
	"strconv"
)

type TripService struct {
	TripRepo *db.TripsRepository
}

func (s *TripService) CreateTrip(trip models.Trip) (any, error) {
	result, err := s.TripRepo.CreateTrip(trip)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *TripService) GetTripByID(tripID string) (any, error) {
	// Convert string ID to integer
	id, err := strconv.Atoi(tripID)
	if err != nil {
		return nil, err
	}

	trip, err := s.TripRepo.GetTripByID(id)
	if err != nil {
		return nil, err
	}

	return trip, nil
}
