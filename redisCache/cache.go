package redisCache

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// SetCacheRedis sets a value in Redis with the specified key and expiration time.
func SetCacheRedis(ctx context.Context, redisClient *redis.Client, key string, value any, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		log.Printf("Failed to marshal data for caching: %v", err)
		return err
	}

	if err := redisClient.Set(ctx, key, data, expiration).Err(); err != nil {
		log.Printf("Failed to set data in cache: %v", err)
		return err
	}

	return nil
}

// GetCacheRedis retrieves a value from Redis and unmarshals it into the provided destination.
func GetCacheRedis(ctx context.Context, redisClient *redis.Client, key string, dest any) error {
	data, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			// Key does not exist
			return nil
		}
		log.Printf("Failed to get data from cache: %v", err)
		return err
	}

	if err := json.Unmarshal([]byte(data), dest); err != nil {
		log.Printf("Failed to unmarshal cached data: %v", err)
		return err
	}

	return nil
}
