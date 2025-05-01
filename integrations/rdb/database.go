package rdb

import (
	"context"
	"log"

	"sync"

	"github.com/go-redis/redis/v8"
)

var (
	redisClient *redis.Client
	redisOnce   sync.Once
)

func InitRedis() *redis.Client {
	redisOnce.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr:         "localhost:6379",
			Password:     "",
			DB:           0,
			MaxRetries:   3,
			MinIdleConns: 10,
		})

		ctx := context.Background()
		if err := redisClient.Ping(ctx).Err(); err != nil {
			log.Fatalf("Failed to connect to Redis: %v", err)
		}
		log.Println("Successfully connected to Redis")
	})

	return redisClient
}

func GetRedisClient() *redis.Client {
	if redisClient == nil {
		InitRedis()
	}
	return redisClient
}

func CloseRedis() error {
	if redisClient != nil {
		return redisClient.Close()
	}
	return nil
}
