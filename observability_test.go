package observability

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Dorico-Dynamics/txova-go-observability/health"
	"github.com/Dorico-Dynamics/txova-go-observability/metrics"
	"github.com/Dorico-Dynamics/txova-go-observability/tracing"
	"github.com/prometheus/client_golang/prometheus"
)

func TestDefaultConfig(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig()

	if !cfg.MetricsEnabled {
		t.Error("MetricsEnabled should be true by default")
	}
	if !cfg.TracingEnabled {
		t.Error("TracingEnabled should be true by default")
	}
	if !cfg.HealthEnabled {
		t.Error("HealthEnabled should be true by default")
	}
}

func TestNew_AllEnabled(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	registry := prometheus.NewRegistry()
	cfg := DefaultConfig()
	cfg.Metrics = metrics.DefaultConfig().WithRegistry(registry).WithSubsystem("test_all")
	cfg.Tracing = tracing.DefaultConfig().
		WithServiceName("test-service").
		WithExporter(tracing.ExporterNone)

	obs, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer obs.Close(ctx)

	if obs.Tracer == nil {
		t.Error("Tracer should not be nil")
	}
	if obs.HealthManager == nil {
		t.Error("HealthManager should not be nil")
	}
	if obs.HealthHandler == nil {
		t.Error("HealthHandler should not be nil")
	}
	if obs.HTTPCollector == nil {
		t.Error("HTTPCollector should not be nil")
	}
	if obs.DBCollector == nil {
		t.Error("DBCollector should not be nil")
	}
	if obs.RedisCollector == nil {
		t.Error("RedisCollector should not be nil")
	}
	if obs.KafkaCollector == nil {
		t.Error("KafkaCollector should not be nil")
	}
	if obs.RideCollector == nil {
		t.Error("RideCollector should not be nil")
	}
	if obs.DriverCollector == nil {
		t.Error("DriverCollector should not be nil")
	}
	if obs.PaymentCollector == nil {
		t.Error("PaymentCollector should not be nil")
	}
	if obs.SafetyCollector == nil {
		t.Error("SafetyCollector should not be nil")
	}
}

func TestNew_AllDisabled(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg := DefaultConfig()
	cfg.MetricsEnabled = false
	cfg.TracingEnabled = false
	cfg.HealthEnabled = false

	obs, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer obs.Close(ctx)

	if obs.Tracer != nil {
		t.Error("Tracer should be nil when disabled")
	}
	if obs.HealthManager != nil {
		t.Error("HealthManager should be nil when disabled")
	}
	if obs.HTTPCollector != nil {
		t.Error("HTTPCollector should be nil when disabled")
	}
}

func TestNew_InvalidTracingConfig(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg := DefaultConfig()
	cfg.Tracing = tracing.Config{ServiceName: ""} // Invalid
	cfg.MetricsEnabled = false
	cfg.HealthEnabled = false

	_, err := New(ctx, cfg)
	if err == nil {
		t.Error("New() should return error for invalid tracing config")
	}
}

func TestObservability_Initialize(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg := DefaultConfig()
	cfg.MetricsEnabled = false
	cfg.TracingEnabled = false

	obs, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if err := obs.Initialize(ctx); err != nil {
		t.Errorf("Initialize() error = %v", err)
	}

	if err := obs.Close(ctx); err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

func TestObservability_Close(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	registry := prometheus.NewRegistry()
	cfg := DefaultConfig()
	cfg.Metrics = metrics.DefaultConfig().WithRegistry(registry).WithSubsystem("test_close")
	cfg.Tracing = tracing.DefaultConfig().
		WithServiceName("test-service").
		WithExporter(tracing.ExporterNone)

	obs, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if err := obs.Initialize(ctx); err != nil {
		t.Errorf("Initialize() error = %v", err)
	}

	if err := obs.Close(ctx); err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

func TestObservability_HealthCheck_NoManager(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg := DefaultConfig()
	cfg.MetricsEnabled = false
	cfg.TracingEnabled = false
	cfg.HealthEnabled = false

	obs, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	err = obs.HealthCheck(ctx)
	if err != nil {
		t.Errorf("HealthCheck() error = %v, want nil", err)
	}
}

func TestObservability_HealthCheck_Healthy(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg := DefaultConfig()
	cfg.MetricsEnabled = false
	cfg.TracingEnabled = false
	cfg.Health = health.DefaultManagerConfig().WithCacheTTL(0)

	obs, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	obs.RegisterHealthChecker(health.NewFuncChecker("test", func(ctx context.Context) error {
		return nil
	}, true))

	err = obs.HealthCheck(ctx)
	if err != nil {
		t.Errorf("HealthCheck() error = %v, want nil", err)
	}
}

func TestObservability_HealthCheck_Unhealthy(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg := DefaultConfig()
	cfg.MetricsEnabled = false
	cfg.TracingEnabled = false
	cfg.Health = health.DefaultManagerConfig().WithCacheTTL(0)

	obs, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	obs.RegisterHealthChecker(health.NewFuncChecker("failing", func(ctx context.Context) error {
		return errors.New("check failed")
	}, true))

	err = obs.HealthCheck(ctx)
	if err == nil {
		t.Error("HealthCheck() should return error when unhealthy")
	}
}

func TestObservability_HTTPMiddleware(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	registry := prometheus.NewRegistry()
	cfg := DefaultConfig()
	cfg.Metrics = metrics.DefaultConfig().WithRegistry(registry).WithSubsystem("test_mw")
	cfg.Tracing = tracing.DefaultConfig().
		WithServiceName("test-service").
		WithExporter(tracing.ExporterNone)

	obs, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer obs.Close(ctx)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	wrapped := obs.HTTPMiddleware()(handler)

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestObservability_HTTPMiddleware_NoMetrics(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg := DefaultConfig()
	cfg.MetricsEnabled = false
	cfg.TracingEnabled = false
	cfg.HealthEnabled = false

	obs, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrapped := obs.HTTPMiddleware()(handler)

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestObservability_RegisterHealthChecker(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg := DefaultConfig()
	cfg.MetricsEnabled = false
	cfg.TracingEnabled = false
	cfg.Health = health.DefaultManagerConfig().WithCacheTTL(0)

	obs, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	checker := health.NewFuncChecker("custom", func(ctx context.Context) error {
		return nil
	}, true)

	obs.RegisterHealthChecker(checker)

	report := obs.HealthManager.Check(ctx)
	if _, ok := report.Checks["custom"]; !ok {
		t.Error("Registered checker not found in report")
	}
}

func TestObservability_RegisterHealthChecker_NoManager(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg := DefaultConfig()
	cfg.MetricsEnabled = false
	cfg.TracingEnabled = false
	cfg.HealthEnabled = false

	obs, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	checker := health.NewFuncChecker("custom", func(ctx context.Context) error {
		return nil
	}, true)

	// Should not panic.
	obs.RegisterHealthChecker(checker)
}

func TestObservability_HTTPRoundTripper_WithTracer(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg := DefaultConfig()
	cfg.MetricsEnabled = false
	cfg.HealthEnabled = false
	cfg.Tracing = tracing.DefaultConfig().
		WithServiceName("test-service").
		WithExporter(tracing.ExporterNone)

	obs, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer obs.Close(ctx)

	rt := obs.HTTPRoundTripper(nil)
	if rt == nil {
		t.Error("HTTPRoundTripper() should not return nil")
	}
}

func TestObservability_HTTPRoundTripper_NoTracer(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg := DefaultConfig()
	cfg.MetricsEnabled = false
	cfg.TracingEnabled = false
	cfg.HealthEnabled = false

	obs, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	rt := obs.HTTPRoundTripper(nil)
	if rt != http.DefaultTransport {
		t.Error("HTTPRoundTripper() should return DefaultTransport when no tracer")
	}
}

func TestObservability_HTTPRoundTripper_CustomBase(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg := DefaultConfig()
	cfg.MetricsEnabled = false
	cfg.TracingEnabled = false
	cfg.HealthEnabled = false

	obs, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	customTransport := &http.Transport{}
	rt := obs.HTTPRoundTripper(customTransport)
	if rt != customTransport {
		t.Error("HTTPRoundTripper() should return custom transport when no tracer")
	}
}

func TestResponseWriter_WriteHeader(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	rw := &responseWriter{ResponseWriter: rec, statusCode: http.StatusOK}

	rw.WriteHeader(http.StatusCreated)

	if rw.statusCode != http.StatusCreated {
		t.Errorf("statusCode = %d, want %d", rw.statusCode, http.StatusCreated)
	}
	if !rw.written {
		t.Error("written should be true")
	}
}

func TestResponseWriter_Write(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	rw := &responseWriter{ResponseWriter: rec, statusCode: http.StatusOK}

	n, err := rw.Write([]byte("test"))
	if err != nil {
		t.Errorf("Write() error = %v", err)
	}
	if n != 4 {
		t.Errorf("Write() n = %d, want 4", n)
	}
	if !rw.written {
		t.Error("written should be true after Write")
	}
}
