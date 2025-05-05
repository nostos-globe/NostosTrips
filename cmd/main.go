package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	controller "main/internal/api"
	dbRepo "main/internal/db"
	"main/internal/service"
	"main/pkg/config"
	"main/pkg/db"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or error loading it: %v", err)
	}

	minioManager := config.InitMinIO()
	if minioManager == nil {
		log.Println("Faliing to init MinIO")
	}

	secretsManager := config.GetSecretsManager()
	if secretsManager != nil {
		secrets := secretsManager.LoadSecrets()
		for key, value := range secrets {
			os.Setenv(key, value)
		}
	} else {
		log.Println("Falling back to environment variables")
	}
}

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	database, err := db.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Initialize repositories
	tripRepo := &dbRepo.TripsRepository{DB: database}
	mediaRepo := &dbRepo.MediaRepository{DB: database}
	albumsTripsRepo := &dbRepo.AlbumsTripsRepository{DB: database}

	// Initialize authClient
	authClient := &service.AuthClient{BaseURL: cfg.AuthServiceUrl}
	profileClient := &service.ProfileClient{BaseURL: cfg.ProfileServiceUrl}

	// Initialize MinioService
	minioService := service.NewMinioService()

	// Initialize services
	albumsTripsService := &service.AlbumsTripsService{AlbumsTripsRepo: albumsTripsRepo}
	tripService := &service.TripService{TripRepo: tripRepo}
	mediaService := &service.MediaService{
		MediaRepo:    mediaRepo,
		MinioService: minioService,
	}
	geocodingService := &service.GeocodingService{}

	// Initialize controllers
	tripHandler := &controller.TripController{
		TripService:      tripService,
		MediaService:     mediaService,
		AuthClient:       authClient,
		ProfileClient:    profileClient,
		AlbumTripService: albumsTripsService,
		LikesClient:      &service.LikesClient{BaseURL: "https://actions.nostos-globe.me"},  // Add this line
	}
	mediaHandler := &controller.MediaController{
		MediaService:     mediaService,
		AuthClient:       authClient,
		GeocodingService: geocodingService,
	}

	// Initialize Gin
	r := gin.Default()

	// Trip routes
	api := r.Group("/api/trips")
	{
		api.POST("/", tripHandler.CreateTrip)
		api.GET("/", tripHandler.GetAllTrips)
		api.POST("/search", tripHandler.SearchTrips)
		api.GET("/public", tripHandler.GetPublicTrips)
		api.GET("/myTrips", tripHandler.GetMyTrips)
		api.GET("/following", tripHandler.GetFollowedUsersTrips)
		api.GET("/user/:id", tripHandler.GetTripsByUserID)
		api.GET("/myLikedTrips", tripHandler.GetMyLikedTrips)
		api.GET("/:id", tripHandler.GetTripByID)
		api.GET("/:id/locations", tripHandler.GetLocationsByTripID)
		api.PUT("/update", tripHandler.UpdateTrip)
		api.DELETE("/delete/:id", tripHandler.DeleteTrip)
	}

	// Media routes in separate group
	mediaApi := r.Group("/api/media")
	{
		mediaApi.POST("/trip/:trip_id", mediaHandler.UploadMedia)
		mediaApi.GET("/:media_id", mediaHandler.GetMediaURL)
		mediaApi.DELETE("/:media_id", mediaHandler.DeleteMedia)
		mediaApi.POST("/:media_id/metadata", mediaHandler.AddMetadataToMedia)
		//mediaApi.GET("/:media_id/metadata", mediaHandler.GetMediaMetadata)
		mediaApi.GET("/:media_id/visibility", mediaHandler.GetMediaVisibility)
		mediaApi.PUT("/:media_id/visibility", mediaHandler.ChangeMediaVisibility)
		mediaApi.GET("/:media_id/location", mediaHandler.GetLocationByMediaID)
		mediaApi.GET("/trip/:trip_id", mediaHandler.GetMediaByTripID)
	}

	// Start server
	log.Println("Server running on http://localhost:8084")
	if err := r.Run(":8084"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
