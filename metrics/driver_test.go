package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestNewDriverCollector(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_driver")

	collector, err := NewDriverCollector(cfg)
	if err != nil {
		t.Fatalf("NewDriverCollector() error = %v", err)
	}
	if collector == nil {
		t.Fatal("NewDriverCollector() returned nil collector")
	}
}

func TestNewDriverCollector_DuplicateRegistration(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_driver_dup")

	collector1, err := NewDriverCollector(cfg)
	if err != nil {
		t.Fatalf("First NewDriverCollector() error = %v", err)
	}

	collector2, err := NewDriverCollector(cfg)
	if err != nil {
		t.Fatalf("Second NewDriverCollector() error = %v", err)
	}

	if collector1 == nil || collector2 == nil {
		t.Fatal("NewDriverCollector() returned nil collectors")
	}
}

func TestDriverCollector_SetDriversOnline(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_driver_online")

	collector, err := NewDriverCollector(cfg)
	if err != nil {
		t.Fatalf("NewDriverCollector() error = %v", err)
	}

	collector.SetDriversOnline("maputo", "standard", 50)
	collector.SetDriversOnline("maputo", "premium", 15)
	collector.SetDriversOnline("beira", "standard", 25)

	count := testutil.ToFloat64(collector.onlineTotal.WithLabelValues("maputo", "standard"))
	if count != 50 {
		t.Errorf("onlineTotal maputo/standard = %v, want 50", count)
	}

	count = testutil.ToFloat64(collector.onlineTotal.WithLabelValues("maputo", "premium"))
	if count != 15 {
		t.Errorf("onlineTotal maputo/premium = %v, want 15", count)
	}

	count = testutil.ToFloat64(collector.onlineTotal.WithLabelValues("beira", "standard"))
	if count != 25 {
		t.Errorf("onlineTotal beira/standard = %v, want 25", count)
	}

	// Test updating a value
	collector.SetDriversOnline("maputo", "standard", 45)
	count = testutil.ToFloat64(collector.onlineTotal.WithLabelValues("maputo", "standard"))
	if count != 45 {
		t.Errorf("onlineTotal maputo/standard after update = %v, want 45", count)
	}
}

func TestDriverCollector_SetAcceptanceRate(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_driver_accept")

	collector, err := NewDriverCollector(cfg)
	if err != nil {
		t.Fatalf("NewDriverCollector() error = %v", err)
	}

	collector.SetAcceptanceRate("driver_001", 0.95)
	collector.SetAcceptanceRate("driver_002", 0.80)
	collector.SetAcceptanceRate("driver_003", 0.65)

	rate := testutil.ToFloat64(collector.acceptanceRate.WithLabelValues("driver_001"))
	if rate != 0.95 {
		t.Errorf("acceptanceRate driver_001 = %v, want 0.95", rate)
	}

	rate = testutil.ToFloat64(collector.acceptanceRate.WithLabelValues("driver_002"))
	if rate != 0.80 {
		t.Errorf("acceptanceRate driver_002 = %v, want 0.80", rate)
	}

	// Test updating a value
	collector.SetAcceptanceRate("driver_001", 0.90)
	rate = testutil.ToFloat64(collector.acceptanceRate.WithLabelValues("driver_001"))
	if rate != 0.90 {
		t.Errorf("acceptanceRate driver_001 after update = %v, want 0.90", rate)
	}
}

func TestDriverCollector_SetRatingAverage(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_driver_rating")

	collector, err := NewDriverCollector(cfg)
	if err != nil {
		t.Fatalf("NewDriverCollector() error = %v", err)
	}

	collector.SetRatingAverage(4.5)

	rating := testutil.ToFloat64(collector.ratingAverage)
	if rating != 4.5 {
		t.Errorf("ratingAverage = %v, want 4.5", rating)
	}

	// Test updating
	collector.SetRatingAverage(4.7)
	rating = testutil.ToFloat64(collector.ratingAverage)
	if rating != 4.7 {
		t.Errorf("ratingAverage after update = %v, want 4.7", rating)
	}
}

func TestDriverCollector_AddEarnings(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_driver_earn")

	collector, err := NewDriverCollector(cfg)
	if err != nil {
		t.Fatalf("NewDriverCollector() error = %v", err)
	}

	collector.AddEarnings("driver_001", 500)
	collector.AddEarnings("driver_001", 750)
	collector.AddEarnings("driver_002", 1000)

	earnings := testutil.ToFloat64(collector.earnings.WithLabelValues("driver_001"))
	if earnings != 1250 {
		t.Errorf("earnings driver_001 = %v, want 1250", earnings)
	}

	earnings = testutil.ToFloat64(collector.earnings.WithLabelValues("driver_002"))
	if earnings != 1000 {
		t.Errorf("earnings driver_002 = %v, want 1000", earnings)
	}
}

func TestDriverCollector_Describe(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_driver_desc")

	collector, err := NewDriverCollector(cfg)
	if err != nil {
		t.Fatalf("NewDriverCollector() error = %v", err)
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

func TestDriverCollector_Collect(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_driver_collect")

	collector, err := NewDriverCollector(cfg)
	if err != nil {
		t.Fatalf("NewDriverCollector() error = %v", err)
	}

	// Record some metrics first
	collector.SetDriversOnline("maputo", "standard", 50)
	collector.SetRatingAverage(4.5)
	collector.AddEarnings("driver_001", 500)

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

func TestDriverCollector_MultipleCities(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_driver_cities")

	collector, err := NewDriverCollector(cfg)
	if err != nil {
		t.Fatalf("NewDriverCollector() error = %v", err)
	}

	cities := []string{"maputo", "beira", "nampula", "matola"}
	for _, city := range cities {
		collector.SetDriversOnline(city, "standard", 20)
		collector.SetDriversOnline(city, "premium", 5)
	}

	for _, city := range cities {
		stdCount := testutil.ToFloat64(collector.onlineTotal.WithLabelValues(city, "standard"))
		if stdCount != 20 {
			t.Errorf("onlineTotal %s/standard = %v, want 20", city, stdCount)
		}

		premCount := testutil.ToFloat64(collector.onlineTotal.WithLabelValues(city, "premium"))
		if premCount != 5 {
			t.Errorf("onlineTotal %s/premium = %v, want 5", city, premCount)
		}
	}
}

func TestDriverCollector_AcceptanceRateRange(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_driver_rate_range")

	collector, err := NewDriverCollector(cfg)
	if err != nil {
		t.Fatalf("NewDriverCollector() error = %v", err)
	}

	// Test boundary values
	collector.SetAcceptanceRate("driver_perfect", 1.0)
	collector.SetAcceptanceRate("driver_zero", 0.0)
	collector.SetAcceptanceRate("driver_half", 0.5)

	if rate := testutil.ToFloat64(collector.acceptanceRate.WithLabelValues("driver_perfect")); rate != 1.0 {
		t.Errorf("acceptanceRate driver_perfect = %v, want 1.0", rate)
	}

	if rate := testutil.ToFloat64(collector.acceptanceRate.WithLabelValues("driver_zero")); rate != 0.0 {
		t.Errorf("acceptanceRate driver_zero = %v, want 0.0", rate)
	}

	if rate := testutil.ToFloat64(collector.acceptanceRate.WithLabelValues("driver_half")); rate != 0.5 {
		t.Errorf("acceptanceRate driver_half = %v, want 0.5", rate)
	}
}
