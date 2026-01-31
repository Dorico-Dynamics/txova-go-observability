package health

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockPinger is a mock implementation of Pinger.
type mockPinger struct {
	err error
}

func (m *mockPinger) Ping(ctx context.Context) error {
	return m.err
}

// mockKafkaClient is a mock implementation of KafkaMetadataFetcher.
type mockKafkaClient struct {
	err error
}

func (m *mockKafkaClient) GetMetadata(ctx context.Context) error {
	return m.err
}

func TestRedisChecker_Name(t *testing.T) {
	t.Parallel()

	checker := NewRedisChecker("redis-primary", &mockPinger{}, true)

	if checker.Name() != "redis-primary" {
		t.Errorf("Name() = %v, want redis-primary", checker.Name())
	}
}

func TestRedisChecker_Required(t *testing.T) {
	t.Parallel()

	required := NewRedisChecker("redis", &mockPinger{}, true)
	optional := NewRedisChecker("redis", &mockPinger{}, false)

	if !required.Required() {
		t.Error("Required() should be true")
	}
	if optional.Required() {
		t.Error("Required() should be false")
	}
}

func TestRedisChecker_Check_Healthy(t *testing.T) {
	t.Parallel()

	checker := NewRedisChecker("redis", &mockPinger{err: nil}, true)
	result := checker.Check(context.Background())

	if result.Status != StatusHealthy {
		t.Errorf("Status = %v, want %v", result.Status, StatusHealthy)
	}
}

func TestRedisChecker_Check_Unhealthy(t *testing.T) {
	t.Parallel()

	checker := NewRedisChecker("redis", &mockPinger{err: errors.New("connection refused")}, true)
	result := checker.Check(context.Background())

	if result.Status != StatusUnhealthy {
		t.Errorf("Status = %v, want %v", result.Status, StatusUnhealthy)
	}
	if result.Error != "connection refused" {
		t.Errorf("Error = %v, want 'connection refused'", result.Error)
	}
}

func TestKafkaChecker_Name(t *testing.T) {
	t.Parallel()

	checker := NewKafkaChecker("kafka-cluster", &mockKafkaClient{}, true)

	if checker.Name() != "kafka-cluster" {
		t.Errorf("Name() = %v, want kafka-cluster", checker.Name())
	}
}

func TestKafkaChecker_Required(t *testing.T) {
	t.Parallel()

	required := NewKafkaChecker("kafka", &mockKafkaClient{}, true)
	optional := NewKafkaChecker("kafka", &mockKafkaClient{}, false)

	if !required.Required() {
		t.Error("Required() should be true")
	}
	if optional.Required() {
		t.Error("Required() should be false")
	}
}

func TestKafkaChecker_Check_Healthy(t *testing.T) {
	t.Parallel()

	checker := NewKafkaChecker("kafka", &mockKafkaClient{err: nil}, true)
	result := checker.Check(context.Background())

	if result.Status != StatusHealthy {
		t.Errorf("Status = %v, want %v", result.Status, StatusHealthy)
	}
}

func TestKafkaChecker_Check_Unhealthy(t *testing.T) {
	t.Parallel()

	checker := NewKafkaChecker("kafka", &mockKafkaClient{err: errors.New("broker unavailable")}, true)
	result := checker.Check(context.Background())

	if result.Status != StatusUnhealthy {
		t.Errorf("Status = %v, want %v", result.Status, StatusUnhealthy)
	}
	if result.Error != "broker unavailable" {
		t.Errorf("Error = %v, want 'broker unavailable'", result.Error)
	}
}

func TestHTTPChecker_Name(t *testing.T) {
	t.Parallel()

	checker := NewHTTPChecker("external-api", "http://example.com/health", nil, false)

	if checker.Name() != "external-api" {
		t.Errorf("Name() = %v, want external-api", checker.Name())
	}
}

func TestHTTPChecker_Required(t *testing.T) {
	t.Parallel()

	required := NewHTTPChecker("api", "http://example.com", nil, true)
	optional := NewHTTPChecker("api", "http://example.com", nil, false)

	if !required.Required() {
		t.Error("Required() should be true")
	}
	if optional.Required() {
		t.Error("Required() should be false")
	}
}

func TestHTTPChecker_Check_Healthy(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	checker := NewHTTPChecker("test-api", server.URL, nil, true)
	result := checker.Check(context.Background())

	if result.Status != StatusHealthy {
		t.Errorf("Status = %v, want %v", result.Status, StatusHealthy)
	}
	if result.Details["status_code"] != 200 {
		t.Errorf("Details[status_code] = %v, want 200", result.Details["status_code"])
	}
}

func TestHTTPChecker_Check_Unhealthy(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	checker := NewHTTPChecker("test-api", server.URL, nil, true)
	result := checker.Check(context.Background())

	if result.Status != StatusUnhealthy {
		t.Errorf("Status = %v, want %v", result.Status, StatusUnhealthy)
	}
}

func TestHTTPChecker_Check_ConnectionError(t *testing.T) {
	t.Parallel()

	// Use a valid but unreachable address (non-routable IP)
	checker := NewHTTPChecker("test-api", "http://10.255.255.1:12345", nil, true)
	result := checker.Check(context.Background())

	if result.Status != StatusUnhealthy {
		t.Errorf("Status = %v, want %v", result.Status, StatusUnhealthy)
	}
}

func TestHTTPChecker_WithExpectedStatusCode(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	checker := NewHTTPChecker("test-api", server.URL, nil, true).WithExpectedStatusCode(http.StatusAccepted)
	result := checker.Check(context.Background())

	if result.Status != StatusHealthy {
		t.Errorf("Status = %v, want %v", result.Status, StatusHealthy)
	}
}

func TestFuncChecker_Name(t *testing.T) {
	t.Parallel()

	checker := NewFuncChecker("custom-check", func(ctx context.Context) error {
		return nil
	}, true)

	if checker.Name() != "custom-check" {
		t.Errorf("Name() = %v, want custom-check", checker.Name())
	}
}

func TestFuncChecker_Required(t *testing.T) {
	t.Parallel()

	required := NewFuncChecker("check", func(ctx context.Context) error { return nil }, true)
	optional := NewFuncChecker("check", func(ctx context.Context) error { return nil }, false)

	if !required.Required() {
		t.Error("Required() should be true")
	}
	if optional.Required() {
		t.Error("Required() should be false")
	}
}

func TestFuncChecker_Check_Healthy(t *testing.T) {
	t.Parallel()

	checker := NewFuncChecker("check", func(ctx context.Context) error {
		return nil
	}, true)
	result := checker.Check(context.Background())

	if result.Status != StatusHealthy {
		t.Errorf("Status = %v, want %v", result.Status, StatusHealthy)
	}
}

func TestFuncChecker_Check_Unhealthy(t *testing.T) {
	t.Parallel()

	checker := NewFuncChecker("check", func(ctx context.Context) error {
		return errors.New("custom error")
	}, true)
	result := checker.Check(context.Background())

	if result.Status != StatusUnhealthy {
		t.Errorf("Status = %v, want %v", result.Status, StatusUnhealthy)
	}
	if result.Error != "custom error" {
		t.Errorf("Error = %v, want 'custom error'", result.Error)
	}
}

func TestFuncChecker_Check_ContextCancellation(t *testing.T) {
	t.Parallel()

	checker := NewFuncChecker("check", func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return nil
		}
	}, true)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result := checker.Check(ctx)

	if result.Status != StatusUnhealthy {
		t.Errorf("Status = %v, want %v", result.Status, StatusUnhealthy)
	}
}
