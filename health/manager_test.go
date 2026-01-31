package health

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDefaultManagerConfig(t *testing.T) {
	t.Parallel()

	cfg := DefaultManagerConfig()

	if cfg.Timeout != 5*time.Second {
		t.Errorf("Timeout = %v, want 5s", cfg.Timeout)
	}
	if cfg.CacheTTL != 30*time.Second {
		t.Errorf("CacheTTL = %v, want 30s", cfg.CacheTTL)
	}
	if cfg.BackgroundInterval != 30*time.Second {
		t.Errorf("BackgroundInterval = %v, want 30s", cfg.BackgroundInterval)
	}
	if cfg.FailureThreshold != 3 {
		t.Errorf("FailureThreshold = %v, want 3", cfg.FailureThreshold)
	}
}

func TestManagerConfig_Chaining(t *testing.T) {
	t.Parallel()

	cfg := DefaultManagerConfig().
		WithTimeout(10 * time.Second).
		WithCacheTTL(60 * time.Second).
		WithBackgroundInterval(15 * time.Second).
		WithFailureThreshold(5)

	if cfg.Timeout != 10*time.Second {
		t.Errorf("Timeout = %v, want 10s", cfg.Timeout)
	}
	if cfg.CacheTTL != 60*time.Second {
		t.Errorf("CacheTTL = %v, want 60s", cfg.CacheTTL)
	}
	if cfg.BackgroundInterval != 15*time.Second {
		t.Errorf("BackgroundInterval = %v, want 15s", cfg.BackgroundInterval)
	}
	if cfg.FailureThreshold != 5 {
		t.Errorf("FailureThreshold = %v, want 5", cfg.FailureThreshold)
	}
}

func TestNewManager(t *testing.T) {
	t.Parallel()

	cfg := DefaultManagerConfig()
	manager := NewManager(cfg)

	if manager == nil {
		t.Fatal("NewManager() returned nil")
	}
}

func TestManager_Register(t *testing.T) {
	t.Parallel()

	manager := NewManager(DefaultManagerConfig())
	checker := NewFuncChecker("test", func(ctx context.Context) error {
		return nil
	}, true)

	manager.Register(checker)

	// Verify by running a check.
	report := manager.Check(context.Background())
	if _, ok := report.Checks["test"]; !ok {
		t.Error("Registered checker not found in report")
	}
}

func TestManager_Check_AllHealthy(t *testing.T) {
	t.Parallel()

	manager := NewManager(DefaultManagerConfig())
	manager.Register(NewFuncChecker("check1", func(ctx context.Context) error {
		return nil
	}, true))
	manager.Register(NewFuncChecker("check2", func(ctx context.Context) error {
		return nil
	}, true))

	report := manager.Check(context.Background())

	if report.Status != StatusHealthy {
		t.Errorf("Status = %v, want %v", report.Status, StatusHealthy)
	}
	if len(report.Checks) != 2 {
		t.Errorf("Checks count = %d, want 2", len(report.Checks))
	}
}

func TestManager_Check_OneUnhealthy(t *testing.T) {
	t.Parallel()

	// Use FailureThreshold(1) so first failure is immediately reported
	manager := NewManager(DefaultManagerConfig().WithFailureThreshold(1))
	manager.Register(NewFuncChecker("healthy", func(ctx context.Context) error {
		return nil
	}, true))
	manager.Register(NewFuncChecker("unhealthy", func(ctx context.Context) error {
		return errors.New("check failed")
	}, true))

	report := manager.Check(context.Background())

	if report.Status != StatusUnhealthy {
		t.Errorf("Status = %v, want %v", report.Status, StatusUnhealthy)
	}
}

func TestManager_Check_Caching(t *testing.T) {
	t.Parallel()

	callCount := 0
	manager := NewManager(DefaultManagerConfig().WithCacheTTL(1 * time.Second))
	manager.Register(NewFuncChecker("test", func(ctx context.Context) error {
		callCount++
		return nil
	}, true))

	// First call.
	manager.Check(context.Background())
	if callCount != 1 {
		t.Errorf("Call count after first check = %d, want 1", callCount)
	}

	// Second call should use cache.
	manager.Check(context.Background())
	if callCount != 1 {
		t.Errorf("Call count after cached check = %d, want 1", callCount)
	}

	// Wait for cache to expire.
	time.Sleep(1100 * time.Millisecond)

	// Third call should not use cache.
	manager.Check(context.Background())
	if callCount != 2 {
		t.Errorf("Call count after cache expiry = %d, want 2", callCount)
	}
}

func TestManager_IsReady(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		checkErr error
		expected bool
	}{
		{"healthy", nil, true},
		{"unhealthy", errors.New("failed"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Use FailureThreshold(1) so first failure is immediately reported
			manager := NewManager(DefaultManagerConfig().WithCacheTTL(0).WithFailureThreshold(1))
			manager.Register(NewFuncChecker("test", func(ctx context.Context) error {
				return tt.checkErr
			}, true))

			if got := manager.IsReady(context.Background()); got != tt.expected {
				t.Errorf("IsReady() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestManager_IsLive(t *testing.T) {
	t.Parallel()

	manager := NewManager(DefaultManagerConfig())

	// IsLive should always return true.
	if !manager.IsLive() {
		t.Error("IsLive() should return true")
	}
}

func TestManager_IsStarted(t *testing.T) {
	t.Parallel()

	manager := NewManager(DefaultManagerConfig().WithCacheTTL(0))
	manager.Register(NewFuncChecker("test", func(ctx context.Context) error {
		return nil
	}, true))

	if !manager.IsStarted(context.Background()) {
		t.Error("IsStarted() should return true when all checks pass")
	}
}

func TestManager_GetFailureCount(t *testing.T) {
	t.Parallel()

	checkErr := errors.New("check failed")
	manager := NewManager(DefaultManagerConfig().WithCacheTTL(0))
	manager.Register(NewFuncChecker("failing", func(ctx context.Context) error {
		return checkErr
	}, true))

	// First check.
	manager.Check(context.Background())
	if count := manager.GetFailureCount("failing"); count != 1 {
		t.Errorf("Failure count = %d, want 1", count)
	}

	// Second check.
	manager.Check(context.Background())
	if count := manager.GetFailureCount("failing"); count != 2 {
		t.Errorf("Failure count = %d, want 2", count)
	}
}

func TestManager_GetFailureCount_Resets(t *testing.T) {
	t.Parallel()

	shouldFail := true
	manager := NewManager(DefaultManagerConfig().WithCacheTTL(0))
	manager.Register(NewFuncChecker("test", func(ctx context.Context) error {
		if shouldFail {
			return errors.New("failed")
		}
		return nil
	}, true))

	// Fail once.
	manager.Check(context.Background())
	if count := manager.GetFailureCount("test"); count != 1 {
		t.Errorf("Failure count = %d, want 1", count)
	}

	// Pass.
	shouldFail = false
	manager.Check(context.Background())
	if count := manager.GetFailureCount("test"); count != 0 {
		t.Errorf("Failure count after success = %d, want 0", count)
	}
}

func TestManager_BackgroundChecks(t *testing.T) {
	t.Parallel()

	callCount := 0
	manager := NewManager(DefaultManagerConfig().
		WithBackgroundInterval(100 * time.Millisecond).
		WithCacheTTL(50 * time.Millisecond))
	manager.Register(NewFuncChecker("test", func(ctx context.Context) error {
		callCount++
		return nil
	}, true))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	manager.StartBackground(ctx)

	// Wait for a few intervals.
	time.Sleep(350 * time.Millisecond)

	manager.StopBackground()

	// Should have been called at least 3 times (initial + 2-3 intervals).
	if callCount < 3 {
		t.Errorf("Call count = %d, want at least 3", callCount)
	}
}
