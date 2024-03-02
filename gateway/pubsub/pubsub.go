package pubsub

import (
	"context"
	"fmt"

	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func SendMessage(ctx context.Context, message []byte, isError bool) error {
	_, span := otel.Tracer("send-to-pubsub").Start(ctx, "pubsub-span")
	defer span.End()
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

	// Message properties, including HTTP headers
	headers := make(amqp.Table)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Bearer your_access_token"
	headers["traceparent"] = fmt.Sprintf("00-%s-%s-01", span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String())
	headers["x-error-case"] = isError

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
		return fmt.Errorf("failed to publish a message: %v", err)
	}

	fmt.Println("Message sent successfully!")
	span.SetAttributes(attribute.String("pubsub isSuccess", "true"))
	return nil
}
