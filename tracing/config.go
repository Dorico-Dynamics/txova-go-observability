package tracing

import (
	"fmt"
)

// PropagationType defines the trace context propagation format.
type PropagationType string

const (
	// PropagationW3C uses W3C trace context format (traceparent, tracestate).
	PropagationW3C PropagationType = "w3c"
	// PropagationB3 uses B3 propagation format.
	PropagationB3 PropagationType = "b3"
)

// ExporterType defines the trace exporter type.
type ExporterType string

const (
	// ExporterOTLPHTTP exports traces via OTLP over HTTP.
	ExporterOTLPHTTP ExporterType = "otlp-http"
	// ExporterOTLPGRPC exports traces via OTLP over gRPC.
	ExporterOTLPGRPC ExporterType = "otlp-grpc"
	// ExporterNone disables trace exporting (for testing).
	ExporterNone ExporterType = "none"
)

// Config holds configuration for the tracing setup.
type Config struct {
	// ServiceName is the name of the service being traced.
	ServiceName string

	// ServiceVersion is the version of the service.
	ServiceVersion string

	// Endpoint is the OTLP collector endpoint.
	Endpoint string

	// SampleRate is the sampling rate (0.0-1.0).
	// 0.0 means no traces, 1.0 means all traces.
	SampleRate float64

	// Propagation defines the trace context propagation format.
	Propagation PropagationType

	// Exporter defines the trace exporter type.
	Exporter ExporterType

	// Insecure disables TLS for the exporter connection.
	Insecure bool

	// Headers are additional headers to send with exports.
	Headers map[string]string
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		ServiceName:    "unknown-service",
		ServiceVersion: "unknown",
		Endpoint:       "localhost:4318",
		SampleRate:     1.0,
		Propagation:    PropagationW3C,
		Exporter:       ExporterOTLPHTTP,
		Insecure:       true,
		Headers:        make(map[string]string),
	}
}

// WithServiceName sets the service name.
func (c Config) WithServiceName(name string) Config {
	c.ServiceName = name
	return c
}

// WithServiceVersion sets the service version.
func (c Config) WithServiceVersion(version string) Config {
	c.ServiceVersion = version
	return c
}

// WithEndpoint sets the OTLP collector endpoint.
func (c Config) WithEndpoint(endpoint string) Config {
	c.Endpoint = endpoint
	return c
}

// WithSampleRate sets the sampling rate.
func (c Config) WithSampleRate(rate float64) Config {
	c.SampleRate = rate
	return c
}

// WithPropagation sets the trace context propagation format.
func (c Config) WithPropagation(propagation PropagationType) Config {
	c.Propagation = propagation
	return c
}

// WithExporter sets the trace exporter type.
func (c Config) WithExporter(exporter ExporterType) Config {
	c.Exporter = exporter
	return c
}

// WithInsecure sets whether to use insecure connection.
func (c Config) WithInsecure(insecure bool) Config {
	c.Insecure = insecure
	return c
}

// WithHeaders sets additional headers for the exporter.
func (c Config) WithHeaders(headers map[string]string) Config {
	c.Headers = headers
	return c
}

// Validate validates the configuration.
func (c Config) Validate() error {
	if c.ServiceName == "" {
		return fmt.Errorf("service name is required")
	}

	if c.SampleRate < 0 || c.SampleRate > 1 {
		return fmt.Errorf("sample rate must be between 0.0 and 1.0, got %f", c.SampleRate)
	}

	switch c.Propagation {
	case PropagationW3C, PropagationB3:
		// Valid
	case "":
		// Default to W3C
	default:
		return fmt.Errorf("invalid propagation type: %s", c.Propagation)
	}

	switch c.Exporter {
	case ExporterOTLPHTTP, ExporterOTLPGRPC, ExporterNone:
		// Valid
	case "":
		// Default to OTLP HTTP
	default:
		return fmt.Errorf("invalid exporter type: %s", c.Exporter)
	}

	return nil
}
