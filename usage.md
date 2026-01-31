# Usage Guide

This guide provides comprehensive examples for using `txova-go-observability` in your Txova services.

## Table of Contents

- [Unified Observability](#unified-observability)
- [Metrics](#metrics)
- [Tracing](#tracing)
- [Health Checks](#health-checks)
- [Integration Patterns](#integration-patterns)

## Unified Observability

The `Observability` struct provides a single entry point for all observability features.

### Basic Setup

```go
package main

import (
    "context"
    "log"
    "net/http"

    "github.com/Dorico-Dynamics/txova-go-observability"
    "github.com/Dorico-Dynamics/txova-go-observability/metrics"
    "github.com/Dorico-Dynamics/txova-go-observability/tracing"
    "github.com/Dorico-Dynamics/txova-go-observability/health"
)

func main() {
    ctx := context.Background()

    obs, err := observability.New(ctx, &observability.Config{
        MetricsEnabled: true,
        TracingEnabled: true,
        HealthEnabled:  true,
        Metrics: metrics.Config{
            Namespace: "txova",
            Subsystem: "ride_service",
        },
        Tracing: tracing.Config{
            ServiceName:    "ride-service",
            ServiceVersion: "1.0.0",
            Endpoint:       "otel-collector:4318",
            SampleRate:     0.1,
            Exporter:       tracing.ExporterOTLPHTTP,
        },
        Health: health.DefaultManagerConfig(),
    })
    if err != nil {
        log.Fatal(err)
    }
    defer obs.Close(ctx)

    // Start background health checks
    if err := obs.Initialize(ctx); err != nil {
        log.Fatal(err)
    }

    // Set up HTTP server with observability
    mux := http.NewServeMux()
    
    // Register health endpoints
    obs.HealthHandler.RegisterRoutes(mux)
    
    // Apply tracing and metrics middleware
    handler := obs.HTTPMiddleware()(mux)

    log.Println("Starting server on :8080")
    http.ListenAndServe(":8080", handler)
}
```

### Using Default Configuration

```go
// Use default configuration (all features enabled)
obs, err := observability.New(ctx, nil)
```

### Selective Features

```go
// Only enable metrics and health checks
obs, err := observability.New(ctx, &observability.Config{
    MetricsEnabled: true,
    TracingEnabled: false,
    HealthEnabled:  true,
})
```

## Metrics

### Using Collectors via Observability

```go
// HTTP metrics are automatically collected via HTTPMiddleware()
// For business metrics, access collectors directly:

// Record a ride request
obs.RideCollector.RecordRideRequested("standard", "maputo")

// Record a completed ride with details
obs.RideCollector.RecordRideCompleted("standard", "maputo")
obs.RideCollector.RecordRideDuration("standard", 15*time.Minute)
obs.RideCollector.RecordRideDistance("standard", 8.5)
obs.RideCollector.RecordRideFare("standard", 250.00)

// Record a cancelled ride
obs.RideCollector.RecordRideCancelled("rider", "changed_mind")

// Record driver metrics
obs.DriverCollector.RecordDriverOnline("maputo", "standard")
obs.DriverCollector.RecordDriverOffline("maputo", "standard")
obs.DriverCollector.RecordDriverAvailable("maputo", "standard")
obs.DriverCollector.RecordDriverBusy("maputo", "standard")

// Record payment metrics
obs.PaymentCollector.RecordPaymentAttempt("mpesa", "success")
obs.PaymentCollector.RecordPaymentAmount("mpesa", 500.00)

// Record database metrics
obs.DBCollector.RecordQuery("select_ride", 5*time.Millisecond)
obs.DBCollector.RecordQueryError("select_ride", "timeout")

// Record Redis metrics
obs.RedisCollector.RecordCommand("GET", 1*time.Millisecond)
obs.RedisCollector.RecordCacheHit("user_profile")
obs.RedisCollector.RecordCacheMiss("ride_cache")

// Record Kafka metrics
obs.KafkaCollector.RecordMessageProduced("ride-events")
obs.KafkaCollector.RecordMessageConsumed("ride-events", "ride-processor")
```

### Standalone Collector Usage

For services that only need specific metrics:

```go
import "github.com/Dorico-Dynamics/txova-go-observability/metrics"

// Create a standalone HTTP collector
httpCollector, err := metrics.NewHTTPCollector(metrics.Config{
    Namespace: "txova",
    Subsystem: "api_gateway",
})
if err != nil {
    log.Fatal(err)
}

// Record metrics manually
httpCollector.RecordRequest("POST", "/v1/rides", 201, 50*time.Millisecond)
httpCollector.RecordRequestSize("POST", "/v1/rides", 256)
httpCollector.RecordResponseSize("POST", "/v1/rides", 512)
```

### Custom Registry

```go
import "github.com/prometheus/client_golang/prometheus"

// Use a custom registry for testing or isolation
registry := prometheus.NewRegistry()

obs, err := observability.New(ctx, &observability.Config{
    Metrics: metrics.Config{
        Namespace: "txova",
        Subsystem: "test_service",
        Registry:  registry,
    },
})
```

### Exposing Metrics

```go
import "github.com/prometheus/client_golang/prometheus/promhttp"

// Add Prometheus metrics endpoint
mux.Handle("/metrics", promhttp.Handler())
```

## Tracing

### Creating Spans

```go
// Using the observability tracer
ctx, span := obs.Tracer.Start(ctx, "process-ride-request")
defer span.End()

// Add attributes to span
span.SetAttributes(
    attribute.String("ride.id", rideID),
    attribute.String("user.id", userID),
    attribute.String("service.type", "standard"),
)

// Record errors
if err != nil {
    span.RecordError(err)
    span.SetStatus(codes.Error, err.Error())
    return err
}
```

### Nested Spans

```go
func ProcessRide(ctx context.Context, obs *observability.Observability, rideID string) error {
    ctx, span := obs.Tracer.Start(ctx, "process-ride")
    defer span.End()

    // Nested span for validation
    if err := validateRide(ctx, obs, rideID); err != nil {
        return err
    }

    // Nested span for matching
    if err := matchDriver(ctx, obs, rideID); err != nil {
        return err
    }

    return nil
}

func validateRide(ctx context.Context, obs *observability.Observability, rideID string) error {
    ctx, span := obs.Tracer.Start(ctx, "validate-ride")
    defer span.End()
    
    // validation logic...
    return nil
}
```

### HTTP Client with Tracing

```go
// Create an HTTP client with trace propagation
client := &http.Client{
    Transport: obs.HTTPRoundTripper(http.DefaultTransport),
}

// Traces are automatically propagated to downstream services
resp, err := client.Do(req.WithContext(ctx))
```

### Standalone Tracer

```go
import "github.com/Dorico-Dynamics/txova-go-observability/tracing"

// Create standalone tracer
tracer, err := tracing.New(ctx, tracing.Config{
    ServiceName:    "payment-service",
    ServiceVersion: "2.0.0",
    Endpoint:       "localhost:4318",
    SampleRate:     1.0, // Sample all traces in dev
    Exporter:       tracing.ExporterOTLPHTTP,
    Insecure:       true,
})
if err != nil {
    log.Fatal(err)
}
defer tracer.Shutdown(ctx)
```

### Context Propagation

```go
import "github.com/Dorico-Dynamics/txova-go-observability/tracing"

// Extract trace context from incoming HTTP request
ctx := tracing.Extract(r.Context(), r.Header)

// Inject trace context into outgoing HTTP request
tracing.Inject(ctx, outgoingReq.Header)

// Helper for Kafka headers (convert to/from map)
kafkaHeaders := tracing.HeadersToMap(ctx)
ctx = tracing.ExtractFromMap(ctx, kafkaHeaders)
```

### Utility Functions

```go
import "github.com/Dorico-Dynamics/txova-go-observability/tracing"

// Get current span from context
span := tracing.SpanFromContext(ctx)

// Add attributes to current span without getting the span
tracing.AddSpanAttributes(ctx,
    attribute.String("key", "value"),
)

// Record error on current span
tracing.RecordError(ctx, err)
```

## Health Checks

### Registering Health Checkers

```go
import (
    "database/sql"
    "github.com/Dorico-Dynamics/txova-go-observability/health"
)

// Register PostgreSQL checker (required dependency)
obs.RegisterHealthChecker(health.NewPostgresChecker("postgres", db, true))

// Register Redis checker (required dependency)
obs.RegisterHealthChecker(health.NewRedisChecker("redis", redisClient, true))

// Register Kafka checker (optional dependency)
obs.RegisterHealthChecker(health.NewKafkaChecker("kafka", kafkaClient, false))

// Register external service checker
obs.RegisterHealthChecker(health.NewHTTPChecker(
    "payment-gateway",
    "https://payment.example.com/health",
    nil, // use default HTTP client
    false, // not required
))
```

### Custom Health Checkers

```go
// Using FuncChecker for custom logic
obs.RegisterHealthChecker(health.NewFuncChecker(
    "feature-flags",
    func(ctx context.Context) error {
        // Check feature flag service
        return featureFlagClient.Ping(ctx)
    },
    false, // not required
))

// Implementing the Checker interface
type CacheChecker struct {
    cache *cache.Client
}

func (c *CacheChecker) Name() string { return "cache" }
func (c *CacheChecker) Required() bool { return true }
func (c *CacheChecker) Check(ctx context.Context) health.Result {
    start := time.Now()
    if err := c.cache.Ping(ctx); err != nil {
        return health.NewUnhealthyResult(time.Since(start), err)
    }
    return health.NewHealthyResult(time.Since(start))
}

obs.RegisterHealthChecker(&CacheChecker{cache: cacheClient})
```

### Health Endpoints

The health handler provides these endpoints:

| Endpoint | Purpose | Kubernetes Probe |
|----------|---------|------------------|
| `/health/live` | Liveness check | `livenessProbe` |
| `/health/ready` | Readiness check | `readinessProbe` |
| `/health/startup` | Startup check | `startupProbe` |
| `/health` | Full health report | Debugging |

### Health Response Format

```json
{
  "status": "healthy",
  "checks": {
    "postgres": {
      "status": "healthy",
      "duration_ms": 2,
      "timestamp": "2024-01-15T10:30:00Z"
    },
    "redis": {
      "status": "healthy",
      "duration_ms": 1,
      "timestamp": "2024-01-15T10:30:00Z"
    },
    "kafka": {
      "status": "degraded",
      "duration_ms": 50,
      "error": "high latency",
      "timestamp": "2024-01-15T10:30:00Z"
    }
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Standalone Health Manager

```go
import "github.com/Dorico-Dynamics/txova-go-observability/health"

// Create standalone health manager
manager := health.NewManager(health.ManagerConfig{
    Timeout:            5 * time.Second,
    CacheTTL:           30 * time.Second,
    BackgroundInterval: 30 * time.Second,
    FailureThreshold:   3,
})

// Register checkers
manager.Register(health.NewPostgresChecker("postgres", db, true))

// Start background checks
manager.StartBackground(ctx)
defer manager.StopBackground()

// Create HTTP handler
handler := health.NewHandler(manager)
handler.RegisterRoutes(mux)
```

## Integration Patterns

### With txova-go-core

```go
import (
    "github.com/Dorico-Dynamics/txova-go-core/app"
    "github.com/Dorico-Dynamics/txova-go-observability"
)

func main() {
    ctx := context.Background()
    
    obs, err := observability.New(ctx, &observability.Config{
        // ... config
    })
    if err != nil {
        log.Fatal(err)
    }

    // Observability implements app interfaces
    application := app.New(
        app.WithInitializer(obs),  // calls obs.Initialize()
        app.WithCloser(obs),       // calls obs.Close()
        app.WithHealthChecker(obs), // calls obs.HealthCheck()
    )

    application.Run(ctx)
}
```

### Middleware Chain

```go
// Combine observability middleware with other middleware
handler := obs.HTTPMiddleware()(
    authMiddleware(
        rateLimitMiddleware(
            yourHandler,
        ),
    ),
)
```

### Testing

```go
import (
    "testing"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/Dorico-Dynamics/txova-go-observability"
    "github.com/Dorico-Dynamics/txova-go-observability/metrics"
    "github.com/Dorico-Dynamics/txova-go-observability/tracing"
)

func TestWithObservability(t *testing.T) {
    ctx := context.Background()

    // Use isolated registry for tests
    registry := prometheus.NewRegistry()

    obs, err := observability.New(ctx, &observability.Config{
        MetricsEnabled: true,
        TracingEnabled: true,
        HealthEnabled:  true,
        Metrics: metrics.Config{
            Registry: registry,
        },
        Tracing: tracing.Config{
            ServiceName: "test-service",
            Exporter:    tracing.ExporterNone, // No export in tests
        },
    })
    if err != nil {
        t.Fatal(err)
    }
    defer obs.Close(ctx)

    // Run your tests...
}
```

### Environment-Based Configuration

```go
func loadConfig() *observability.Config {
    cfg := observability.DefaultConfig()

    // Override from environment
    if endpoint := os.Getenv("OTEL_EXPORTER_ENDPOINT"); endpoint != "" {
        cfg.Tracing.Endpoint = endpoint
    }

    if rate := os.Getenv("TRACE_SAMPLE_RATE"); rate != "" {
        if r, err := strconv.ParseFloat(rate, 64); err == nil {
            cfg.Tracing.SampleRate = r
        }
    }

    if os.Getenv("METRICS_DISABLED") == "true" {
        cfg.MetricsEnabled = false
    }

    return &cfg
}
```

## Best Practices

1. **Use the unified Observability struct** - It ensures all components are properly coordinated.

2. **Register health checkers for all dependencies** - Mark critical dependencies as required.

3. **Set appropriate sample rates** - Use 1.0 in development, lower in production (0.01-0.1).

4. **Add context to spans** - Include relevant business context as attributes.

5. **Use meaningful span names** - Follow the pattern `verb-noun` (e.g., `process-payment`).

6. **Handle errors properly** - Record errors on spans before returning.

7. **Expose /metrics endpoint** - Required for Prometheus scraping.

8. **Configure health check timeouts** - Avoid blocking health checks with slow dependencies.
