package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/pkbhowmick/dev-hack/product/presenter/healthcheck"
	"github.com/pkbhowmick/dev-hack/product/presenter/product"
	"github.com/pkbhowmick/dev-hack/product/pubsub"
	productuc "github.com/pkbhowmick/dev-hack/product/usecase/product"
)

type TraceConfig struct {
	Protocal    string
	Endpoint    string
	SampleRate  float64
	ServiceName string
	Environment string
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
			semconv.ServiceNameKey.String(conf.ServiceName),
			semconv.DeploymentEnvironmentKey.String(conf.Environment),
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
		Protocal:    "grpc",
		Endpoint:    "localhost:4317",
		SampleRate:  1.0,
		ServiceName: "devops-hackathon",
		Environment: "production",
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

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	go func() {
		if err := pubsub.ListenForMessage(); err != nil {
			slog.Error(err.Error())
		}
	}()

	server, err := newServer()
	if err != nil {
		return err
	}

	log.Println("server is running...")
	return server.ListenAndServe()
}

func newServer() (*http.Server, error) {
	c, err := newDIContainer()
	if err != nil {
		return nil, err
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	err = c.Invoke(func(
		productUC *productuc.Usecase,
	) {
		r.GET("/", healthcheck.Handler())

		r.GET("/products/:userId", product.GetListHandler(productUC))
		r.POST("/products", product.CreationHandler(productUC))
	})
	if err != nil {
		return nil, err
	}

	return &http.Server{Addr: ":8080", Handler: r}, nil
}
