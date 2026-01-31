package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// RedisCollector collects Redis metrics.
type RedisCollector struct {
	commandsTotal   *prometheus.CounterVec
	commandDuration *prometheus.HistogramVec
	cacheHitsTotal  *prometheus.CounterVec
	cacheMissTotal  *prometheus.CounterVec
}

// NewRedisCollector creates a new RedisCollector with the given configuration.
func NewRedisCollector(cfg Config) (*RedisCollector, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	c := &RedisCollector{
		commandsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: cfg.Namespace,
				Subsystem: cfg.Subsystem,
				Name:      "redis_commands_total",
				Help:      "Total number of Redis commands executed.",
			},
			[]string{"command"},
		),
		commandDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: cfg.Namespace,
				Subsystem: cfg.Subsystem,
				Name:      "redis_command_duration_seconds",
				Help:      "Redis command latency in seconds.",
				Buckets:   DBLatencyBuckets,
			},
			[]string{"command"},
		),
		cacheHitsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: cfg.Namespace,
				Subsystem: cfg.Subsystem,
				Name:      "redis_cache_hits_total",
				Help:      "Total number of cache hits.",
			},
			[]string{"cache"},
		),
		cacheMissTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: cfg.Namespace,
				Subsystem: cfg.Subsystem,
				Name:      "redis_cache_misses_total",
				Help:      "Total number of cache misses.",
			},
			[]string{"cache"},
		),
	}

	// Register all metrics with the registry.
	collectors := []prometheus.Collector{
		c.commandsTotal,
		c.commandDuration,
		c.cacheHitsTotal,
		c.cacheMissTotal,
	}

	for _, collector := range collectors {
		if err := cfg.Registry.Register(collector); err != nil {
			if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
				switch existing := are.ExistingCollector.(type) {
				case *prometheus.CounterVec:
					if collector == c.commandsTotal {
						c.commandsTotal = existing
					} else if collector == c.cacheHitsTotal {
						c.cacheHitsTotal = existing
					} else if collector == c.cacheMissTotal {
						c.cacheMissTotal = existing
					}
				case *prometheus.HistogramVec:
					if collector == c.commandDuration {
						c.commandDuration = existing
					}
				}
			} else {
				return nil, err
			}
		}
	}

	return c, nil
}

// RecordCommand records a Redis command execution.
// command: Redis command name (e.g., "GET", "SET", "HGET").
func (c *RedisCollector) RecordCommand(command string, duration time.Duration) {
	c.commandsTotal.WithLabelValues(command).Inc()
	c.commandDuration.WithLabelValues(command).Observe(duration.Seconds())
}

// RecordCacheHit records a cache hit.
// cache: cache name or key pattern (e.g., "user_session", "ride_status").
func (c *RedisCollector) RecordCacheHit(cache string) {
	c.cacheHitsTotal.WithLabelValues(cache).Inc()
}

// RecordCacheMiss records a cache miss.
// cache: cache name or key pattern (e.g., "user_session", "ride_status").
func (c *RedisCollector) RecordCacheMiss(cache string) {
	c.cacheMissTotal.WithLabelValues(cache).Inc()
}

// Describe implements prometheus.Collector.
func (c *RedisCollector) Describe(ch chan<- *prometheus.Desc) {
	c.commandsTotal.Describe(ch)
	c.commandDuration.Describe(ch)
	c.cacheHitsTotal.Describe(ch)
	c.cacheMissTotal.Describe(ch)
}

// Collect implements prometheus.Collector.
func (c *RedisCollector) Collect(ch chan<- prometheus.Metric) {
	c.commandsTotal.Collect(ch)
	c.commandDuration.Collect(ch)
	c.cacheHitsTotal.Collect(ch)
	c.cacheMissTotal.Collect(ch)
}
