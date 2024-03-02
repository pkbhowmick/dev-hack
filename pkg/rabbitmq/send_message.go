package rabbitmq

import (
	"context"
	"fmt"

	"github.com/streadway/amqp"
)

func SendMessage(ctx context.Context, message []byte, headers amqp.Table) error {
	// Connect to RabbitMQ server
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare a queue
	queueName := "demo-queue"
	_, err = ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %v", err)
	}

	// Publish the message with custom headers
	err = ch.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         message,
			DeliveryMode: amqp.Transient,
			Headers:      headers,
		})
	if err != nil {
		return fmt.Errorf("failed to publish a message: %v", err)
	}

	fmt.Println("Message sent successfully!")
	return nil
}
