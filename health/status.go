package health

// Status represents the health status of a component.
type Status string

const (
	// StatusHealthy indicates the component is functioning normally.
	StatusHealthy Status = "healthy"

	// StatusUnhealthy indicates the component is not functioning properly.
	StatusUnhealthy Status = "unhealthy"

	// StatusDegraded indicates the component is functioning but with issues.
	StatusDegraded Status = "degraded"
)

// IsHealthy returns true if the status is healthy.
func (s Status) IsHealthy() bool {
	return s == StatusHealthy
}

// String returns the string representation of the status.
func (s Status) String() string {
	return string(s)
}
