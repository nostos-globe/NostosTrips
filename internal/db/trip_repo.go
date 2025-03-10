package db

import (
	"main/internal/models"

	"gorm.io/gorm"
)

type FollowRepository struct {
	DB *gorm.DB
}

func (repo *FollowRepository) GetFollowByIDs(followerID uint, followedID uint) (any, error) {
	var follow models.Follow
	err := repo.DB.Table("auth.followers").Where("follower_id = ? AND followed_id = ?", followerID, followedID).First(&follow).Error
	if err != nil {
		return nil, err
	}
	return &follow, nil
}

func (repo *FollowRepository) FollowUser(follow *models.Follow) error {
	return repo.DB.Table("auth.followers").Create(follow).Error
}

func (repo *FollowRepository) UnFollowUser(follow *models.Follow) error {
	return repo.DB.Table("auth.followers").Delete(follow).Error
}

func (repo *FollowRepository) ListFollowers(profileID uint) (any, error) {
	var followers int64
	err := repo.DB.Table("auth.followers").Where("followed_id = ?", profileID).Count(&followers).Error
	if err != nil {
		return nil, err
	}
	return followers, nil
}

func (repo *FollowRepository) ListFollowing(profileID uint) (any, error) {
	var following int64
	err := repo.DB.Table("auth.followers").Where("follower_id = ?", profileID).Count(&following).Error
	if err != nil {
		return nil, err
	}
	return following, nil
}
