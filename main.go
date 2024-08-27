package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

// Redis client
var rdb *redis.Client
var ctx = context.Background()

func main() {
	// Initialize Redis client
	rdb = redis.NewClient(&redis.Options{
		Addr: "hogwartsredis:6379", // Redis address
		DB:   0,                    // use default DB
	})

	// Connect to RabbitMQ
	conn, err := amqp091.Dial("amqp://guest:guest@hogwartsmq:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare a queue to receive messages
	q, err := ch.QueueDeclare(
		"crawl_queue", // queue name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Consume messages from the queue
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	// Process each message
	for msg := range msgs {
		taskID := msg.MessageId
		messageContent := string(msg.Body) + " - from Redis"
		fmt.Printf("Received task: %s\n", taskID)

		// Simulate task processing and generate a result
		result := fmt.Sprintf("%s 테스크가 성공적으로 처리 됨. 내용은 -> %s", taskID, messageContent)

		// Simulate processing time
		time.Sleep(5 * time.Second)

		// Store the result and update status in Redis hash
		err := rdb.HSet(ctx, taskID, map[string]interface{}{
			"result":    result,
			"processed": "true",
		}).Err()
		if err != nil {
			log.Fatalf("Failed to store result in Redis: %v", err)
		}

		// Notify completion
		err = rdb.Publish(ctx, "task_completions", taskID).Err()
		if err != nil {
			log.Fatalf("Failed to publish completion message to Redis: %v", err)
		}

		fmt.Printf("Task '%s' completed and stored in Redis\n", taskID)
	}
}
