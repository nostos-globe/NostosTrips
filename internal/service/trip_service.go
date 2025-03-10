package service

import (
	"main/internal/db"
	"main/internal/models"
)

type FollowService struct {
	FollowRepo *db.FollowRepository
}

func (s *FollowService) GetFollowByIDs(followerID uint, followedID uint) (any, error) {
	return s.FollowRepo.GetFollowByIDs(followerID, followedID)
}

func (s *FollowService) FollowUser(follow *models.Follow) error {
	return s.FollowRepo.FollowUser(follow)
}

func (s *FollowService) UnFollowUser(follow *models.Follow) error {
	return s.FollowRepo.UnFollowUser(follow)
}

func (s *FollowService) ListFollowers(profileID uint) (any, error) {
	return s.FollowRepo.ListFollowers(profileID)
}

func (s *FollowService) ListFollowing(profileID uint) (any, error) {
	return s.FollowRepo.ListFollowing(profileID)
}
