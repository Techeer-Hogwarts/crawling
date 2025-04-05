package main

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/Techeer-Hogwarts/crawling/cmd"
	"github.com/Techeer-Hogwarts/crawling/cmd/blogs"
	"github.com/Techeer-Hogwarts/crawling/cmd/rabbitmq"
	"github.com/Techeer-Hogwarts/crawling/cmd/redisInteractor"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	ctx := context.Background()
	tracerProvider, err := cmd.InitTracer(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize tracing: %v", err)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Fatalf("트레이서 종료 중 오류 발생: %v", err)
		}
	}()

	newConnection := rabbitmq.NewConnection()
	defer newConnection.Close()

	newRedisClient, err := redisInteractor.NewClient()
	if err != nil {
		log.Fatalf("Failed to create a new Redis client: %v", err)
	}
	defer newRedisClient.Close()

	newChannel := rabbitmq.NewChannel(newConnection)
	defer newChannel.Close()

	queue1 := rabbitmq.DeclareQueue(newChannel, "crawl_queue")
	consumedMessages := rabbitmq.ConsumeMessages(newChannel, queue1.Name)

	const numWorkers = 5 // Number of concurrent workers
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	tracer := otel.Tracer("worker")

	for i := 0; i < numWorkers; i++ {
		log.Printf("Starting worker %d", i)
		go func(workerID int) {
			defer wg.Done()
			for msg := range consumedMessages {
				msgctx := cmd.ExtractTraceContext(msg)
				ctx, span := tracer.Start(msgctx, "ReceiveMessage",
					trace.WithAttributes(attribute.String("worker_id", strconv.Itoa(workerID))))
				log.Printf("Worker %d processing message: %s", workerID, msg.Body)
				processMessage(ctx, msg, newRedisClient)
				span.End()
			}
		}(i)
	}
	wg.Wait()
}

func processMessage(ctx context.Context, msg amqp091.Delivery, newRedisClient *redis.Client) {
	tracer := otel.Tracer("worker")
	ctx, span := tracer.Start(ctx, "DecodeAndValidate",
		trace.WithAttributes(attribute.String("message_id", msg.MessageId)),
	)
	defer span.End()
	var blogRequest blogs.BlogRequest
	err := json.Unmarshal(msg.Body, &blogRequest)
	if err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return
	}

	// signUp_blog_fetch, blogs_daily_update, shared_post_fetch
	crawlingType := msg.Type

	url, host, err := cmd.ValidateAndSanitizeURL(string(blogRequest.Data))
	if err != nil {
		log.Printf("Invalid or unsafe URL: %v", err)
		return
	}
	log.Printf("Processing URL: %v", url)

	blogRequest.UserID = cmd.ExtractUserID(msg.MessageId)
	ctx, crawlSpan := tracer.Start(ctx, "CrawlBlog",
		trace.WithAttributes(attribute.String("url", url), attribute.String("host", host)),
	)
	blogPosts, err := cmd.CrawlBlog(url, host, crawlingType)
	crawlSpan.End()

	if err != nil {
		log.Printf("Failed to crawl blog: %v, userID: %v", err, blogRequest.UserID)
		return
	}
	log.Printf("User ID: %v, Blog Posts: %v", blogRequest.UserID, blogPosts)
	blogPosts.UserID = blogRequest.UserID
	// responseJSON, _ := json.MarshalIndent(blogPosts, "", "  ")
	// fmt.Println(string(responseJSON))

	ctx, redisSpan := tracer.Start(ctx, "SetRedisData",
		trace.WithAttributes(attribute.String("message_id", msg.MessageId)),
	)
	err = redisInteractor.SetData(ctx, newRedisClient, msg.MessageId, blogPosts)
	defer redisSpan.End()
	if err != nil {
		redisSpan.RecordError(err)
		redisSpan.SetStatus(codes.Error, err.Error())
		log.Printf("Failed to set data: %v", err)
		return
	}

	ctx, notifySpan := tracer.Start(ctx, "NotifyCompletion",
		trace.WithAttributes(attribute.String("message_id", msg.MessageId)),
	)
	err = redisInteractor.NotifyCompletion(ctx, newRedisClient, msg.MessageId, crawlingType)
	notifySpan.End()

	if err != nil {
		log.Printf("Failed to notify completion: %v", err)
		return
	}

	log.Printf("Successfully processed and stored blog data. Time: %v", time.Now())
}
