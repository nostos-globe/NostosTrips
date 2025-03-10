package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost            string
	DBUser            string
	DBPassword        string
	DBName            string
	DBPort            string
	JWTSecret         string
	AuthServiceUrl    string
	ProfileServiceUrl string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	return &Config{
		DBHost:            "localhost",
		DBUser:            "root",
		DBPassword:        "root",
		DBName:            os.Getenv("DB_NAME"),
		DBPort:            os.Getenv("DB_PORT"),
		JWTSecret:         os.Getenv("JWT_SECRET"),
		AuthServiceUrl:    "http://localhost:8081",
		ProfileServiceUrl: "http://localhost:8082",
	}
}
