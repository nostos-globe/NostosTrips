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
	// Initialize authClient
	authClient := &service.AuthClient{BaseURL: cfg.AuthServiceUrl}
	// Initialize services
	tripService := &service.TripService{TripRepo: tripRepo}

	// Initialize controllers
	tripHandler := &controller.TripController{TripService: tripService, AuthClient: authClient}

	// Initialize Gin
	r := gin.Default()

	// Profile routes
	api := r.Group("/api/trips")
	{
		api.POST("/", tripHandler.CreateTrip)

		api.GET("/", tripHandler.GetAllTrips)
		api.GET("/myTrips", tripHandler.GetMyTrips)
		api.GET("/:id", tripHandler.GetTripByID)

	}

	// Start server
	log.Println("Server running on http://localhost:8083")
	if err := r.Run(":8083"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
