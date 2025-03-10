package service

import (
	"main/internal/db"
	"main/internal/models"
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
