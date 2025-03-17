package db

import (
	"gorm.io/gorm"
)

type MediaRepository struct {
	DB *gorm.DB
}

func (repo *MediaRepository) DeleteMediaByTripID(tripID int) error {
	result := repo.DB.Table("media.media").Where("trip_id =?", tripID).Delete("media.media")
	if result.Error != nil {
		return result.Error
	}

	return nil
}
