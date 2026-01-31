// Package observability provides a unified observability solution for Txova services,
// including metrics collection, distributed tracing, and health checks.
package observability

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Dorico-Dynamics/txova-go-observability/health"
	"github.com/Dorico-Dynamics/txova-go-observability/metrics"
	"github.com/Dorico-Dynamics/txova-go-observability/tracing"
)

// Config holds configuration for the observability setup.
type Config struct {
	// Metrics configuration.
	Metrics metrics.Config

	// Tracing configuration.
	Tracing tracing.Config

	// Health configuration.
	Health health.ManagerConfig

	// Enabled flags for each subsystem.
	MetricsEnabled bool
	TracingEnabled bool
	HealthEnabled  bool
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Metrics:        metrics.DefaultConfig(),
		Tracing:        tracing.DefaultConfig(),
		Health:         health.DefaultManagerConfig(),
		MetricsEnabled: true,
		TracingEnabled: true,
		HealthEnabled:  true,
	}
}

// Observability provides a central entry point for all observability features.
type Observability struct {
	config Config

	// Tracer is the OpenTelemetry tracer.
	Tracer *tracing.Tracer

	// HealthManager manages health checks.
	HealthManager *health.Manager

	// HealthHandler provides HTTP handlers for health endpoints.
	HealthHandler *health.Handler

	// HTTPCollector collects HTTP metrics.
	HTTPCollector *metrics.HTTPCollector

	// DBCollector collects database metrics.
	DBCollector *metrics.DBCollector

	// RedisCollector collects Redis metrics.
	RedisCollector *metrics.RedisCollector

	// KafkaCollector collects Kafka metrics.
	KafkaCollector *metrics.KafkaCollector

	// RideCollector collects ride metrics.
	RideCollector *metrics.RideCollector

	// DriverCollector collects driver metrics.
	DriverCollector *metrics.DriverCollector

	// PaymentCollector collects payment metrics.
	PaymentCollector *metrics.PaymentCollector

	// SafetyCollector collects safety metrics.
	SafetyCollector *metrics.SafetyCollector
}

// New creates a new Observability instance with the given configuration.
func New(ctx context.Context, cfg Config) (*Observability, error) {
	obs := &Observability{
		config: cfg,
	}

	// Initialize tracing.
	if cfg.TracingEnabled {
		tracer, err := tracing.New(ctx, cfg.Tracing)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize tracing: %w", err)
		}
		obs.Tracer = tracer
	}

	// Initialize health manager.
	if cfg.HealthEnabled {
		obs.HealthManager = health.NewManager(cfg.Health)
		obs.HealthHandler = health.NewHandler(obs.HealthManager)
	}

	// Initialize metrics collectors.
	if cfg.MetricsEnabled {
		var err error

		obs.HTTPCollector, err = metrics.NewHTTPCollector(cfg.Metrics)
		if err != nil {
			return nil, fmt.Errorf("failed to create HTTP collector: %w", err)
		}

		obs.DBCollector, err = metrics.NewDBCollector(cfg.Metrics)
		if err != nil {
			return nil, fmt.Errorf("failed to create DB collector: %w", err)
		}

		obs.RedisCollector, err = metrics.NewRedisCollector(cfg.Metrics)
		if err != nil {
			return nil, fmt.Errorf("failed to create Redis collector: %w", err)
		}

		obs.KafkaCollector, err = metrics.NewKafkaCollector(cfg.Metrics)
		if err != nil {
			return nil, fmt.Errorf("failed to create Kafka collector: %w", err)
		}

		obs.RideCollector, err = metrics.NewRideCollector(cfg.Metrics)
		if err != nil {
			return nil, fmt.Errorf("failed to create Ride collector: %w", err)
		}

		obs.DriverCollector, err = metrics.NewDriverCollector(cfg.Metrics)
		if err != nil {
			return nil, fmt.Errorf("failed to create Driver collector: %w", err)
		}

		obs.PaymentCollector, err = metrics.NewPaymentCollector(cfg.Metrics)
		if err != nil {
			return nil, fmt.Errorf("failed to create Payment collector: %w", err)
		}

		obs.SafetyCollector, err = metrics.NewSafetyCollector(cfg.Metrics)
		if err != nil {
			return nil, fmt.Errorf("failed to create Safety collector: %w", err)
		}
	}

	return obs, nil
}

// Initialize starts all observability subsystems.
// This implements the app.Initializer interface from txova-go-core.
func (o *Observability) Initialize(ctx context.Context) error {
	if o.HealthManager != nil && o.config.HealthEnabled {
		o.HealthManager.StartBackground(ctx)
	}
	return nil
}

// Close shuts down all observability subsystems.
// This implements the app.Closer interface from txova-go-core.
func (o *Observability) Close(ctx context.Context) error {
	if o.HealthManager != nil {
		o.HealthManager.StopBackground()
	}

	if o.Tracer != nil {
		if err := o.Tracer.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown tracer: %w", err)
		}
	}

	return nil
}

// HealthCheck returns the current health status.
// This implements the app.HealthChecker interface from txova-go-core.
func (o *Observability) HealthCheck(ctx context.Context) error {
	if o.HealthManager == nil {
		return nil
	}

	report := o.HealthManager.Check(ctx)
	if report.Status == health.StatusUnhealthy {
		return fmt.Errorf("health check failed: %s", report.Status)
	}
	return nil
}

// HTTPMiddleware returns an HTTP middleware that adds tracing and metrics.
func (o *Observability) HTTPMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		handler := next

		// Apply metrics middleware.
		if o.HTTPCollector != nil {
			handler = o.metricsMiddleware(handler)
		}

		// Apply tracing middleware.
		if o.Tracer != nil {
			handler = tracing.Middleware(o.Tracer)(handler)
		}

		return handler
	}
}

// metricsMiddleware wraps an HTTP handler to collect metrics.
func (o *Observability) metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		o.HTTPCollector.IncRequestsInFlight()
		defer o.HTTPCollector.DecRequestsInFlight()

		// Create a response writer wrapper to capture the status code.
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		start := time.Now()

		next.ServeHTTP(rw, r)

		// Record metrics after request completes.
		duration := time.Since(start)
		o.HTTPCollector.RecordRequest(r.Method, r.URL.Path, rw.statusCode, duration)
	})
}

// responseWriter wraps http.ResponseWriter to capture the status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.written {
		rw.statusCode = code
		rw.written = true
	}
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.statusCode = http.StatusOK
		rw.written = true
	}
	return rw.ResponseWriter.Write(b)
}

// RegisterHealthChecker registers a health checker with the manager.
func (o *Observability) RegisterHealthChecker(checker health.Checker) {
	if o.HealthManager != nil {
		o.HealthManager.Register(checker)
	}
}

// HTTPRoundTripper returns an HTTP RoundTripper with tracing.
func (o *Observability) HTTPRoundTripper(base http.RoundTripper) http.RoundTripper {
	if o.Tracer != nil {
		return tracing.RoundTripper(o.Tracer, base)
	}
	if base == nil {
		return http.DefaultTransport
	}
	return base
}
