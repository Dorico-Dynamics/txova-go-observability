package health

import (
	"encoding/json"
	"net/http"
)

// Handler provides HTTP handlers for health check endpoints.
type Handler struct {
	manager *Manager
}

// NewHandler creates a new health check HTTP handler.
func NewHandler(manager *Manager) *Handler {
	return &Handler{
		manager: manager,
	}
}

// LiveHandler handles the /health/live endpoint (liveness probe).
// Returns 200 if the service is alive, regardless of dependency health.
func (h *Handler) LiveHandler(w http.ResponseWriter, r *http.Request) {
	if h.manager.IsLive() {
		h.writeJSON(w, http.StatusOK, map[string]any{
			"status": StatusHealthy,
		})
	} else {
		h.writeJSON(w, http.StatusServiceUnavailable, map[string]any{
			"status": StatusUnhealthy,
		})
	}
}

// ReadyHandler handles the /health/ready endpoint (readiness probe).
// Returns 200 if the service can accept traffic, 503 if not.
func (h *Handler) ReadyHandler(w http.ResponseWriter, r *http.Request) {
	report := h.manager.Check(r.Context())

	statusCode := http.StatusOK
	if report.Status == StatusUnhealthy {
		statusCode = http.StatusServiceUnavailable
	}

	h.writeJSON(w, statusCode, report)
}

// StartupHandler handles the /health/startup endpoint (startup probe).
// Returns 200 if the service has completed startup, 503 if not.
func (h *Handler) StartupHandler(w http.ResponseWriter, r *http.Request) {
	if h.manager.IsStarted(r.Context()) {
		report := h.manager.Check(r.Context())
		h.writeJSON(w, http.StatusOK, report)
	} else {
		h.writeJSON(w, http.StatusServiceUnavailable, map[string]any{
			"status": StatusUnhealthy,
			"error":  "service not yet started",
		})
	}
}

// FullHandler handles a full health check endpoint with all details.
func (h *Handler) FullHandler(w http.ResponseWriter, r *http.Request) {
	report := h.manager.Check(r.Context())

	statusCode := http.StatusOK
	if report.Status == StatusUnhealthy {
		statusCode = http.StatusServiceUnavailable
	}

	h.writeJSON(w, statusCode, report)
}

// writeJSON writes a JSON response.
func (h *Handler) writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data) //nolint:errcheck // best effort write to response
}

// RegisterRoutes registers health check routes on an http.ServeMux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health/live", h.LiveHandler)
	mux.HandleFunc("GET /health/ready", h.ReadyHandler)
	mux.HandleFunc("GET /health/startup", h.StartupHandler)
	mux.HandleFunc("GET /health", h.FullHandler)
}

// Routes returns a map of routes to handlers for custom routers.
func (h *Handler) Routes() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"/health/live":    h.LiveHandler,
		"/health/ready":   h.ReadyHandler,
		"/health/startup": h.StartupHandler,
		"/health":         h.FullHandler,
	}
}
