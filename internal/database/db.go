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

	// err = DB.AutoMigrate(&models.Video{})
	// if err != nil {
	// 	log.Fatalf("Failed to migrate database: %v", err)
	// }
}

// SaveVideo inserts a new video record into the database.
func SaveVideo(video *models.Video) error {
	// The clause.OnConflict allows us to define what to do in case of a duplicate entry,
	// here it's set to do nothing, but it can be set to update the record if needed.
	result := DB.Clauses(clause.OnConflict{DoNothing: true}).Create(video)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func ListVideos(pageSize, pageNumber int) ([]models.Video, error) {
	var videos []models.Video

	offset := (pageNumber - 1) * pageSize
	result := DB.Limit(pageSize).
		Offset(offset).
		Order("publish_date desc").
		Find(&videos)

	if result.Error != nil {
		return nil, result.Error
	}
	return videos, nil
}

func GetAllVideos() ([]models.Video, error) {
	var videos []models.Video

	result := DB.Find(&videos)
	if result.Error != nil {
		return nil, result.Error
	}

	return videos, nil
}
