package redisInteractor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Techeer-Hogwarts/crawling/cmd/blogs"
	"github.com/Techeer-Hogwarts/crawling/config"
	"github.com/redis/go-redis/v9"
)

func NewClient() (*redis.Client, error) {
	host := config.GetEnv("REDIS_HOST", "localhost")
	port := config.GetEnv("REDIS_PORT", "6379")
	password := config.GetEnv("REDIS_PASSWORD", "test")
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       0,
	})
	return rdb, nil
}

func SetData(ctx context.Context, rdb *redis.Client, key string, value blogs.BlogResponse) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		log.Printf("Failed to marshal data: %v", err)
		return err
	}
	err = rdb.HSet(ctx, key, map[string]interface{}{
		"result":    string(jsonValue),
		"processed": "true",
		"userId":    value.UserID,
	}).Err()
	if err != nil {
		log.Printf("Failed to set data: %v", err)
		return err
	}
	log.Printf("Successfully set data for key: %s", key)
	return nil
}

func NotifyCompletion(ctx context.Context, rdb *redis.Client, key, msgType string) error {
	err := rdb.Publish(ctx, msgType, key).Err()
	if err != nil {
		log.Printf("Failed to publish completion message to Redis: %v", err)
		return err
	}
	return nil
}
