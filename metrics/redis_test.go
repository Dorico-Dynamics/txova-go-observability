package metrics

import (
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestNewRedisCollector(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_redis")

	collector, err := NewRedisCollector(cfg)
	if err != nil {
		t.Fatalf("NewRedisCollector() error = %v", err)
	}
	if collector == nil {
		t.Fatal("NewRedisCollector() returned nil collector")
	}
}

func TestNewRedisCollector_DuplicateRegistration(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_redis_dup")

	collector1, err := NewRedisCollector(cfg)
	if err != nil {
		t.Fatalf("First NewRedisCollector() error = %v", err)
	}

	collector2, err := NewRedisCollector(cfg)
	if err != nil {
		t.Fatalf("Second NewRedisCollector() error = %v", err)
	}

	if collector1 == nil || collector2 == nil {
		t.Fatal("NewRedisCollector() returned nil collectors")
	}
}

func TestRedisCollector_RecordCommand(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_redis_cmd")

	collector, err := NewRedisCollector(cfg)
	if err != nil {
		t.Fatalf("NewRedisCollector() error = %v", err)
	}

	collector.RecordCommand("GET", 1*time.Millisecond)
	collector.RecordCommand("GET", 2*time.Millisecond)
	collector.RecordCommand("SET", 3*time.Millisecond)

	count := testutil.ToFloat64(collector.commandsTotal.WithLabelValues("GET"))
	if count != 2 {
		t.Errorf("commandsTotal GET = %v, want 2", count)
	}

	count = testutil.ToFloat64(collector.commandsTotal.WithLabelValues("SET"))
	if count != 1 {
		t.Errorf("commandsTotal SET = %v, want 1", count)
	}

	histCount := testutil.CollectAndCount(collector.commandDuration)
	if histCount == 0 {
		t.Error("commandDuration histogram has no metrics")
	}
}

func TestRedisCollector_RecordCacheHit(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_redis_hit")

	collector, err := NewRedisCollector(cfg)
	if err != nil {
		t.Fatalf("NewRedisCollector() error = %v", err)
	}

	collector.RecordCacheHit("user_session")
	collector.RecordCacheHit("user_session")
	collector.RecordCacheHit("ride_status")

	count := testutil.ToFloat64(collector.cacheHitsTotal.WithLabelValues("user_session"))
	if count != 2 {
		t.Errorf("cacheHitsTotal user_session = %v, want 2", count)
	}

	count = testutil.ToFloat64(collector.cacheHitsTotal.WithLabelValues("ride_status"))
	if count != 1 {
		t.Errorf("cacheHitsTotal ride_status = %v, want 1", count)
	}
}

func TestRedisCollector_RecordCacheMiss(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_redis_miss")

	collector, err := NewRedisCollector(cfg)
	if err != nil {
		t.Fatalf("NewRedisCollector() error = %v", err)
	}

	collector.RecordCacheMiss("user_session")
	collector.RecordCacheMiss("driver_location")
	collector.RecordCacheMiss("driver_location")
	collector.RecordCacheMiss("driver_location")

	count := testutil.ToFloat64(collector.cacheMissTotal.WithLabelValues("user_session"))
	if count != 1 {
		t.Errorf("cacheMissTotal user_session = %v, want 1", count)
	}

	count = testutil.ToFloat64(collector.cacheMissTotal.WithLabelValues("driver_location"))
	if count != 3 {
		t.Errorf("cacheMissTotal driver_location = %v, want 3", count)
	}
}

func TestRedisCollector_Describe(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_redis_desc")

	collector, err := NewRedisCollector(cfg)
	if err != nil {
		t.Fatalf("NewRedisCollector() error = %v", err)
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

func TestRedisCollector_Collect(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_redis_collect")

	collector, err := NewRedisCollector(cfg)
	if err != nil {
		t.Fatalf("NewRedisCollector() error = %v", err)
	}

	// Record some metrics first
	collector.RecordCommand("GET", 1*time.Millisecond)
	collector.RecordCacheHit("test")

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

func TestRedisCollector_AllCommands(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_redis_all_cmd")

	collector, err := NewRedisCollector(cfg)
	if err != nil {
		t.Fatalf("NewRedisCollector() error = %v", err)
	}

	commands := []string{"GET", "SET", "HGET", "HSET", "LPUSH", "RPOP", "SADD", "DEL", "EXPIRE"}
	for _, cmd := range commands {
		collector.RecordCommand(cmd, 1*time.Millisecond)
	}

	for _, cmd := range commands {
		count := testutil.ToFloat64(collector.commandsTotal.WithLabelValues(cmd))
		if count != 1 {
			t.Errorf("commandsTotal %s = %v, want 1", cmd, count)
		}
	}
}

func TestRedisCollector_CacheHitMissRatio(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	cfg := DefaultConfig().WithRegistry(registry).WithSubsystem("test_redis_ratio")

	collector, err := NewRedisCollector(cfg)
	if err != nil {
		t.Fatalf("NewRedisCollector() error = %v", err)
	}

	// Simulate 80% hit rate
	for i := 0; i < 80; i++ {
		collector.RecordCacheHit("sessions")
	}
	for i := 0; i < 20; i++ {
		collector.RecordCacheMiss("sessions")
	}

	hits := testutil.ToFloat64(collector.cacheHitsTotal.WithLabelValues("sessions"))
	misses := testutil.ToFloat64(collector.cacheMissTotal.WithLabelValues("sessions"))

	if hits != 80 {
		t.Errorf("cacheHitsTotal = %v, want 80", hits)
	}
	if misses != 20 {
		t.Errorf("cacheMissTotal = %v, want 20", misses)
	}
}
