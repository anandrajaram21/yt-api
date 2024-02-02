package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	db "github.com/anandrajaram21/yt-api/database"
	models "github.com/anandrajaram21/yt-api/models"
	"github.com/joho/godotenv"
	redis "github.com/redis/go-redis/v9"
	"gorm.io/gorm/clause"
)

var (
	ctx = context.Background()
	rdb *redis.Client
)

func initCache() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}

	rdb = redis.NewClient(opt)

	// Ping the Redis server to check connectivity
	_, redisError := rdb.Ping(ctx).Result()
	if redisError != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
}

func setVideoInCache(video models.Video) {
	// Serialize video data (e.g., as JSON)
	jsonData, err := json.Marshal(video)
	if err != nil {
		log.Printf("Error serializing video data: %v", err)
		return
	}

	// Set the data in the cache, consider using video ID or title as the key
	err = rdb.Set(ctx, video.VideoID, jsonData, 0).Err()
	if err != nil {
		log.Printf("Failed to set video in cache: %v", err)
	}
}

func getVideoFromCache(videoId string) *models.Video {
	val, err := rdb.Get(ctx, videoId).Result()
	if err == redis.Nil {
		log.Printf("Video %s does not exist in cache", videoId)
		return nil
	} else if err != nil {
		log.Printf("Error retrieving video from cache: %v", err)
		return nil
	}

	var video models.Video
	err = json.Unmarshal([]byte(val), &video)
	if err != nil {
		log.Printf("Error unmarshaling video data: %v", err)
		return nil
	}

	return &video
}

func main() {
	initCache()
	database := db.SetupDatabase()
	service, err := initYoutubeClient()
	listArr := []string{"id", "snippet"}

	if err != nil {
		log.Fatalf("Youtube Client could not be initialized")
	}

	searchTerm := "football"

	daysOld := 7

	call := service.Search.List(listArr).
		Q(searchTerm).
		MaxResults(20).
		Type("video").
		PublishedAfter(time.Now().AddDate(0, 0, -daysOld).Format(time.RFC3339))

	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error calling Youtube API: %v", err)
	}

	for _, item := range response.Items {
		log.Println(item.Snippet.Title)
	}

	for _, item := range response.Items {
		publishedAtTime, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
		if err != nil {
			log.Fatalf("Error parsing time: %v", err)
		}

		video := models.Video{
			VideoID:      item.Id.VideoId,
			Title:        item.Snippet.Title,
			Description:  item.Snippet.Description,
			PublishDate:  publishedAtTime,
			ThumbnailURL: item.Snippet.Thumbnails.Default.Url,
		}

		cachedVideo := getVideoFromCache(video.VideoID)

		if cachedVideo == nil {
			result := database.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "video_id"}},                                                         // Use VideoID as the conflict target
				DoUpdates: clause.AssignmentColumns([]string{"title", "description", "publish_date", "thumbnail_url"}), // columns that should be updated
			}).Create(&video)

			if result.Error != nil {
				log.Printf("Error upserting video: %v", result.Error)
				continue
			}

			setVideoInCache(video)
		} else {
			log.Printf("Video already exists in cache: %s", video.Title)
		}

	}
}
