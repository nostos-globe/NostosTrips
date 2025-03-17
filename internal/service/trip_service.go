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

func (s *TripService) UpdateTrip(trip models.Trip) (any, error) {
	// Convert string ID to integer
	result, err := s.TripRepo.UpdateTrip(trip)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *TripService) DeleteTrip(tripID string) error {
	// Convert string ID to integer
	id, err := strconv.Atoi(tripID)
	if err != nil {
		return err
	}
	err = s.TripRepo.DeleteTrip(id)
	if err != nil {
		return err
	}

	return nil
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

func (s *TripService) GetAllTrips() ([]models.Trip, error) {
	trips, err := s.TripRepo.GetAllTrips()
	if err != nil {
		return nil, err
	}

	return trips, nil
}

func (s *TripService) GetMyTrips(userID uint) ([]models.Trip, error) {
	trips, err := s.TripRepo.GetMyTrips(userID)
	if err != nil {
		return nil, err

	}
	return trips, nil
}
