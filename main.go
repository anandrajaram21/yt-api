package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	youtube "google.golang.org/api/youtube/v3"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("API_KEY")

	ctx := context.Background()

	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	log.Println(youtubeService)
}
