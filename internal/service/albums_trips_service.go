package service

import (
	"main/internal/db"
	"strconv"
	"fmt"
)

type AlbumsTripsService struct {
	AlbumsTripsRepo *db.AlbumsTripsRepository
}

func (s *AlbumsTripsService) CreateAlbumTrip(albumID string, tripID uint) error {
    // Convert albumID from string to uint
    id, err := strconv.ParseUint(albumID, 10, 32)
    if err != nil {
        return fmt.Errorf("invalid album ID format: %v", err)
    }
    
    return s.AlbumsTripsRepo.CreateAlbumTrip(uint(id), tripID)
}