package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestNewKafkaCollector(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_kafka")

	collector, err := NewKafkaCollector(cfg)
	if err != nil {
		t.Fatalf("NewKafkaCollector() error = %v", err)
	}
	if collector == nil {
		t.Fatal("NewKafkaCollector() returned nil collector")
	}
}

func TestNewKafkaCollector_DuplicateRegistration(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_kafka_dup")

	collector1, err := NewKafkaCollector(cfg)
	if err != nil {
		t.Fatalf("First NewKafkaCollector() error = %v", err)
	}

	collector2, err := NewKafkaCollector(cfg)
	if err != nil {
		t.Fatalf("Second NewKafkaCollector() error = %v", err)
	}

	if collector1 == nil || collector2 == nil {
		t.Fatal("NewKafkaCollector() returned nil collectors")
	}
}

func TestKafkaCollector_RecordMessageProduced(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_kafka_prod")

	collector, err := NewKafkaCollector(cfg)
	if err != nil {
		t.Fatalf("NewKafkaCollector() error = %v", err)
	}

	collector.RecordMessageProduced("ride_events")
	collector.RecordMessageProduced("ride_events")
	collector.RecordMessageProduced("payment_events")

	count := testutil.ToFloat64(collector.messagesProducedTotal.WithLabelValues("ride_events"))
	if count != 2 {
		t.Errorf("messagesProducedTotal ride_events = %v, want 2", count)
	}

	count = testutil.ToFloat64(collector.messagesProducedTotal.WithLabelValues("payment_events"))
	if count != 1 {
		t.Errorf("messagesProducedTotal payment_events = %v, want 1", count)
	}
}

func TestKafkaCollector_RecordMessageConsumed(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_kafka_cons")

	collector, err := NewKafkaCollector(cfg)
	if err != nil {
		t.Fatalf("NewKafkaCollector() error = %v", err)
	}

	collector.RecordMessageConsumed("ride_events", "ride_processor")
	collector.RecordMessageConsumed("ride_events", "ride_processor")
	collector.RecordMessageConsumed("ride_events", "analytics_processor")

	count := testutil.ToFloat64(collector.messagesConsumedTotal.WithLabelValues("ride_events", "ride_processor"))
	if count != 2 {
		t.Errorf("messagesConsumedTotal ride_events/ride_processor = %v, want 2", count)
	}

	count = testutil.ToFloat64(collector.messagesConsumedTotal.WithLabelValues("ride_events", "analytics_processor"))
	if count != 1 {
		t.Errorf("messagesConsumedTotal ride_events/analytics_processor = %v, want 1", count)
	}
}

func TestKafkaCollector_SetConsumerLag(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_kafka_lag")

	collector, err := NewKafkaCollector(cfg)
	if err != nil {
		t.Fatalf("NewKafkaCollector() error = %v", err)
	}

	collector.SetConsumerLag("ride_events", "0", "ride_processor", 100)
	collector.SetConsumerLag("ride_events", "1", "ride_processor", 50)
	collector.SetConsumerLag("payment_events", "0", "payment_processor", 25)

	lag := testutil.ToFloat64(collector.consumerLag.WithLabelValues("ride_events", "0", "ride_processor"))
	if lag != 100 {
		t.Errorf("consumerLag ride_events/0/ride_processor = %v, want 100", lag)
	}

	lag = testutil.ToFloat64(collector.consumerLag.WithLabelValues("ride_events", "1", "ride_processor"))
	if lag != 50 {
		t.Errorf("consumerLag ride_events/1/ride_processor = %v, want 50", lag)
	}

	// Test updating lag value
	collector.SetConsumerLag("ride_events", "0", "ride_processor", 0)
	lag = testutil.ToFloat64(collector.consumerLag.WithLabelValues("ride_events", "0", "ride_processor"))
	if lag != 0 {
		t.Errorf("consumerLag after update = %v, want 0", lag)
	}
}

func TestKafkaCollector_RecordProduceError(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_kafka_prod_err")

	collector, err := NewKafkaCollector(cfg)
	if err != nil {
		t.Fatalf("NewKafkaCollector() error = %v", err)
	}

	collector.RecordProduceError("ride_events")
	collector.RecordProduceError("ride_events")
	collector.RecordProduceError("payment_events")

	count := testutil.ToFloat64(collector.produceErrorsTotal.WithLabelValues("ride_events"))
	if count != 2 {
		t.Errorf("produceErrorsTotal ride_events = %v, want 2", count)
	}

	count = testutil.ToFloat64(collector.produceErrorsTotal.WithLabelValues("payment_events"))
	if count != 1 {
		t.Errorf("produceErrorsTotal payment_events = %v, want 1", count)
	}
}

func TestKafkaCollector_RecordConsumeError(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_kafka_cons_err")

	collector, err := NewKafkaCollector(cfg)
	if err != nil {
		t.Fatalf("NewKafkaCollector() error = %v", err)
	}

	collector.RecordConsumeError("ride_events")
	collector.RecordConsumeError("notification_events")
	collector.RecordConsumeError("notification_events")

	count := testutil.ToFloat64(collector.consumeErrorsTotal.WithLabelValues("ride_events"))
	if count != 1 {
		t.Errorf("consumeErrorsTotal ride_events = %v, want 1", count)
	}

	count = testutil.ToFloat64(collector.consumeErrorsTotal.WithLabelValues("notification_events"))
	if count != 2 {
		t.Errorf("consumeErrorsTotal notification_events = %v, want 2", count)
	}
}

func TestKafkaCollector_Describe(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_kafka_desc")

	collector, err := NewKafkaCollector(cfg)
	if err != nil {
		t.Fatalf("NewKafkaCollector() error = %v", err)
	}

	ch := make(chan *prometheus.Desc, 100)
	collector.Describe(ch)
	close(ch)

	count := 0
	for range ch {
		count++
	}

	if count == 0 {
		t.Error("Describe() produced no descriptors")
	}
}

func TestKafkaCollector_Collect(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_kafka_collect")

	collector, err := NewKafkaCollector(cfg)
	if err != nil {
		t.Fatalf("NewKafkaCollector() error = %v", err)
	}

	// Record some metrics first
	collector.RecordMessageProduced("test_topic")
	collector.SetConsumerLag("test_topic", "0", "test_group", 10)

	ch := make(chan prometheus.Metric, 100)
	collector.Collect(ch)
	close(ch)

	count := 0
	for range ch {
		count++
	}

	if count == 0 {
		t.Error("Collect() produced no metrics")
	}
}

func TestKafkaCollector_MultipleTopics(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_kafka_multi")

	collector, err := NewKafkaCollector(cfg)
	if err != nil {
		t.Fatalf("NewKafkaCollector() error = %v", err)
	}

	topics := []string{"ride_events", "payment_events", "notification_events", "driver_events", "user_events"}
	for _, topic := range topics {
		collector.RecordMessageProduced(topic)
		collector.RecordMessageConsumed(topic, "processor")
	}

	for _, topic := range topics {
		prodCount := testutil.ToFloat64(collector.messagesProducedTotal.WithLabelValues(topic))
		if prodCount != 1 {
			t.Errorf("messagesProducedTotal %s = %v, want 1", topic, prodCount)
		}

		consCount := testutil.ToFloat64(collector.messagesConsumedTotal.WithLabelValues(topic, "processor"))
		if consCount != 1 {
			t.Errorf("messagesConsumedTotal %s/processor = %v, want 1", topic, consCount)
		}
	}
}

func TestKafkaCollector_HighThroughput(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_kafka_throughput")

	collector, err := NewKafkaCollector(cfg)
	if err != nil {
		t.Fatalf("NewKafkaCollector() error = %v", err)
	}

	// Simulate high throughput
	for i := 0; i < 10000; i++ {
		collector.RecordMessageProduced("high_volume_topic")
		collector.RecordMessageConsumed("high_volume_topic", "fast_consumer")
	}

	prodCount := testutil.ToFloat64(collector.messagesProducedTotal.WithLabelValues("high_volume_topic"))
	if prodCount != 10000 {
		t.Errorf("messagesProducedTotal = %v, want 10000", prodCount)
	}

	consCount := testutil.ToFloat64(collector.messagesConsumedTotal.WithLabelValues("high_volume_topic", "fast_consumer"))
	if consCount != 10000 {
		t.Errorf("messagesConsumedTotal = %v, want 10000", consCount)
	}
}
