package db

import (
    "fmt"
    "log"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"

    "main/pkg/config"
)

func ConnectDB(cfg *config.Config) (*gorm.DB, error) {
    log.Printf("Attempting to connect to database: host=%s user=%s dbname=%s port=%s",
        cfg.DBHost, cfg.DBUser, cfg.DBName, cfg.DBPort)

    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Madrid",
        cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    log.Println("Database connected successfully.")
    return db, nil
}
