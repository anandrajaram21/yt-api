package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/anandrajaram21/yt-api/internal/cache"
	"github.com/anandrajaram21/yt-api/internal/config"
	"github.com/anandrajaram21/yt-api/internal/database"
	"github.com/anandrajaram21/yt-api/internal/messaging"
	"github.com/anandrajaram21/yt-api/internal/models"
	"github.com/aws/aws-lambda-go/lambda"
)

var ctx = context.Background()

func init() {
	cfg := config.LoadConfig()

	// Initialize the cache
	_, err := cache.InitializeCache(cfg.RedisURL)
	if err != nil {
		log.Fatalf("Failed to initialize cache: %v", err)
	}

	// Initialize the database
	database.InitializeDatabase(cfg.DBConfig)
}

func Handler(ctx context.Context) {
	cfg := config.LoadConfig()

	// Continuously receive and process messages
	for {
		message, err := messaging.ReceiveMessage(cfg.AWSConfig.SQSUrl)
		if err != nil {
			log.Printf("Failed to receive messages: %v", err)
			continue
		}

		if message == nil {
			// No message received, continue to check for new messages
			continue
		}

		var video models.Video
		err = json.Unmarshal([]byte(*message.Body), &video)
		if err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			continue
		}

		// Upsert the video data to the database
		err = database.UpsertVideo(video)
		if err != nil {
			log.Printf("Failed to upsert video to database: %v", err)
			continue
		}

		// Also, write to Redis cache
		videoJSON, err := json.Marshal(video)
		if err != nil {
			log.Printf("Error marshalling video: %v", err)
			continue
		}

		err = cache.Set(video.VideoID, videoJSON, 24*time.Hour) // Set an expiration as needed
		if err != nil {
			log.Printf("Failed to set video in cache: %v", err)
			continue
		}

		log.Printf("Video %s upserted to DB and cached", video.VideoID)
	}
}

func main() {
	lambda.Start(Handler)
}
