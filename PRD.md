# txova-go-observability

## Overview
Observability library providing Prometheus metrics, OpenTelemetry tracing, and health check utilities for monitoring Txova services.

**Module:** `github.com/txova/txova-go-observability`

---

## Packages

### `metrics` - Prometheus Metrics

#### HTTP Metrics
| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| http_requests_total | Counter | method, path, status | Total requests |
| http_request_duration_seconds | Histogram | method, path | Request latency |
| http_request_size_bytes | Histogram | method, path | Request body size |
| http_response_size_bytes | Histogram | method, path | Response body size |
| http_requests_in_flight | Gauge | - | Current active requests |

#### Database Metrics
| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| db_connections_total | Gauge | pool, state | Connection pool stats |
| db_query_duration_seconds | Histogram | operation | Query latency |
| db_query_errors_total | Counter | operation, error | Query failures |
| db_transaction_duration_seconds | Histogram | - | Transaction latency |

#### Redis Metrics
| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| redis_commands_total | Counter | command | Commands executed |
| redis_command_duration_seconds | Histogram | command | Command latency |
| redis_cache_hits_total | Counter | cache | Cache hits |
| redis_cache_misses_total | Counter | cache | Cache misses |

#### Kafka Metrics
| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| kafka_messages_produced_total | Counter | topic | Messages published |
| kafka_messages_consumed_total | Counter | topic, group | Messages consumed |
| kafka_consumer_lag | Gauge | topic, partition, group | Consumer lag |
| kafka_produce_errors_total | Counter | topic | Publish failures |
| kafka_consume_errors_total | Counter | topic | Consume failures |

---

#### Business Metrics

**Ride Metrics:**
| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| rides_requested_total | Counter | service_type, city | Rides requested |
| rides_completed_total | Counter | service_type, city | Rides completed |
| rides_cancelled_total | Counter | cancelled_by, reason | Rides cancelled |
| ride_duration_seconds | Histogram | service_type | Trip duration |
| ride_distance_km | Histogram | service_type | Trip distance |
| ride_fare_mzn | Histogram | service_type | Trip fare |
| ride_wait_time_seconds | Histogram | service_type | Time to match |

**Driver Metrics:**
| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| drivers_online_total | Gauge | city, service_type | Online drivers |
| driver_acceptance_rate | Gauge | driver_id | Acceptance rate |
| driver_rating_average | Gauge | - | Average rating |
| driver_earnings_mzn | Counter | driver_id | Total earnings |

**Payment Metrics:**
| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| payments_total | Counter | method, status | Payment attempts |
| payment_amount_mzn | Histogram | method | Payment amounts |
| payment_processing_seconds | Histogram | method | Processing time |
| refunds_total | Counter | reason | Refunds issued |

**Safety Metrics:**
| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| emergencies_triggered_total | Counter | type, city | SOS activations |
| incidents_reported_total | Counter | severity | Incidents |
| trip_shares_total | Counter | - | Trip sharing |

---

#### Metric Registration
| Requirement | Description |
|-------------|-------------|
| Namespace | "txova" prefix for all metrics |
| Subsystem | Service name (e.g., "ride_service") |
| Labels | Keep cardinality low (<100 unique values) |
| Buckets | Appropriate for expected distributions |

**Histogram Buckets:**
| Metric Type | Buckets |
|-------------|---------|
| HTTP latency | 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10 |
| DB latency | 0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1 |
| Fare amount | 50, 100, 250, 500, 1000, 2500, 5000, 10000, 25000 |
| Duration | 60, 300, 600, 900, 1800, 3600 |

---

### `tracing` - OpenTelemetry Tracing

#### Tracer Configuration
| Setting | Description |
|---------|-------------|
| service_name | Name for traces (e.g., "ride-service") |
| endpoint | OTLP collector endpoint |
| sample_rate | Sampling rate (0.0-1.0) |
| propagation | W3C trace context |

#### Span Attributes
| Attribute | Description |
|-----------|-------------|
| service.name | Service identifier |
| service.version | Service version |
| user.id | Authenticated user |
| request.id | Correlation ID |
| http.method | HTTP method |
| http.route | Route pattern |
| http.status_code | Response status |
| db.system | Database type |
| db.operation | Query type |
| messaging.system | Kafka |
| messaging.destination | Topic name |

#### Span Creation
| Requirement | Description |
|-------------|-------------|
| HTTP handlers | Auto-instrument via middleware |
| Database queries | Wrap with span |
| Kafka publish | Include in producer |
| Kafka consume | Extract from message |
| HTTP clients | Propagate context |

---

#### Trace Context Propagation
| Header | Description |
|--------|-------------|
| traceparent | W3C trace context |
| tracestate | W3C trace state |
| X-Request-ID | Correlation ID (also propagated) |

**Requirements:**
- Inject trace context in outgoing HTTP requests
- Extract trace context from incoming HTTP requests
- Propagate via Kafka message headers
- Support baggage for cross-service data

---

### `health` - Health Checks

#### Health Check Types
| Type | Endpoint | Description |
|------|----------|-------------|
| Liveness | /health/live | Is process running? |
| Readiness | /health/ready | Can accept traffic? |
| Startup | /health/startup | Has started successfully? |

#### Check Components
| Component | Check | Required |
|-----------|-------|----------|
| PostgreSQL | Ping with timeout | Yes |
| Redis | PING command | Yes |
| Kafka | Metadata request | Yes |
| External APIs | Health endpoint | No |

#### Response Format
| Field | Description |
|-------|-------------|
| status | "healthy" or "unhealthy" |
| checks | Map of component checks |
| checks[].status | Component status |
| checks[].duration_ms | Check duration |
| checks[].error | Error message (if failed) |
| timestamp | Check time |

**HTTP Status:**
| Scenario | Status |
|----------|--------|
| All checks pass | 200 |
| Any required check fails | 503 |
| Only optional check fails | 200 |

---

#### Health Check Configuration
| Setting | Default | Description |
|---------|---------|-------------|
| timeout | 5s | Check timeout |
| interval | 30s | Background check interval |
| threshold | 3 | Failures before unhealthy |

**Requirements:**
- Cache results to avoid overload
- Background checking for slow dependencies
- Configurable per-component timeouts
- Log all health check failures

---

### `alerts` - Alert Definitions

**Critical Alerts (Page):**
| Alert | Condition | Description |
|-------|-----------|-------------|
| ServiceDown | up == 0 for 2m | Service not responding |
| HighErrorRate | 5xx > 5% for 5m | High error rate |
| DatabaseDown | db_up == 0 for 1m | Database unavailable |
| KafkaLag | lag > 10000 for 5m | Consumer falling behind |
| PaymentFailures | failures > 10% for 5m | Payment issues |

**Warning Alerts (Notify):**
| Alert | Condition | Description |
|-------|-----------|-------------|
| HighLatency | p99 > 5s for 10m | Slow responses |
| LowCacheHit | hit_rate < 50% for 30m | Cache ineffective |
| HighCPU | cpu > 80% for 15m | Resource pressure |
| HighMemory | memory > 85% for 15m | Memory pressure |
| LowDrivers | online < 10 for 30m | Driver shortage |

---

## Dashboard Requirements

**Service Dashboard:**
| Panel | Description |
|-------|-------------|
| Request Rate | Requests per second |
| Error Rate | Percentage of 5xx |
| Latency | P50, P95, P99 |
| Active Requests | In-flight count |
| Resource Usage | CPU, Memory |

**Business Dashboard:**
| Panel | Description |
|-------|-------------|
| Active Rides | Current in-progress |
| Drivers Online | By city, service type |
| Ride Completion Rate | Success percentage |
| Average Wait Time | Time to driver match |
| Revenue | Hourly, daily totals |

---

## Dependencies

**Internal:**
- `txova-go-core`

**External:**
- `github.com/prometheus/client_golang` — Prometheus client
- `go.opentelemetry.io/otel` — OpenTelemetry SDK
- `go.opentelemetry.io/otel/exporters/otlp` — OTLP exporter

---

## Success Metrics
| Metric | Target |
|--------|--------|
| Test coverage | > 80% |
| Metric collection overhead | < 1% CPU |
| Trace sampling accuracy | Within 5% |
| Health check latency | < 100ms |
| Alert accuracy | > 95% |
