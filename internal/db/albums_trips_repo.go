package db

import (
	"gorm.io/gorm"
)

type AlbumsTripsRepository struct {
    DB *gorm.DB
}

type AlbumTrip struct {
	AlbumID uint `gorm:"primaryKey"`
	TripID  uint `gorm:"primaryKey"`
}

func (repo *AlbumsTripsRepository) CreateAlbumTrip(albumID uint, tripID uint) error {
	albumTrip := AlbumTrip{
		AlbumID: albumID,
		TripID:  tripID,
	}
	
	result := repo.DB.Table("albums.album_trips").Create(&albumTrip)
	return result.Error
}