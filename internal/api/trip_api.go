package controller

import (
	"main/internal/models"
	"main/internal/service"
	"net/http"
	"github.com/gin-gonic/gin"
)

type TripController struct {
	TripService  *service.TripService
	MediaService *service.MediaService
	AuthClient   *service.AuthClient
}

func (c *TripController) CreateTrip(ctx *gin.Context) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Visibility  string `json:"visibility"`
		StartDate   string `json:"start_date"`
		EndDate     string `json:"end_date"`
	}

	// Get user ID from authenticated context
	tokenCookie, err := ctx.Cookie("auth_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no token found"})
		return
	}

	TokenResponse, err := c.AuthClient.GetUserID(tokenCookie)
	if err != nil || TokenResponse == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "failed to find this user"})
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tripMapper := &models.TripMapper{}
	trip := tripMapper.ToTrip(req, TokenResponse)

	createdTrip, err := c.TripService.CreateTrip(trip)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create trip"})
		return
	}
	trip = createdTrip.(models.Trip)

	ctx.JSON(http.StatusCreated, trip)
}

func (c *TripController) UpdateTrip(ctx *gin.Context) {
	var req struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Visibility  string `json:"visibility"`
		StartDate   string `json:"start_date"`
		EndDate     string `json:"end_date"`
	}
	// Get user ID from authenticated context
	tokenCookie, err := ctx.Cookie("auth_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no token found"})
		return
	}

	TokenResponse, err := c.AuthClient.GetUserID(tokenCookie)
	if err != nil || TokenResponse == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "failed to find this user"})
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tripMapper := &models.TripMapper{}
	trip := tripMapper.ToTripUpdate(req, TokenResponse)
	result, err := c.TripService.UpdateTrip(trip)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update trip"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "trip updated successfully", "trip": result})
}

func (c *TripController) DeleteTrip(ctx *gin.Context) {
	// Get user ID from authenticated context
	tokenCookie, err := ctx.Cookie("auth_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no token found"})
		return
	}

	TokenResponse, err := c.AuthClient.GetUserID(tokenCookie)
	if err != nil || TokenResponse == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "failed to find this user"})
		return
	}

	tripID := ctx.Param("id")

	deleteMedia := ctx.DefaultQuery("delete_media", "false")
	shouldDeleteMedia := deleteMedia == "true"

	if shouldDeleteMedia {
		err := c.MediaService.DeleteMediaByTripID(tripID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete media"})
			return
		}
	}

	err = c.TripService.DeleteTrip(tripID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete trip"})
		return
	}
	if shouldDeleteMedia {
		ctx.JSON(http.StatusOK, gin.H{"message": "trip and media deleted successfully"})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "trip deleted successfully"})
	}
}

func (c *TripController) GetTripByID(ctx *gin.Context) {
	tripID := ctx.Param("id")
	trip, err := c.TripService.GetTripByID(tripID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "trip not found"})
		return
	}

	ctx.JSON(http.StatusOK, trip)
}

func (c *TripController) GetAllTrips(ctx *gin.Context) {
	trips, err := c.TripService.GetAllTrips()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve trips"})
		return
	}

	ctx.JSON(http.StatusOK, trips)
}

func (c *TripController) GetAllPublicTrips(ctx *gin.Context) {
	trips, err := c.TripService.GetAllPublicTrips()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve public trips"})
		return
	}

	var tripsWithMedia []gin.H
	for _, trip := range trips {
		media, err := c.MediaService.GetMediaByTripID(int64(trip.TripID), int64(trip.UserID))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve media"})
			return
		}

		tripWithMedia := gin.H{
			"trip":  trip,
			"media": media,
		}
		tripsWithMedia = append(tripsWithMedia, tripWithMedia)
	}

	ctx.JSON(http.StatusOK, tripsWithMedia)
}

func (c *TripController) GetMyTrips(ctx *gin.Context) {
	// Get user ID from authenticated context
	tokenCookie, err := ctx.Cookie("auth_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no token found"})
		return
	}

	TokenResponse, err := c.AuthClient.GetUserID(tokenCookie)
	if err != nil || TokenResponse == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "failed to find this user"})
		return
	}

	// Get trips with their associated media
	trips, err := c.TripService.GetMyTrips(TokenResponse)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve trips"})
		return
	}

	// For each trip, fetch its associated media
	var tripsWithMedia []gin.H
	for _, trip := range trips {
		media, err := c.MediaService.GetMediaByTripID(int64(trip.TripID), int64(TokenResponse))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve media"})
			return
		}

		tripWithMedia := gin.H{
			"trip":  trip,
			"media": media,
		}
		tripsWithMedia = append(tripsWithMedia, tripWithMedia)
	}

	ctx.JSON(http.StatusOK, tripsWithMedia)
}
