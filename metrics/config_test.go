package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestDefaultConfig(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig()

	if cfg.Namespace != DefaultNamespace {
		t.Errorf("Namespace = %v, want %v", cfg.Namespace, DefaultNamespace)
	}
	if cfg.Subsystem != "" {
		t.Errorf("Subsystem = %v, want empty string", cfg.Subsystem)
	}
	if cfg.Registry != prometheus.DefaultRegisterer {
		t.Error("Registry should be prometheus.DefaultRegisterer")
	}
}

func TestConfig_WithNamespace(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig().WithNamespace("custom")

	if cfg.Namespace != "custom" {
		t.Errorf("Namespace = %v, want custom", cfg.Namespace)
	}
}

func TestConfig_WithSubsystem(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig().WithSubsystem("ride_service")

	if cfg.Subsystem != "ride_service" {
		t.Errorf("Subsystem = %v, want ride_service", cfg.Subsystem)
	}
}

func TestConfig_WithRegistry(t *testing.T) {
	t.Parallel()

	customRegistry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(customRegistry)

	if cfg.Registry != customRegistry {
		t.Error("Registry should be the custom registry")
	}
}

func TestConfig_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		cfg  Config
	}{
		{
			name: "default config",
			cfg:  DefaultConfig(),
		},
		{
			name: "empty namespace gets default",
			cfg:  Config{Namespace: "", Subsystem: "test"},
		},
		{
			name: "nil registry gets default",
			cfg:  Config{Namespace: "test", Registry: nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := tt.cfg.Validate()
			if err != nil {
				t.Errorf("Validate() error = %v, want nil", err)
			}
		})
	}
}

func TestConfig_Chaining(t *testing.T) {
	t.Parallel()

	customRegistry := prometheus.NewRegistry()
	cfg := DefaultConfig().
		WithNamespace("myapp").
		WithSubsystem("api").
		WithRegistry(customRegistry)

	if cfg.Namespace != "myapp" {
		t.Errorf("Namespace = %v, want myapp", cfg.Namespace)
	}
	if cfg.Subsystem != "api" {
		t.Errorf("Subsystem = %v, want api", cfg.Subsystem)
	}
	if cfg.Registry != customRegistry {
		t.Error("Registry should be the custom registry")
	}
}
