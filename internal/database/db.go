package database

import (
	"log"

	"github.com/anandrajaram21/yt-api/internal/config"
	"github.com/anandrajaram21/yt-api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var DB *gorm.DB

func InitializeDatabase(cfg config.DatabaseConfig) {
	var err error
	DB, err = gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = DB.AutoMigrate(&models.Video{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}

func UpsertVideo(video models.Video) error {
	result := DB.Clauses(clause.OnConflict{
		UpdateAll: true, // Updates all fields to new values if conflict occurs
	}).Create(&video)

	if result.Error != nil {
		return result.Error
	}
	return nil
}
