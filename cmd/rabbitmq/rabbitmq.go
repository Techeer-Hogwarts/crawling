package rabbitmq

import (
	"fmt"
	"log"

	"github.com/Techeer-Hogwarts/crawling/config"
	"github.com/rabbitmq/amqp091-go"
)

func NewConnection() *amqp091.Connection {
	// Connect to RabbitMQ
	user := config.GetEnv("RABBITMQ_USER", "guest")
	password := config.GetEnv("RABBITMQ_PASSWORD", "guest")
	host := config.GetEnv("RABBITMQ_HOST", "localhost")
	conn, err := amqp091.Dial(fmt.Sprintf("amqp://%s:%s@%s:5672/", user, password, host))
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	log.Println("Connected to RabbitMQ")
	// defer conn.Close()
	return conn
}

func NewChannel(conn *amqp091.Connection) *amqp091.Channel {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	log.Println("Opened a channel")
	return ch
}

func DeclareQueue(ch *amqp091.Channel, name string) amqp091.Queue {
	queue, err := ch.QueueDeclare(
		name,  // queue name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}
	log.Println("Declared a queue")
	return queue
}

func ConsumeMessages(ch *amqp091.Channel, queue string) <-chan amqp091.Delivery {
	msgs, err := ch.Consume(
		queue, // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}
	log.Println("Consumed messages")
	return msgs
}
