package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Techeer-Hogwarts/crawling/cmd/rabbitmq"
	"github.com/Techeer-Hogwarts/crawling/cmd/redisInteractor"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

var (
	redisClient *redis.Client
	rabbitChan  *amqp091.Channel
)

func main() {
	// Connect to Redis
	var err error
	redisClient, err = redisInteractor.NewClient()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Connect to RabbitMQ
	rabbitConn := rabbitmq.NewConnection()
	defer rabbitConn.Close()
	rabbitChan = rabbitmq.NewChannel(rabbitConn)
	defer rabbitChan.Close()

	// Declare a queue in RabbitMQ
	queue := rabbitmq.DeclareQueue(rabbitChan, "crawl_queue")

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
			return
		}

		var msg Message
		err := json.NewDecoder(r.Body).Decode(&msg)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			log.Printf("Failed to decode request body: %v", err)
			return
		}

		// Encode the message to JSON []byte
		messageBytes, err := json.Marshal(msg)
		if err != nil {
			http.Error(w, "Failed to encode message", http.StatusInternalServerError)
			log.Printf("Failed to encode message to JSON: %v", err)
			return
		}

		// Publish message to RabbitMQ
		err = PublishMessage(rabbitChan, queue.Name, messageBytes)
		if err != nil {
			http.Error(w, "Failed to send message", http.StatusInternalServerError)
			log.Printf("Failed to publish message to RabbitMQ: %v", err)
			return
		}

		log.Printf("Message sent to RabbitMQ: %s", string(messageBytes))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Message sent successfully"))
	})

	// Start the HTTP server
	log.Println("Server is running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

// PublishMessage sends a message to the specified queue using RabbitMQ's default exchange.
func PublishMessage(ch *amqp091.Channel, queueName string, message []byte) error {
	// Put the channel in confirm mode to ensure that all publishings are acknowledged
	if err := ch.Confirm(false); err != nil {
		log.Printf("Channel could not be put into confirm mode: %v", err)
		return err
	}

	// Create a notification listener for undeliverable messages
	returns := ch.NotifyReturn(make(chan amqp091.Return))
	go func() {
		for ret := range returns {
			log.Printf("Message returned: %s", string(ret.Body))
		}
	}()

	err := ch.Publish(
		"",        // exchange - default exchange
		queueName, // routing key - queue name
		true,      // mandatory - ensures message is returned if no queue is bound
		false,     // immediate - ensures the message is returned if no consumer is available
		amqp091.Publishing{
			ContentType: "plain/text", // Using plain text for message content
			Body:        message,      // Using byte array for JSON-encoded message
		},
	)
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
		return err
	}

	// Wait for confirmation of the published message
	confirms := ch.NotifyPublish(make(chan amqp091.Confirmation))
	for confirm := range confirms {
		if confirm.Ack {
			log.Printf("Message delivery confirmed (ack): %d", confirm.DeliveryTag)
			break
		} else {
			log.Printf("Message delivery not confirmed (nack): %d", confirm.DeliveryTag)
			break
		}
	}

	return nil
}
