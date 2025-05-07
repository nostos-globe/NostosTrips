package controller

import (
	"fmt"
	"main/internal/models"
	"main/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TripController struct {
	TripService      *service.TripService
	MediaService     *service.MediaService
	AuthClient       *service.AuthClient
	ProfileClient    *service.ProfileClient
	AlbumTripService *service.AlbumsTripsService
	LikesClient      *service.LikesClient  // Add this line
}
func (c *TripController) CreateTrip(ctx *gin.Context) {
	fmt.Printf("Starting CreateTrip request\n")
	
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Visibility  string `json:"visibility"`
		StartDate   string `json:"start_date"`
		EndDate     string `json:"end_date"`
		AlbumID     string   `json:"album_id"` 
	}

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

	if err := ctx.ShouldBindJSON(&req); err != nil {
		fmt.Printf("Error: Failed to bind JSON request - %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tripMapper := &models.TripMapper{}
	trip := tripMapper.ToTripRequest(req, TokenResponse)

	createdTrip, err := c.TripService.CreateTrip(trip)
	if err != nil {
		fmt.Printf("Error: Failed to create trip - %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create trip"})
		return
	}
	trip = createdTrip.(models.Trip)

	if req.AlbumID != 0 {
		albumIDStr := strconv.FormatUint(uint64(req.AlbumID), 10)
        err = c.AlbumTripService.CreateAlbumTrip(albumIDStr, uint(trip.TripID))
        if err != nil {
            fmt.Printf("Error: Failed to create album-trip association - %v\n", err)
        }
	}

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

	media, err := c.MediaService.GetMediaByTripID(int64(trip.TripID), int64(trip.UserID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve media"})
		return
	}

	tripWithMedia := gin.H{
		"trip":  trip,
		"media": media,
	}

	ctx.JSON(http.StatusOK, tripWithMedia)
}

func (c *TripController) SearchTrips(ctx *gin.Context) {
	var searchRequest struct {
		Query string `json:"query"`
	}
	if err := ctx.ShouldBindJSON(&searchRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenCookie, err := ctx.Cookie("auth_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no token found"})
		return
	}

	TokenResponse, err := c.AuthClient.GetUserID(tokenCookie)
	if err != nil || TokenResponse == 0 {
		fmt.Printf("Error: Failed to get user ID - %v\n", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "failed to find this user"})
		return
	}

	trips, err := c.TripService.SearchTrips(searchRequest.Query, TokenResponse)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

func (c *TripController) GetAllTrips(ctx *gin.Context) {
	trips, err := c.TripService.GetAllTrips()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve trips"})
		return
	}

	ctx.JSON(http.StatusOK, trips)
}

/*
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
*/

func (c *TripController) GetPublicTrips(ctx *gin.Context) {
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

	trips, err := c.TripService.GetPublicTripsForEveryone(TokenResponse)
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

func (c *TripController) GetTripsByUserID(ctx *gin.Context) {
	userID := ctx.Param("id")

	// Get user's trips with their associated media
	trips, err := c.TripService.GetTripsByUserID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve user's trips"})
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
func (c *TripController) GetLocationsByTripID(ctx *gin.Context) {
	fmt.Printf("Starting GetLocationsByTripID request\n")

	tripID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	fmt.Printf("Processing request for trip ID: %d\n", tripID)

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
	fmt.Printf("Authenticated user ID: %d\n", TokenResponse)

	// Check if the trip belongs to the user
	trip, err := c.TripService.GetTripByID(strconv.FormatUint(uint64(tripID), 10))
	if err != nil {
		fmt.Printf("Error: Trip not found - %v\n", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "trip not found"})
		return
	}
	fmt.Printf("Found trip with ID: %d\n", trip.TripID)

	media, err := c.MediaService.GetMediaByTripID(int64(trip.TripID), int64(trip.UserID))
	if err != nil {
		fmt.Printf("Error: Failed to retrieve media - %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve media"})
		return
	}
	fmt.Printf("Retrieved %d media items\n", len(media))

	var locations []models.Location
	for _, m := range media {
		fmt.Printf("Fetching location for media ID: %d\n", m.MediaID)
		location, err := c.MediaService.GetLocationByMediaID(m.MediaID)
		if err != nil {
			fmt.Printf("Error: Failed to retrieve location for media ID %d - %v\n", m.MediaID, err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve locations"})
			return
		}
		locations = append(locations, *location)
	}

	fmt.Printf("Successfully retrieved %d locations\n", len(locations))
	ctx.JSON(http.StatusOK, locations)
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
	fmt.Printf("Processing request for user ID: %d\n", userID)

	// Get followed users and followers from Profile service
	followedUsers, err := c.ProfileClient.GetFollowing(tokenCookie, userID)
	if err != nil {
		fmt.Printf("Error: Failed to get followed users - %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve followed users"})
		return
	}
	fmt.Printf("Retrieved %d followed users\n", len(followedUsers))

	followers, err := c.ProfileClient.GetFollowers(tokenCookie, userID)
	if err != nil {
		fmt.Printf("Error: Failed to get followers - %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve followers"})
		return
	}
	fmt.Printf("Retrieved %d followers\n", len(followers))

	// Create a map of mutual follows (friends)
	mutualFollows := make(map[int]bool)
	for _, follower := range followers {
		mutualFollows[follower] = true
	}
	fmt.Printf("Identified %d mutual follows\n", len(mutualFollows))

	var allTripsWithMedia []gin.H
	for _, followedID := range followedUsers {
		fmt.Printf("Processing trips for followed user ID: %d\n", followedID)

		var trips []models.Trip
		if mutualFollows[followedID] {
			fmt.Printf("Fetching public and friends trips for mutual follow %d\n", followedID)
			trips, err = c.TripService.GetPublicAndFriendsTripsForUser(uint(followedID))
		} else {
			fmt.Printf("Fetching public trips for followed user %d\n", followedID)
			trips, err = c.TripService.GetPublicTripsForUser(uint(followedID))
		}

		if err != nil {
			fmt.Printf("Error: Failed to retrieve trips for user %d - %v\n", followedID, err)
			continue
		}
		fmt.Printf("Retrieved %d trips for user %d\n", len(trips), followedID)

		for _, trip := range trips {
			fmt.Printf("Processing trip ID: %d\n", trip.TripID)
			media, err := c.MediaService.GetMediaByTripID(int64(trip.TripID), int64(userID))
			if err != nil {
				fmt.Printf("Error: Failed to retrieve media for trip %d - %v\n", trip.TripID, err)
				continue
			}

			if len(media) == 0 {
				fmt.Printf("Skipping trip %d as it has no media\n", trip.TripID)
				continue
			}
			fmt.Printf("Found %d media items for trip %d\n", len(media), trip.TripID)

			country := "Unknown"
			tripWithMedia := gin.H{
				"trip":    trip,
				"media":   media,
				"user_id": followedID,
				"country": country,
			}
			allTripsWithMedia = append(allTripsWithMedia, tripWithMedia)
		}
	}

	fmt.Printf("Successfully completed GetFollowedUsersTrips request. Returning %d trips\n", len(allTripsWithMedia))
	ctx.JSON(http.StatusOK, allTripsWithMedia)
}

func (c *TripController) GetMyLikedTrips(ctx *gin.Context) {
	fmt.Printf("Starting GetMyLikedTrips request\n")

	// Get auth token
	tokenCookie, err := ctx.Cookie("auth_token")
	if err != nil {
		fmt.Printf("Error: No auth token found - %v\n", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no token found"})
		return
	}

	// Get liked trip IDs
	likedTripIDs, err := c.LikesClient.GetMyLikes(tokenCookie)
	if err != nil {
		fmt.Printf("Error: Failed to get liked trips - %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve liked trips"})
		return
	}

	var tripsWithMedia []gin.H
	for _, tripID := range likedTripIDs {
		// Get trip details
		trip, err := c.TripService.GetTripByID(fmt.Sprintf("%d", tripID))
		if err != nil {
			fmt.Printf("Error: Failed to get trip %d - %v\n", tripID, err)
			continue
		}

		// Get media for trip
		media, err := c.MediaService.GetMediaByTripID(int64(trip.TripID), int64(trip.UserID))
		if err != nil {
			fmt.Printf("Error: Failed to get media for trip %d - %v\n", tripID, err)
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