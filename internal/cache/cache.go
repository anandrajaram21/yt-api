package cache

import (
	"context"
	"log"
	"time"

	redis "github.com/redis/go-redis/v9"
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

func GetClient() *redis.Client {
	return rdb
}

func Set(key string, value interface{}, expiration time.Duration) error {
	return rdb.Set(ctx, key, value, expiration).Err()
}

func Get(key string) (string, error) {
	return rdb.Get(ctx, key).Result()
}
