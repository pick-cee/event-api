package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pick-cee/events-api/internal/database"
)

// Set a value in cache
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return database.RedisClient.Set(ctx, key, data, expiration).Err()
}

// Get a value from cache
func Get(ctx context.Context, key string, dest interface{}) error {
	data, err := database.RedisClient.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dest)
}

// Delete a key
func Delete(ctx context.Context, key ...string) error {
	return database.RedisClient.Del(ctx, key...).Err()
}

// Check if a key exists
func Exists(ctx context.Context, key string) (bool, error) {
	count, err := database.RedisClient.Exists(ctx, key).Result()
	return count > 0, err
}
