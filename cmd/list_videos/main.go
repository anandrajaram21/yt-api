package main

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/anandrajaram21/yt-api/internal/cache"
	"github.com/anandrajaram21/yt-api/internal/config"
	"github.com/anandrajaram21/yt-api/internal/database"
	"github.com/anandrajaram21/yt-api/internal/youtube"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

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

func Handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	pageSize := request.QueryStringParameters["pageSize"]
	pageNumber := request.QueryStringParameters["pageNumber"]

	// Checking for the presence of the pageSize and pageNumber query parameters
	if pageSize == "" || pageNumber == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "pageSize and pageNumber query parameters are required",
		}, nil
	}

	// Convert the pageSize and pageNumber strings to integers
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		mes := `{ "message": "pageSize must be a number" }`
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       mes,
		}, nil
	}

	pageNumberInt, err := strconv.Atoi(pageNumber)
	if err != nil {
		mes := `{ "message": "pageNumber must be a number" }`
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       mes,
		}, nil
	}

	// Perform the pagination logic
	videos, err := database.ListVideos(pageSizeInt, pageNumberInt)
	if err != nil {
		mes := `{ "message": "Failed to list videos" }`
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       mes,
		}, nil
	}

	// Convert the videos to JSON
	videosJSON, err := json.Marshal(videos)
	if err != nil {
		mes := `{ "message": "Failed to marshal videos to JSON" }`
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       mes,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(videosJSON),
	}, nil
}

func main() {
	lambda.Start(Handler)
}
