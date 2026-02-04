package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pick-cee/events-api/internal/config"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func ConnectRedis(cfg *config.Config) error {
	opt, err := redis.ParseURL(cfg.RedisURL)

	if err != nil {
		return fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	RedisClient = redis.NewClient(opt)

	// Verify connection with ping
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("✅ Connected to Redis")
	return nil
}

func GetClient() (*redis.Client, error) {
	if RedisClient == nil {
		return nil, fmt.Errorf("Redis client not initialized")
	}
	return RedisClient, nil
}

func DisconnectRedis() error {
	if RedisClient == nil {
		return nil
	}

	if err := RedisClient.Close(); err != nil {
		return fmt.Errorf("failed to close Redis connection: %w", err)
	}

	log.Println("🔌 Disconnected from Redis")
	return nil
}
