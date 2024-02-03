package main

import (
	"log"
	"time"

	"github.com/anandrajaram21/yt-api/internal/cache"
	"github.com/anandrajaram21/yt-api/internal/config"
	"github.com/anandrajaram21/yt-api/internal/database"
	"github.com/anandrajaram21/yt-api/internal/models"
	"github.com/anandrajaram21/yt-api/internal/youtube"
	"github.com/aws/aws-lambda-go/lambda"
)

type VideoData struct {
	VideoID      string `json:"video_id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	PublishDate  string `json:"publish_date"`
	ThumbnailURL string `json:"thumbnail_url"`
}

func init() {
	cfg := config.LoadConfig()

	// Initialize Redis
	_, err := cache.InitializeCache(cfg.RedisURL)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Initialize Database
	database.InitializeDatabase(cfg.DBConfig)

	// Initialize YouTube API
	youtube.InitializeYouTubeAPI(cfg.YouTubeAPIKey)
}

func Handler() {
	searchQuery := "football"
	maxResults := int64(20)

	// List videos from YouTube API
	videos, err := youtube.ListVideos(searchQuery, maxResults)
	if err != nil {
		log.Fatalf("Failed to list videos: %v", err)
	}

	for _, item := range videos.Items {
		publishDate, _ := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
		video := models.Video{
			VideoID:      item.Id.VideoId,
			Title:        item.Snippet.Title,
			Description:  item.Snippet.Description,
			PublishDate:  publishDate,
			ThumbnailURL: item.Snippet.Thumbnails.High.Url,
		}

		cachedVideo, err := cache.RetrieveVideo(video.VideoID)
		if err != nil {
			log.Printf("Failed to retrieve video from cache: %v", err)
			continue
		}

		if cachedVideo == nil {
			log.Println("Saving video")
			err = database.SaveVideo(&video)
			if err != nil {
				log.Printf("Failed to save video: %v", err)
				continue
			}

			err = cache.StoreVideo(&video)
			if err != nil {
				log.Printf("Failed to store video in cache: %v", err)
				continue
			}
		}
	}
}

func main() {
	lambda.Start(Handler)
}
