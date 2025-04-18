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
			semconv.ServiceNameKey.String("go-consumer-crawler"),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	return tp, nil
}

func ExtractTraceContext(msg amqp091.Delivery) context.Context {
	propagator := otel.GetTextMapPropagator()
	ctx := context.Background()
	carrier := propagation.MapCarrier{}
	for key, value := range msg.Headers {
		log.Printf("Key: %s, Value: %v (%T)", key, value, value)
		switch v := value.(type) {
		case string:
			carrier.Set(key, v)
		case []byte:
			carrier.Set(key, string(v))
		default:
			log.Printf("Unsupported header type: %T", v)
		}
	}
	ctx = propagator.Extract(ctx, carrier)
	// logTraceContext(ctx)
	return ctx
}

// func logTraceContext(ctx context.Context) {
// 	spanContext := oteltrace.SpanContextFromContext(ctx)
// 	if spanContext.IsValid() {
// 		// Log the trace ID and span ID
// 		log.Printf("Trace ID: %s, Span ID: %s", spanContext.TraceID().String(), spanContext.SpanID().String())
// 	} else {
// 		log.Println("Invalid Span Context")
// 	}
// }
