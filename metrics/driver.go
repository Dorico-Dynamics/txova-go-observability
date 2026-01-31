package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// DriverCollector collects driver-related business metrics.
type DriverCollector struct {
	onlineTotal    *prometheus.GaugeVec
	acceptanceRate *prometheus.GaugeVec
	ratingAverage  prometheus.Gauge
	earnings       *prometheus.CounterVec
}

// NewDriverCollector creates a new DriverCollector with the given configuration.
func NewDriverCollector(cfg Config) (*DriverCollector, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	c := &DriverCollector{
		onlineTotal: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: cfg.Namespace,
				Subsystem: cfg.Subsystem,
				Name:      "drivers_online_total",
				Help:      "Current number of online drivers.",
			},
			[]string{"city", "service_type"},
		),
		acceptanceRate: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: cfg.Namespace,
				Subsystem: cfg.Subsystem,
				Name:      "driver_acceptance_rate",
				Help:      "Driver acceptance rate (0.0-1.0).",
			},
			[]string{"driver_id"},
		),
		ratingAverage: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: cfg.Namespace,
				Subsystem: cfg.Subsystem,
				Name:      "driver_rating_average",
				Help:      "Average driver rating across all drivers.",
			},
		),
		earnings: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: cfg.Namespace,
				Subsystem: cfg.Subsystem,
				Name:      "driver_earnings_mzn",
				Help:      "Total driver earnings in MZN (smallest currency unit).",
			},
			[]string{"driver_id"},
		),
	}

	// Register all metrics with the registry.
	collectors := []prometheus.Collector{
		c.onlineTotal,
		c.acceptanceRate,
		c.ratingAverage,
		c.earnings,
	}

	for _, collector := range collectors {
		if err := cfg.Registry.Register(collector); err != nil {
			if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
				switch existing := are.ExistingCollector.(type) {
				case *prometheus.GaugeVec:
					if collector == c.onlineTotal {
						c.onlineTotal = existing
					} else if collector == c.acceptanceRate {
						c.acceptanceRate = existing
					}
				case prometheus.Gauge:
					if collector == c.ratingAverage {
						c.ratingAverage = existing
					}
				case *prometheus.CounterVec:
					if collector == c.earnings {
						c.earnings = existing
					}
				}
			} else {
				return nil, err
			}
		}
	}

	return c, nil
}

// SetDriversOnline sets the number of online drivers.
// city: city name
// serviceType: type of service (e.g., "standard", "premium", "moto")
// count: number of online drivers.
func (c *DriverCollector) SetDriversOnline(city, serviceType string, count float64) {
	c.onlineTotal.WithLabelValues(city, serviceType).Set(count)
}

// SetAcceptanceRate sets the acceptance rate for a driver.
// driverID: driver identifier
// rate: acceptance rate (0.0-1.0).
func (c *DriverCollector) SetAcceptanceRate(driverID string, rate float64) {
	c.acceptanceRate.WithLabelValues(driverID).Set(rate)
}

// SetRatingAverage sets the average driver rating.
// rating: average rating value.
func (c *DriverCollector) SetRatingAverage(rating float64) {
	c.ratingAverage.Set(rating)
}

// AddEarnings adds to a driver's total earnings.
// driverID: driver identifier
// amountMZN: amount to add in MZN (smallest currency unit).
func (c *DriverCollector) AddEarnings(driverID string, amountMZN float64) {
	c.earnings.WithLabelValues(driverID).Add(amountMZN)
}

// Describe implements prometheus.Collector.
func (c *DriverCollector) Describe(ch chan<- *prometheus.Desc) {
	c.onlineTotal.Describe(ch)
	c.acceptanceRate.Describe(ch)
	c.ratingAverage.Describe(ch)
	c.earnings.Describe(ch)
}

// Collect implements prometheus.Collector.
func (c *DriverCollector) Collect(ch chan<- prometheus.Metric) {
	c.onlineTotal.Collect(ch)
	c.acceptanceRate.Collect(ch)
	c.ratingAverage.Collect(ch)
	c.earnings.Collect(ch)
}
