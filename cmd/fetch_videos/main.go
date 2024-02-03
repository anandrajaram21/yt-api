package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/anandrajaram21/yt-api/internal/config"
	"github.com/anandrajaram21/yt-api/internal/messaging"
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

func Handler(ctx context.Context) {
	cfg := config.LoadConfig()

	log.Println("YOUTUBE API KEY", os.Getenv("API_KEY"))

	yt := youtube.NewYouTubeService(cfg.YouTubeAPIKey)

	searchQuery := "football"
	maxResults := int64(10)

	videos, err := yt.ListVideos(searchQuery, maxResults)
	if err != nil {
		log.Fatalf("Error fetching videos: %v", err)
	}

	for _, item := range videos.Items {
		videoData := VideoData{
			VideoID:      item.Id.VideoId,
			Title:        item.Snippet.Title,
			Description:  item.Snippet.Description,
			PublishDate:  item.Snippet.PublishedAt,
			ThumbnailURL: item.Snippet.Thumbnails.Default.Url,
		}

		data, err := json.Marshal(videoData)
		if err != nil {
			log.Printf("Error marshalling video data: %v", err)
			continue
		}

		err = messaging.SendMessage(cfg.AWSConfig.SQSUrl, string(data))

		if err != nil {
			log.Printf("Error sending message: %v", err)
			continue
		}

	}
}

func main() {
	lambda.Start(Handler)
}
