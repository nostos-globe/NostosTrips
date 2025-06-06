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
    fmt.Println("Starting media upload process")

    if err := ctx.Request.ParseMultipartForm(1000 << 20); err != nil {
        fmt.Printf("Failed to parse form: %v\n", err)
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form"})
        return
    }

    tripID, err := strconv.ParseInt(ctx.Param("trip_id"), 10, 64)
    if err != nil {
        fmt.Printf("Error: Invalid trip ID - %v\n", err)
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid trip ID"})
        return
    }
    fmt.Printf("Processing upload for trip ID: %d\n", tripID)

    // Get user ID from authenticated context
    tokenCookie, err := ctx.Cookie("auth_token")
    if err != nil {
        fmt.Println("No auth token found in cookies")
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no token found"})
        return
    }

    userID, err := c.AuthClient.GetUserID(tokenCookie)
    if err != nil || userID == 0 {
        fmt.Printf("Failed to get user ID from token: %v\n", err)
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "failed to authenticate user"})
        return
    }
    fmt.Printf("Authenticated user ID: %d\n", userID)

    // Get form values
    visibility := models.VisibilityEnum(ctx.Request.FormValue("visibility"))
    if visibility == "" {
        visibility = models.Public
        fmt.Println("No visibility specified, defaulting to PUBLIC")
    }
    fmt.Printf("Media visibility set to: %s\n", visibility)

    // Handle file upload
    file, header, err := ctx.Request.FormFile("media")
    if err != nil {
        fmt.Printf("Error: No file provided - %v\n", err)
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "no file provided"})
        return
    }
    defer file.Close()
    fmt.Printf("Processing media file: %s (size: %d bytes)\n", header.Filename, header.Size)

    // Extract metadata
    fmt.Println("Attempting to extract metadata from file")
    metadata, err := c.MediaService.ExtractMetadata(file, header)
    requiresManualLocation := false
    if err != nil {
        if err.Error() == "MANUAL_LOCATION_REQUIRED" {
            requiresManualLocation = true
            fmt.Println("Media requires manual location input")
        } else {
            fmt.Printf("Warning: Failed to extract metadata: %v\n", err)
        }
    } else {
        fmt.Println("Successfully extracted metadata")
    }

    // Upload to MinIO
    fmt.Println("Uploading file to MinIO")
    objectName, err := c.MediaService.UploadMedia(int64(userID), file, header, visibility)
    if err != nil {
        fmt.Printf("Error: Failed to upload media to MinIO - %v\n", err)
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload media"})
        return
    }
    fmt.Printf("Successfully uploaded file. Object name: %s\n", objectName)

    // Create media record
    media := models.Media{
        TripID:       tripID,
        UserID:       int64(userID),
        LocationID:   metadata.LocationID,
        Type:        metadata.Type,
        FilePath:     objectName,
        Visibility:   visibility,
        UploadDate:   time.Now(),
        CaptureDate:  metadata.CaptureDate,
        GpsLatitude:  metadata.Latitude,
        GpsLongitude: metadata.Longitude,
        GpsAltitude:  metadata.Altitude,
    }

    // Save media metadata
    fmt.Println("Saving media metadata to database")
    if err = c.MediaService.SaveMedia(&media); err != nil {
        fmt.Printf("Failed to save media metadata: %v\n", err)
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save media metadata"})
        return
    }

    // Prepare response
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

    if requiresManualLocation {
        fmt.Println("Returning response with manual location flag")
        response["requiresManualLocation"] = true
        ctx.JSON(http.StatusAccepted, response)
        return
    }

    fmt.Println("Media upload completed successfully")
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

func (c *MediaController) ChangeMediaVisibility(ctx *gin.Context) {
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

	// Parse new visibility from request body
	var requestBody struct {
		Visibility models.VisibilityEnum `json:"visibility"`
	}

	if err := ctx.BindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body format"})
		return
	}

	// Update media visibility
	err = c.MediaService.ChangeMediaVisibility(mediaID, int64(userID), requestBody.Visibility)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to change media visibility"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "media visibility changed successfully"})
}

func (c *MediaController) GetMediaVisibility(ctx *gin.Context) {
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

	visibility, err := c.MediaService.GetMediaVisibility(mediaID, uint(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get media visibility"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"visibility": visibility})
}

func (c *MediaController) GetMediaByTripID(ctx *gin.Context) {
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

	media, err := c.MediaService.GetMediaByTripID(tripID, int64(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get media"})
	}

	ctx.JSON(http.StatusOK, media)
}

func (c *MediaController) GetMediaByID(ctx *gin.Context) {
    mediaID, err := strconv.ParseInt(ctx.Param("media_id"), 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid media ID"})
        return
    }

    media, err := c.MediaService.GetMediaByID(mediaID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve media"})
        return
    }

    ctx.JSON(http.StatusOK, media)
}

func (c *MediaController) GetLocationByMediaID(ctx *gin.Context) {
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

	location, err := c.MediaService.GetLocationByMediaID(mediaID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get location"})
		return
	}

	ctx.JSON(http.StatusOK, location)
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
