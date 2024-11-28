package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/Techeer-Hogwarts/crawling/cmd"
	"github.com/Techeer-Hogwarts/crawling/cmd/blogs"
	"github.com/Techeer-Hogwarts/crawling/cmd/rabbitmq"
	"github.com/Techeer-Hogwarts/crawling/cmd/redisInteractor"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

func main() {
	newConnection := rabbitmq.NewConnection()
	defer newConnection.Close()
	newRedisClient, err := redisInteractor.NewClient()
	if err != nil {
		log.Fatalf("Failed to create a new Redis client: %v", err)
	}
	redisContext := context.Background()
	newChannel := rabbitmq.NewChannel(newConnection)
	defer newChannel.Close()
	queue1 := rabbitmq.DeclareQueue(newChannel, "crawl_queue")
	consumedMessages := rabbitmq.ConsumeMessages(newChannel, queue1.Name)

	const numWorkers = 5 // Number of concurrent workers
	var wg sync.WaitGroup
	wg.Add(numWorkers) // Add the number of workers to the WaitGroup

	for i := 0; i < numWorkers; i++ {
		log.Printf("Starting worker %d", i)
		go func(workerID int) { // Each worker is a goroutine
			defer wg.Done()
			for msg := range consumedMessages { // Continuously process messages
				log.Printf("Worker %d processing message: %s", workerID, msg.Body)
				processMessage(msg, redisContext, newRedisClient)
			}
		}(i)
	}
	wg.Wait()
}

func processMessage(msg amqp091.Delivery, redisContext context.Context, newRedisClient *redis.Client) {
	var blogRequest blogs.BlogRequest
	err := json.Unmarshal(msg.Body, &blogRequest)
	if err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return
	}

	url, host, err := cmd.ValidateAndSanitizeURL(string(blogRequest.Data))
	if err != nil {
		log.Printf("Invalid or unsafe URL: %v", err)
		return
	}
	log.Printf("Processing URL: %v", url)

	blogPosts, err := cmd.CrawlBlog(url, host)
	if err != nil {
		log.Printf("Failed to crawl blog: %v, userID: %v", err, blogRequest.UserID)
		return
	}
	blogPosts.UserID = blogRequest.UserID
	// responseJSON, _ := json.MarshalIndent(blogPosts, "", "  ")
	// fmt.Println(string(responseJSON))

	err = redisInteractor.SetData(redisContext, newRedisClient, msg.MessageId, blogPosts)
	if err != nil {
		log.Printf("Failed to set data: %v", err)
		return
	}

	err = redisInteractor.NotifyCompletion(redisContext, newRedisClient, msg.MessageId)
	if err != nil {
		log.Printf("Failed to notify completion: %v", err)
		return
	}

	log.Printf("Successfully processed and stored blog data. Time: %v", time.Now())
}
