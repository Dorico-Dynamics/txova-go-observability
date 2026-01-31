package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// SafetyCollector collects safety-related business metrics.
type SafetyCollector struct {
	emergenciesTotal *prometheus.CounterVec
	incidentsTotal   *prometheus.CounterVec
	tripSharesTotal  prometheus.Counter
}

// NewSafetyCollector creates a new SafetyCollector with the given configuration.
func NewSafetyCollector(cfg Config) (*SafetyCollector, error) {
	cfg, err := cfg.Validate()
	if err != nil {
		return nil, err
	}

	c := &SafetyCollector{}

	c.emergenciesTotal, err = registerCollector(cfg.Registry, prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: cfg.Namespace,
			Subsystem: cfg.Subsystem,
			Name:      "emergencies_triggered_total",
			Help:      "Total number of emergency (SOS) activations.",
		},
		[]string{"type", "city"},
	))
	if err != nil {
		return nil, err
	}

	c.incidentsTotal, err = registerCollector(cfg.Registry, prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: cfg.Namespace,
			Subsystem: cfg.Subsystem,
			Name:      "incidents_reported_total",
			Help:      "Total number of incidents reported.",
		},
		[]string{"severity"},
	))
	if err != nil {
		return nil, err
	}

	c.tripSharesTotal, err = registerCollector(cfg.Registry, prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: cfg.Namespace,
			Subsystem: cfg.Subsystem,
			Name:      "trip_shares_total",
			Help:      "Total number of trip sharing activations.",
		},
	))
	if err != nil {
		return nil, err
	}

	return c, nil
}

// RecordEmergency records an emergency (SOS) activation.
// emergencyType: type of emergency (e.g., "sos_button", "auto_detected", "police_request")
// city: city where the emergency was triggered.
func (c *SafetyCollector) RecordEmergency(emergencyType, city string) {
	c.emergenciesTotal.WithLabelValues(emergencyType, city).Inc()
}

// RecordIncident records an incident report.
// severity: incident severity (e.g., "low", "medium", "high", "critical").
func (c *SafetyCollector) RecordIncident(severity string) {
	c.incidentsTotal.WithLabelValues(severity).Inc()
}

// RecordTripShare records a trip sharing activation.
func (c *SafetyCollector) RecordTripShare() {
	c.tripSharesTotal.Inc()
}

// Describe implements prometheus.Collector.
func (c *SafetyCollector) Describe(ch chan<- *prometheus.Desc) {
	c.emergenciesTotal.Describe(ch)
	c.incidentsTotal.Describe(ch)
	c.tripSharesTotal.Describe(ch)
}

// Collect implements prometheus.Collector.
func (c *SafetyCollector) Collect(ch chan<- prometheus.Metric) {
	c.emergenciesTotal.Collect(ch)
	c.incidentsTotal.Collect(ch)
	c.tripSharesTotal.Collect(ch)
}
