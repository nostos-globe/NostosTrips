package models

import "time"

type Profile struct {
	ProfileID uint `gorm:"primaryKey;column:profile_id"`

	UserID uint `gorm:"column:user_id;uniqueIndex:idx_user_profile"`
	User   User `gorm:"references:UserID"`

	Username        string     `gorm:"size:100;unique;not null"`
	Bio             *string    `gorm:"type:text"`
	ProfilePicture  *string    `gorm:"size:255"`
	Theme           *string    `gorm:"size:50"`
	Location        *string    `gorm:"size:255"`
	Website         *string    `gorm:"size:255"`
	Birthdate       *string    `gorm:"type:date"`
	Language        *string    `gorm:"size:10"`
	PrivacySettings *string    `gorm:"type:jsonb"`
	UpdatedAt       *time.Time `gorm:"type:timestamp;default:current_timestamp"`
	CreatedAt       *time.Time `gorm:"type:timestamp;default:current_timestamp"`
}
