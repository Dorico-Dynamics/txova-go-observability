package tracing

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddleware(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg := DefaultConfig().
		WithServiceName("test-service").
		WithExporter(ExporterNone)

	tracer, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer tracer.Shutdown(ctx)

	// Create a simple handler.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Wrap with middleware.
	wrapped := Middleware(tracer)(handler)

	// Create a test request.
	req := httptest.NewRequest("GET", "/api/v1/rides", nil)
	req.Header.Set("User-Agent", "TestAgent/1.0")
	req.Header.Set("X-Request-ID", "req-123")

	// Create a response recorder.
	rec := httptest.NewRecorder()

	// Serve the request.
	wrapped.ServeHTTP(rec, req)

	// Verify response.
	if rec.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", rec.Code, http.StatusOK)
	}
	if rec.Body.String() != "OK" {
		t.Errorf("Body = %s, want OK", rec.Body.String())
	}
}

func TestMiddleware_ErrorStatusCode(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg := DefaultConfig().
		WithServiceName("test-service").
		WithExporter(ExporterNone)

	tracer, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer tracer.Shutdown(ctx)

	// Create a handler that returns 500.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	})

	wrapped := Middleware(tracer)(handler)

	req := httptest.NewRequest("GET", "/api/v1/error", nil)
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("Status code = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

func TestMiddleware_ExtractsTraceContext(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg := DefaultConfig().
		WithServiceName("test-service").
		WithExporter(ExporterNone)

	tracer, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer tracer.Shutdown(ctx)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify span is in context.
		span := SpanFromContext(r.Context())
		if span == nil {
			t.Error("Span should be in request context")
		}
		w.WriteHeader(http.StatusOK)
	})

	wrapped := Middleware(tracer)(handler)

	req := httptest.NewRequest("GET", "/api/v1/rides", nil)
	req.Header.Set("traceparent", "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01")

	rec := httptest.NewRecorder()
	wrapped.ServeHTTP(rec, req)
}

func TestMiddleware_ClientIP_XForwardedFor(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg := DefaultConfig().
		WithServiceName("test-service").
		WithExporter(ExporterNone)

	tracer, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer tracer.Shutdown(ctx)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrapped := Middleware(tracer)(handler)

	req := httptest.NewRequest("GET", "/api/v1/rides", nil)
	req.Header.Set("X-Forwarded-For", "192.168.1.1")

	rec := httptest.NewRecorder()
	wrapped.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestMiddleware_ClientIP_XRealIP(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg := DefaultConfig().
		WithServiceName("test-service").
		WithExporter(ExporterNone)

	tracer, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer tracer.Shutdown(ctx)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrapped := Middleware(tracer)(handler)

	req := httptest.NewRequest("GET", "/api/v1/rides", nil)
	req.Header.Set("X-Real-IP", "10.0.0.1")

	rec := httptest.NewRecorder()
	wrapped.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestResponseWriter_WriteHeader(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	rw := &responseWriter{ResponseWriter: rec, statusCode: http.StatusOK}

	rw.WriteHeader(http.StatusCreated)

	if rw.statusCode != http.StatusCreated {
		t.Errorf("statusCode = %d, want %d", rw.statusCode, http.StatusCreated)
	}

	// Second call should not change status.
	rw.WriteHeader(http.StatusBadRequest)
	if rw.statusCode != http.StatusCreated {
		t.Errorf("statusCode after second WriteHeader = %d, want %d", rw.statusCode, http.StatusCreated)
	}
}

func TestResponseWriter_Write(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	rw := &responseWriter{ResponseWriter: rec, statusCode: http.StatusOK}

	n, err := rw.Write([]byte("test"))
	if err != nil {
		t.Errorf("Write() error = %v", err)
	}
	if n != 4 {
		t.Errorf("Write() n = %d, want 4", n)
	}
	if !rw.written {
		t.Error("written should be true after Write()")
	}
}

func TestResponseWriter_Unwrap(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	rw := &responseWriter{ResponseWriter: rec}

	if rw.Unwrap() != rec {
		t.Error("Unwrap() should return the wrapped ResponseWriter")
	}
}

func TestRoundTripper(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg := DefaultConfig().
		WithServiceName("test-service").
		WithExporter(ExporterNone)

	tracer, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer tracer.Shutdown(ctx)

	// Create a test server.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify trace context was propagated.
		if r.Header.Get("traceparent") == "" {
			t.Error("traceparent header should be set")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	// Create client with tracing round tripper.
	client := &http.Client{
		Transport: RoundTripper(tracer, nil),
	}

	// Start a span to have trace context.
	spanCtx, span := tracer.Start(ctx, "client-call")
	defer span.End()

	req, err := http.NewRequestWithContext(spanCtx, "GET", server.URL, nil)
	if err != nil {
		t.Fatalf("NewRequestWithContext() error = %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestRoundTripper_NilBase(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg := DefaultConfig().
		WithServiceName("test-service").
		WithExporter(ExporterNone)

	tracer, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer tracer.Shutdown(ctx)

	// Create with nil base, should use DefaultTransport.
	rt := RoundTripper(tracer, nil)

	if rt == nil {
		t.Error("RoundTripper() returned nil")
	}
}

func TestRoundTripper_ErrorResponse(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg := DefaultConfig().
		WithServiceName("test-service").
		WithExporter(ExporterNone)

	tracer, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer tracer.Shutdown(ctx)

	// Create a test server that returns 500.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := &http.Client{
		Transport: RoundTripper(tracer, nil),
	}

	spanCtx, span := tracer.Start(ctx, "client-call")
	defer span.End()

	req, err := http.NewRequestWithContext(spanCtx, "GET", server.URL, nil)
	if err != nil {
		t.Fatalf("NewRequestWithContext() error = %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Status code = %d, want %d", resp.StatusCode, http.StatusInternalServerError)
	}
}

func TestSpanFromContext(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg := DefaultConfig().
		WithServiceName("test-service").
		WithExporter(ExporterNone)

	tracer, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer tracer.Shutdown(ctx)

	spanCtx, span := tracer.Start(ctx, "test-span")
	defer span.End()

	retrievedSpan := SpanFromContext(spanCtx)

	if retrievedSpan != span {
		t.Error("SpanFromContext() should return the same span")
	}
}

func TestAddSpanAttributes(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg := DefaultConfig().
		WithServiceName("test-service").
		WithExporter(ExporterNone)

	tracer, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer tracer.Shutdown(ctx)

	spanCtx, span := tracer.Start(ctx, "test-span")
	defer span.End()

	// This should not panic.
	AddSpanAttributes(spanCtx, RideID("ride-123"), DriverID("driver-456"))
}

func TestRecordError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg := DefaultConfig().
		WithServiceName("test-service").
		WithExporter(ExporterNone)

	tracer, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer tracer.Shutdown(ctx)

	spanCtx, span := tracer.Start(ctx, "test-span")
	defer span.End()

	testErr := errors.New("test error")

	// This should not panic.
	RecordError(spanCtx, testErr)
}

func TestGetClientIP(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		headers  map[string]string
		remoteIP string
		expected string
	}{
		{
			name:     "X-Forwarded-For",
			headers:  map[string]string{"X-Forwarded-For": "192.168.1.1"},
			remoteIP: "10.0.0.1:1234",
			expected: "192.168.1.1",
		},
		{
			name:     "X-Real-IP",
			headers:  map[string]string{"X-Real-IP": "192.168.1.2"},
			remoteIP: "10.0.0.1:1234",
			expected: "192.168.1.2",
		},
		{
			name:     "RemoteAddr fallback",
			headers:  map[string]string{},
			remoteIP: "10.0.0.1:1234",
			expected: "10.0.0.1:1234",
		},
		{
			name:     "X-Forwarded-For takes precedence",
			headers:  map[string]string{"X-Forwarded-For": "192.168.1.1", "X-Real-IP": "192.168.1.2"},
			remoteIP: "10.0.0.1:1234",
			expected: "192.168.1.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest("GET", "/", nil)
			req.RemoteAddr = tt.remoteIP
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}

			got := getClientIP(req)
			if got != tt.expected {
				t.Errorf("getClientIP() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestMiddleware_WriteWithoutHeader(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg := DefaultConfig().
		WithServiceName("test-service").
		WithExporter(ExporterNone)

	tracer, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer tracer.Shutdown(ctx)

	// Handler that writes without explicitly calling WriteHeader.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	wrapped := Middleware(tracer)(handler)

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	// Should default to 200.
	if rec.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", rec.Code, http.StatusOK)
	}
	body, _ := io.ReadAll(rec.Body)
	if string(body) != "OK" {
		t.Errorf("Body = %s, want OK", string(body))
	}
}
