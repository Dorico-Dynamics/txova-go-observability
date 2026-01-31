package tracing

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig()

	if cfg.ServiceName != "unknown-service" {
		t.Errorf("ServiceName = %v, want unknown-service", cfg.ServiceName)
	}
	if cfg.ServiceVersion != "unknown" {
		t.Errorf("ServiceVersion = %v, want unknown", cfg.ServiceVersion)
	}
	if cfg.Endpoint != "localhost:4318" {
		t.Errorf("Endpoint = %v, want localhost:4318", cfg.Endpoint)
	}
	if cfg.SampleRate != 1.0 {
		t.Errorf("SampleRate = %v, want 1.0", cfg.SampleRate)
	}
	if cfg.Propagation != PropagationW3C {
		t.Errorf("Propagation = %v, want %v", cfg.Propagation, PropagationW3C)
	}
	if cfg.Exporter != ExporterOTLPHTTP {
		t.Errorf("Exporter = %v, want %v", cfg.Exporter, ExporterOTLPHTTP)
	}
	if !cfg.Insecure {
		t.Error("Insecure should be true by default")
	}
	if cfg.Headers == nil {
		t.Error("Headers should not be nil")
	}
}

func TestConfig_WithServiceName(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig().WithServiceName("ride-service")

	if cfg.ServiceName != "ride-service" {
		t.Errorf("ServiceName = %v, want ride-service", cfg.ServiceName)
	}
}

func TestConfig_WithServiceVersion(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig().WithServiceVersion("v1.2.3")

	if cfg.ServiceVersion != "v1.2.3" {
		t.Errorf("ServiceVersion = %v, want v1.2.3", cfg.ServiceVersion)
	}
}

func TestConfig_WithEndpoint(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig().WithEndpoint("otel-collector:4317")

	if cfg.Endpoint != "otel-collector:4317" {
		t.Errorf("Endpoint = %v, want otel-collector:4317", cfg.Endpoint)
	}
}

func TestConfig_WithSampleRate(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig().WithSampleRate(0.5)

	if cfg.SampleRate != 0.5 {
		t.Errorf("SampleRate = %v, want 0.5", cfg.SampleRate)
	}
}

func TestConfig_WithPropagation(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig().WithPropagation(PropagationB3)

	if cfg.Propagation != PropagationB3 {
		t.Errorf("Propagation = %v, want %v", cfg.Propagation, PropagationB3)
	}
}

func TestConfig_WithExporter(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig().WithExporter(ExporterOTLPGRPC)

	if cfg.Exporter != ExporterOTLPGRPC {
		t.Errorf("Exporter = %v, want %v", cfg.Exporter, ExporterOTLPGRPC)
	}
}

func TestConfig_WithInsecure(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig().WithInsecure(false)

	if cfg.Insecure {
		t.Error("Insecure should be false")
	}
}

func TestConfig_WithHeaders(t *testing.T) {
	t.Parallel()

	headers := map[string]string{
		"Authorization": "Bearer token",
		"X-Custom":      "value",
	}
	cfg := DefaultConfig().WithHeaders(headers)

	if len(cfg.Headers) != 2 {
		t.Errorf("Headers length = %d, want 2", len(cfg.Headers))
	}
	if cfg.Headers["Authorization"] != "Bearer token" {
		t.Errorf("Headers[Authorization] = %v, want Bearer token", cfg.Headers["Authorization"])
	}
}

func TestConfig_Chaining(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig().
		WithServiceName("test-service").
		WithServiceVersion("v1.0.0").
		WithEndpoint("collector:4317").
		WithSampleRate(0.8).
		WithPropagation(PropagationW3C).
		WithExporter(ExporterOTLPGRPC).
		WithInsecure(true)

	if cfg.ServiceName != "test-service" {
		t.Errorf("ServiceName = %v, want test-service", cfg.ServiceName)
	}
	if cfg.ServiceVersion != "v1.0.0" {
		t.Errorf("ServiceVersion = %v, want v1.0.0", cfg.ServiceVersion)
	}
	if cfg.Endpoint != "collector:4317" {
		t.Errorf("Endpoint = %v, want collector:4317", cfg.Endpoint)
	}
	if cfg.SampleRate != 0.8 {
		t.Errorf("SampleRate = %v, want 0.8", cfg.SampleRate)
	}
}

func TestConfig_Validate_Valid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		cfg  Config
	}{
		{
			name: "default config with service name",
			cfg:  DefaultConfig().WithServiceName("test"),
		},
		{
			name: "full config",
			cfg: Config{
				ServiceName:    "test-service",
				ServiceVersion: "v1.0.0",
				Endpoint:       "localhost:4318",
				SampleRate:     0.5,
				Propagation:    PropagationW3C,
				Exporter:       ExporterOTLPHTTP,
			},
		},
		{
			name: "sample rate 0",
			cfg:  DefaultConfig().WithServiceName("test").WithSampleRate(0),
		},
		{
			name: "sample rate 1",
			cfg:  DefaultConfig().WithServiceName("test").WithSampleRate(1),
		},
		{
			name: "empty propagation defaults",
			cfg: Config{
				ServiceName: "test",
				Propagation: "",
			},
		},
		{
			name: "empty exporter defaults",
			cfg: Config{
				ServiceName: "test",
				Exporter:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.cfg.Validate(); err != nil {
				t.Errorf("Validate() error = %v, want nil", err)
			}
		})
	}
}

func TestConfig_Validate_Invalid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cfg     Config
		wantErr string
	}{
		{
			name:    "empty service name",
			cfg:     Config{ServiceName: ""},
			wantErr: "service name is required",
		},
		{
			name:    "negative sample rate",
			cfg:     Config{ServiceName: "test", SampleRate: -0.1},
			wantErr: "sample rate must be between 0.0 and 1.0",
		},
		{
			name:    "sample rate greater than 1",
			cfg:     Config{ServiceName: "test", SampleRate: 1.5},
			wantErr: "sample rate must be between 0.0 and 1.0",
		},
		{
			name:    "invalid propagation type",
			cfg:     Config{ServiceName: "test", Propagation: "invalid"},
			wantErr: "invalid propagation type",
		},
		{
			name:    "invalid exporter type",
			cfg:     Config{ServiceName: "test", Exporter: "invalid"},
			wantErr: "invalid exporter type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.cfg.Validate()
			if err == nil {
				t.Errorf("Validate() error = nil, want error containing %q", tt.wantErr)
				return
			}
			if !contains(err.Error(), tt.wantErr) {
				t.Errorf("Validate() error = %v, want error containing %q", err, tt.wantErr)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsAt(s, substr, 0))
}

func containsAt(s, substr string, start int) bool {
	if start+len(substr) > len(s) {
		return false
	}
	for i := start; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestPropagationType_Constants(t *testing.T) {
	t.Parallel()

	if PropagationW3C != "w3c" {
		t.Errorf("PropagationW3C = %v, want w3c", PropagationW3C)
	}
	if PropagationB3 != "b3" {
		t.Errorf("PropagationB3 = %v, want b3", PropagationB3)
	}
}

func TestExporterType_Constants(t *testing.T) {
	t.Parallel()

	if ExporterOTLPHTTP != "otlp-http" {
		t.Errorf("ExporterOTLPHTTP = %v, want otlp-http", ExporterOTLPHTTP)
	}
	if ExporterOTLPGRPC != "otlp-grpc" {
		t.Errorf("ExporterOTLPGRPC = %v, want otlp-grpc", ExporterOTLPGRPC)
	}
	if ExporterNone != "none" {
		t.Errorf("ExporterNone = %v, want none", ExporterNone)
	}
}
