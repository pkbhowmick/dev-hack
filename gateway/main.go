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

	"github.com/pkbhowmick/dev-hack/gateway/middleware"
	"github.com/pkbhowmick/dev-hack/gateway/presenter/healthcheck"
	"github.com/pkbhowmick/dev-hack/gateway/pubsub"
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
	server, err := newServer()
	if err != nil {
		return err
	}

	log.Println("server is running")
	return server.ListenAndServe()
}

func productHandler(c *gin.Context) {
	if err := pubsub.SendMessage(c.Request.Context(), []byte("Hello World!")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func newServer() (*http.Server, error) {
	// c, err := newDIContainer()
	// if err != nil {
	// 	return nil, err
	// }

	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/", middleware.AddTracing("healthHandler", "healthSpan"), healthcheck.Handler())
	r.POST("/product", middleware.AddTracing("gatewayHandler", "gatewaySpan"), productHandler)

	// err = c.Invoke(func(
	// 	signupuc *signupuc.Usecase,
	// 	productuc *productuc.Usecase,
	// ) {
	// 	r.GET("/", middleware.AddTracing("healthHandler", "healthSpan"), healthcheck.Handler())
	// 	r.POST("/signup", signup.Handler(signupuc))
	// 	r.POST("/users/:userId/products", product.CreationHandler(productuc))
	// 	r.GET("/users/:userId/products", product.ListHandler(productuc))
	// })
	// if err != nil {
	// 	return nil, err
	// }

	return &http.Server{Addr: ":8080", Handler: r}, nil
}
