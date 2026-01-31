# txova-go-observability

Unified observability library providing Prometheus metrics, OpenTelemetry tracing, and health check utilities for Txova services.

## Overview

`txova-go-observability` provides comprehensive observability capabilities for Txova services through a single, cohesive API. It includes Prometheus metrics collection, OpenTelemetry distributed tracing, and Kubernetes-compatible health checks.

**Module:** `github.com/Dorico-Dynamics/txova-go-observability`

## Features

- **Unified API** - Single entry point for all observability features
- **Prometheus Metrics** - HTTP, database, Redis, Kafka, and business metrics
- **OpenTelemetry Tracing** - Distributed tracing with W3C/B3 propagation
- **Health Checks** - Liveness, readiness, and startup probes
- **txova-go-core Integration** - Implements `app.Initializer`, `app.Closer`, and `app.HealthChecker` interfaces

## Packages

| Package | Description |
|---------|-------------|
| `observability` | Unified entry point combining all features |
| `metrics` | Prometheus metric collectors |
| `tracing` | OpenTelemetry tracer setup and middleware |
| `health` | Health check manager and HTTP handlers |

## Installation

```bash
go get github.com/Dorico-Dynamics/txova-go-observability
```

## Quick Start

The recommended approach is to use the unified `Observability` struct:

```go
package main

import (
    "context"
    "net/http"

    "github.com/Dorico-Dynamics/txova-go-observability"
)

func main() {
    ctx := context.Background()

    // Create observability with configuration
    obs, err := observability.New(ctx, &observability.Config{
        MetricsEnabled: true,
        TracingEnabled: true,
        HealthEnabled:  true,
        Tracing: tracing.Config{
            ServiceName:    "ride-service",
            ServiceVersion: "1.0.0",
            Endpoint:       "otel-collector:4318",
            SampleRate:     0.1,
        },
    })
    if err != nil {
        log.Fatal(err)
    }
    defer obs.Close(ctx)

    // Start background health checks
    obs.Initialize(ctx)

    // Use HTTP middleware (adds tracing + metrics)
    mux := http.NewServeMux()
    handler := obs.HTTPMiddleware()(mux)

    // Register health endpoints
    obs.HealthHandler.RegisterRoutes(mux)

    http.ListenAndServe(":8080", handler)
}
```

## Detailed Usage

See [usage.md](./usage.md) for comprehensive examples and best practices.

## Metric Collectors

All collectors are created via the unified `Observability` struct or individually:

| Collector | Description |
|-----------|-------------|
| `HTTPCollector` | HTTP request metrics |
| `DBCollector` | Database query metrics |
| `RedisCollector` | Redis command metrics |
| `KafkaCollector` | Kafka producer/consumer metrics |
| `RideCollector` | Ride business metrics |
| `DriverCollector` | Driver business metrics |
| `PaymentCollector` | Payment business metrics |
| `SafetyCollector` | Safety incident metrics |

## Health Checkers

Built-in health checkers for common dependencies:

| Checker | Description |
|---------|-------------|
| `PostgresChecker` | PostgreSQL database health |
| `RedisChecker` | Redis connection health |
| `KafkaChecker` | Kafka broker health |
| `HTTPChecker` | External HTTP service health |
| `FuncChecker` | Custom function-based checks |

## Histogram Buckets

| Metric Type | Buckets |
|-------------|---------|
| HTTP latency | 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10 |
| DB latency | 0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1 |
| Fare (MZN) | 50, 100, 250, 500, 1000, 2500, 5000, 10000, 25000 |
| Distance (km) | 1, 2, 5, 10, 20, 50, 100, 200 |

## Dependencies

**Internal:**
- `txova-go-core` (interfaces only)

**External:**
- `github.com/prometheus/client_golang` v1.22+
- `go.opentelemetry.io/otel` v1.35+
- `go.opentelemetry.io/otel/exporters/otlp/otlptrace` v1.35+

## Development

### Requirements

- Go 1.24+
- golangci-lint v2.8+

### Testing

```bash
go test ./...
```

### Linting

```bash
golangci-lint run
```

### Test Coverage Target

> 80%

## License

Proprietary - Dorico Dynamics
