package health

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewHandler(t *testing.T) {
	t.Parallel()

	manager := NewManager(DefaultManagerConfig())
	handler := NewHandler(manager)

	if handler == nil {
		t.Fatal("NewHandler() returned nil")
	}
}

func TestHandler_LiveHandler_Healthy(t *testing.T) {
	t.Parallel()

	manager := NewManager(DefaultManagerConfig())
	handler := NewHandler(manager)

	req := httptest.NewRequest(http.MethodGet, "/health/live", nil)
	rec := httptest.NewRecorder()

	handler.LiveHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", rec.Code, http.StatusOK)
	}

	var response map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Status = %v, want healthy", response["status"])
	}
}

func TestHandler_ReadyHandler_Healthy(t *testing.T) {
	t.Parallel()

	manager := NewManager(DefaultManagerConfig().WithCacheTTL(0))
	manager.Register(NewFuncChecker("test", func(ctx context.Context) error {
		return nil
	}, true))
	handler := NewHandler(manager)

	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	rec := httptest.NewRecorder()

	handler.ReadyHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", rec.Code, http.StatusOK)
	}

	var response Report
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Status != StatusHealthy {
		t.Errorf("Status = %v, want healthy", response.Status)
	}
}

func TestHandler_ReadyHandler_Unhealthy(t *testing.T) {
	t.Parallel()

	manager := NewManager(DefaultManagerConfig().WithCacheTTL(0))
	manager.Register(NewFuncChecker("failing", func(ctx context.Context) error {
		return errors.New("check failed")
	}, true))
	handler := NewHandler(manager)

	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	rec := httptest.NewRecorder()

	handler.ReadyHandler(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("Status code = %d, want %d", rec.Code, http.StatusServiceUnavailable)
	}

	var response Report
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Status != StatusUnhealthy {
		t.Errorf("Status = %v, want unhealthy", response.Status)
	}
}

func TestHandler_StartupHandler_Started(t *testing.T) {
	t.Parallel()

	manager := NewManager(DefaultManagerConfig().WithCacheTTL(0))
	manager.Register(NewFuncChecker("test", func(ctx context.Context) error {
		return nil
	}, true))
	handler := NewHandler(manager)

	req := httptest.NewRequest(http.MethodGet, "/health/startup", nil)
	rec := httptest.NewRecorder()

	handler.StartupHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestHandler_StartupHandler_NotStarted(t *testing.T) {
	t.Parallel()

	manager := NewManager(DefaultManagerConfig().WithCacheTTL(0))
	manager.Register(NewFuncChecker("failing", func(ctx context.Context) error {
		return errors.New("not ready")
	}, true))
	handler := NewHandler(manager)

	req := httptest.NewRequest(http.MethodGet, "/health/startup", nil)
	rec := httptest.NewRecorder()

	handler.StartupHandler(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("Status code = %d, want %d", rec.Code, http.StatusServiceUnavailable)
	}
}

func TestHandler_FullHandler(t *testing.T) {
	t.Parallel()

	manager := NewManager(DefaultManagerConfig().WithCacheTTL(0))
	manager.Register(NewFuncChecker("test1", func(ctx context.Context) error {
		return nil
	}, true))
	manager.Register(NewFuncChecker("test2", func(ctx context.Context) error {
		return nil
	}, false))
	handler := NewHandler(manager)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	handler.FullHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", rec.Code, http.StatusOK)
	}

	var response Report
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(response.Checks) != 2 {
		t.Errorf("Checks count = %d, want 2", len(response.Checks))
	}
}

func TestHandler_Routes(t *testing.T) {
	t.Parallel()

	manager := NewManager(DefaultManagerConfig())
	handler := NewHandler(manager)

	routes := handler.Routes()

	expectedRoutes := []string{
		"/health/live",
		"/health/ready",
		"/health/startup",
		"/health",
	}

	for _, route := range expectedRoutes {
		if _, ok := routes[route]; !ok {
			t.Errorf("Route %s not found", route)
		}
	}
}

func TestHandler_ContentType(t *testing.T) {
	t.Parallel()

	manager := NewManager(DefaultManagerConfig())
	handler := NewHandler(manager)

	req := httptest.NewRequest(http.MethodGet, "/health/live", nil)
	rec := httptest.NewRecorder()

	handler.LiveHandler(rec, req)

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Content-Type = %v, want application/json", contentType)
	}
}

func TestHandler_ReadyHandler_Degraded(t *testing.T) {
	t.Parallel()

	manager := NewManager(DefaultManagerConfig().WithCacheTTL(0))
	manager.Register(NewFuncChecker("required", func(ctx context.Context) error {
		return nil
	}, true))
	manager.Register(NewFuncChecker("optional", func(ctx context.Context) error {
		return errors.New("optional failed")
	}, false))
	handler := NewHandler(manager)

	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	rec := httptest.NewRecorder()

	handler.ReadyHandler(rec, req)

	// Degraded should still return 200.
	if rec.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", rec.Code, http.StatusOK)
	}

	var response Report
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Status != StatusDegraded {
		t.Errorf("Status = %v, want degraded", response.Status)
	}
}
