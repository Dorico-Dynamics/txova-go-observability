package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
)

// Tracer wraps an OpenTelemetry tracer with configuration and lifecycle management.
type Tracer struct {
	provider *sdktrace.TracerProvider
	tracer   trace.Tracer
	config   Config
}

// New creates a new Tracer with the given configuration.
func New(ctx context.Context, cfg Config) (*Tracer, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// Set defaults for empty values.
	if cfg.Propagation == "" {
		cfg.Propagation = PropagationW3C
	}
	if cfg.Exporter == "" {
		cfg.Exporter = ExporterOTLPHTTP
	}

	// Create resource with service information.
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String(cfg.ServiceVersion),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create exporter.
	var exporter sdktrace.SpanExporter
	switch cfg.Exporter {
	case ExporterOTLPHTTP:
		opts := []otlptracehttp.Option{
			otlptracehttp.WithEndpoint(cfg.Endpoint),
		}
		if cfg.Insecure {
			opts = append(opts, otlptracehttp.WithInsecure())
		}
		if len(cfg.Headers) > 0 {
			opts = append(opts, otlptracehttp.WithHeaders(cfg.Headers))
		}
		exporter, err = otlptracehttp.New(ctx, opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create OTLP HTTP exporter: %w", err)
		}

	case ExporterOTLPGRPC:
		opts := []otlptracegrpc.Option{
			otlptracegrpc.WithEndpoint(cfg.Endpoint),
		}
		if cfg.Insecure {
			opts = append(opts, otlptracegrpc.WithInsecure())
		}
		if len(cfg.Headers) > 0 {
			opts = append(opts, otlptracegrpc.WithHeaders(cfg.Headers))
		}
		exporter, err = otlptracegrpc.New(ctx, opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create OTLP gRPC exporter: %w", err)
		}

	case ExporterNone:
		// No exporter, useful for testing.
		exporter = nil

	default:
		return nil, fmt.Errorf("unsupported exporter type: %s", cfg.Exporter)
	}

	// Create sampler based on sample rate.
	var sampler sdktrace.Sampler
	switch {
	case cfg.SampleRate <= 0:
		sampler = sdktrace.NeverSample()
	case cfg.SampleRate >= 1:
		sampler = sdktrace.AlwaysSample()
	default:
		sampler = sdktrace.TraceIDRatioBased(cfg.SampleRate)
	}

	// Create tracer provider options.
	providerOpts := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sampler),
	}

	if exporter != nil {
		providerOpts = append(providerOpts, sdktrace.WithBatcher(exporter))
	}

	// Create tracer provider.
	provider := sdktrace.NewTracerProvider(providerOpts...)

	// Set global tracer provider.
	otel.SetTracerProvider(provider)

	// Set global propagator based on configuration.
	var propagator propagation.TextMapPropagator
	switch cfg.Propagation {
	case PropagationW3C:
		propagator = propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		)
	case PropagationB3:
		// B3 propagation would require additional import.
		// For now, default to W3C.
		propagator = propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		)
	default:
		propagator = propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		)
	}
	otel.SetTextMapPropagator(propagator)

	return &Tracer{
		provider: provider,
		tracer:   provider.Tracer(cfg.ServiceName),
		config:   cfg,
	}, nil
}

// Tracer returns the underlying OpenTelemetry tracer.
func (t *Tracer) Tracer() trace.Tracer {
	return t.tracer
}

// Provider returns the underlying tracer provider.
func (t *Tracer) Provider() *sdktrace.TracerProvider {
	return t.provider
}

// Config returns the tracer configuration.
func (t *Tracer) Config() Config {
	return t.config
}

// Start creates a new span with the given name.
func (t *Tracer) Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, name, opts...)
}

// Shutdown gracefully shuts down the tracer provider.
func (t *Tracer) Shutdown(ctx context.Context) error {
	return t.provider.Shutdown(ctx)
}

// ForceFlush flushes any pending spans.
func (t *Tracer) ForceFlush(ctx context.Context) error {
	return t.provider.ForceFlush(ctx)
}
