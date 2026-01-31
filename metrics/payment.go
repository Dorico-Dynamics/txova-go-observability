package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// PaymentCollector collects payment-related business metrics.
type PaymentCollector struct {
	paymentsTotal  *prometheus.CounterVec
	paymentAmount  *prometheus.HistogramVec
	processingTime *prometheus.HistogramVec
	refundsTotal   *prometheus.CounterVec
}

// NewPaymentCollector creates a new PaymentCollector with the given configuration.
func NewPaymentCollector(cfg Config) (*PaymentCollector, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	c := &PaymentCollector{
		paymentsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: cfg.Namespace,
				Subsystem: cfg.Subsystem,
				Name:      "payments_total",
				Help:      "Total number of payment attempts.",
			},
			[]string{"method", "status"},
		),
		paymentAmount: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: cfg.Namespace,
				Subsystem: cfg.Subsystem,
				Name:      "payment_amount_mzn",
				Help:      "Payment amounts in MZN (smallest currency unit).",
				Buckets:   PaymentAmountBuckets,
			},
			[]string{"method"},
		),
		processingTime: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: cfg.Namespace,
				Subsystem: cfg.Subsystem,
				Name:      "payment_processing_seconds",
				Help:      "Payment processing time in seconds.",
				Buckets:   HTTPLatencyBuckets,
			},
			[]string{"method"},
		),
		refundsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: cfg.Namespace,
				Subsystem: cfg.Subsystem,
				Name:      "refunds_total",
				Help:      "Total number of refunds issued.",
			},
			[]string{"reason"},
		),
	}

	// Register all metrics with the registry.
	collectors := []prometheus.Collector{
		c.paymentsTotal,
		c.paymentAmount,
		c.processingTime,
		c.refundsTotal,
	}

	for _, collector := range collectors {
		if err := cfg.Registry.Register(collector); err != nil {
			if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
				switch existing := are.ExistingCollector.(type) {
				case *prometheus.CounterVec:
					if collector == c.paymentsTotal {
						c.paymentsTotal = existing
					} else if collector == c.refundsTotal {
						c.refundsTotal = existing
					}
				case *prometheus.HistogramVec:
					if collector == c.paymentAmount {
						c.paymentAmount = existing
					} else if collector == c.processingTime {
						c.processingTime = existing
					}
				}
			} else {
				return nil, err
			}
		}
	}

	return c, nil
}

// RecordPayment records a payment attempt.
// method: payment method (e.g., "mpesa", "card", "cash")
// status: payment status (e.g., "success", "failed", "pending").
func (c *PaymentCollector) RecordPayment(method, status string) {
	c.paymentsTotal.WithLabelValues(method, status).Inc()
}

// RecordPaymentAmount records the amount of a payment.
// method: payment method (e.g., "mpesa", "card", "cash")
// amountMZN: payment amount in MZN (smallest currency unit).
func (c *PaymentCollector) RecordPaymentAmount(method string, amountMZN float64) {
	c.paymentAmount.WithLabelValues(method).Observe(amountMZN)
}

// RecordProcessingTime records the time taken to process a payment.
// method: payment method (e.g., "mpesa", "card", "cash")
// duration: processing time.
func (c *PaymentCollector) RecordProcessingTime(method string, duration time.Duration) {
	c.processingTime.WithLabelValues(method).Observe(duration.Seconds())
}

// RecordRefund records a refund.
// reason: reason for refund (e.g., "ride_cancelled", "overcharge", "dispute").
func (c *PaymentCollector) RecordRefund(reason string) {
	c.refundsTotal.WithLabelValues(reason).Inc()
}

// Describe implements prometheus.Collector.
func (c *PaymentCollector) Describe(ch chan<- *prometheus.Desc) {
	c.paymentsTotal.Describe(ch)
	c.paymentAmount.Describe(ch)
	c.processingTime.Describe(ch)
	c.refundsTotal.Describe(ch)
}

// Collect implements prometheus.Collector.
func (c *PaymentCollector) Collect(ch chan<- prometheus.Metric) {
	c.paymentsTotal.Collect(ch)
	c.paymentAmount.Collect(ch)
	c.processingTime.Collect(ch)
	c.refundsTotal.Collect(ch)
}
