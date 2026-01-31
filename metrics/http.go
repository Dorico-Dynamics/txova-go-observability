package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// HTTPCollector collects HTTP request metrics.
// It implements the server.MetricsCollector interface from txova-go-core.
type HTTPCollector struct {
	requestsTotal    *prometheus.CounterVec
	requestDuration  *prometheus.HistogramVec
	requestSize      *prometheus.HistogramVec
	responseSize     *prometheus.HistogramVec
	requestsInFlight prometheus.Gauge
	panicsTotal      *prometheus.CounterVec
}

// NewHTTPCollector creates a new HTTPCollector with the given configuration.
func NewHTTPCollector(cfg Config) (*HTTPCollector, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	c := &HTTPCollector{
		requestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: cfg.Namespace,
				Subsystem: cfg.Subsystem,
				Name:      "http_requests_total",
				Help:      "Total number of HTTP requests.",
			},
			[]string{"method", "path", "status"},
		),
		requestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: cfg.Namespace,
				Subsystem: cfg.Subsystem,
				Name:      "http_request_duration_seconds",
				Help:      "HTTP request latency in seconds.",
				Buckets:   HTTPLatencyBuckets,
			},
			[]string{"method", "path"},
		),
		requestSize: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: cfg.Namespace,
				Subsystem: cfg.Subsystem,
				Name:      "http_request_size_bytes",
				Help:      "HTTP request body size in bytes.",
				Buckets:   RequestSizeBuckets,
			},
			[]string{"method", "path"},
		),
		responseSize: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: cfg.Namespace,
				Subsystem: cfg.Subsystem,
				Name:      "http_response_size_bytes",
				Help:      "HTTP response body size in bytes.",
				Buckets:   RequestSizeBuckets,
			},
			[]string{"method", "path"},
		),
		requestsInFlight: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: cfg.Namespace,
				Subsystem: cfg.Subsystem,
				Name:      "http_requests_in_flight",
				Help:      "Current number of HTTP requests being processed.",
			},
		),
		panicsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: cfg.Namespace,
				Subsystem: cfg.Subsystem,
				Name:      "http_panics_total",
				Help:      "Total number of panics during HTTP request handling.",
			},
			[]string{"method", "path"},
		),
	}

	// Register all metrics with the registry.
	collectors := []prometheus.Collector{
		c.requestsTotal,
		c.requestDuration,
		c.requestSize,
		c.responseSize,
		c.requestsInFlight,
		c.panicsTotal,
	}

	for _, collector := range collectors {
		if err := cfg.Registry.Register(collector); err != nil {
			// If already registered, try to unregister and re-register.
			// This can happen in tests or when recreating collectors.
			if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
				// Use the existing collector.
				switch existing := are.ExistingCollector.(type) {
				case *prometheus.CounterVec:
					if collector == c.requestsTotal {
						c.requestsTotal = existing
					} else if collector == c.panicsTotal {
						c.panicsTotal = existing
					}
				case *prometheus.HistogramVec:
					if collector == c.requestDuration {
						c.requestDuration = existing
					} else if collector == c.requestSize {
						c.requestSize = existing
					} else if collector == c.responseSize {
						c.responseSize = existing
					}
				case prometheus.Gauge:
					if collector == c.requestsInFlight {
						c.requestsInFlight = existing
					}
				}
			} else {
				return nil, err
			}
		}
	}

	return c, nil
}

// RecordRequest implements server.MetricsCollector.
// It records metrics for a completed HTTP request.
func (c *HTTPCollector) RecordRequest(method, path string, statusCode int, duration time.Duration) {
	status := strconv.Itoa(statusCode)
	c.requestsTotal.WithLabelValues(method, path, status).Inc()
	c.requestDuration.WithLabelValues(method, path).Observe(duration.Seconds())
}

// RecordPanic implements server.MetricsCollector.
// It records when a panic occurs during request handling.
func (c *HTTPCollector) RecordPanic(method, path string) {
	c.panicsTotal.WithLabelValues(method, path).Inc()
}

// RecordRequestSize records the size of an HTTP request body.
func (c *HTTPCollector) RecordRequestSize(method, path string, size int64) {
	c.requestSize.WithLabelValues(method, path).Observe(float64(size))
}

// RecordResponseSize records the size of an HTTP response body.
func (c *HTTPCollector) RecordResponseSize(method, path string, size int64) {
	c.responseSize.WithLabelValues(method, path).Observe(float64(size))
}

// IncRequestsInFlight increments the in-flight requests gauge.
func (c *HTTPCollector) IncRequestsInFlight() {
	c.requestsInFlight.Inc()
}

// DecRequestsInFlight decrements the in-flight requests gauge.
func (c *HTTPCollector) DecRequestsInFlight() {
	c.requestsInFlight.Dec()
}

// Describe implements prometheus.Collector.
func (c *HTTPCollector) Describe(ch chan<- *prometheus.Desc) {
	c.requestsTotal.Describe(ch)
	c.requestDuration.Describe(ch)
	c.requestSize.Describe(ch)
	c.responseSize.Describe(ch)
	c.requestsInFlight.Describe(ch)
	c.panicsTotal.Describe(ch)
}

// Collect implements prometheus.Collector.
func (c *HTTPCollector) Collect(ch chan<- prometheus.Metric) {
	c.requestsTotal.Collect(ch)
	c.requestDuration.Collect(ch)
	c.requestSize.Collect(ch)
	c.responseSize.Collect(ch)
	c.requestsInFlight.Collect(ch)
	c.panicsTotal.Collect(ch)
}
