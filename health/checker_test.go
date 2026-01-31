package health

import (
	"errors"
	"testing"
	"time"
)

func TestNewHealthyResult(t *testing.T) {
	t.Parallel()

	result := NewHealthyResult(100 * time.Millisecond)

	if result.Status != StatusHealthy {
		t.Errorf("Status = %v, want %v", result.Status, StatusHealthy)
	}
	if result.DurationMS != 100 {
		t.Errorf("DurationMS = %v, want 100", result.DurationMS)
	}
	if result.Error != "" {
		t.Errorf("Error = %v, want empty", result.Error)
	}
	if result.Timestamp.IsZero() {
		t.Error("Timestamp should not be zero")
	}
}

func TestNewUnhealthyResult(t *testing.T) {
	t.Parallel()

	err := errors.New("connection failed")
	result := NewUnhealthyResult(50*time.Millisecond, err)

	if result.Status != StatusUnhealthy {
		t.Errorf("Status = %v, want %v", result.Status, StatusUnhealthy)
	}
	if result.DurationMS != 50 {
		t.Errorf("DurationMS = %v, want 50", result.DurationMS)
	}
	if result.Error != "connection failed" {
		t.Errorf("Error = %v, want 'connection failed'", result.Error)
	}
}

func TestNewUnhealthyResult_NilError(t *testing.T) {
	t.Parallel()

	result := NewUnhealthyResult(50*time.Millisecond, nil)

	if result.Status != StatusUnhealthy {
		t.Errorf("Status = %v, want %v", result.Status, StatusUnhealthy)
	}
	if result.Error != "" {
		t.Errorf("Error = %v, want empty", result.Error)
	}
}

func TestNewDegradedResult(t *testing.T) {
	t.Parallel()

	result := NewDegradedResult(75*time.Millisecond, "high latency")

	if result.Status != StatusDegraded {
		t.Errorf("Status = %v, want %v", result.Status, StatusDegraded)
	}
	if result.DurationMS != 75 {
		t.Errorf("DurationMS = %v, want 75", result.DurationMS)
	}
	if result.Error != "high latency" {
		t.Errorf("Error = %v, want 'high latency'", result.Error)
	}
}

func TestResult_WithDetails(t *testing.T) {
	t.Parallel()

	result := NewHealthyResult(100 * time.Millisecond).WithDetails(map[string]any{
		"connections": 10,
		"version":     "14.0",
	})

	if result.Details == nil {
		t.Fatal("Details should not be nil")
	}
	if result.Details["connections"] != 10 {
		t.Errorf("Details[connections] = %v, want 10", result.Details["connections"])
	}
	if result.Details["version"] != "14.0" {
		t.Errorf("Details[version] = %v, want '14.0'", result.Details["version"])
	}
}

func TestNewReport_AllHealthy(t *testing.T) {
	t.Parallel()

	checks := map[string]Result{
		"postgres": NewHealthyResult(10 * time.Millisecond),
		"redis":    NewHealthyResult(5 * time.Millisecond),
	}
	required := map[string]bool{
		"postgres": true,
		"redis":    true,
	}

	report := NewReport(checks, required)

	if report.Status != StatusHealthy {
		t.Errorf("Status = %v, want %v", report.Status, StatusHealthy)
	}
	if len(report.Checks) != 2 {
		t.Errorf("Checks count = %d, want 2", len(report.Checks))
	}
}

func TestNewReport_RequiredUnhealthy(t *testing.T) {
	t.Parallel()

	checks := map[string]Result{
		"postgres": NewUnhealthyResult(10*time.Millisecond, errors.New("connection failed")),
		"redis":    NewHealthyResult(5 * time.Millisecond),
	}
	required := map[string]bool{
		"postgres": true,
		"redis":    true,
	}

	report := NewReport(checks, required)

	if report.Status != StatusUnhealthy {
		t.Errorf("Status = %v, want %v", report.Status, StatusUnhealthy)
	}
}

func TestNewReport_OptionalUnhealthy(t *testing.T) {
	t.Parallel()

	checks := map[string]Result{
		"postgres":     NewHealthyResult(10 * time.Millisecond),
		"external_api": NewUnhealthyResult(5*time.Millisecond, errors.New("timeout")),
	}
	required := map[string]bool{
		"postgres":     true,
		"external_api": false,
	}

	report := NewReport(checks, required)

	// Should be degraded, not unhealthy, because only optional check failed.
	if report.Status != StatusDegraded {
		t.Errorf("Status = %v, want %v", report.Status, StatusDegraded)
	}
}

func TestNewReport_DegradedCheck(t *testing.T) {
	t.Parallel()

	checks := map[string]Result{
		"postgres": NewHealthyResult(10 * time.Millisecond),
		"redis":    NewDegradedResult(5*time.Millisecond, "high latency"),
	}
	required := map[string]bool{
		"postgres": true,
		"redis":    true,
	}

	report := NewReport(checks, required)

	if report.Status != StatusDegraded {
		t.Errorf("Status = %v, want %v", report.Status, StatusDegraded)
	}
}

func TestNewReport_Timestamp(t *testing.T) {
	t.Parallel()

	checks := map[string]Result{}
	required := map[string]bool{}

	before := time.Now()
	report := NewReport(checks, required)
	after := time.Now()

	if report.Timestamp.Before(before) || report.Timestamp.After(after) {
		t.Error("Timestamp should be between before and after")
	}
}
