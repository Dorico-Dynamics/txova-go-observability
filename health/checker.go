package health

import (
	"context"
	"time"
)

// Checker defines the interface for health check components.
type Checker interface {
	// Name returns the name of the component being checked.
	Name() string

	// Check performs the health check and returns the result.
	Check(ctx context.Context) Result

	// Required returns true if this check must pass for the service to be healthy.
	Required() bool
}

// Result represents the result of a health check.
type Result struct {
	// Status is the health status of the component.
	Status Status `json:"status"`

	// DurationMS is the time taken to perform the check in milliseconds.
	DurationMS int64 `json:"duration_ms"`

	// Error contains the error message if the check failed.
	Error string `json:"error,omitempty"`

	// Details contains additional information about the check.
	Details map[string]any `json:"details,omitempty"`

	// Timestamp is when the check was performed.
	Timestamp time.Time `json:"timestamp"`
}

// NewHealthyResult creates a healthy result with the given duration.
func NewHealthyResult(duration time.Duration) Result {
	return Result{
		Status:     StatusHealthy,
		DurationMS: duration.Milliseconds(),
		Timestamp:  time.Now(),
	}
}

// NewUnhealthyResult creates an unhealthy result with the given error.
func NewUnhealthyResult(duration time.Duration, err error) Result {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	return Result{
		Status:     StatusUnhealthy,
		DurationMS: duration.Milliseconds(),
		Error:      errMsg,
		Timestamp:  time.Now(),
	}
}

// NewDegradedResult creates a degraded result with the given message.
func NewDegradedResult(duration time.Duration, message string) Result {
	return Result{
		Status:     StatusDegraded,
		DurationMS: duration.Milliseconds(),
		Error:      message,
		Timestamp:  time.Now(),
	}
}

// WithDetails adds details to the result.
func (r Result) WithDetails(details map[string]any) Result {
	r.Details = details
	return r
}

// Report represents the overall health report.
type Report struct {
	// Status is the overall health status.
	Status Status `json:"status"`

	// Checks contains the results of individual health checks.
	Checks map[string]Result `json:"checks"`

	// Timestamp is when the report was generated.
	Timestamp time.Time `json:"timestamp"`
}

// NewReport creates a new health report from check results.
func NewReport(checks map[string]Result, requiredChecks map[string]bool) Report {
	status := StatusHealthy

	for name, result := range checks {
		if result.Status == StatusUnhealthy {
			// If a required check is unhealthy, the whole service is unhealthy.
			if requiredChecks[name] {
				status = StatusUnhealthy
				break
			}
			// If an optional check is unhealthy, mark as degraded at most.
			if status == StatusHealthy {
				status = StatusDegraded
			}
		} else if result.Status == StatusDegraded && status == StatusHealthy {
			status = StatusDegraded
		}
	}

	return Report{
		Status:    status,
		Checks:    checks,
		Timestamp: time.Now(),
	}
}
