package service

import (
	"fmt"
	"main/internal/db"
	"strconv"
)

type MediaService struct {
	MediaRepo *db.MediaRepository
}

func (s *MediaService) DeleteMediaByTripID(tripID string) error {
	id, err := strconv.Atoi(tripID)
	if err != nil {
		return err
	}

	result := s.MediaRepo.DeleteMediaByTripID(id)
	if result != nil {
		return result
	}

	return nil
}

func (s *MediaService) GetMediaByID(id int) (any, error) {
	return nil, fmt.Errorf("not implemented")
}
