# txova-go-observability

Observability library providing Prometheus metrics, OpenTelemetry tracing, and health check utilities for monitoring Txova services.

## Overview

`txova-go-observability` provides comprehensive observability capabilities for Txova services, including Prometheus metrics for monitoring, OpenTelemetry distributed tracing, and Kubernetes-compatible health checks.

**Module:** `github.com/txova/txova-go-observability`

## Features

- **Prometheus Metrics** - HTTP, database, Redis, Kafka, and business metrics
- **OpenTelemetry Tracing** - Distributed tracing with context propagation
- **Health Checks** - Liveness, readiness, and startup probes
- **Alert Definitions** - Pre-configured alert rules

## Packages

| Package | Description |
|---------|-------------|
| `metrics` | Prometheus metric collectors |
| `tracing` | OpenTelemetry tracer setup |
| `health` | Health check endpoints |
| `alerts` | Alert rule definitions |

## Installation

```bash
go get github.com/txova/txova-go-observability
```

## Usage

### Prometheus Metrics

```go
import "github.com/txova/txova-go-observability/metrics"

// Initialize metrics
m := metrics.New(metrics.Config{
    Namespace: "txova",
    Subsystem: "ride_service",
})

// HTTP metrics (use as middleware)
r.Use(m.HTTPMiddleware())

// Record business metrics
m.RidesRequested.WithLabelValues("standard", "maputo").Inc()
m.RideCompleted.WithLabelValues("standard", "maputo").Inc()
m.RideFare.WithLabelValues("standard").Observe(float64(fare.Amount()))
m.DriversOnline.WithLabelValues("maputo", "standard").Set(float64(count))

// Database metrics
m.RecordQuery("select_user", duration)
m.RecordQueryError("select_user", "timeout")

// Cache metrics
m.RecordCacheHit("user_profile")
m.RecordCacheMiss("user_profile")
```

### Available Metrics

**HTTP Metrics:**
- `http_requests_total` - Total requests (method, path, status)
- `http_request_duration_seconds` - Request latency histogram
- `http_requests_in_flight` - Current active requests

**Database Metrics:**
- `db_connections_total` - Connection pool stats
- `db_query_duration_seconds` - Query latency
- `db_query_errors_total` - Query failures

**Business Metrics:**
- `rides_requested_total` - Rides requested (service_type, city)
- `rides_completed_total` - Rides completed
- `rides_cancelled_total` - Rides cancelled (cancelled_by, reason)
- `drivers_online_total` - Online drivers (city, service_type)
- `payments_total` - Payment attempts (method, status)

### OpenTelemetry Tracing

```go
import "github.com/txova/txova-go-observability/tracing"

// Initialize tracer
tracer, err := tracing.New(tracing.Config{
    ServiceName: "ride-service",
    Endpoint:    "otel-collector:4317",
    SampleRate:  0.1, // 10% sampling
})
defer tracer.Shutdown(ctx)

// Create span
ctx, span := tracer.Start(ctx, "process-ride-request")
defer span.End()

// Add attributes
span.SetAttributes(
    attribute.String("ride.id", rideID.String()),
    attribute.String("user.id", userID.String()),
)

// Record errors
if err != nil {
    span.RecordError(err)
    span.SetStatus(codes.Error, err.Error())
}
```

### Trace Context Propagation

```go
import "github.com/txova/txova-go-observability/tracing"

// HTTP middleware automatically propagates trace context
r.Use(tracing.HTTPMiddleware(tracer))

// For outgoing HTTP requests
req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
tracing.InjectHTTP(ctx, req) // Adds traceparent header

// For Kafka messages
tracing.InjectKafka(ctx, message) // Adds trace headers
```

### Health Checks

```go
import "github.com/txova/txova-go-observability/health"

checker := health.New(health.Config{
    Timeout:  5 * time.Second,
    Interval: 30 * time.Second,
})

// Register checks
checker.AddCheck("postgres", health.PostgresCheck(pool))
checker.AddCheck("redis", health.RedisCheck(redisClient))
checker.AddCheck("kafka", health.KafkaCheck(kafkaClient))

// Mount endpoints
r.Get("/health/live", checker.LivenessHandler())
r.Get("/health/ready", checker.ReadinessHandler())
r.Get("/health/startup", checker.StartupHandler())
```

### Health Response Format

```json
{
  "status": "healthy",
  "checks": {
    "postgres": {
      "status": "healthy",
      "duration_ms": 2
    },
    "redis": {
      "status": "healthy", 
      "duration_ms": 1
    }
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Alert Rules

**Critical Alerts (Page):**
- `ServiceDown` - Service not responding for 2m
- `HighErrorRate` - 5xx rate > 5% for 5m
- `DatabaseDown` - Database unavailable for 1m
- `KafkaLag` - Consumer lag > 10,000 for 5m
- `PaymentFailures` - Payment failure rate > 10% for 5m

**Warning Alerts (Notify):**
- `HighLatency` - P99 latency > 5s for 10m
- `LowCacheHit` - Cache hit rate < 50% for 30m
- `HighCPU` - CPU > 80% for 15m
- `LowDrivers` - Online drivers < 10 for 30m

## Histogram Buckets

| Metric Type | Buckets |
|-------------|---------|
| HTTP latency | 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10 |
| DB latency | 0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1 |
| Fare amount | 50, 100, 250, 500, 1000, 2500, 5000, 10000, 25000 |

## Dependencies

**Internal:**
- `txova-go-core`

**External:**
- `github.com/prometheus/client_golang` - Prometheus client
- `go.opentelemetry.io/otel` - OpenTelemetry SDK
- `go.opentelemetry.io/otel/exporters/otlp` - OTLP exporter

## Development

### Requirements

- Go 1.25+

### Testing

```bash
go test ./...
```

### Test Coverage Target

> 80%

## License

Proprietary - Dorico Dynamics
