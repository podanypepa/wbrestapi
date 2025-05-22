package main

import (
	"fmt"
	"os"

	"github.com/podanypepa/wbrestapi/pkg/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSL"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to database: %w", err)
	}

	if err = db.AutoMigrate(&repository.User{}); err != nil {
		return nil, fmt.Errorf("AutoMigrate failed: %w", err)
	}

	return db, nil
}
