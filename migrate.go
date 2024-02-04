package main

import (
	"log"
	"os"

	"github.com/anandrajaram21/yt-api/internal/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var DB *gorm.DB

	DB, err = gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatalln("Failed to connect to database")
	}

	err = DB.AutoMigrate(&models.Video{})
	if err != nil {
		log.Fatalln("Failed to migrate database")
	}
}
