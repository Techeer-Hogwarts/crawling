package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

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
		taskID := string(msg.MessageId)
		taskBody := string(msg.Body)
		fmt.Printf("Received task: %s with task ID: %s\n", taskBody, taskID)

		// Simulate task processing and generate a result
		result := fmt.Sprintf("Processed: %s", taskID)

		// Simulate processing time
		time.Sleep(5 * time.Second)

		// Store the result in Redis and update status to "completed"
		err := rdb.Set(ctx, taskID+":result", result, 0).Err()
		if err != nil {
			log.Fatalf("Failed to store result in Redis: %v", err)
		}

		// Update the task status to "completed"
		err = rdb.Set(ctx, taskID, "completed", 0).Err()
		if err != nil {
			log.Fatalf("Failed to update task status in Redis: %v", err)
		}

		fmt.Printf("Task '%s' completed and stored in Redis\n", taskID)
	}
}
