package pubsub

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func ListenForMessage() error {
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

	// Consume messages
	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %v", err)
	}

	slog.Info("listening for messages")
	// Process incoming messages
	for msg := range msgs {
		// Access custom headers
		contentType, ok := msg.Headers["Content-Type"].(string)
		if !ok {
			contentType = "unknown"
		}

		authorization, ok := msg.Headers["Authorization"].(string)
		if !ok {
			authorization = "unknown"
		}

		traceParent, ok := msg.Headers["traceparent"].(string)
		if !ok {
			traceParent = "unknown"
		}

		fmt.Printf("Received a message: %s\n", msg.Body)
		fmt.Printf("Content-Type: %s\n", contentType)
		fmt.Printf("Authorization: %s\n", authorization)
		fmt.Printf("TraceContext: %s\n", traceParent)

		req, err := http.NewRequest("", "", nil)
		if err != nil {
			slog.Error(err.Error())
		}
		req.Header.Add("traceparent", traceParent)

		prop := otel.GetTextMapPropagator()
		propCtx := prop.Extract(context.TODO(), propagation.HeaderCarrier(req.Header))

		_, span := otel.Tracer("product-service").Start(propCtx, "product-span")

		time.Sleep(time.Second * 3)

		span.End()
	}
	return nil
}
