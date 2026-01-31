package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestNewSafetyCollector(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_safety")

	collector, err := NewSafetyCollector(cfg)
	if err != nil {
		t.Fatalf("NewSafetyCollector() error = %v", err)
	}
	if collector == nil {
		t.Fatal("NewSafetyCollector() returned nil collector")
	}
}

func TestNewSafetyCollector_DuplicateRegistration(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_safety_dup")

	collector1, err := NewSafetyCollector(cfg)
	if err != nil {
		t.Fatalf("First NewSafetyCollector() error = %v", err)
	}

	collector2, err := NewSafetyCollector(cfg)
	if err != nil {
		t.Fatalf("Second NewSafetyCollector() error = %v", err)
	}

	if collector1 == nil || collector2 == nil {
		t.Fatal("NewSafetyCollector() returned nil collectors")
	}
}

func TestSafetyCollector_RecordEmergency(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_safety_emer")

	collector, err := NewSafetyCollector(cfg)
	if err != nil {
		t.Fatalf("NewSafetyCollector() error = %v", err)
	}

	collector.RecordEmergency("sos_button", "maputo")
	collector.RecordEmergency("sos_button", "maputo")
	collector.RecordEmergency("auto_detected", "maputo")
	collector.RecordEmergency("police_request", "beira")

	count := testutil.ToFloat64(collector.emergenciesTotal.WithLabelValues("sos_button", "maputo"))
	if count != 2 {
		t.Errorf("emergenciesTotal sos_button/maputo = %v, want 2", count)
	}

	count = testutil.ToFloat64(collector.emergenciesTotal.WithLabelValues("auto_detected", "maputo"))
	if count != 1 {
		t.Errorf("emergenciesTotal auto_detected/maputo = %v, want 1", count)
	}

	count = testutil.ToFloat64(collector.emergenciesTotal.WithLabelValues("police_request", "beira"))
	if count != 1 {
		t.Errorf("emergenciesTotal police_request/beira = %v, want 1", count)
	}
}

func TestSafetyCollector_RecordIncident(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_safety_inc")

	collector, err := NewSafetyCollector(cfg)
	if err != nil {
		t.Fatalf("NewSafetyCollector() error = %v", err)
	}

	collector.RecordIncident("low")
	collector.RecordIncident("low")
	collector.RecordIncident("medium")
	collector.RecordIncident("high")
	collector.RecordIncident("critical")

	count := testutil.ToFloat64(collector.incidentsTotal.WithLabelValues("low"))
	if count != 2 {
		t.Errorf("incidentsTotal low = %v, want 2", count)
	}

	count = testutil.ToFloat64(collector.incidentsTotal.WithLabelValues("medium"))
	if count != 1 {
		t.Errorf("incidentsTotal medium = %v, want 1", count)
	}

	count = testutil.ToFloat64(collector.incidentsTotal.WithLabelValues("high"))
	if count != 1 {
		t.Errorf("incidentsTotal high = %v, want 1", count)
	}

	count = testutil.ToFloat64(collector.incidentsTotal.WithLabelValues("critical"))
	if count != 1 {
		t.Errorf("incidentsTotal critical = %v, want 1", count)
	}
}

func TestSafetyCollector_RecordTripShare(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_safety_share")

	collector, err := NewSafetyCollector(cfg)
	if err != nil {
		t.Fatalf("NewSafetyCollector() error = %v", err)
	}

	collector.RecordTripShare()
	collector.RecordTripShare()
	collector.RecordTripShare()

	count := testutil.ToFloat64(collector.tripSharesTotal)
	if count != 3 {
		t.Errorf("tripSharesTotal = %v, want 3", count)
	}
}

func TestSafetyCollector_Describe(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_safety_desc")

	collector, err := NewSafetyCollector(cfg)
	if err != nil {
		t.Fatalf("NewSafetyCollector() error = %v", err)
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

func TestSafetyCollector_Collect(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_safety_collect")

	collector, err := NewSafetyCollector(cfg)
	if err != nil {
		t.Fatalf("NewSafetyCollector() error = %v", err)
	}

	// Record some metrics first
	collector.RecordEmergency("sos_button", "maputo")
	collector.RecordIncident("high")
	collector.RecordTripShare()

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

func TestSafetyCollector_AllEmergencyTypes(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_safety_emer_types")

	collector, err := NewSafetyCollector(cfg)
	if err != nil {
		t.Fatalf("NewSafetyCollector() error = %v", err)
	}

	emergencyTypes := []string{"sos_button", "auto_detected", "police_request", "driver_alert"}
	for _, et := range emergencyTypes {
		collector.RecordEmergency(et, "maputo")
	}

	for _, et := range emergencyTypes {
		count := testutil.ToFloat64(collector.emergenciesTotal.WithLabelValues(et, "maputo"))
		if count != 1 {
			t.Errorf("emergenciesTotal %s/maputo = %v, want 1", et, count)
		}
	}
}

func TestSafetyCollector_AllSeverityLevels(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_safety_sev")

	collector, err := NewSafetyCollector(cfg)
	if err != nil {
		t.Fatalf("NewSafetyCollector() error = %v", err)
	}

	severities := []string{"low", "medium", "high", "critical"}
	for _, sev := range severities {
		collector.RecordIncident(sev)
	}

	for _, sev := range severities {
		count := testutil.ToFloat64(collector.incidentsTotal.WithLabelValues(sev))
		if count != 1 {
			t.Errorf("incidentsTotal %s = %v, want 1", sev, count)
		}
	}
}

func TestSafetyCollector_MultipleCities(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_safety_cities")

	collector, err := NewSafetyCollector(cfg)
	if err != nil {
		t.Fatalf("NewSafetyCollector() error = %v", err)
	}

	cities := []string{"maputo", "beira", "nampula", "matola"}
	for _, city := range cities {
		collector.RecordEmergency("sos_button", city)
	}

	for _, city := range cities {
		count := testutil.ToFloat64(collector.emergenciesTotal.WithLabelValues("sos_button", city))
		if count != 1 {
			t.Errorf("emergenciesTotal sos_button/%s = %v, want 1", city, count)
		}
	}
}

func TestSafetyCollector_TripSharesHighVolume(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_safety_volume")

	collector, err := NewSafetyCollector(cfg)
	if err != nil {
		t.Fatalf("NewSafetyCollector() error = %v", err)
	}

	// Simulate high volume of trip shares
	for i := 0; i < 1000; i++ {
		collector.RecordTripShare()
	}

	count := testutil.ToFloat64(collector.tripSharesTotal)
	if count != 1000 {
		t.Errorf("tripSharesTotal = %v, want 1000", count)
	}
}
