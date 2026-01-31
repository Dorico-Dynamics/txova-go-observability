package metrics

import (
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestNewPaymentCollector(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_payment")

	collector, err := NewPaymentCollector(cfg)
	if err != nil {
		t.Fatalf("NewPaymentCollector() error = %v", err)
	}
	if collector == nil {
		t.Fatal("NewPaymentCollector() returned nil collector")
	}
}

func TestNewPaymentCollector_DuplicateRegistration(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_payment_dup")

	collector1, err := NewPaymentCollector(cfg)
	if err != nil {
		t.Fatalf("First NewPaymentCollector() error = %v", err)
	}

	collector2, err := NewPaymentCollector(cfg)
	if err != nil {
		t.Fatalf("Second NewPaymentCollector() error = %v", err)
	}

	if collector1 == nil || collector2 == nil {
		t.Fatal("NewPaymentCollector() returned nil collectors")
	}
}

func TestPaymentCollector_RecordPayment(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_payment_rec")

	collector, err := NewPaymentCollector(cfg)
	if err != nil {
		t.Fatalf("NewPaymentCollector() error = %v", err)
	}

	collector.RecordPayment("mpesa", "success")
	collector.RecordPayment("mpesa", "success")
	collector.RecordPayment("mpesa", "failed")
	collector.RecordPayment("card", "success")
	collector.RecordPayment("cash", "success")

	count := testutil.ToFloat64(collector.paymentsTotal.WithLabelValues("mpesa", "success"))
	if count != 2 {
		t.Errorf("paymentsTotal mpesa/success = %v, want 2", count)
	}

	count = testutil.ToFloat64(collector.paymentsTotal.WithLabelValues("mpesa", "failed"))
	if count != 1 {
		t.Errorf("paymentsTotal mpesa/failed = %v, want 1", count)
	}

	count = testutil.ToFloat64(collector.paymentsTotal.WithLabelValues("card", "success"))
	if count != 1 {
		t.Errorf("paymentsTotal card/success = %v, want 1", count)
	}

	count = testutil.ToFloat64(collector.paymentsTotal.WithLabelValues("cash", "success"))
	if count != 1 {
		t.Errorf("paymentsTotal cash/success = %v, want 1", count)
	}
}

func TestPaymentCollector_RecordPaymentAmount(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_payment_amt")

	collector, err := NewPaymentCollector(cfg)
	if err != nil {
		t.Fatalf("NewPaymentCollector() error = %v", err)
	}

	collector.RecordPaymentAmount("mpesa", 150)
	collector.RecordPaymentAmount("mpesa", 500)
	collector.RecordPaymentAmount("card", 2500)
	collector.RecordPaymentAmount("cash", 100)

	histCount := testutil.CollectAndCount(collector.paymentAmount)
	if histCount == 0 {
		t.Error("paymentAmount histogram has no metrics")
	}
}

func TestPaymentCollector_RecordProcessingTime(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_payment_proc")

	collector, err := NewPaymentCollector(cfg)
	if err != nil {
		t.Fatalf("NewPaymentCollector() error = %v", err)
	}

	collector.RecordProcessingTime("mpesa", 2*time.Second)
	collector.RecordProcessingTime("mpesa", 3*time.Second)
	collector.RecordProcessingTime("card", 500*time.Millisecond)
	collector.RecordProcessingTime("cash", 100*time.Millisecond)

	histCount := testutil.CollectAndCount(collector.processingTime)
	if histCount == 0 {
		t.Error("processingTime histogram has no metrics")
	}
}

func TestPaymentCollector_RecordRefund(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_payment_refund")

	collector, err := NewPaymentCollector(cfg)
	if err != nil {
		t.Fatalf("NewPaymentCollector() error = %v", err)
	}

	collector.RecordRefund("ride_cancelled")
	collector.RecordRefund("ride_cancelled")
	collector.RecordRefund("overcharge")
	collector.RecordRefund("dispute")

	count := testutil.ToFloat64(collector.refundsTotal.WithLabelValues("ride_cancelled"))
	if count != 2 {
		t.Errorf("refundsTotal ride_cancelled = %v, want 2", count)
	}

	count = testutil.ToFloat64(collector.refundsTotal.WithLabelValues("overcharge"))
	if count != 1 {
		t.Errorf("refundsTotal overcharge = %v, want 1", count)
	}

	count = testutil.ToFloat64(collector.refundsTotal.WithLabelValues("dispute"))
	if count != 1 {
		t.Errorf("refundsTotal dispute = %v, want 1", count)
	}
}

func TestPaymentCollector_Describe(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_payment_desc")

	collector, err := NewPaymentCollector(cfg)
	if err != nil {
		t.Fatalf("NewPaymentCollector() error = %v", err)
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

func TestPaymentCollector_Collect(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_payment_collect")

	collector, err := NewPaymentCollector(cfg)
	if err != nil {
		t.Fatalf("NewPaymentCollector() error = %v", err)
	}

	// Record some metrics first
	collector.RecordPayment("mpesa", "success")
	collector.RecordPaymentAmount("mpesa", 500)

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

func TestPaymentCollector_AllPaymentMethods(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_payment_methods")

	collector, err := NewPaymentCollector(cfg)
	if err != nil {
		t.Fatalf("NewPaymentCollector() error = %v", err)
	}

	methods := []string{"mpesa", "emola", "card", "cash", "wallet"}
	for _, method := range methods {
		collector.RecordPayment(method, "success")
		collector.RecordPaymentAmount(method, 500)
		collector.RecordProcessingTime(method, time.Second)
	}

	for _, method := range methods {
		count := testutil.ToFloat64(collector.paymentsTotal.WithLabelValues(method, "success"))
		if count != 1 {
			t.Errorf("paymentsTotal %s/success = %v, want 1", method, count)
		}
	}
}

func TestPaymentCollector_AllPaymentStatuses(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_payment_statuses")

	collector, err := NewPaymentCollector(cfg)
	if err != nil {
		t.Fatalf("NewPaymentCollector() error = %v", err)
	}

	statuses := []string{"success", "failed", "pending", "cancelled", "refunded"}
	for _, status := range statuses {
		collector.RecordPayment("mpesa", status)
	}

	for _, status := range statuses {
		count := testutil.ToFloat64(collector.paymentsTotal.WithLabelValues("mpesa", status))
		if count != 1 {
			t.Errorf("paymentsTotal mpesa/%s = %v, want 1", status, count)
		}
	}
}

func TestPaymentCollector_HighThroughput(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_payment_throughput")

	collector, err := NewPaymentCollector(cfg)
	if err != nil {
		t.Fatalf("NewPaymentCollector() error = %v", err)
	}

	// Simulate high transaction volume
	for i := 0; i < 1000; i++ {
		collector.RecordPayment("mpesa", "success")
		collector.RecordPaymentAmount("mpesa", float64(100+i))
		collector.RecordProcessingTime("mpesa", time.Duration(100+i)*time.Millisecond)
	}

	count := testutil.ToFloat64(collector.paymentsTotal.WithLabelValues("mpesa", "success"))
	if count != 1000 {
		t.Errorf("paymentsTotal = %v, want 1000", count)
	}
}
