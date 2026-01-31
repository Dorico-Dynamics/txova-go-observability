package health

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

// ManagerConfig holds configuration for the health check manager.
type ManagerConfig struct {
	// Timeout is the default timeout for health checks.
	Timeout time.Duration

	// CacheTTL is how long to cache health check results.
	CacheTTL time.Duration

	// BackgroundInterval is how often to run background checks.
	BackgroundInterval time.Duration

	// FailureThreshold is how many consecutive failures before marking unhealthy.
	FailureThreshold int

	// Logger is the logger to use for health check events.
	Logger *slog.Logger
}

// DefaultManagerConfig returns the default manager configuration.
func DefaultManagerConfig() ManagerConfig {
	return ManagerConfig{
		Timeout:            5 * time.Second,
		CacheTTL:           30 * time.Second,
		BackgroundInterval: 30 * time.Second,
		FailureThreshold:   3,
		Logger:             slog.Default(),
	}
}

// WithTimeout sets the check timeout.
func (c ManagerConfig) WithTimeout(timeout time.Duration) ManagerConfig {
	c.Timeout = timeout
	return c
}

// WithCacheTTL sets the cache TTL.
func (c ManagerConfig) WithCacheTTL(ttl time.Duration) ManagerConfig {
	c.CacheTTL = ttl
	return c
}

// WithBackgroundInterval sets the background check interval.
func (c ManagerConfig) WithBackgroundInterval(interval time.Duration) ManagerConfig {
	c.BackgroundInterval = interval
	return c
}

// WithFailureThreshold sets the failure threshold.
func (c ManagerConfig) WithFailureThreshold(threshold int) ManagerConfig {
	c.FailureThreshold = threshold
	return c
}

// WithLogger sets the logger.
func (c ManagerConfig) WithLogger(logger *slog.Logger) ManagerConfig {
	c.Logger = logger
	return c
}

// Manager manages health checks for multiple components.
type Manager struct {
	config         ManagerConfig
	checkers       []Checker
	requiredChecks map[string]bool

	mu            sync.RWMutex
	cachedReport  *Report
	cacheTime     time.Time
	failureCounts map[string]int

	stopCh chan struct{}
	doneCh chan struct{}
}

// NewManager creates a new health check manager.
func NewManager(config ManagerConfig) *Manager {
	if config.Logger == nil {
		config.Logger = slog.Default()
	}
	return &Manager{
		config:         config,
		checkers:       make([]Checker, 0),
		requiredChecks: make(map[string]bool),
		failureCounts:  make(map[string]int),
		stopCh:         make(chan struct{}),
		doneCh:         make(chan struct{}),
	}
}

// Register registers a health checker.
func (m *Manager) Register(checker Checker) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.checkers = append(m.checkers, checker)
	m.requiredChecks[checker.Name()] = checker.Required()
}

// Check performs all health checks and returns a report.
func (m *Manager) Check(ctx context.Context) Report {
	// Check if we have a valid cached report.
	m.mu.RLock()
	if m.cachedReport != nil && time.Since(m.cacheTime) < m.config.CacheTTL {
		report := *m.cachedReport
		m.mu.RUnlock()
		return report
	}
	m.mu.RUnlock()

	// Perform checks.
	return m.runChecks(ctx)
}

// runChecks performs all health checks.
func (m *Manager) runChecks(ctx context.Context) Report {
	m.mu.RLock()
	checkers := make([]Checker, len(m.checkers))
	copy(checkers, m.checkers)
	requiredChecks := make(map[string]bool)
	for k, v := range m.requiredChecks {
		requiredChecks[k] = v
	}
	m.mu.RUnlock()

	results := make(map[string]Result)
	var wg sync.WaitGroup

	resultCh := make(chan struct {
		name   string
		result Result
	}, len(checkers))

	for _, checker := range checkers {
		wg.Add(1)
		go func(c Checker) {
			defer wg.Done()

			// Create a timeout context for this check.
			checkCtx, cancel := context.WithTimeout(ctx, m.config.Timeout)
			defer cancel()

			result := c.Check(checkCtx)

			// Log failures.
			if result.Status != StatusHealthy {
				m.config.Logger.Warn("health check failed",
					"component", c.Name(),
					"status", result.Status,
					"error", result.Error,
					"duration_ms", result.DurationMS,
				)
			}

			resultCh <- struct {
				name   string
				result Result
			}{
				name:   c.Name(),
				result: result,
			}
		}(checker)
	}

	// Wait for all checks to complete.
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Collect results.
	for r := range resultCh {
		results[r.name] = r.result

		// Update failure counts.
		m.mu.Lock()
		if r.result.Status != StatusHealthy {
			m.failureCounts[r.name]++
		} else {
			m.failureCounts[r.name] = 0
		}
		m.mu.Unlock()
	}

	report := NewReport(results, requiredChecks)

	// Cache the report.
	m.mu.Lock()
	m.cachedReport = &report
	m.cacheTime = time.Now()
	m.mu.Unlock()

	return report
}

// StartBackground starts background health checks.
func (m *Manager) StartBackground(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(m.config.BackgroundInterval)
		defer ticker.Stop()
		defer close(m.doneCh)

		// Run initial check.
		m.runChecks(ctx)

		for {
			select {
			case <-ticker.C:
				m.runChecks(ctx)
			case <-m.stopCh:
				return
			case <-ctx.Done():
				return
			}
		}
	}()
}

// StopBackground stops background health checks.
func (m *Manager) StopBackground() {
	close(m.stopCh)
	<-m.doneCh
}

// GetFailureCount returns the consecutive failure count for a component.
func (m *Manager) GetFailureCount(name string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.failureCounts[name]
}

// IsReady returns true if the service is ready to accept traffic.
func (m *Manager) IsReady(ctx context.Context) bool {
	report := m.Check(ctx)
	return report.Status == StatusHealthy || report.Status == StatusDegraded
}

// IsLive returns true if the service is alive.
func (m *Manager) IsLive() bool {
	// For liveness, we just check if the manager is running.
	// The actual health of dependencies doesn't matter for liveness.
	return true
}

// IsStarted returns true if the service has completed startup.
func (m *Manager) IsStarted(ctx context.Context) bool {
	// Check that all required checks have passed at least once.
	report := m.Check(ctx)
	return report.Status == StatusHealthy
}
