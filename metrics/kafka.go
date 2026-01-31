package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// KafkaCollector collects Kafka metrics.
type KafkaCollector struct {
	messagesProducedTotal *prometheus.CounterVec
	messagesConsumedTotal *prometheus.CounterVec
	consumerLag           *prometheus.GaugeVec
	produceErrorsTotal    *prometheus.CounterVec
	consumeErrorsTotal    *prometheus.CounterVec
}

// NewKafkaCollector creates a new KafkaCollector with the given configuration.
func NewKafkaCollector(cfg Config) (*KafkaCollector, error) {
	cfg, err := cfg.Validate()
	if err != nil {
		return nil, err
	}

	c := &KafkaCollector{}

	c.messagesProducedTotal, err = registerCollector(cfg.Registry, prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: cfg.Namespace,
			Subsystem: cfg.Subsystem,
			Name:      "kafka_messages_produced_total",
			Help:      "Total number of Kafka messages produced.",
		},
		[]string{"topic"},
	))
	if err != nil {
		return nil, err
	}

	c.messagesConsumedTotal, err = registerCollector(cfg.Registry, prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: cfg.Namespace,
			Subsystem: cfg.Subsystem,
			Name:      "kafka_messages_consumed_total",
			Help:      "Total number of Kafka messages consumed.",
		},
		[]string{"topic", "group"},
	))
	if err != nil {
		return nil, err
	}

	c.consumerLag, err = registerCollector(cfg.Registry, prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: cfg.Namespace,
			Subsystem: cfg.Subsystem,
			Name:      "kafka_consumer_lag",
			Help:      "Current consumer lag by topic, partition, and consumer group.",
		},
		[]string{"topic", "partition", "group"},
	))
	if err != nil {
		return nil, err
	}

	c.produceErrorsTotal, err = registerCollector(cfg.Registry, prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: cfg.Namespace,
			Subsystem: cfg.Subsystem,
			Name:      "kafka_produce_errors_total",
			Help:      "Total number of Kafka produce errors.",
		},
		[]string{"topic"},
	))
	if err != nil {
		return nil, err
	}

	c.consumeErrorsTotal, err = registerCollector(cfg.Registry, prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: cfg.Namespace,
			Subsystem: cfg.Subsystem,
			Name:      "kafka_consume_errors_total",
			Help:      "Total number of Kafka consume errors.",
		},
		[]string{"topic"},
	))
	if err != nil {
		return nil, err
	}

	return c, nil
}

// RecordMessageProduced records a successfully produced message.
// topic: Kafka topic name.
func (c *KafkaCollector) RecordMessageProduced(topic string) {
	c.messagesProducedTotal.WithLabelValues(topic).Inc()
}

// RecordMessageConsumed records a successfully consumed message.
// topic: Kafka topic name
// group: consumer group name.
func (c *KafkaCollector) RecordMessageConsumed(topic, group string) {
	c.messagesConsumedTotal.WithLabelValues(topic, group).Inc()
}

// SetConsumerLag sets the current consumer lag for a topic/partition/group.
// topic: Kafka topic name
// partition: partition number (as string)
// group: consumer group name
// lag: current lag (number of messages behind).
func (c *KafkaCollector) SetConsumerLag(topic, partition, group string, lag float64) {
	c.consumerLag.WithLabelValues(topic, partition, group).Set(lag)
}

// RecordProduceError records a produce error.
// topic: Kafka topic name.
func (c *KafkaCollector) RecordProduceError(topic string) {
	c.produceErrorsTotal.WithLabelValues(topic).Inc()
}

// RecordConsumeError records a consume error.
// topic: Kafka topic name.
func (c *KafkaCollector) RecordConsumeError(topic string) {
	c.consumeErrorsTotal.WithLabelValues(topic).Inc()
}

// Describe implements prometheus.Collector.
func (c *KafkaCollector) Describe(ch chan<- *prometheus.Desc) {
	c.messagesProducedTotal.Describe(ch)
	c.messagesConsumedTotal.Describe(ch)
	c.consumerLag.Describe(ch)
	c.produceErrorsTotal.Describe(ch)
	c.consumeErrorsTotal.Describe(ch)
}

// Collect implements prometheus.Collector.
func (c *KafkaCollector) Collect(ch chan<- prometheus.Metric) {
	c.messagesProducedTotal.Collect(ch)
	c.messagesConsumedTotal.Collect(ch)
	c.consumerLag.Collect(ch)
	c.produceErrorsTotal.Collect(ch)
	c.consumeErrorsTotal.Collect(ch)
}
