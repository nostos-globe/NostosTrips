package controller

import (
	"fmt"
	"main/internal/models"
	"main/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type MediaController struct {
	MediaService     *service.MediaService
	AuthClient       *service.AuthClient
	GeocodingService *service.GeocodingService
}

func (c *MediaController) UploadMedia(ctx *gin.Context) {
	tripID, err := strconv.ParseInt(ctx.Param("trip_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid trip ID"})
		return
	}

	tokenCookie, err := ctx.Cookie("auth_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no token found"})
		return
	}

	userID, err := c.AuthClient.GetUserID(tokenCookie)

	if err != nil || userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "failed to authenticate user"})
		return
	}

	file, header, err := ctx.Request.FormFile("media")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "no file provided"})
		return
	}
	defer file.Close()

	visibility := models.VisibilityEnum(ctx.PostForm("visibility"))
	if visibility == "" {
		visibility = models.Public
	}

	// Extract metadata
	metadata, err := c.MediaService.ExtractMetadata(file, header)
	requiresManualLocation := false
	if err != nil {
		if err.Error() == "MANUAL_LOCATION_REQUIRED" {
			requiresManualLocation = true
			fmt.Printf("Media requires manual location input\n")
		} else {
			fmt.Printf("Warning: Failed to extract metadata: %v\n", err)
		}
	}

	objectName, err := c.MediaService.UploadMedia(int64(userID), file, header, visibility)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload media"})
		return
	}

	// Create media record with metadata
	media := models.Media{
		TripID:       tripID,
		UserID:       int64(userID),
		LocationID:   func() int64 { if metadata.LocationID == 0 { return 0 }; return metadata.LocationID }(),
		Type:         func() string { if metadata.Type == "" { return "" }; return metadata.Type }(),
		FilePath:     objectName,
		Visibility:   visibility,
		UploadDate:   time.Now(),
		CaptureDate:  metadata.CaptureDate,
		GpsLatitude:  metadata.Latitude,
		GpsLongitude: metadata.Longitude,
		GpsAltitude:  metadata.Altitude,
	}

	err = c.MediaService.SaveMedia(&media)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save media metadata"})
		return
	}

	response := gin.H{
		"message": "media uploaded successfully",
		"path":    objectName,
		"mediaID": media.MediaID,
		"metadata": gin.H{
			"type":        metadata.Type,
			"captureDate": metadata.CaptureDate,
			"location": gin.H{
				"latitude":  metadata.Latitude,
				"longitude": metadata.Longitude,
				"altitude":  metadata.Altitude,
				"city":      metadata.City,
				"country":   metadata.Country,
			},
		},
	}

	// Add flag for manual location if needed
	if requiresManualLocation {
		response["requiresManualLocation"] = true
		ctx.JSON(http.StatusAccepted, response) // Status 202 Accepted
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *MediaController) GetMediaURL(ctx *gin.Context) {
	mediaID, err := strconv.ParseInt(ctx.Param("media_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid media ID"})
		return
	}

	tokenCookie, err := ctx.Cookie("auth_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no token found"})
		return
	}

	userID, err := c.AuthClient.GetUserID(tokenCookie)
	if err != nil || userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "failed to authenticate user"})
		return
	}

	url, err := c.MediaService.GetMediaURL(mediaID, int64(userID))
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"url": url})
}

func (c *MediaController) AddMetadataToMedia(ctx *gin.Context) {
	mediaID, err := strconv.ParseInt(ctx.Param("media_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid media ID"})
		return
	}

	// Authenticate user
	tokenCookie, err := ctx.Cookie("auth_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no token found"})
		return
	}

	userID, err := c.AuthClient.GetUserID(tokenCookie)
	if err != nil || userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "failed to authenticate user"})
		return
	}

	// Parse metadata from request body
	var metadata struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Altitude  float64 `json:"altitude"`
	}

	if err := ctx.BindJSON(&metadata); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid metadata format"})
		return
	}

	// Update media with new metadata
	err = c.MediaService.UpdateMediaMetadata(mediaID, int64(userID), metadata.Latitude,
		metadata.Longitude, metadata.Altitude)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update metadata"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "metadata updated successfully",
		"metadata": metadata,
	})
}

func (c *MediaController) DeleteMedia(ctx *gin.Context) {
	mediaID, err := strconv.ParseInt(ctx.Param("media_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid media ID"})
		return
	}

	tokenCookie, err := ctx.Cookie("auth_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no token found"})
		return
	}

	userID, err := c.AuthClient.GetUserID(tokenCookie)
	if err != nil || userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "failed to authenticate user"})
		return
	}

	// Delete the media from both MinIO and database
	err = c.MediaService.DeleteMediaCompletely(mediaID, int64(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete media: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "media deleted successfully"})
}
