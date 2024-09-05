package redisInteractor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Techeer-Hogwarts/crawling/cmd"
	"github.com/Techeer-Hogwarts/crawling/config"
	"github.com/redis/go-redis/v9"
)

func NewClient() (*redis.Client, error) {
	host := config.GetEnv("REDIS_HOST", "localhost")
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:6379", host), // Redis address
		DB:   0,                            // use default DB
	})
	return rdb, nil
}

func SetData(ctx context.Context, rdb *redis.Client, key string, value cmd.BlogResponse) error {
	jsonValue, err := json.Marshal(value)
	err = rdb.HSet(ctx, key, map[string]interface{}{
		"result":    string(jsonValue),
		"processed": "true",
	}).Err()
	if err != nil {
		log.Printf("Failed to set data: %v", err)
		return err
	}
	return nil
}

func NotifyCompletion(ctx context.Context, rdb *redis.Client, key string) error {
	err := rdb.Publish(ctx, "task_completions", key).Err()
	if err != nil {
		log.Fatalf("Failed to publish completion message to Redis: %v", err)
	}
	return nil
}
