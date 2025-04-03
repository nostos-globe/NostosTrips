package db

import (
	"main/internal/models"

	"gorm.io/gorm"
)

type MediaRepository struct {
	DB *gorm.DB
}

func (repo *MediaRepository) UpdateMedia(d int64, media *models.Media) error {
	result := repo.DB.Table("media.media").Where("media_id = ?", d).Updates(media)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *MediaRepository) GetMediaByTripID(tripID int64) ([]*models.Media, error) {
	var media []*models.Media
	result := repo.DB.Table("media.media").Where("trip_id =?", tripID).Find(&media)
	if result.Error != nil {
		return nil, result.Error
	}
	return media, nil
}


func (repo *MediaRepository) SaveMedia(media *models.Media) error {
	result := repo.DB.Table("media.media").Create(media)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *MediaRepository) DeleteMediaByTripID(tripID int) error {
	result := repo.DB.Table("media.media").Where("trip_id =?", tripID).Delete("media.media")
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *MediaRepository) DeleteMedia(mediaID int64) error {
	result := repo.DB.Table("media.media").Where("media_id = ?", mediaID).Delete(&models.Media{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *MediaRepository) GetMediaByID(mediaID int64) (*models.Media, error) {
	var media models.Media
	result := r.DB.Table("media.media").First(&media, mediaID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &media, nil
}

func (r *MediaRepository) AreFriends(userID1, userID2 int64) bool {
	var count int64
	r.DB.Table("friendships").
		Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)",
			userID1, userID2, userID2, userID1).
		Count(&count)
	return count > 0
}

func (r *MediaRepository) GetLocationByCountryAndCity(location *models.Location) (*models.Location, error) {
	var result models.Location
	dbResult := r.DB.Table("locations.locations").
		Where("country = ? AND city = ?", location.Country, location.City).
		First(&result)

	if dbResult.Error != nil {
		return nil, dbResult.Error
	}

	return &result, nil
}

func (r *MediaRepository) SaveLocationInfo(location *models.Location) error {
	result := r.DB.Table("locations.locations").Create(location)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
