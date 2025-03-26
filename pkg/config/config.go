package config

import (
	"log"

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
		DBName:            "nostos",
		DBPort:            "5432",
		JWTSecret:         "SECRET",
		AuthServiceUrl:    "http://localhost:8082",
		ProfileServiceUrl: "http://localhost:8083",
	}
}
