# txova-go-observability Execution Plan

**Version:** 1.1  
**Module:** `github.com/Dorico-Dynamics/txova-go-observability`  
**Target Test Coverage:** >80%  
**Internal Dependencies:** txova-go-core  
**External Dependencies:** prometheus/client_golang, go.opentelemetry.io/otel, go.opentelemetry.io/otel/exporters/otlp

---

## Current Status

| Phase | Status | Coverage |
|-------|--------|----------|
| Phase 1: Project Setup | ✅ Complete | - |
| Phase 2: metrics (HTTP & Infrastructure) | Not Started | - |
| Phase 3: metrics (Business) | Not Started | - |
| Phase 4: tracing | Not Started | - |
| Phase 5: health | Not Started | - |
| Phase 6: Integration & QA | Not Started | - |

**Overall Test Coverage: 0%**

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
- [ ] Define `Config` struct with namespace, subsystem, and bucket configurations
- [ ] Implement histogram bucket presets (HTTP, DB, duration, fare)
- [ ] Implement label validation (cardinality guards < 100 unique values)
- [ ] Write tests for configuration and bucket presets

### 2.2 Package: `metrics` - HTTP Metrics
- [ ] Implement `http_requests_total` counter (method, path, status labels)
- [ ] Implement `http_request_duration_seconds` histogram (buckets: 0.005-10s)
- [ ] Implement `http_request_size_bytes` histogram
- [ ] Implement `http_response_size_bytes` histogram
- [ ] Implement `http_requests_in_flight` gauge
- [ ] Implement `HTTPCollector` struct implementing `server.MetricsCollector`
- [ ] Write tests using prometheus/testutil

### 2.3 Package: `metrics` - Database Metrics
- [ ] Implement `db_connections_total` gauge (pool, state labels)
- [ ] Implement `db_query_duration_seconds` histogram (buckets: 0.001-1s)
- [ ] Implement `db_query_errors_total` counter (operation, error labels)
- [ ] Implement `db_transaction_duration_seconds` histogram
- [ ] Implement `DBCollector` struct with query instrumentation wrapper
- [ ] Write tests for all database metrics

### 2.4 Package: `metrics` - Redis Metrics
- [ ] Implement `redis_commands_total` counter (command label)
- [ ] Implement `redis_command_duration_seconds` histogram (command label)
- [ ] Implement `redis_cache_hits_total` counter (cache label)
- [ ] Implement `redis_cache_misses_total` counter (cache label)
- [ ] Implement `RedisCollector` struct with command instrumentation
- [ ] Write tests for all Redis metrics

### 2.5 Package: `metrics` - Kafka Metrics
- [ ] Implement `kafka_messages_produced_total` counter (topic label)
- [ ] Implement `kafka_messages_consumed_total` counter (topic, group labels)
- [ ] Implement `kafka_consumer_lag` gauge (topic, partition, group labels)
- [ ] Implement `kafka_produce_errors_total` counter (topic label)
- [ ] Implement `kafka_consume_errors_total` counter (topic label)
- [ ] Implement `KafkaCollector` struct with producer/consumer hooks
- [ ] Write tests for all Kafka metrics

**Deliverables:**
- [ ] `metrics/` package with HTTP, DB, Redis, and Kafka collectors
- [ ] All collectors registered with "txova" namespace
- [ ] Tests verifying metric registration and label handling (>80% coverage)

---

## Phase 3: Metrics - Business (Week 3)

### 3.1 Package: `metrics` - Ride Metrics
- [ ] Implement `rides_requested_total` counter (service_type, city labels)
- [ ] Implement `rides_completed_total` counter (service_type, city labels)
- [ ] Implement `rides_cancelled_total` counter (cancelled_by, reason labels)
- [ ] Implement `ride_duration_seconds` histogram (service_type label, buckets: 60-3600s)
- [ ] Implement `ride_distance_km` histogram (service_type label)
- [ ] Implement `ride_fare_mzn` histogram (service_type label, buckets: 50-25000 MZN)
- [ ] Implement `ride_wait_time_seconds` histogram (service_type label)
- [ ] Implement `RideCollector` struct with recording methods
- [ ] Write tests for all ride metrics

### 3.2 Package: `metrics` - Driver Metrics
- [ ] Implement `drivers_online_total` gauge (city, service_type labels)
- [ ] Implement `driver_acceptance_rate` gauge (driver_id label)
- [ ] Implement `driver_rating_average` gauge
- [ ] Implement `driver_earnings_mzn` counter (driver_id label)
- [ ] Implement `DriverCollector` struct with recording methods
- [ ] Write tests for all driver metrics

### 3.3 Package: `metrics` - Payment Metrics
- [ ] Implement `payments_total` counter (method, status labels)
- [ ] Implement `payment_amount_mzn` histogram (method label)
- [ ] Implement `payment_processing_seconds` histogram (method label)
- [ ] Implement `refunds_total` counter (reason label)
- [ ] Implement `PaymentCollector` struct with recording methods
- [ ] Write tests for all payment metrics

### 3.4 Package: `metrics` - Safety Metrics
- [ ] Implement `emergencies_triggered_total` counter (type, city labels)
- [ ] Implement `incidents_reported_total` counter (severity label)
- [ ] Implement `trip_shares_total` counter
- [ ] Implement `SafetyCollector` struct with recording methods
- [ ] Write tests for all safety metrics

**Deliverables:**
- [ ] Business metrics collectors (ride, driver, payment, safety)
- [ ] All histograms using appropriate bucket distributions
- [ ] Tests verifying metric recording and label cardinality (>80% coverage)

---

## Phase 4: Tracing (Week 4)

### 4.1 Package: `tracing` - Tracer Setup
- [ ] Define `Config` struct (service_name, endpoint, sample_rate, propagation)
- [ ] Implement tracer initialization with OTLP exporter
- [ ] Implement configurable sampling (0.0-1.0)
- [ ] Implement graceful shutdown for tracer provider
- [ ] Write tests for tracer initialization

### 4.2 Package: `tracing` - Span Attributes
- [ ] Define standard attribute keys (service.name, service.version, user.id, request.id)
- [ ] Define HTTP attribute keys (http.method, http.route, http.status_code)
- [ ] Define DB attribute keys (db.system, db.operation)
- [ ] Define messaging attribute keys (messaging.system, messaging.destination)
- [ ] Implement attribute helper functions
- [ ] Write tests for attribute creation

### 4.3 Package: `tracing` - Context Propagation
- [ ] Implement W3C trace context extraction from HTTP headers (traceparent, tracestate)
- [ ] Implement W3C trace context injection into HTTP headers
- [ ] Implement Kafka header injection for producer
- [ ] Implement Kafka header extraction for consumer
- [ ] Implement X-Request-ID correlation with trace context
- [ ] Implement baggage support for cross-service data
- [ ] Write tests for header injection/extraction

### 4.4 Package: `tracing` - HTTP Middleware
- [ ] Implement tracing middleware for Chi router
- [ ] Auto-create spans for incoming requests
- [ ] Extract context from incoming headers
- [ ] Inject context into outgoing requests
- [ ] Write tests for middleware

**Deliverables:**
- [ ] `tracing/` package with OpenTelemetry integration
- [ ] W3C trace context propagation (HTTP and Kafka)
- [ ] Chi-compatible middleware
- [ ] Tests for tracer setup, propagation, and middleware (>80% coverage)

---

## Phase 5: Health Checks (Week 5)

### 5.1 Package: `health` - Core Infrastructure
- [ ] Define `Checker` interface compatible with `app.HealthChecker`
- [ ] Define `CheckResult` struct (status, duration_ms, error)
- [ ] Define `Report` struct (status, checks map, timestamp)
- [ ] Implement `Status` type (healthy, unhealthy)
- [ ] Implement `Manager` for registering and running checks
- [ ] Write tests for manager and result aggregation

### 5.2 Package: `health` - Component Checks
- [ ] Implement `PostgresChecker` with PING and timeout
- [ ] Implement `RedisChecker` with PING command and timeout
- [ ] Implement `KafkaChecker` with metadata request and timeout
- [ ] Implement `HTTPChecker` for external API health endpoints
- [ ] Mark checks as required or optional
- [ ] Write tests for all checkers (with mocks)

### 5.3 Package: `health` - Caching & Background Checks
- [ ] Implement result caching with configurable TTL
- [ ] Implement background checking at configurable interval (default: 30s)
- [ ] Implement failure threshold before marking unhealthy (default: 3)
- [ ] Implement per-component timeout configuration (default: 5s)
- [ ] Log all health check failures using `txova-go-core/logging`
- [ ] Write tests for caching and background checks

### 5.4 Package: `health` - HTTP Handlers
- [ ] Implement `/health/live` endpoint (liveness probe)
- [ ] Implement `/health/ready` endpoint (readiness probe)
- [ ] Implement `/health/startup` endpoint (startup probe)
- [ ] Return 200 when healthy, 503 when unhealthy
- [ ] Return 200 if only optional checks fail
- [ ] Implement JSON response format matching PRD spec
- [ ] Write tests for all endpoints

**Deliverables:**
- [ ] `health/` package with component checks and HTTP handlers
- [ ] Caching and background check mechanism
- [ ] Integration with txova-go-core app lifecycle
- [ ] Tests for all components (>80% coverage)

---

## Phase 6: Integration & Quality Assurance (Week 6)

### 6.1 Integration with txova-go-core
- [ ] Implement `Observability` struct as central entry point
- [ ] Implement `app.Initializer` interface for startup
- [ ] Implement `app.Closer` interface for shutdown
- [ ] Implement `app.HealthChecker` interface for health reporting
- [ ] Implement `server.MetricsCollector` interface for HTTP metrics
- [ ] Provide `Register(app *app.App, server *server.Server)` convenience function
- [ ] Write integration tests with txova-go-core

### 6.2 Cross-Package Integration
- [ ] Verify all packages work together without circular dependencies
- [ ] Ensure consistent error handling using `txova-go-core/errors`
- [ ] Validate context propagation through middleware chain
- [ ] Test full request flow (metrics + tracing + health)

### 6.3 Quality Assurance
- [ ] Run full test suite and verify >80% coverage
- [ ] Run golangci-lint and fix all issues
- [ ] Run `go vet` and address all warnings
- [ ] Test with `go build` for all target platforms
- [ ] Verify metric collection overhead < 1% CPU
- [ ] Verify health check latency < 100ms
- [ ] Verify trace sampling accuracy within 5%

### 6.4 Documentation
- [ ] Ensure all exported types and functions have godoc comments

### 6.5 Release
- [ ] Tag release as v1.0.0 (via CI/CD, not manual)
- [ ] Push to GitHub
- [ ] Verify module is accessible via `go get`

**Deliverables:**
- [ ] Complete, tested library
- [ ] >80% test coverage verified
- [ ] v1.0.0 release tagged and published (via CI/CD)

---

## Success Criteria

| Criteria | Target | Current |
|----------|--------|---------|
| Test Coverage | >80% | 0% |
| Metric Collection Overhead | <1% CPU | Not measured |
| Trace Sampling Accuracy | Within 5% | Not measured |
| Health Check Latency | <100ms | Not measured |
| Linting Errors | 0 | 0 |
| go vet Warnings | 0 | 0 |

---

## Package Dependency Order

```
metrics (uses core/errors for error handling)
    ↓
tracing (uses core/errors, core/context for request IDs)
    ↓
health (uses core/errors, core/logging, implements core/app interfaces)
    ↓
root package (integrates all, provides convenience functions)
```

---

## Integration Points with txova-go-core

| Core Package | Integration |
|--------------|-------------|
| `app` | Implement Initializer, Closer, HealthChecker interfaces |
| `server` | Implement MetricsCollector interface |
| `context` | Extract request_id, user_id, correlation_id for spans/metrics |
| `logging` | Use Logger for observability component logging |
| `errors` | Use AppError for all error handling |

---

## Risk Mitigation

| Risk | Mitigation |
|------|------------|
| High label cardinality causing Prometheus OOM | Implement cardinality guards, validate labels |
| OpenTelemetry version compatibility | Pin specific versions, test with integration |
| Health check timeouts under load | Implement caching, background checks |
| Context propagation edge cases in Kafka | Comprehensive tests for header injection/extraction |
| Metric collection overhead | Benchmark tests, optimize hot paths |

---

## Histogram Bucket Reference

| Metric Type | Buckets |
|-------------|---------|
| HTTP latency (seconds) | 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10 |
| DB latency (seconds) | 0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1 |
| Duration (seconds) | 60, 300, 600, 900, 1800, 3600 |
| Fare amount (MZN) | 50, 100, 250, 500, 1000, 2500, 5000, 10000, 25000 |

---

## Alert Definitions Reference

**Critical Alerts (Page):**
| Alert | Condition |
|-------|-----------|
| ServiceDown | up == 0 for 2m |
| HighErrorRate | 5xx > 5% for 5m |
| DatabaseDown | db_up == 0 for 1m |
| KafkaLag | lag > 10000 for 5m |
| PaymentFailures | failures > 10% for 5m |

**Warning Alerts (Notify):**
| Alert | Condition |
|-------|-----------|
| HighLatency | p99 > 5s for 10m |
| LowCacheHit | hit_rate < 50% for 30m |
| HighCPU | cpu > 80% for 15m |
| HighMemory | memory > 85% for 15m |
| LowDrivers | online < 10 for 30m |
