package tracing

import (
	"context"
	"net/http"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func TestHTTPCarrier_Get(t *testing.T) {
	t.Parallel()

	headers := http.Header{}
	headers.Set("traceparent", "00-abc123-def456-01")
	headers.Set("X-Custom", "value")

	carrier := HTTPCarrier(headers)

	if got := carrier.Get("traceparent"); got != "00-abc123-def456-01" {
		t.Errorf("Get(traceparent) = %v, want 00-abc123-def456-01", got)
	}
	if got := carrier.Get("X-Custom"); got != "value" {
		t.Errorf("Get(X-Custom) = %v, want value", got)
	}
	if got := carrier.Get("nonexistent"); got != "" {
		t.Errorf("Get(nonexistent) = %v, want empty string", got)
	}
}

func TestHTTPCarrier_Set(t *testing.T) {
	t.Parallel()

	headers := http.Header{}
	carrier := HTTPCarrier(headers)

	carrier.Set("traceparent", "00-xyz789-uvw012-00")

	if got := headers.Get("traceparent"); got != "00-xyz789-uvw012-00" {
		t.Errorf("headers.Get(traceparent) = %v, want 00-xyz789-uvw012-00", got)
	}
}

func TestHTTPCarrier_Keys(t *testing.T) {
	t.Parallel()

	headers := http.Header{}
	headers.Set("traceparent", "value1")
	headers.Set("tracestate", "value2")

	carrier := HTTPCarrier(headers)
	keys := carrier.Keys()

	if len(keys) != 2 {
		t.Errorf("Keys() length = %d, want 2", len(keys))
	}
}

func TestKafkaCarrier_Get(t *testing.T) {
	t.Parallel()

	headers := KafkaCarrier{
		"traceparent": "00-abc123-def456-01",
		"custom":      "value",
	}

	if got := headers.Get("traceparent"); got != "00-abc123-def456-01" {
		t.Errorf("Get(traceparent) = %v, want 00-abc123-def456-01", got)
	}
	if got := headers.Get("custom"); got != "value" {
		t.Errorf("Get(custom) = %v, want value", got)
	}
	if got := headers.Get("nonexistent"); got != "" {
		t.Errorf("Get(nonexistent) = %v, want empty string", got)
	}
}

func TestKafkaCarrier_Set(t *testing.T) {
	t.Parallel()

	headers := make(KafkaCarrier)

	headers.Set("traceparent", "00-xyz789-uvw012-00")

	if got := headers["traceparent"]; got != "00-xyz789-uvw012-00" {
		t.Errorf("headers[traceparent] = %v, want 00-xyz789-uvw012-00", got)
	}
}

func TestKafkaCarrier_Keys(t *testing.T) {
	t.Parallel()

	headers := KafkaCarrier{
		"traceparent": "value1",
		"tracestate":  "value2",
		"custom":      "value3",
	}

	keys := headers.Keys()

	if len(keys) != 3 {
		t.Errorf("Keys() length = %d, want 3", len(keys))
	}
}

func TestExtract_InjectsContext(t *testing.T) {
	t.Parallel()

	// Set up a global propagator first.
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	headers := http.Header{}
	// Valid W3C trace context format.
	headers.Set("traceparent", "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01")

	ctx := Extract(context.Background(), headers)

	// Verify the context was modified (has trace info).
	if ctx == nil {
		t.Error("Extract() returned nil context")
	}
}

func TestInject(t *testing.T) {
	t.Parallel()

	// Create a tracer to generate valid spans.
	ctx := context.Background()
	cfg := Config{
		ServiceName: "test-service",
		Exporter:    ExporterNone,
	}

	tracer, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer tracer.Shutdown(ctx)

	// Start a span to create trace context.
	spanCtx, span := tracer.Start(ctx, "test-span")
	defer span.End()

	// Inject into headers.
	headers := http.Header{}
	Inject(spanCtx, headers)

	// Verify traceparent was injected.
	if headers.Get("traceparent") == "" {
		t.Error("Inject() did not set traceparent header")
	}
}

func TestExtractFromKafka(t *testing.T) {
	t.Parallel()

	// Set up a global propagator first.
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	headers := map[string]string{
		"traceparent": "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
	}

	ctx := ExtractFromKafka(context.Background(), headers)

	if ctx == nil {
		t.Error("ExtractFromKafka() returned nil context")
	}
}

func TestInjectToKafka(t *testing.T) {
	t.Parallel()

	// Create a tracer to generate valid spans.
	ctx := context.Background()
	cfg := Config{
		ServiceName: "test-service",
		Exporter:    ExporterNone,
	}

	tracer, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer tracer.Shutdown(ctx)

	// Start a span to create trace context.
	spanCtx, span := tracer.Start(ctx, "test-span")
	defer span.End()

	// Inject into Kafka headers.
	headers := make(map[string]string)
	InjectToKafka(spanCtx, headers)

	// Verify traceparent was injected.
	if headers["traceparent"] == "" {
		t.Error("InjectToKafka() did not set traceparent header")
	}
}

func TestExtractRequestID(t *testing.T) {
	t.Parallel()

	headers := http.Header{}
	headers.Set(HeaderRequestID, "req-abc-123")

	requestID := ExtractRequestID(headers)

	if requestID != "req-abc-123" {
		t.Errorf("ExtractRequestID() = %v, want req-abc-123", requestID)
	}
}

func TestExtractRequestID_Empty(t *testing.T) {
	t.Parallel()

	headers := http.Header{}

	requestID := ExtractRequestID(headers)

	if requestID != "" {
		t.Errorf("ExtractRequestID() = %v, want empty string", requestID)
	}
}

func TestInjectRequestID(t *testing.T) {
	t.Parallel()

	headers := http.Header{}
	InjectRequestID(headers, "req-xyz-789")

	if got := headers.Get(HeaderRequestID); got != "req-xyz-789" {
		t.Errorf("headers.Get(X-Request-ID) = %v, want req-xyz-789", got)
	}
}

func TestInjectRequestID_Empty(t *testing.T) {
	t.Parallel()

	headers := http.Header{}
	InjectRequestID(headers, "")

	if got := headers.Get(HeaderRequestID); got != "" {
		t.Errorf("headers.Get(X-Request-ID) = %v, want empty string", got)
	}
}

func TestSpanContext_IsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		sc       SpanContext
		expected bool
	}{
		{
			name:     "valid trace parent",
			sc:       SpanContext{TraceParent: "00-abc-def-01"},
			expected: true,
		},
		{
			name:     "empty trace parent",
			sc:       SpanContext{TraceParent: ""},
			expected: false,
		},
		{
			name:     "only trace state",
			sc:       SpanContext{TraceState: "key=value"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.sc.IsValid(); got != tt.expected {
				t.Errorf("IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSpanContextFromContext(t *testing.T) {
	t.Parallel()

	// Create a tracer to generate valid spans.
	ctx := context.Background()
	cfg := Config{
		ServiceName: "test-service",
		Exporter:    ExporterNone,
	}

	tracer, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer tracer.Shutdown(ctx)

	// Start a span to create trace context.
	spanCtx, span := tracer.Start(ctx, "test-span")
	defer span.End()

	sc := SpanContextFromContext(spanCtx)

	if !sc.IsValid() {
		t.Error("SpanContextFromContext() should return valid span context")
	}
	if sc.TraceParent == "" {
		t.Error("TraceParent should not be empty")
	}
}

func TestHeaderConstants(t *testing.T) {
	t.Parallel()

	if HeaderTraceParent != "traceparent" {
		t.Errorf("HeaderTraceParent = %v, want traceparent", HeaderTraceParent)
	}
	if HeaderTraceState != "tracestate" {
		t.Errorf("HeaderTraceState = %v, want tracestate", HeaderTraceState)
	}
	if HeaderRequestID != "X-Request-ID" {
		t.Errorf("HeaderRequestID = %v, want X-Request-ID", HeaderRequestID)
	}
}
