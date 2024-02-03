package cache

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/anandrajaram21/yt-api/internal/models"
	"github.com/redis/go-redis/v9"
)

var (
	ctx = context.Background()
	rdb *redis.Client
)

func InitializeCache(redisUrl string) (*redis.Client, error) {
	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}

	rdb = redis.NewClient(opt)

	// Ping the Redis server to check connectivity
	_, redisError := rdb.Ping(ctx).Result()
	if redisError != nil {
		return nil, err
	}

	return rdb, nil
}

// StoreVideo stores a Video object in Redis cache.
func StoreVideo(video *models.Video) error {
	videoJSON, err := json.Marshal(video)
	if err != nil {
		return err
	}

	// Use the video ID as the key
	return rdb.Set(ctx, video.VideoID, videoJSON, 24*time.Hour).Err() // Setting a 24-hour TTL for the cache
}

// RetrieveVideo retrieves a Video object from Redis cache by its ID.
func RetrieveVideo(videoID string) (*models.Video, error) {
	result, err := rdb.Get(ctx, videoID).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	} else if err != nil {
		return nil, err
	}

	var video models.Video
	err = json.Unmarshal([]byte(result), &video)
	if err != nil {
		return nil, err
	}

	return &video, nil
}
