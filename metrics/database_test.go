package metrics

import (
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestNewDBCollector(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_db")

	collector, err := NewDBCollector(cfg)
	if err != nil {
		t.Fatalf("NewDBCollector() error = %v", err)
	}
	if collector == nil {
		t.Fatal("NewDBCollector() returned nil collector")
	}
}

func TestNewDBCollector_DuplicateRegistration(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_db_dup")

	collector1, err := NewDBCollector(cfg)
	if err != nil {
		t.Fatalf("First NewDBCollector() error = %v", err)
	}

	collector2, err := NewDBCollector(cfg)
	if err != nil {
		t.Fatalf("Second NewDBCollector() error = %v", err)
	}

	if collector1 == nil || collector2 == nil {
		t.Fatal("NewDBCollector() returned nil collectors")
	}
}

func TestDBCollector_SetConnections(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_db_conn")

	collector, err := NewDBCollector(cfg)
	if err != nil {
		t.Fatalf("NewDBCollector() error = %v", err)
	}

	collector.SetConnections("primary", "idle", 5)
	collector.SetConnections("primary", "in_use", 10)
	collector.SetConnections("replica", "idle", 3)

	count := testutil.ToFloat64(collector.connectionsTotal.WithLabelValues("primary", "idle"))
	if count != 5 {
		t.Errorf("connections primary/idle = %v, want 5", count)
	}

	count = testutil.ToFloat64(collector.connectionsTotal.WithLabelValues("primary", "in_use"))
	if count != 10 {
		t.Errorf("connections primary/in_use = %v, want 10", count)
	}

	count = testutil.ToFloat64(collector.connectionsTotal.WithLabelValues("replica", "idle"))
	if count != 3 {
		t.Errorf("connections replica/idle = %v, want 3", count)
	}

	// Test updating a value
	collector.SetConnections("primary", "idle", 8)
	count = testutil.ToFloat64(collector.connectionsTotal.WithLabelValues("primary", "idle"))
	if count != 8 {
		t.Errorf("connections primary/idle after update = %v, want 8", count)
	}
}

func TestDBCollector_RecordQueryDuration(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_db_query")

	collector, err := NewDBCollector(cfg)
	if err != nil {
		t.Fatalf("NewDBCollector() error = %v", err)
	}

	collector.RecordQueryDuration("select", 10*time.Millisecond)
	collector.RecordQueryDuration("select", 20*time.Millisecond)
	collector.RecordQueryDuration("insert", 5*time.Millisecond)

	histCount := testutil.CollectAndCount(collector.queryDuration)
	if histCount == 0 {
		t.Error("queryDuration histogram has no metrics")
	}
}

func TestDBCollector_RecordQueryError(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_db_err")

	collector, err := NewDBCollector(cfg)
	if err != nil {
		t.Fatalf("NewDBCollector() error = %v", err)
	}

	collector.RecordQueryError("insert", "constraint")
	collector.RecordQueryError("insert", "constraint")
	collector.RecordQueryError("select", "timeout")

	count := testutil.ToFloat64(collector.queryErrorsTotal.WithLabelValues("insert", "constraint"))
	if count != 2 {
		t.Errorf("queryErrors insert/constraint = %v, want 2", count)
	}

	count = testutil.ToFloat64(collector.queryErrorsTotal.WithLabelValues("select", "timeout"))
	if count != 1 {
		t.Errorf("queryErrors select/timeout = %v, want 1", count)
	}
}

func TestDBCollector_RecordTransactionDuration(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_db_tx")

	collector, err := NewDBCollector(cfg)
	if err != nil {
		t.Fatalf("NewDBCollector() error = %v", err)
	}

	collector.RecordTransactionDuration(50 * time.Millisecond)
	collector.RecordTransactionDuration(100 * time.Millisecond)

	histCount := testutil.CollectAndCount(collector.transactionDuration)
	if histCount == 0 {
		t.Error("transactionDuration histogram has no metrics")
	}
}

func TestDBCollector_Describe(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_db_desc")

	collector, err := NewDBCollector(cfg)
	if err != nil {
		t.Fatalf("NewDBCollector() error = %v", err)
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

func TestDBCollector_Collect(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_db_collect")

	collector, err := NewDBCollector(cfg)
	if err != nil {
		t.Fatalf("NewDBCollector() error = %v", err)
	}

	// Record some metrics first
	collector.SetConnections("primary", "idle", 5)
	collector.RecordQueryDuration("select", 10*time.Millisecond)

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

func TestDBCollector_AllOperationTypes(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_db_ops")

	collector, err := NewDBCollector(cfg)
	if err != nil {
		t.Fatalf("NewDBCollector() error = %v", err)
	}

	operations := []string{"select", "insert", "update", "delete"}
	for _, op := range operations {
		collector.RecordQueryDuration(op, 10*time.Millisecond)
	}

	histCount := testutil.CollectAndCount(collector.queryDuration)
	if histCount == 0 {
		t.Error("queryDuration histogram has no metrics after recording all operation types")
	}
}

func TestDBCollector_AllErrorTypes(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_db_all_err")

	collector, err := NewDBCollector(cfg)
	if err != nil {
		t.Fatalf("NewDBCollector() error = %v", err)
	}

	errorTypes := []string{"timeout", "connection", "constraint", "deadlock"}
	for _, errType := range errorTypes {
		collector.RecordQueryError("select", errType)
	}

	for _, errType := range errorTypes {
		count := testutil.ToFloat64(collector.queryErrorsTotal.WithLabelValues("select", errType))
		if count != 1 {
			t.Errorf("queryErrors select/%s = %v, want 1", errType, count)
		}
	}
}
