package metrics

import (
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestNewHTTPCollector(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_http")

	collector, err := NewHTTPCollector(cfg)
	if err != nil {
		t.Fatalf("NewHTTPCollector() error = %v", err)
	}
	if collector == nil {
		t.Fatal("NewHTTPCollector() returned nil collector")
	}
}

func TestNewHTTPCollector_DuplicateRegistration(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_http_dup")

	collector1, err := NewHTTPCollector(cfg)
	if err != nil {
		t.Fatalf("First NewHTTPCollector() error = %v", err)
	}

	collector2, err := NewHTTPCollector(cfg)
	if err != nil {
		t.Fatalf("Second NewHTTPCollector() error = %v", err)
	}

	if collector1 == nil || collector2 == nil {
		t.Fatal("NewHTTPCollector() returned nil collectors")
	}
}

func TestHTTPCollector_RecordRequest(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_record_req")

	collector, err := NewHTTPCollector(cfg)
	if err != nil {
		t.Fatalf("NewHTTPCollector() error = %v", err)
	}

	collector.RecordRequest("GET", "/api/v1/rides", 200, 50*time.Millisecond)
	collector.RecordRequest("POST", "/api/v1/rides", 201, 100*time.Millisecond)
	collector.RecordRequest("GET", "/api/v1/rides", 500, 200*time.Millisecond)

	// Verify counter was incremented
	count := testutil.ToFloat64(collector.requestsTotal.WithLabelValues("GET", "/api/v1/rides", "200"))
	if count != 1 {
		t.Errorf("requestsTotal for GET 200 = %v, want 1", count)
	}

	count = testutil.ToFloat64(collector.requestsTotal.WithLabelValues("POST", "/api/v1/rides", "201"))
	if count != 1 {
		t.Errorf("requestsTotal for POST 201 = %v, want 1", count)
	}

	count = testutil.ToFloat64(collector.requestsTotal.WithLabelValues("GET", "/api/v1/rides", "500"))
	if count != 1 {
		t.Errorf("requestsTotal for GET 500 = %v, want 1", count)
	}
}

func TestHTTPCollector_RecordPanic(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_panic")

	collector, err := NewHTTPCollector(cfg)
	if err != nil {
		t.Fatalf("NewHTTPCollector() error = %v", err)
	}

	collector.RecordPanic("GET", "/api/v1/rides")
	collector.RecordPanic("GET", "/api/v1/rides")

	count := testutil.ToFloat64(collector.panicsTotal.WithLabelValues("GET", "/api/v1/rides"))
	if count != 2 {
		t.Errorf("panicsTotal = %v, want 2", count)
	}
}

func TestHTTPCollector_RecordRequestSize(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_req_size")

	collector, err := NewHTTPCollector(cfg)
	if err != nil {
		t.Fatalf("NewHTTPCollector() error = %v", err)
	}

	collector.RecordRequestSize("POST", "/api/v1/rides", 1024)
	collector.RecordRequestSize("POST", "/api/v1/rides", 2048)

	// Histogram sum should be 3072
	histCount := testutil.CollectAndCount(collector.requestSize)
	if histCount == 0 {
		t.Error("requestSize histogram has no metrics")
	}
}

func TestHTTPCollector_RecordResponseSize(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_resp_size")

	collector, err := NewHTTPCollector(cfg)
	if err != nil {
		t.Fatalf("NewHTTPCollector() error = %v", err)
	}

	collector.RecordResponseSize("GET", "/api/v1/rides", 4096)

	histCount := testutil.CollectAndCount(collector.responseSize)
	if histCount == 0 {
		t.Error("responseSize histogram has no metrics")
	}
}

func TestHTTPCollector_RequestsInFlight(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_in_flight")

	collector, err := NewHTTPCollector(cfg)
	if err != nil {
		t.Fatalf("NewHTTPCollector() error = %v", err)
	}

	collector.IncRequestsInFlight()
	collector.IncRequestsInFlight()
	collector.IncRequestsInFlight()

	count := testutil.ToFloat64(collector.requestsInFlight)
	if count != 3 {
		t.Errorf("requestsInFlight = %v, want 3", count)
	}

	collector.DecRequestsInFlight()
	count = testutil.ToFloat64(collector.requestsInFlight)
	if count != 2 {
		t.Errorf("requestsInFlight after dec = %v, want 2", count)
	}
}

func TestHTTPCollector_Describe(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_describe")

	collector, err := NewHTTPCollector(cfg)
	if err != nil {
		t.Fatalf("NewHTTPCollector() error = %v", err)
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

func TestHTTPCollector_Collect(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_collect")

	collector, err := NewHTTPCollector(cfg)
	if err != nil {
		t.Fatalf("NewHTTPCollector() error = %v", err)
	}

	// Record some metrics first
	collector.RecordRequest("GET", "/test", 200, 10*time.Millisecond)

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

func TestHTTPCollector_MultipleRequests(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_multi")

	collector, err := NewHTTPCollector(cfg)
	if err != nil {
		t.Fatalf("NewHTTPCollector() error = %v", err)
	}

	// Simulate multiple requests
	for i := 0; i < 100; i++ {
		collector.RecordRequest("GET", "/health", 200, time.Duration(i)*time.Millisecond)
	}

	count := testutil.ToFloat64(collector.requestsTotal.WithLabelValues("GET", "/health", "200"))
	if count != 100 {
		t.Errorf("requestsTotal = %v, want 100", count)
	}
}
