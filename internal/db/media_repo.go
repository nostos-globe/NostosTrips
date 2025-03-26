package db

import (
	"gorm.io/gorm"
	"main/internal/models"
)

type MediaRepository struct {
	DB *gorm.DB
}

func (repo *MediaRepository) SaveMedia(media *models.Media) error {
	result := repo.DB.Table("media.media").Create(media)
	if result.Error!= nil {
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
    result := r.DB.First(&media, mediaID)
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
