package cmd

import (
	"context"
	"log"

	"github.com/Techeer-Hogwarts/crawling/config"
	"github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func InitTracer(ctx context.Context) (*trace.TracerProvider, error) {
	otelEndpoint := config.GetEnv("TRACING_GRPC", "otel-collector:4317")
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(otelEndpoint),
	)
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("crawler-worker"),
		)),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}

func ExtractTraceContext(msg amqp091.Delivery) context.Context {
	propagator := otel.GetTextMapPropagator()
	ctx := context.Background()
	carrier := propagation.MapCarrier{}
	log.Printf("Headers: %v", msg.Headers)
	for key, value := range msg.Headers {
		if strValue, ok := value.(string); ok {
			carrier.Set(key, strValue)
		}
	}
	return propagator.Extract(ctx, carrier)
}
