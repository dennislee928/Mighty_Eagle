package config

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client

// InitRedis initializes the Redis connection
func InitRedis() (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s",
		getEnv("REDIS_HOST", "localhost"),
		getEnv("REDIS_PORT", "6379"),
	)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       0, // use default DB
	})

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	Redis = client
	log.Println("✅ Redis connection established")
	return client, nil
}
