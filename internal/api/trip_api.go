package controller

import (
	"fmt"
	"main/internal/models"
	"main/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TripController struct {
	TripService   *service.TripService
	MediaService  *service.MediaService
	AuthClient    *service.AuthClient
	ProfileClient *service.ProfileClient
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

		// Skip trips that have no media
		if len(media) == 0 {
			continue
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
	fmt.Printf("Starting GetMyTrips request\n")

	// Get user ID from authenticated context
	tokenCookie, err := ctx.Cookie("auth_token")
	if err != nil {
		fmt.Printf("Error: No auth token found - %v\n", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no token found"})
		return
	}

	TokenResponse, err := c.AuthClient.GetUserID(tokenCookie)
	if err != nil || TokenResponse == 0 {
		fmt.Printf("Error: Failed to get user ID - %v\n", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "failed to find this user"})
		return
	}

	fmt.Printf("Fetching trips for user ID: %d\n", TokenResponse)

	// Get trips with their associated media
	trips, err := c.TripService.GetMyTrips(TokenResponse)
	if err != nil {
		fmt.Printf("Error: Failed to retrieve trips - %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve trips"})
		return
	}

	fmt.Printf("Successfully retrieved %d trips\n", len(trips))

	// For each trip, fetch its associated media
	var tripsWithMedia []gin.H
	for _, trip := range trips {
		fmt.Printf("Fetching media for trip ID: %d\n", trip.TripID)
		media, err := c.MediaService.GetMediaByTripID(int64(trip.TripID), int64(TokenResponse))
		if err != nil {
			fmt.Printf("Error: Failed to retrieve media for trip ID %d - %v\n", trip.TripID, err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve media"})
			return
		}

		tripWithMedia := gin.H{
			"trip":  trip,
			"media": media,
		}
		tripsWithMedia = append(tripsWithMedia, tripWithMedia)
	}

	fmt.Printf("Successfully completed GetMyTrips request\n")
	ctx.JSON(http.StatusOK, tripsWithMedia)
}

func (c *TripController) GetFollowedUsersTrips(ctx *gin.Context) {
	fmt.Printf("Starting GetFollowedUsersTrips request\n")

	tokenCookie, err := ctx.Cookie("auth_token")
	if err != nil {
		fmt.Printf("Error: No auth token found - %v\n", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no token found"})
		return
	}

	userID, err := c.AuthClient.GetUserID(tokenCookie)
	if err != nil || userID == 0 {
		fmt.Printf("Error: Failed to get user ID - %v\n", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "failed to find this user"})
		return
	}

	// Get followed users and followers from Profile service
	followedUsers, err := c.ProfileClient.GetFollowing(tokenCookie, userID)
	if err != nil {
		fmt.Printf("Error: Failed to get followed users - %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve followed users"})
		return
	}

	fmt.Printf("Success: Following - %v\n", followedUsers)

	followers, err := c.ProfileClient.GetFollowers(tokenCookie, userID)
	if err != nil {
		fmt.Printf("Error: Failed to get followers - %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve followers"})
		return
	}

	fmt.Printf("Success: Followers - %v\n", followers)

	// Create a map of mutual follows (friends)
	mutualFollows := make(map[int]bool)
	for _, follower := range followers {
		mutualFollows[follower] = true
	}

	var allTripsWithMedia []gin.H
	for _, followedID := range followedUsers {
		fmt.Printf("Fetching trips for followed user ID: %d\n", followedID)

		var trips []models.Trip
		if mutualFollows[followedID] {
			// Get both public and friends-only trips for mutual follows
			trips, err = c.TripService.GetPublicAndFriendsTripsForUser(uint(followedID))
			fmt.Printf("Found trips for user %d - %v\n", followedID, trips, err)
		} else {
			// Get only public trips for non-mutual follows
			trips, err = c.TripService.GetPublicTripsForUser(uint(followedID))
			fmt.Printf("Found trips for user %d - %v\n", followedID, trips, err)
		}

		if err != nil {
			fmt.Printf("Error: Failed to retrieve trips for user %d - %v\n", followedID, err)
			continue
		}

		fmt.Printf("Found %d trips for user %d\n", len(trips), followedID)

		for _, trip := range trips {
			fmt.Printf("Fetching media for trip ID: %d\n", trip.TripID)
			media, err := c.MediaService.GetMediaByTripID(int64(trip.TripID), int64(userID))
			if err != nil {
				fmt.Printf("Error: Failed to retrieve media for trip %d - %v\n", trip.TripID, err)
				continue
			}

			// Skip trips with no media
			if len(media) == 0 {
				fmt.Printf("Skipping trip %d as it has no media\n", trip.TripID)
				continue
			}

			mediaData, err := c.MediaService.GetMediaDataByTripID(int64(trip.TripID), int64(followedID))
			if err != nil {
				fmt.Printf("Error: Failed to retrieve media for trip %d - %v\n", trip.TripID, err)
				continue
			}

			// Skip if mediaData is empty
			if len(mediaData) == 0 {
				fmt.Printf("Skipping trip %d as it has no media data\n", trip.TripID)
				continue
			}

			var country string
			for _, m := range mediaData {
				if m.GpsLatitude != 0 && m.GpsLongitude != 0 {
					locationInfo, err := c.MediaService.GetLocationInfo(m.GpsLatitude, m.GpsLongitude)
					if err == nil && locationInfo.Country != "" {
						country = locationInfo.Country
						break
					}
				}
			}

			tripWithMedia := gin.H{
				"trip":    trip,
				"media":   media,
				"user_id": followedID,
				"country": country,
			}
			allTripsWithMedia = append(allTripsWithMedia, tripWithMedia)
		}
	}

	fmt.Printf("Successfully completed GetFollowedUsersTrips request with %d trips\n", len(allTripsWithMedia))
	ctx.JSON(http.StatusOK, allTripsWithMedia)
}
