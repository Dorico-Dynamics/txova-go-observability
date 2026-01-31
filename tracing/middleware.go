package tracing

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// sanitizeURL returns a sanitized URL string with query parameters and fragment removed
// to prevent PII leakage in traces.
func sanitizeURL(u *url.URL) string {
	sanitized := *u
	sanitized.RawQuery = ""
	sanitized.Fragment = ""
	return sanitized.String()
}

// Middleware returns an HTTP middleware that creates spans for incoming requests.
func Middleware(tracer *Tracer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract trace context from incoming headers.
			ctx := Extract(r.Context(), r.Header)

			// Get the route pattern if available, otherwise use the URL path.
			route := r.URL.Path

			// Start a new span for this request.
			spanName := fmt.Sprintf("%s %s", r.Method, route)
			ctx, span := tracer.Start(ctx, spanName,
				trace.WithSpanKind(trace.SpanKindServer),
			)
			defer span.End()

			// Add standard HTTP attributes.
			// Use sanitized URL to prevent PII leakage from query parameters.
			span.SetAttributes(
				HTTPMethod(r.Method),
				HTTPRoute(route),
				HTTPURL(sanitizeURL(r.URL)),
				HTTPScheme(r.URL.Scheme),
				HTTPHost(r.Host),
			)

			// Add optional attributes.
			if userAgent := r.Header.Get("User-Agent"); userAgent != "" {
				span.SetAttributes(HTTPUserAgent(userAgent))
			}

			if clientIP := getClientIP(r); clientIP != "" {
				span.SetAttributes(HTTPClientIP(clientIP))
			}

			// Extract and propagate request ID.
			if requestID := ExtractRequestID(r.Header); requestID != "" {
				span.SetAttributes(RequestID(requestID))
			}

			// Create a response writer wrapper to capture the status code.
			rw := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// Serve the request with the updated context.
			next.ServeHTTP(rw, r.WithContext(ctx))

			// Record the status code.
			span.SetAttributes(HTTPStatusCode(rw.statusCode))

			// Mark the span as error if status code indicates an error.
			if rw.statusCode >= 400 {
				span.SetStatus(codes.Error, http.StatusText(rw.statusCode))
			} else {
				span.SetStatus(codes.Ok, "")
			}
		})
	}
}

// responseWriter wraps http.ResponseWriter to capture the status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

// WriteHeader captures the status code and writes it.
func (rw *responseWriter) WriteHeader(statusCode int) {
	if !rw.written {
		rw.statusCode = statusCode
		rw.written = true
	}
	rw.ResponseWriter.WriteHeader(statusCode)
}

// Write writes the response body.
func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.statusCode = http.StatusOK
		rw.written = true
	}
	return rw.ResponseWriter.Write(b)
}

// Unwrap returns the wrapped ResponseWriter for http.ResponseController.
func (rw *responseWriter) Unwrap() http.ResponseWriter {
	return rw.ResponseWriter
}

// getClientIP extracts the client IP from the request.
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first.
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}

	// Check X-Real-IP header.
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr.
	return r.RemoteAddr
}

// RoundTripper returns an http.RoundTripper that propagates trace context.
func RoundTripper(tracer *Tracer, base http.RoundTripper) http.RoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}
	return &tracingRoundTripper{
		tracer: tracer,
		base:   base,
	}
}

// tracingRoundTripper is an http.RoundTripper that adds tracing.
type tracingRoundTripper struct {
	tracer *Tracer
	base   http.RoundTripper
}

// RoundTrip implements http.RoundTripper.
func (rt *tracingRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	ctx := r.Context()

	// Start a new span for the outgoing request.
	spanName := fmt.Sprintf("HTTP %s %s", r.Method, r.URL.Host)
	ctx, span := rt.tracer.Start(ctx, spanName,
		trace.WithSpanKind(trace.SpanKindClient),
	)
	defer span.End()

	// Add HTTP attributes.
	// Use sanitized URL to prevent PII leakage from query parameters.
	span.SetAttributes(
		HTTPMethod(r.Method),
		HTTPURL(sanitizeURL(r.URL)),
		HTTPHost(r.Host),
	)

	// Clone the request to avoid modifying the original.
	req := r.Clone(ctx)

	// Inject trace context into outgoing headers.
	Inject(ctx, req.Header)

	// Make the request.
	resp, err := rt.base.RoundTrip(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Record the response status.
	span.SetAttributes(HTTPStatusCode(resp.StatusCode))

	if resp.StatusCode >= 400 {
		span.SetStatus(codes.Error, http.StatusText(resp.StatusCode))
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return resp, nil
}

// SpanFromContext returns the current span from context.
func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// AddSpanAttributes adds attributes to the current span in context.
func AddSpanAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attrs...)
}

// RecordError records an error on the current span in context.
func RecordError(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}
