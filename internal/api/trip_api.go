package controller

import (
	"main/internal/models"
	"main/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TripController struct {
	TripService *service.TripService
	AuthClient  *service.AuthClient
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

	trip := models.Trip{
		Name:        req.Name,
		Description: req.Description,
		Visibility:  req.Visibility,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		UserID:      TokenResponse,
	}

	createdTrip, err := c.TripService.CreateTrip(trip)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create trip"})
		return
	}
	trip = createdTrip.(models.Trip)

	ctx.JSON(http.StatusCreated, trip)
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
