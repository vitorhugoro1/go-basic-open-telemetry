package main

import (
	"context"

	"log"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

func main() {
	// Inicializamos o TracerProvider para que possamos integrar ao nosso GracefulShutdown
	tp, err := setupOpenTelemetry("local", "go-otel-basic")

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	r := gin.New()
	// Injetamos o Middleware que vai fazer a instrumentação dos nossos requests
	r.Use(otelgin.Middleware("http-server"))

	r.GET("/up", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"hello": "world",
		})
	})

	r.Run()
}

func setupOpenTelemetry(environment string, serviceName string) (*sdktrace.TracerProvider, error) {
	var tp *sdktrace.TracerProvider
	var err error
	ctx := context.Background()

	tp, err = setupProvider(ctx, environment, serviceName)
	if err != nil {
		return nil, err
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
	)

	return tp, nil
}

func setupProvider(ctx context.Context, environment string, serviceName string) (*sdktrace.TracerProvider, error) {
	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)

	if err != nil {
		return nil, err
	}

	var traceExporter sdktrace.SpanExporter

	// Configuração do Exporter em gRPC, que é um requisito padrão em nossa arquitetura com OpenTelemetry
	traceExporter, err = otlptracegrpc.New(ctx)
	if err != nil {
		return nil, err
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)

	tp := sdktrace.NewTracerProvider(
		// Configuração do SampleRate via ParentBasedTraceIDRatio
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.DeploymentEnvironmentKey.String(environment),
		)),
	)

	return tp, nil
}
