package metrics

import (
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestNewRideCollector(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_ride")

	collector, err := NewRideCollector(cfg)
	if err != nil {
		t.Fatalf("NewRideCollector() error = %v", err)
	}
	if collector == nil {
		t.Fatal("NewRideCollector() returned nil collector")
	}
}

func TestNewRideCollector_DuplicateRegistration(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_ride_dup")

	collector1, err := NewRideCollector(cfg)
	if err != nil {
		t.Fatalf("First NewRideCollector() error = %v", err)
	}

	collector2, err := NewRideCollector(cfg)
	if err != nil {
		t.Fatalf("Second NewRideCollector() error = %v", err)
	}

	if collector1 == nil || collector2 == nil {
		t.Fatal("NewRideCollector() returned nil collectors")
	}
}

func TestRideCollector_RecordRideRequested(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_ride_req")

	collector, err := NewRideCollector(cfg)
	if err != nil {
		t.Fatalf("NewRideCollector() error = %v", err)
	}

	collector.RecordRideRequested("standard", "maputo")
	collector.RecordRideRequested("standard", "maputo")
	collector.RecordRideRequested("premium", "maputo")
	collector.RecordRideRequested("standard", "beira")

	count := testutil.ToFloat64(collector.requestedTotal.WithLabelValues("standard", "maputo"))
	if count != 2 {
		t.Errorf("requestedTotal standard/maputo = %v, want 2", count)
	}

	count = testutil.ToFloat64(collector.requestedTotal.WithLabelValues("premium", "maputo"))
	if count != 1 {
		t.Errorf("requestedTotal premium/maputo = %v, want 1", count)
	}

	count = testutil.ToFloat64(collector.requestedTotal.WithLabelValues("standard", "beira"))
	if count != 1 {
		t.Errorf("requestedTotal standard/beira = %v, want 1", count)
	}
}

func TestRideCollector_RecordRideCompleted(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_ride_comp")

	collector, err := NewRideCollector(cfg)
	if err != nil {
		t.Fatalf("NewRideCollector() error = %v", err)
	}

	collector.RecordRideCompleted("standard", "maputo")
	collector.RecordRideCompleted("standard", "maputo")
	collector.RecordRideCompleted("moto", "maputo")

	count := testutil.ToFloat64(collector.completedTotal.WithLabelValues("standard", "maputo"))
	if count != 2 {
		t.Errorf("completedTotal standard/maputo = %v, want 2", count)
	}

	count = testutil.ToFloat64(collector.completedTotal.WithLabelValues("moto", "maputo"))
	if count != 1 {
		t.Errorf("completedTotal moto/maputo = %v, want 1", count)
	}
}

func TestRideCollector_RecordRideCancelled(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_ride_cancel")

	collector, err := NewRideCollector(cfg)
	if err != nil {
		t.Fatalf("NewRideCollector() error = %v", err)
	}

	collector.RecordRideCancelled("rider", "changed_mind")
	collector.RecordRideCancelled("rider", "changed_mind")
	collector.RecordRideCancelled("driver", "no_show")
	collector.RecordRideCancelled("system", "no_drivers")

	count := testutil.ToFloat64(collector.cancelledTotal.WithLabelValues("rider", "changed_mind"))
	if count != 2 {
		t.Errorf("cancelledTotal rider/changed_mind = %v, want 2", count)
	}

	count = testutil.ToFloat64(collector.cancelledTotal.WithLabelValues("driver", "no_show"))
	if count != 1 {
		t.Errorf("cancelledTotal driver/no_show = %v, want 1", count)
	}

	count = testutil.ToFloat64(collector.cancelledTotal.WithLabelValues("system", "no_drivers"))
	if count != 1 {
		t.Errorf("cancelledTotal system/no_drivers = %v, want 1", count)
	}
}

func TestRideCollector_RecordRideDuration(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_ride_dur")

	collector, err := NewRideCollector(cfg)
	if err != nil {
		t.Fatalf("NewRideCollector() error = %v", err)
	}

	collector.RecordRideDuration("standard", 15*time.Minute)
	collector.RecordRideDuration("standard", 30*time.Minute)
	collector.RecordRideDuration("premium", 45*time.Minute)

	histCount := testutil.CollectAndCount(collector.duration)
	if histCount == 0 {
		t.Error("duration histogram has no metrics")
	}
}

func TestRideCollector_RecordRideDistance(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_ride_dist")

	collector, err := NewRideCollector(cfg)
	if err != nil {
		t.Fatalf("NewRideCollector() error = %v", err)
	}

	collector.RecordRideDistance("standard", 5.5)
	collector.RecordRideDistance("standard", 12.3)
	collector.RecordRideDistance("premium", 25.0)

	histCount := testutil.CollectAndCount(collector.distance)
	if histCount == 0 {
		t.Error("distance histogram has no metrics")
	}
}

func TestRideCollector_RecordRideFare(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_ride_fare")

	collector, err := NewRideCollector(cfg)
	if err != nil {
		t.Fatalf("NewRideCollector() error = %v", err)
	}

	collector.RecordRideFare("standard", 150)
	collector.RecordRideFare("standard", 350)
	collector.RecordRideFare("premium", 1500)

	histCount := testutil.CollectAndCount(collector.fare)
	if histCount == 0 {
		t.Error("fare histogram has no metrics")
	}
}

func TestRideCollector_RecordRideWaitTime(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_ride_wait")

	collector, err := NewRideCollector(cfg)
	if err != nil {
		t.Fatalf("NewRideCollector() error = %v", err)
	}

	collector.RecordRideWaitTime("standard", 2*time.Minute)
	collector.RecordRideWaitTime("standard", 5*time.Minute)
	collector.RecordRideWaitTime("moto", 1*time.Minute)

	histCount := testutil.CollectAndCount(collector.waitTime)
	if histCount == 0 {
		t.Error("waitTime histogram has no metrics")
	}
}

func TestRideCollector_Describe(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_ride_desc")

	collector, err := NewRideCollector(cfg)
	if err != nil {
		t.Fatalf("NewRideCollector() error = %v", err)
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

func TestRideCollector_Collect(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_ride_collect")

	collector, err := NewRideCollector(cfg)
	if err != nil {
		t.Fatalf("NewRideCollector() error = %v", err)
	}

	// Record some metrics first
	collector.RecordRideRequested("standard", "maputo")
	collector.RecordRideCompleted("standard", "maputo")

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

func TestRideCollector_AllServiceTypes(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_ride_types")

	collector, err := NewRideCollector(cfg)
	if err != nil {
		t.Fatalf("NewRideCollector() error = %v", err)
	}

	serviceTypes := []string{"standard", "premium", "moto", "comfort"}
	for _, st := range serviceTypes {
		collector.RecordRideRequested(st, "maputo")
		collector.RecordRideCompleted(st, "maputo")
		collector.RecordRideDuration(st, 10*time.Minute)
		collector.RecordRideDistance(st, 5.0)
		collector.RecordRideFare(st, 200)
		collector.RecordRideWaitTime(st, 3*time.Minute)
	}

	for _, st := range serviceTypes {
		reqCount := testutil.ToFloat64(collector.requestedTotal.WithLabelValues(st, "maputo"))
		if reqCount != 1 {
			t.Errorf("requestedTotal %s/maputo = %v, want 1", st, reqCount)
		}

		compCount := testutil.ToFloat64(collector.completedTotal.WithLabelValues(st, "maputo"))
		if compCount != 1 {
			t.Errorf("completedTotal %s/maputo = %v, want 1", st, compCount)
		}
	}
}
