package tracing

import (
	"context"
	"testing"
)

func TestNew_ValidConfig(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg := DefaultConfig().
		WithServiceName("test-service").
		WithExporter(ExporterNone) // Use no exporter for tests

	tracer, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if tracer == nil {
		t.Fatal("New() returned nil tracer")
	}

	defer func() {
		if err := tracer.Shutdown(ctx); err != nil {
			t.Errorf("Shutdown() error = %v", err)
		}
	}()

	if tracer.Tracer() == nil {
		t.Error("Tracer() returned nil")
	}
	if tracer.Provider() == nil {
		t.Error("Provider() returned nil")
	}
	if tracer.Config().ServiceName != "test-service" {
		t.Errorf("Config().ServiceName = %v, want test-service", tracer.Config().ServiceName)
	}
}

func TestNew_InvalidConfig(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg := Config{ServiceName: ""} // Invalid: empty service name

	tracer, err := New(ctx, cfg)
	if err == nil {
		t.Error("New() error = nil, want error")
		if tracer != nil {
			tracer.Shutdown(ctx)
		}
	}
}

func TestTracer_Start(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg := DefaultConfig().
		WithServiceName("test-service").
		WithExporter(ExporterNone)

	tracer, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer tracer.Shutdown(ctx)

	spanCtx, span := tracer.Start(ctx, "test-span")
	if span == nil {
		t.Fatal("Start() returned nil span")
	}
	defer span.End()

	if spanCtx == ctx {
		t.Error("Start() should return a new context")
	}

	if !span.SpanContext().IsValid() {
		t.Error("Span context should be valid")
	}
}

func TestTracer_Shutdown(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg := DefaultConfig().
		WithServiceName("test-service").
		WithExporter(ExporterNone)

	tracer, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Shutdown should not error.
	if err := tracer.Shutdown(ctx); err != nil {
		t.Errorf("Shutdown() error = %v", err)
	}
}

func TestTracer_ForceFlush(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg := DefaultConfig().
		WithServiceName("test-service").
		WithExporter(ExporterNone)

	tracer, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer tracer.Shutdown(ctx)

	// ForceFlush should not error.
	if err := tracer.ForceFlush(ctx); err != nil {
		t.Errorf("ForceFlush() error = %v", err)
	}
}

func TestNew_SamplerConfigurations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		sampleRate float64
	}{
		{"always sample", 1.0},
		{"never sample", 0.0},
		{"half sample", 0.5},
		{"low sample", 0.1},
		{"high sample", 0.9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			cfg := DefaultConfig().
				WithServiceName("test-service").
				WithExporter(ExporterNone).
				WithSampleRate(tt.sampleRate)

			tracer, err := New(ctx, cfg)
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}
			defer tracer.Shutdown(ctx)

			// Create a span to verify it works.
			_, span := tracer.Start(ctx, "test-span")
			span.End()
		})
	}
}

func TestNew_PropagationConfigurations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		propagation PropagationType
	}{
		{"w3c propagation", PropagationW3C},
		{"b3 propagation", PropagationB3},
		{"empty propagation defaults to w3c", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			cfg := DefaultConfig().
				WithServiceName("test-service").
				WithExporter(ExporterNone).
				WithPropagation(tt.propagation)

			tracer, err := New(ctx, cfg)
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}
			defer tracer.Shutdown(ctx)
		})
	}
}

func TestNew_ExporterNone(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg := DefaultConfig().
		WithServiceName("test-service").
		WithExporter(ExporterNone)

	tracer, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer tracer.Shutdown(ctx)

	// Create spans to verify it works without exporter.
	_, span := tracer.Start(ctx, "test-span")
	span.End()
}

func TestTracer_NestedSpans(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg := DefaultConfig().
		WithServiceName("test-service").
		WithExporter(ExporterNone)

	tracer, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer tracer.Shutdown(ctx)

	// Create parent span.
	ctx1, parentSpan := tracer.Start(ctx, "parent-span")
	defer parentSpan.End()

	// Create child span.
	_, childSpan := tracer.Start(ctx1, "child-span")
	defer childSpan.End()

	// Child span should have parent.
	parentCtx := parentSpan.SpanContext()
	childCtx := childSpan.SpanContext()

	if parentCtx.TraceID() != childCtx.TraceID() {
		t.Error("Child span should have same TraceID as parent")
	}
}

func TestNew_WithHeaders(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	headers := map[string]string{
		"Authorization": "Bearer token",
	}
	cfg := DefaultConfig().
		WithServiceName("test-service").
		WithExporter(ExporterNone).
		WithHeaders(headers)

	tracer, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer tracer.Shutdown(ctx)

	if tracer.Config().Headers["Authorization"] != "Bearer token" {
		t.Error("Headers not preserved in config")
	}
}
