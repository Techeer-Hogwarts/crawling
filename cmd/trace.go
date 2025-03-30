package cmd

import (
	"context"

	"github.com/Techeer-Hogwarts/crawling/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
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
