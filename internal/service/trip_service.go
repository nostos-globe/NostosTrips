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

func (s *TripService) GetTripByID(tripID string) (models.Trip, error) {
	// Convert string ID to integer
	id, err := strconv.Atoi(tripID)
	if err != nil {
		return models.Trip{}, err
	}

	trip, err := s.TripRepo.GetTripByID(id)
	if err != nil {
		return models.Trip{}, err
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

func (s *TripService) GetAllPublicTrips() ([]models.Trip, error) {
    trips, err := s.TripRepo.GetAllPublicTrips()
    if err != nil {
        return nil, err
    }
    return trips, nil
}

func (s *TripService) GetPublicTripsForEveryone(userID uint) ([]models.Trip, error) {
    return s.TripRepo.GetPublicTripsForEveryone(userID)
}

func (s *TripService) GetPublicTripsForUser(userID uint) ([]models.Trip, error) {
    return s.TripRepo.GetPublicTripsForUser(userID)
}

func (s *TripService) GetPublicAndFriendsTripsForUser(userID uint) ([]models.Trip, error) {
    return s.TripRepo.GetPublicAndFriendsTripsForUser(userID)
}

func (s *TripService) GetTripsByUserID(userID string) ([]models.Trip, error) {

	id, err := strconv.Atoi(userID)
	if err!= nil {
		return nil, err
	}

	return s.TripRepo.GetTripsByUserID(uint(id))
}

func (s *TripService) SearchTrips(query string, userID uint) ([]models.Trip, error) {
    return s.TripRepo.SearchTrips(query, userID)
}