package health

import (
	"testing"
)

func TestStatus_IsHealthy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		status   Status
		expected bool
	}{
		{"healthy", StatusHealthy, true},
		{"unhealthy", StatusUnhealthy, false},
		{"degraded", StatusDegraded, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.status.IsHealthy(); got != tt.expected {
				t.Errorf("IsHealthy() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestStatus_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		status   Status
		expected string
	}{
		{"healthy", StatusHealthy, "healthy"},
		{"unhealthy", StatusUnhealthy, "unhealthy"},
		{"degraded", StatusDegraded, "degraded"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.status.String(); got != tt.expected {
				t.Errorf("String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestStatusConstants(t *testing.T) {
	t.Parallel()

	if StatusHealthy != "healthy" {
		t.Errorf("StatusHealthy = %v, want healthy", StatusHealthy)
	}
	if StatusUnhealthy != "unhealthy" {
		t.Errorf("StatusUnhealthy = %v, want unhealthy", StatusUnhealthy)
	}
	if StatusDegraded != "degraded" {
		t.Errorf("StatusDegraded = %v, want degraded", StatusDegraded)
	}
}
