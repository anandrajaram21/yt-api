package db

import (
	"log"

	config "github.com/anandrajaram21/yt-api/internal/config"

	models "github.com/anandrajaram21/yt-api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDatabase(cfg config.DatabaseConfig) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run the auto migration
	err = db.AutoMigrate(&models.Video{})
	if err != nil {
		log.Fatalf("Failed to run auto migration: %v", err)
	}

	return db
}
