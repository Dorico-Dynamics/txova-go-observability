package tracing

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/contrib/propagators/b3"
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
func New(ctx context.Context, cfg Config) (*Tracer, error) { //nolint:gocritic // cfg passed by value for API simplicity
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	applyConfigDefaults(&cfg)

	res, err := createResource(&cfg)
	if err != nil {
		return nil, err
	}

	exporter, err := createExporter(ctx, &cfg)
	if err != nil && !errors.Is(err, errNoExporter) {
		return nil, err
	}
	// If errNoExporter, exporter is nil which is handled by createProvider

	sampler := createSampler(cfg.SampleRate)
	provider := createProvider(res, sampler, exporter)

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(createPropagator(cfg.Propagation))

	return &Tracer{
		provider: provider,
		tracer:   provider.Tracer(cfg.ServiceName),
		config:   cfg,
	}, nil
}

// applyConfigDefaults sets default values for empty config fields.
func applyConfigDefaults(cfg *Config) {
	if cfg.Propagation == "" {
		cfg.Propagation = PropagationW3C
	}
	if cfg.Exporter == "" {
		cfg.Exporter = ExporterOTLPHTTP
	}
}

// createResource creates an OpenTelemetry resource with service information.
func createResource(cfg *Config) (*resource.Resource, error) {
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
	return res, nil
}

// errNoExporter is a sentinel value indicating no exporter is configured.
// This is not an error condition, just indicates ExporterNone was selected.
var errNoExporter = fmt.Errorf("no exporter configured")

// createExporter creates a span exporter based on configuration.
func createExporter(ctx context.Context, cfg *Config) (sdktrace.SpanExporter, error) {
	switch cfg.Exporter {
	case ExporterOTLPHTTP:
		return createHTTPExporter(ctx, cfg)
	case ExporterOTLPGRPC:
		return createGRPCExporter(ctx, cfg)
	case ExporterNone:
		return nil, errNoExporter
	default:
		return nil, fmt.Errorf("unsupported exporter type: %s", cfg.Exporter)
	}
}

// createHTTPExporter creates an OTLP HTTP exporter.
func createHTTPExporter(ctx context.Context, cfg *Config) (sdktrace.SpanExporter, error) {
	opts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(cfg.Endpoint),
	}
	if cfg.Insecure {
		opts = append(opts, otlptracehttp.WithInsecure())
	}
	if len(cfg.Headers) > 0 {
		opts = append(opts, otlptracehttp.WithHeaders(cfg.Headers))
	}
	exporter, err := otlptracehttp.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP HTTP exporter: %w", err)
	}
	return exporter, nil
}

// createGRPCExporter creates an OTLP gRPC exporter.
func createGRPCExporter(ctx context.Context, cfg *Config) (sdktrace.SpanExporter, error) {
	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(cfg.Endpoint),
	}
	if cfg.Insecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}
	if len(cfg.Headers) > 0 {
		opts = append(opts, otlptracegrpc.WithHeaders(cfg.Headers))
	}
	exporter, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP gRPC exporter: %w", err)
	}
	return exporter, nil
}

// createSampler creates a sampler based on sample rate.
func createSampler(sampleRate float64) sdktrace.Sampler {
	switch {
	case sampleRate <= 0:
		return sdktrace.NeverSample()
	case sampleRate >= 1:
		return sdktrace.AlwaysSample()
	default:
		return sdktrace.TraceIDRatioBased(sampleRate)
	}
}

// createProvider creates a tracer provider with the given configuration.
func createProvider(res *resource.Resource, sampler sdktrace.Sampler, exporter sdktrace.SpanExporter) *sdktrace.TracerProvider {
	providerOpts := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sampler),
	}

	if exporter != nil {
		providerOpts = append(providerOpts, sdktrace.WithBatcher(exporter))
	}

	return sdktrace.NewTracerProvider(providerOpts...)
}

// createPropagator creates a trace context propagator based on the configured propagation type.
func createPropagator(propagationType PropagationType) propagation.TextMapPropagator {
	switch propagationType {
	case PropagationB3:
		return propagation.NewCompositeTextMapPropagator(
			b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader)),
			propagation.Baggage{},
		)
	case PropagationW3C:
		return propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		)
	default:
		// Default to W3C trace context
		return propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		)
	}
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
