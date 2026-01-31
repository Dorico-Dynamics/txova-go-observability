package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// DBCollector collects database metrics.
type DBCollector struct {
	connectionsTotal    *prometheus.GaugeVec
	queryDuration       *prometheus.HistogramVec
	queryErrorsTotal    *prometheus.CounterVec
	transactionDuration *prometheus.HistogramVec
}

// NewDBCollector creates a new DBCollector with the given configuration.
func NewDBCollector(cfg Config) (*DBCollector, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	c := &DBCollector{
		connectionsTotal: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: cfg.Namespace,
				Subsystem: cfg.Subsystem,
				Name:      "db_connections_total",
				Help:      "Current number of database connections by pool and state.",
			},
			[]string{"pool", "state"},
		),
		queryDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: cfg.Namespace,
				Subsystem: cfg.Subsystem,
				Name:      "db_query_duration_seconds",
				Help:      "Database query latency in seconds.",
				Buckets:   DBLatencyBuckets,
			},
			[]string{"operation"},
		),
		queryErrorsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: cfg.Namespace,
				Subsystem: cfg.Subsystem,
				Name:      "db_query_errors_total",
				Help:      "Total number of database query errors.",
			},
			[]string{"operation", "error"},
		),
		transactionDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: cfg.Namespace,
				Subsystem: cfg.Subsystem,
				Name:      "db_transaction_duration_seconds",
				Help:      "Database transaction latency in seconds.",
				Buckets:   DBLatencyBuckets,
			},
			[]string{},
		),
	}

	// Register all metrics with the registry.
	collectors := []prometheus.Collector{
		c.connectionsTotal,
		c.queryDuration,
		c.queryErrorsTotal,
		c.transactionDuration,
	}

	for _, collector := range collectors {
		if err := cfg.Registry.Register(collector); err != nil {
			if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
				switch existing := are.ExistingCollector.(type) {
				case *prometheus.GaugeVec:
					if collector == c.connectionsTotal {
						c.connectionsTotal = existing
					}
				case *prometheus.HistogramVec:
					if collector == c.queryDuration {
						c.queryDuration = existing
					} else if collector == c.transactionDuration {
						c.transactionDuration = existing
					}
				case *prometheus.CounterVec:
					if collector == c.queryErrorsTotal {
						c.queryErrorsTotal = existing
					}
				}
			} else {
				return nil, err
			}
		}
	}

	return c, nil
}

// SetConnections sets the number of connections for a pool and state.
// pool: connection pool name (e.g., "primary", "replica")
// state: connection state (e.g., "idle", "in_use", "max_open").
func (c *DBCollector) SetConnections(pool, state string, count float64) {
	c.connectionsTotal.WithLabelValues(pool, state).Set(count)
}

// RecordQueryDuration records the duration of a database query.
// operation: query operation type (e.g., "select", "insert", "update", "delete").
func (c *DBCollector) RecordQueryDuration(operation string, duration time.Duration) {
	c.queryDuration.WithLabelValues(operation).Observe(duration.Seconds())
}

// RecordQueryError records a database query error.
// operation: query operation type (e.g., "select", "insert", "update", "delete")
// errorType: error classification (e.g., "timeout", "connection", "constraint").
func (c *DBCollector) RecordQueryError(operation, errorType string) {
	c.queryErrorsTotal.WithLabelValues(operation, errorType).Inc()
}

// RecordTransactionDuration records the duration of a database transaction.
func (c *DBCollector) RecordTransactionDuration(duration time.Duration) {
	c.transactionDuration.WithLabelValues().Observe(duration.Seconds())
}

// Describe implements prometheus.Collector.
func (c *DBCollector) Describe(ch chan<- *prometheus.Desc) {
	c.connectionsTotal.Describe(ch)
	c.queryDuration.Describe(ch)
	c.queryErrorsTotal.Describe(ch)
	c.transactionDuration.Describe(ch)
}

// Collect implements prometheus.Collector.
func (c *DBCollector) Collect(ch chan<- prometheus.Metric) {
	c.connectionsTotal.Collect(ch)
	c.queryDuration.Collect(ch)
	c.queryErrorsTotal.Collect(ch)
	c.transactionDuration.Collect(ch)
}
