// Package metrics provides Prometheus metric collectors for the Txova platform.
// It includes collectors for HTTP requests, database queries, Redis commands,
// Kafka messages, and business-specific metrics.
package metrics

import (
	"errors"

	"github.com/prometheus/client_golang/prometheus"
)

// Default namespace for all Txova metrics.
const DefaultNamespace = "txova"

// Config holds configuration for metric collectors.
type Config struct {
	// Namespace is the prefix for all metric names. Default: "txova".
	Namespace string

	// Subsystem is the service name (e.g., "ride_service", "payment_service").
	Subsystem string

	// Registry is the Prometheus registry to use. If nil, the default registry is used.
	Registry prometheus.Registerer
}

// DefaultConfig returns a Config with default values.
func DefaultConfig() Config {
	return Config{
		Namespace: DefaultNamespace,
		Subsystem: "",
		Registry:  prometheus.DefaultRegisterer,
	}
}

// WithNamespace returns a new Config with the specified namespace.
func (c Config) WithNamespace(namespace string) Config {
	c.Namespace = namespace
	return c
}

// WithSubsystem returns a new Config with the specified subsystem.
func (c Config) WithSubsystem(subsystem string) Config {
	c.Subsystem = subsystem
	return c
}

// WithRegistry returns a new Config with the specified registry.
func (c Config) WithRegistry(registry prometheus.Registerer) Config {
	c.Registry = registry
	return c
}

// Validate checks that the configuration is valid and returns a validated copy.
func (c Config) Validate() (Config, error) { //nolint:unparam // error kept for API consistency and future validation
	if c.Namespace == "" {
		c.Namespace = DefaultNamespace
	}
	if c.Registry == nil {
		c.Registry = prometheus.DefaultRegisterer
	}
	return c, nil
}

// registerCollector registers a collector with the registry, handling already registered errors.
// If the collector is already registered, it returns the existing collector.
func registerCollector[T prometheus.Collector](registry prometheus.Registerer, collector T) (T, error) {
	if err := registry.Register(collector); err != nil {
		var are prometheus.AlreadyRegisteredError
		if errors.As(err, &are) {
			if existing, ok := are.ExistingCollector.(T); ok {
				return existing, nil
			}
		}
		var zero T
		return zero, err
	}
	return collector, nil
}
