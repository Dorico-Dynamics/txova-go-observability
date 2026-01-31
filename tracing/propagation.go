package tracing

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/propagation"
)

const (
	// HeaderTraceParent is the W3C trace context header.
	HeaderTraceParent = "traceparent"

	// HeaderTraceState is the W3C trace state header.
	HeaderTraceState = "tracestate"

	// HeaderRequestID is the correlation ID header.
	HeaderRequestID = "X-Request-ID"
)

// HTTPCarrier wraps http.Header to implement propagation.TextMapCarrier.
type HTTPCarrier http.Header

// Get returns the value associated with the passed key.
func (c HTTPCarrier) Get(key string) string {
	return http.Header(c).Get(key)
}

// Set stores the key-value pair.
func (c HTTPCarrier) Set(key, value string) {
	http.Header(c).Set(key, value)
}

// Keys returns all keys in the carrier.
func (c HTTPCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}
	return keys
}

// KafkaCarrier wraps Kafka headers to implement propagation.TextMapCarrier.
type KafkaCarrier map[string]string

// Get returns the value associated with the passed key.
func (c KafkaCarrier) Get(key string) string {
	return c[key]
}

// Set stores the key-value pair.
func (c KafkaCarrier) Set(key, value string) {
	c[key] = value
}

// Keys returns all keys in the carrier.
func (c KafkaCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}
	return keys
}

// Extract extracts trace context from HTTP headers.
func Extract(ctx context.Context, headers http.Header) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, HTTPCarrier(headers))
}

// Inject injects trace context into HTTP headers.
func Inject(ctx context.Context, headers http.Header) {
	otel.GetTextMapPropagator().Inject(ctx, HTTPCarrier(headers))
}

// ExtractFromKafka extracts trace context from Kafka headers.
func ExtractFromKafka(ctx context.Context, headers map[string]string) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, KafkaCarrier(headers))
}

// InjectToKafka injects trace context into Kafka headers.
// If headers is nil, returns early to prevent panic.
func InjectToKafka(ctx context.Context, headers map[string]string) {
	if headers == nil {
		return
	}
	otel.GetTextMapPropagator().Inject(ctx, KafkaCarrier(headers))
}

// ExtractRequestID extracts the X-Request-ID from HTTP headers.
// If not present, returns empty string.
func ExtractRequestID(headers http.Header) string {
	return headers.Get(HeaderRequestID)
}

// InjectRequestID injects the X-Request-ID into HTTP headers.
func InjectRequestID(headers http.Header, requestID string) {
	if requestID != "" {
		headers.Set(HeaderRequestID, requestID)
	}
}

// ExtractBaggage extracts baggage items from context.
func ExtractBaggage(ctx context.Context) baggage.Baggage {
	return baggage.FromContext(ctx)
}

// InjectBaggage injects baggage into context.
func InjectBaggage(ctx context.Context, b baggage.Baggage) context.Context {
	return baggage.ContextWithBaggage(ctx, b)
}

// SpanContextFromContext returns the span context from the given context.
func SpanContextFromContext(ctx context.Context) SpanContext {
	prop := otel.GetTextMapPropagator()

	// Create a carrier to hold the extracted values.
	carrier := make(propagation.MapCarrier)
	prop.Inject(ctx, carrier)

	return SpanContext{
		TraceParent: carrier.Get(HeaderTraceParent),
		TraceState:  carrier.Get(HeaderTraceState),
	}
}

// SpanContext holds the W3C trace context values.
type SpanContext struct {
	TraceParent string
	TraceState  string
}

// IsValid returns true if the span context has a valid trace parent.
func (sc SpanContext) IsValid() bool {
	return sc.TraceParent != ""
}
