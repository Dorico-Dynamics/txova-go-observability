package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// RideCollector collects ride-related business metrics.
type RideCollector struct {
	requestedTotal *prometheus.CounterVec
	completedTotal *prometheus.CounterVec
	cancelledTotal *prometheus.CounterVec
	duration       *prometheus.HistogramVec
	distance       *prometheus.HistogramVec
	fare           *prometheus.HistogramVec
	waitTime       *prometheus.HistogramVec
}

// NewRideCollector creates a new RideCollector with the given configuration.
func NewRideCollector(cfg Config) (*RideCollector, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	c := &RideCollector{
		requestedTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: cfg.Namespace,
				Subsystem: cfg.Subsystem,
				Name:      "rides_requested_total",
				Help:      "Total number of ride requests.",
			},
			[]string{"service_type", "city"},
		),
		completedTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: cfg.Namespace,
				Subsystem: cfg.Subsystem,
				Name:      "rides_completed_total",
				Help:      "Total number of completed rides.",
			},
			[]string{"service_type", "city"},
		),
		cancelledTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: cfg.Namespace,
				Subsystem: cfg.Subsystem,
				Name:      "rides_cancelled_total",
				Help:      "Total number of cancelled rides.",
			},
			[]string{"cancelled_by", "reason"},
		),
		duration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: cfg.Namespace,
				Subsystem: cfg.Subsystem,
				Name:      "ride_duration_seconds",
				Help:      "Duration of rides in seconds.",
				Buckets:   DurationBuckets,
			},
			[]string{"service_type"},
		),
		distance: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: cfg.Namespace,
				Subsystem: cfg.Subsystem,
				Name:      "ride_distance_km",
				Help:      "Distance of rides in kilometers.",
				Buckets:   DistanceBuckets,
			},
			[]string{"service_type"},
		),
		fare: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: cfg.Namespace,
				Subsystem: cfg.Subsystem,
				Name:      "ride_fare_mzn",
				Help:      "Fare of rides in MZN (smallest currency unit).",
				Buckets:   FareBuckets,
			},
			[]string{"service_type"},
		),
		waitTime: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: cfg.Namespace,
				Subsystem: cfg.Subsystem,
				Name:      "ride_wait_time_seconds",
				Help:      "Time to match a driver in seconds.",
				Buckets:   DurationBuckets,
			},
			[]string{"service_type"},
		),
	}

	// Register all metrics with the registry.
	collectors := []prometheus.Collector{
		c.requestedTotal,
		c.completedTotal,
		c.cancelledTotal,
		c.duration,
		c.distance,
		c.fare,
		c.waitTime,
	}

	for _, collector := range collectors {
		if err := cfg.Registry.Register(collector); err != nil {
			if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
				switch existing := are.ExistingCollector.(type) {
				case *prometheus.CounterVec:
					if collector == c.requestedTotal {
						c.requestedTotal = existing
					} else if collector == c.completedTotal {
						c.completedTotal = existing
					} else if collector == c.cancelledTotal {
						c.cancelledTotal = existing
					}
				case *prometheus.HistogramVec:
					if collector == c.duration {
						c.duration = existing
					} else if collector == c.distance {
						c.distance = existing
					} else if collector == c.fare {
						c.fare = existing
					} else if collector == c.waitTime {
						c.waitTime = existing
					}
				}
			} else {
				return nil, err
			}
		}
	}

	return c, nil
}

// RecordRideRequested records a ride request.
// serviceType: type of service (e.g., "standard", "premium", "moto")
// city: city where the ride was requested.
func (c *RideCollector) RecordRideRequested(serviceType, city string) {
	c.requestedTotal.WithLabelValues(serviceType, city).Inc()
}

// RecordRideCompleted records a completed ride.
// serviceType: type of service (e.g., "standard", "premium", "moto")
// city: city where the ride was completed.
func (c *RideCollector) RecordRideCompleted(serviceType, city string) {
	c.completedTotal.WithLabelValues(serviceType, city).Inc()
}

// RecordRideCancelled records a cancelled ride.
// cancelledBy: who cancelled the ride (e.g., "rider", "driver", "system")
// reason: reason for cancellation (e.g., "no_drivers", "rider_cancelled", "timeout").
func (c *RideCollector) RecordRideCancelled(cancelledBy, reason string) {
	c.cancelledTotal.WithLabelValues(cancelledBy, reason).Inc()
}

// RecordRideDuration records the duration of a ride.
// serviceType: type of service (e.g., "standard", "premium", "moto")
// duration: duration of the ride.
func (c *RideCollector) RecordRideDuration(serviceType string, duration time.Duration) {
	c.duration.WithLabelValues(serviceType).Observe(duration.Seconds())
}

// RecordRideDistance records the distance of a ride.
// serviceType: type of service (e.g., "standard", "premium", "moto")
// distanceKm: distance of the ride in kilometers.
func (c *RideCollector) RecordRideDistance(serviceType string, distanceKm float64) {
	c.distance.WithLabelValues(serviceType).Observe(distanceKm)
}

// RecordRideFare records the fare of a ride.
// serviceType: type of service (e.g., "standard", "premium", "moto")
// fareMZN: fare amount in MZN (smallest currency unit).
func (c *RideCollector) RecordRideFare(serviceType string, fareMZN float64) {
	c.fare.WithLabelValues(serviceType).Observe(fareMZN)
}

// RecordRideWaitTime records the time to match a driver.
// serviceType: type of service (e.g., "standard", "premium", "moto")
// waitTime: time taken to match a driver.
func (c *RideCollector) RecordRideWaitTime(serviceType string, waitTime time.Duration) {
	c.waitTime.WithLabelValues(serviceType).Observe(waitTime.Seconds())
}

// Describe implements prometheus.Collector.
func (c *RideCollector) Describe(ch chan<- *prometheus.Desc) {
	c.requestedTotal.Describe(ch)
	c.completedTotal.Describe(ch)
	c.cancelledTotal.Describe(ch)
	c.duration.Describe(ch)
	c.distance.Describe(ch)
	c.fare.Describe(ch)
	c.waitTime.Describe(ch)
}

// Collect implements prometheus.Collector.
func (c *RideCollector) Collect(ch chan<- prometheus.Metric) {
	c.requestedTotal.Collect(ch)
	c.completedTotal.Collect(ch)
	c.cancelledTotal.Collect(ch)
	c.duration.Collect(ch)
	c.distance.Collect(ch)
	c.fare.Collect(ch)
	c.waitTime.Collect(ch)
}
