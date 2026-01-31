// Package metrics provides Prometheus metric collectors for the Txova platform.
// It includes collectors for HTTP requests, database queries, Redis commands,
// Kafka messages, and business-specific metrics.
package metrics

import (
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

// Validate checks that the configuration is valid.
func (c Config) Validate() error {
	if c.Namespace == "" {
		c.Namespace = DefaultNamespace
	}
	if c.Registry == nil {
		c.Registry = prometheus.DefaultRegisterer
	}
	return nil
}
