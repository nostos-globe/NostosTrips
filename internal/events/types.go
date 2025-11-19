package events

import "time"

type TripCreatedEvent struct {
	TripID    int       `json:"tripId"`
	OwnerID   uint      `json:"ownerId"`
	CreatedAt time.Time `json:"createdAt"`
}

type TripUpdatedEvent struct {
	TripID    int       `json:"tripId"`
	OwnerID   uint      `json:"ownerId"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type TripDeletedEvent struct {
	TripID    int       `json:"tripId"`
	OwnerID   uint      `json:"ownerId"`
	DeletedAt time.Time `json:"deletedAt"`
}

type MediaUploadedEvent struct {
	MediaID    int64     `json:"mediaId"`
	TripID     int64     `json:"tripId"`
	UserID     int64     `json:"userId"`
	Type       string    `json:"type"` // foto, video...
	UploadedAt time.Time `json:"uploadedAt"`
}
