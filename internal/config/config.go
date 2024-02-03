package config

import (
	"os"
)

type Config struct {
	RedisURL      string
	YouTubeAPIKey string
	DBConfig      DatabaseConfig
}

type DatabaseConfig struct {
	DatabaseURL string
}

func LoadConfig() Config {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	dbConfig := DatabaseConfig{
		DatabaseURL: getEnv("DATABASE_URL", ""),
	}

	config := Config{
		RedisURL:      getEnv("REDIS_URL", ""),
		YouTubeAPIKey: getEnv("API_KEY", ""),
		DBConfig:      dbConfig,
	}

	return config
}

// getEnv reads an environment variable and returns its value.
// If the variable is not set and a default value is provided, it returns the default.
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = defaultValue
	}
	return value
}
