package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/pkbhowmick/dev-hack/pkg/rabbitmq"
)

type TraceConfig struct {
	Protocal   string
	Endpoint   string
	SampleRate float64
}

func getOTELTraceClient(config TraceConfig) (otlptrace.Client, error) {
	switch config.Protocal {
	case "http":
		return otlptracehttp.NewClient(
			otlptracehttp.WithInsecure(),
			otlptracehttp.WithEndpoint(config.Endpoint),
		), nil
	case "grpc":
		return otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(config.Endpoint),
		), nil
	default:
		return nil, errors.New("unknown remote tracing protocal")
	}
}

func getRemoteTraceProvider(conf TraceConfig) (*sdktrace.TracerProvider, error) {
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String("devops-hackathon-project"),
			semconv.DeploymentEnvironmentKey.String("production"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create the otel resource, error: %s", err.Error())
	}

	traceClient, err := getOTELTraceClient(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to create the remote trace client, error: %s", err.Error())
	}

	traceExporter, err := otlptrace.New(context.Background(), traceClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create the remote trace exporter, error: %s", err.Error())
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	return sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(conf.SampleRate))),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	), nil
}

func initTracing(conf TraceConfig) (*sdktrace.TracerProvider, error) {
	return getRemoteTraceProvider(conf)
}

func main() {
	conf := TraceConfig{
		Protocal:   "grpc",
		Endpoint:   "localhost:4317",
		SampleRate: 1.0,
	}
	tp, err := initTracing(conf)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			slog.Error("Failed to shutdown tracing")
		}
	}()
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	slog.Info(fmt.Sprintf("Tracing enabled with config: endpoint: %s, protocal: %s, sample rate: %f", conf.Endpoint, conf.Protocal, conf.SampleRate))

	// Message properties, including HTTP headers
	headers := make(amqp.Table)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Bearer your_access_token"
	headers["traceparent"] = "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"

	ctx, span := otel.Tracer("test").Start(context.Background(), "testSpan")

	message := []byte("Hello from Service 1")
	if err := rabbitmq.SendMessage(ctx, message, headers); err != nil {
		panic(err)
	}
	span.End()

	select {}
}
