package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	// Connect to RabbitMQ server
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare a queue
	queueName := "your_queue_name"
	_, err = ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Message properties, including HTTP headers
	headers := make(amqp.Table)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Bearer your_access_token"

	// Message content
	message := "Hello, RabbitMQ!"

	// Publish the message with custom headers
	err = ch.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(message),
			DeliveryMode: amqp.Transient,
			Headers:      headers,
		})
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}

	fmt.Println("Message sent successfully!")
}
