# txova-go-observability Execution Plan

**Version:** 1.2  
**Module:** `github.com/Dorico-Dynamics/txova-go-observability`  
**Target Test Coverage:** >80%  
**Internal Dependencies:** txova-go-core  
**External Dependencies:** prometheus/client_golang v1.23.2, go.opentelemetry.io/otel v1.39.0

---

## Current Status

| Phase | Status | Coverage |
|-------|--------|----------|
| Phase 1: Project Setup | ✅ Complete | - |
| Phase 2: metrics (HTTP & Infrastructure) | ✅ Complete | 94.3% |
| Phase 3: metrics (Business) | ✅ Complete | 94.3% |
| Phase 4: tracing | ✅ Complete | 84.0% |
| Phase 5: health | ✅ Complete | 88.3% |
| Phase 6: Integration & QA | ✅ Complete | 89.0% |

**Overall Test Coverage: ~89%**

---

## Phase 1: Project Setup (Week 1)

### 1.1 Project Initialization
- [x] Initialize Go module with `go mod init github.com/Dorico-Dynamics/txova-go-observability`
- [x] Add dependency on `txova-go-core`
- [x] Add external dependencies (prometheus, otel, otlp)
- [x] Create directory structure for packages (metrics, tracing, health)
- [x] Set up `.gitignore` for Go projects
- [x] Configure golangci-lint with strict rules
- [x] Set up GitHub Actions workflows (test, release)

**Note:** Error handling uses `txova-go-core/errors` package. No separate errors package needed.

**Deliverables:**
- [x] Project structure and CI/CD configuration

---

## Phase 2: Metrics - HTTP & Infrastructure (Week 2)

### 2.1 Package: `metrics` - Core Infrastructure
- [x] Define `Config` struct with namespace, subsystem, and bucket configurations
- [x] Implement histogram bucket presets (HTTP, DB, duration, fare)
- [x] Implement label validation (cardinality guards < 100 unique values)
- [x] Write tests for configuration and bucket presets

### 2.2 Package: `metrics` - HTTP Metrics
- [x] Implement `http_requests_total` counter (method, path, status labels)
- [x] Implement `http_request_duration_seconds` histogram (buckets: 0.005-10s)
- [x] Implement `http_request_size_bytes` histogram
- [x] Implement `http_response_size_bytes` histogram
- [x] Implement `http_requests_in_flight` gauge
- [x] Implement `HTTPCollector` struct implementing `server.MetricsCollector`
- [x] Write tests using prometheus/testutil

### 2.3 Package: `metrics` - Database Metrics
- [x] Implement `db_connections_total` gauge (pool, state labels)
- [x] Implement `db_query_duration_seconds` histogram (buckets: 0.001-1s)
- [x] Implement `db_query_errors_total` counter (operation, error labels)
- [x] Implement `db_transaction_duration_seconds` histogram
- [x] Implement `DBCollector` struct with query instrumentation wrapper
- [x] Write tests for all database metrics

### 2.4 Package: `metrics` - Redis Metrics
- [x] Implement `redis_commands_total` counter (command label)
- [x] Implement `redis_command_duration_seconds` histogram (command label)
- [x] Implement `redis_cache_hits_total` counter (cache label)
- [x] Implement `redis_cache_misses_total` counter (cache label)
- [x] Implement `RedisCollector` struct with command instrumentation
- [x] Write tests for all Redis metrics

### 2.5 Package: `metrics` - Kafka Metrics
- [x] Implement `kafka_messages_produced_total` counter (topic label)
- [x] Implement `kafka_messages_consumed_total` counter (topic, group labels)
- [x] Implement `kafka_consumer_lag` gauge (topic, partition, group labels)
- [x] Implement `kafka_produce_errors_total` counter (topic label)
- [x] Implement `kafka_consume_errors_total` counter (topic label)
- [x] Implement `KafkaCollector` struct with producer/consumer hooks
- [x] Write tests for all Kafka metrics

**Deliverables:**
- [x] `metrics/` package with HTTP, DB, Redis, and Kafka collectors
- [x] All collectors registered with "txova" namespace
- [x] Tests verifying metric registration and label handling (94.3% coverage)

---

## Phase 3: Metrics - Business (Week 3)

### 3.1 Package: `metrics` - Ride Metrics
- [x] Implement `rides_requested_total` counter (service_type, city labels)
- [x] Implement `rides_completed_total` counter (service_type, city labels)
- [x] Implement `rides_cancelled_total` counter (cancelled_by, reason labels)
- [x] Implement `ride_duration_seconds` histogram (service_type label, buckets: 60-3600s)
- [x] Implement `ride_distance_km` histogram (service_type label)
- [x] Implement `ride_fare_mzn` histogram (service_type label, buckets: 50-25000 MZN)
- [x] Implement `ride_wait_time_seconds` histogram (service_type label)
- [x] Implement `RideCollector` struct with recording methods
- [x] Write tests for all ride metrics

### 3.2 Package: `metrics` - Driver Metrics
- [x] Implement `drivers_online_total` gauge (city, service_type labels)
- [x] Implement `driver_acceptance_rate` gauge (driver_id label)
- [x] Implement `driver_rating_average` gauge
- [x] Implement `driver_earnings_mzn` counter (driver_id label)
- [x] Implement `DriverCollector` struct with recording methods
- [x] Write tests for all driver metrics

### 3.3 Package: `metrics` - Payment Metrics
- [x] Implement `payments_total` counter (method, status labels)
- [x] Implement `payment_amount_mzn` histogram (method label)
- [x] Implement `payment_processing_seconds` histogram (method label)
- [x] Implement `refunds_total` counter (reason label)
- [x] Implement `PaymentCollector` struct with recording methods
- [x] Write tests for all payment metrics

### 3.4 Package: `metrics` - Safety Metrics
- [x] Implement `emergencies_triggered_total` counter (type, city labels)
- [x] Implement `incidents_reported_total` counter (severity label)
- [x] Implement `trip_shares_total` counter
- [x] Implement `SafetyCollector` struct with recording methods
- [x] Write tests for all safety metrics

**Deliverables:**
- [x] Business metrics collectors (ride, driver, payment, safety)
- [x] All histograms using appropriate bucket distributions
- [x] Tests verifying metric recording and label cardinality (94.3% coverage)

---

## Phase 4: Tracing (Week 4)

### 4.1 Package: `tracing` - Tracer Setup
- [x] Define `Config` struct (service_name, endpoint, sample_rate, propagation)
- [x] Implement tracer initialization with OTLP exporter
- [x] Implement configurable sampling (0.0-1.0)
- [x] Implement graceful shutdown for tracer provider
- [x] Write tests for tracer initialization

### 4.2 Package: `tracing` - Span Attributes
- [x] Define standard attribute keys (service.name, service.version, user.id, request.id)
- [x] Define HTTP attribute keys (http.method, http.route, http.status_code)
- [x] Define DB attribute keys (db.system, db.operation)
- [x] Define messaging attribute keys (messaging.system, messaging.destination)
- [x] Implement attribute helper functions
- [x] Write tests for attribute creation

### 4.3 Package: `tracing` - Context Propagation
- [x] Implement W3C trace context extraction from HTTP headers (traceparent, tracestate)
- [x] Implement W3C trace context injection into HTTP headers
- [x] Implement Kafka header injection for producer
- [x] Implement Kafka header extraction for consumer
- [x] Implement X-Request-ID correlation with trace context
- [x] Implement baggage support for cross-service data
- [x] Write tests for header injection/extraction

### 4.4 Package: `tracing` - HTTP Middleware
- [x] Implement tracing middleware for HTTP handlers
- [x] Auto-create spans for incoming requests
- [x] Extract context from incoming headers
- [x] Inject context into outgoing requests
- [x] Write tests for middleware

**Deliverables:**
- [x] `tracing/` package with OpenTelemetry integration
- [x] W3C trace context propagation (HTTP and Kafka)
- [x] HTTP middleware compatible with standard library
- [x] Tests for tracer setup, propagation, and middleware (84.0% coverage)

---

## Phase 5: Health Checks (Week 5)

### 5.1 Package: `health` - Core Infrastructure
- [x] Define `Checker` interface compatible with `app.HealthChecker`
- [x] Define `Result` struct (status, duration_ms, error)
- [x] Define `Report` struct (status, checks map, timestamp)
- [x] Implement `Status` type (healthy, unhealthy, degraded)
- [x] Implement `Manager` for registering and running checks
- [x] Write tests for manager and result aggregation

### 5.2 Package: `health` - Component Checks
- [x] Implement `PostgresChecker` with PING and timeout
- [x] Implement `RedisChecker` with PING command and timeout
- [x] Implement `KafkaChecker` with metadata request and timeout
- [x] Implement `HTTPChecker` for external API health endpoints
- [x] Implement `FuncChecker` for custom health check functions
- [x] Mark checks as required or optional
- [x] Write tests for all checkers (with mocks)

### 5.3 Package: `health` - Caching & Background Checks
- [x] Implement result caching with configurable TTL
- [x] Implement background checking at configurable interval (default: 30s)
- [x] Implement failure threshold before marking unhealthy (default: 3)
- [x] Implement per-component timeout configuration (default: 5s)
- [x] Log all health check failures using slog
- [x] Write tests for caching and background checks

### 5.4 Package: `health` - HTTP Handlers
- [x] Implement `/health/live` endpoint (liveness probe)
- [x] Implement `/health/ready` endpoint (readiness probe)
- [x] Implement `/health/startup` endpoint (startup probe)
- [x] Return 200 when healthy, 503 when unhealthy
- [x] Return 200 if only optional checks fail (degraded status)
- [x] Implement JSON response format matching PRD spec
- [x] Write tests for all endpoints

**Deliverables:**
- [x] `health/` package with component checks and HTTP handlers
- [x] Caching and background check mechanism
- [x] Integration with txova-go-core app lifecycle interfaces
- [x] Tests for all components (88.3% coverage)

---

## Phase 6: Integration & Quality Assurance (Week 6)

### 6.1 Integration with txova-go-core
- [x] Implement `Observability` struct as central entry point
- [x] Implement `Initialize` method for startup (app.Initializer)
- [x] Implement `Close` method for shutdown (app.Closer)
- [x] Implement `HealthCheck` method for health reporting (app.HealthChecker)
- [x] Implement `HTTPMiddleware` for combined metrics and tracing
- [x] Write integration tests

### 6.2 Cross-Package Integration
- [x] Verify all packages work together without circular dependencies
- [x] Ensure consistent error handling using fmt.Errorf with %w
- [x] Validate context propagation through middleware chain
- [x] Test full request flow (metrics + tracing + health)

### 6.3 Quality Assurance
- [x] Run full test suite and verify >80% coverage
- [x] Run `go vet` and address all warnings
- [x] Test with `go build` for all target platforms

### 6.4 Documentation
- [x] Ensure all exported types and functions have godoc comments

**Deliverables:**
- [x] Complete, tested library
- [x] >80% test coverage verified (~89% overall)
- [x] All packages building without errors

---

## Success Criteria

| Criteria | Target | Current |
|----------|--------|---------|
| Test Coverage | >80% | ~89% ✅ |
| metrics coverage | >80% | 94.3% ✅ |
| tracing coverage | >80% | 84.0% ✅ |
| health coverage | >80% | 88.3% ✅ |
| root coverage | >80% | 89.0% ✅ |
| Linting Errors | 0 | 0 ✅ |
| go vet Warnings | 0 | 0 ✅ |

---

## Package Dependency Order

```
metrics (no internal dependencies)
    ↓
tracing (no internal dependencies)
    ↓
health (no internal dependencies)
    ↓
root package (imports metrics, tracing, health)
```

---

## Histogram Bucket Reference

| Metric Type | Buckets |
|-------------|---------|
| HTTP latency (seconds) | 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10 |
| DB latency (seconds) | 0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1 |
| Duration (seconds) | 60, 300, 600, 900, 1800, 3600 |
| Fare amount (MZN) | 50, 100, 250, 500, 1000, 2500, 5000, 10000, 25000 |

---

## Files Created

### metrics/
- `config.go` - Configuration with namespace, subsystem, registry
- `buckets.go` - Histogram bucket presets
- `http.go` - HTTPCollector for HTTP metrics
- `database.go` - DBCollector for database metrics
- `redis.go` - RedisCollector for Redis metrics
- `kafka.go` - KafkaCollector for Kafka metrics
- `ride.go` - RideCollector for ride business metrics
- `driver.go` - DriverCollector for driver business metrics
- `payment.go` - PaymentCollector for payment metrics
- `safety.go` - SafetyCollector for safety metrics
- `*_test.go` - Tests for all collectors

### tracing/
- `config.go` - Configuration for tracer setup
- `tracer.go` - Tracer wrapper with lifecycle management
- `attributes.go` - Span attribute helper functions
- `propagation.go` - W3C context propagation
- `middleware.go` - HTTP middleware and RoundTripper
- `*_test.go` - Tests for all components

### health/
- `status.go` - Status type (healthy, unhealthy, degraded)
- `checker.go` - Checker interface and Result/Report types
- `checkers.go` - Postgres, Redis, Kafka, HTTP, Func checkers
- `manager.go` - Manager for coordinating health checks
- `handler.go` - HTTP handlers for health endpoints
- `*_test.go` - Tests for all components

### root package
- `observability.go` - Central Observability struct
- `observability_test.go` - Integration tests
