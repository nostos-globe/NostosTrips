package models

import "time"

type User struct {
	UserID              uint      `gorm:"primaryKey;autoIncrement" json:"user_id"`
	Email               string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	PasswordHash        string    `gorm:"type:varchar(255);not null" json:"password_hash"`
	FailedLoginAttempts uint      `gorm:"not null;default:0" json:"failed_login_attempts"`
	AccountLocked       bool      `gorm:"type:tinyint(1);not null;default:0" json:"account_locked"`
	RegistrationDate    time.Time `gorm:"type:timestamp;not null;default:current_timestamp()" json:"registration_date"`
	ResetToken          string    `gorm:"type:varchar(255)" json:"-"`
	ResetTokenExpiry    time.Time `gorm:"type:timestamp" json:"-"`
}
