package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestCode_String(t *testing.T) {
	tests := []struct {
		name string
		code Code
		want string
	}{
		{"collector init failed", CodeCollectorInitFailed, "COLLECTOR_INIT_FAILED"},
		{"tracer init failed", CodeTracerInitFailed, "TRACER_INIT_FAILED"},
		{"exporter error", CodeExporterError, "EXPORTER_ERROR"},
		{"health check failed", CodeHealthCheckFailed, "HEALTH_CHECK_FAILED"},
		{"health check timeout", CodeHealthCheckTimeout, "HEALTH_CHECK_TIMEOUT"},
		{"invalid config", CodeInvalidConfig, "INVALID_CONFIG"},
		{"registration failed", CodeRegistrationFailed, "REGISTRATION_FAILED"},
		{"label cardinality exceeded", CodeLabelCardinalityExceeded, "LABEL_CARDINALITY_EXCEEDED"},
		{"span creation failed", CodeSpanCreationFailed, "SPAN_CREATION_FAILED"},
		{"context propagation failed", CodeContextPropagationFailed, "CONTEXT_PROPAGATION_FAILED"},
		{"shutdown failed", CodeShutdownFailed, "SHUTDOWN_FAILED"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.code.String(); got != tt.want {
				t.Errorf("Code.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	err := New(CodeInvalidConfig, "test message")

	if err.Code() != CodeInvalidConfig {
		t.Errorf("Code() = %v, want %v", err.Code(), CodeInvalidConfig)
	}
	if err.Message() != "test message" {
		t.Errorf("Message() = %v, want %v", err.Message(), "test message")
	}
	if err.Unwrap() != nil {
		t.Errorf("Unwrap() = %v, want nil", err.Unwrap())
	}
}

func TestNewf(t *testing.T) {
	err := Newf(CodeInvalidConfig, "test %s %d", "message", 42)

	if err.Code() != CodeInvalidConfig {
		t.Errorf("Code() = %v, want %v", err.Code(), CodeInvalidConfig)
	}
	if err.Message() != "test message 42" {
		t.Errorf("Message() = %v, want %v", err.Message(), "test message 42")
	}
}

func TestWrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := Wrap(CodeExporterError, "wrapper message", cause)

	if err.Code() != CodeExporterError {
		t.Errorf("Code() = %v, want %v", err.Code(), CodeExporterError)
	}
	if err.Message() != "wrapper message" {
		t.Errorf("Message() = %v, want %v", err.Message(), "wrapper message")
	}
	if err.Unwrap() != cause {
		t.Errorf("Unwrap() = %v, want %v", err.Unwrap(), cause)
	}
}

func TestWrapf(t *testing.T) {
	cause := errors.New("underlying error")
	err := Wrapf(CodeExporterError, cause, "wrapper %s", "message")

	if err.Code() != CodeExporterError {
		t.Errorf("Code() = %v, want %v", err.Code(), CodeExporterError)
	}
	if err.Message() != "wrapper message" {
		t.Errorf("Message() = %v, want %v", err.Message(), "wrapper message")
	}
	if err.Unwrap() != cause {
		t.Errorf("Unwrap() = %v, want %v", err.Unwrap(), cause)
	}
}

func TestError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		want string
	}{
		{
			name: "without cause",
			err:  New(CodeInvalidConfig, "test message"),
			want: "INVALID_CONFIG: test message",
		},
		{
			name: "with cause",
			err:  Wrap(CodeExporterError, "wrapper message", errors.New("cause")),
			want: "EXPORTER_ERROR: wrapper message: cause",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Is(t *testing.T) {
	err1 := New(CodeInvalidConfig, "message 1")
	err2 := New(CodeInvalidConfig, "message 2")
	err3 := New(CodeExporterError, "message 3")

	if !err1.Is(err2) {
		t.Error("Expected err1.Is(err2) to be true (same code)")
	}
	if err1.Is(err3) {
		t.Error("Expected err1.Is(err3) to be false (different code)")
	}
	if err1.Is(errors.New("standard error")) {
		t.Error("Expected err1.Is(standard error) to be false")
	}
}

func TestError_WithMessage(t *testing.T) {
	original := Wrap(CodeInvalidConfig, "original", errors.New("cause"))
	modified := original.WithMessage("modified")

	if modified.Code() != original.Code() {
		t.Errorf("Code changed: got %v, want %v", modified.Code(), original.Code())
	}
	if modified.Message() != "modified" {
		t.Errorf("Message() = %v, want %v", modified.Message(), "modified")
	}
	if modified.Unwrap() != original.Unwrap() {
		t.Errorf("Cause changed: got %v, want %v", modified.Unwrap(), original.Unwrap())
	}
}

func TestError_WithCause(t *testing.T) {
	original := New(CodeInvalidConfig, "message")
	newCause := errors.New("new cause")
	modified := original.WithCause(newCause)

	if modified.Code() != original.Code() {
		t.Errorf("Code changed: got %v, want %v", modified.Code(), original.Code())
	}
	if modified.Message() != original.Message() {
		t.Errorf("Message changed: got %v, want %v", modified.Message(), original.Message())
	}
	if modified.Unwrap() != newCause {
		t.Errorf("Unwrap() = %v, want %v", modified.Unwrap(), newCause)
	}
}

func TestCollectorInitFailed(t *testing.T) {
	err := CollectorInitFailed("prometheus")
	if err.Code() != CodeCollectorInitFailed {
		t.Errorf("Code() = %v, want %v", err.Code(), CodeCollectorInitFailed)
	}
	if err.Message() != "prometheus" {
		t.Errorf("Message() = %v, want %v", err.Message(), "prometheus")
	}
}

func TestCollectorInitFailedWrap(t *testing.T) {
	cause := errors.New("connection refused")
	err := CollectorInitFailedWrap("prometheus", cause)
	if err.Code() != CodeCollectorInitFailed {
		t.Errorf("Code() = %v, want %v", err.Code(), CodeCollectorInitFailed)
	}
	if err.Unwrap() != cause {
		t.Errorf("Unwrap() = %v, want %v", err.Unwrap(), cause)
	}
}

func TestTracerInitFailed(t *testing.T) {
	err := TracerInitFailed("otlp exporter failed")
	if err.Code() != CodeTracerInitFailed {
		t.Errorf("Code() = %v, want %v", err.Code(), CodeTracerInitFailed)
	}
}

func TestTracerInitFailedWrap(t *testing.T) {
	cause := errors.New("connection timeout")
	err := TracerInitFailedWrap("otlp exporter failed", cause)
	if err.Code() != CodeTracerInitFailed {
		t.Errorf("Code() = %v, want %v", err.Code(), CodeTracerInitFailed)
	}
	if err.Unwrap() != cause {
		t.Errorf("Unwrap() = %v, want %v", err.Unwrap(), cause)
	}
}

func TestExporterError(t *testing.T) {
	err := ExporterError("failed to export metrics")
	if err.Code() != CodeExporterError {
		t.Errorf("Code() = %v, want %v", err.Code(), CodeExporterError)
	}
}

func TestExporterErrorWrap(t *testing.T) {
	cause := errors.New("network error")
	err := ExporterErrorWrap("failed to export metrics", cause)
	if err.Code() != CodeExporterError {
		t.Errorf("Code() = %v, want %v", err.Code(), CodeExporterError)
	}
	if err.Unwrap() != cause {
		t.Errorf("Unwrap() = %v, want %v", err.Unwrap(), cause)
	}
}

func TestHealthCheckFailed(t *testing.T) {
	err := HealthCheckFailed("postgres")
	if err.Code() != CodeHealthCheckFailed {
		t.Errorf("Code() = %v, want %v", err.Code(), CodeHealthCheckFailed)
	}
	expectedMsg := "health check failed for component: postgres"
	if err.Message() != expectedMsg {
		t.Errorf("Message() = %v, want %v", err.Message(), expectedMsg)
	}
}

func TestHealthCheckFailedWrap(t *testing.T) {
	cause := errors.New("connection refused")
	err := HealthCheckFailedWrap("postgres", cause)
	if err.Code() != CodeHealthCheckFailed {
		t.Errorf("Code() = %v, want %v", err.Code(), CodeHealthCheckFailed)
	}
	if err.Unwrap() != cause {
		t.Errorf("Unwrap() = %v, want %v", err.Unwrap(), cause)
	}
}

func TestHealthCheckTimeout(t *testing.T) {
	err := HealthCheckTimeout("redis")
	if err.Code() != CodeHealthCheckTimeout {
		t.Errorf("Code() = %v, want %v", err.Code(), CodeHealthCheckTimeout)
	}
	expectedMsg := "health check timed out for component: redis"
	if err.Message() != expectedMsg {
		t.Errorf("Message() = %v, want %v", err.Message(), expectedMsg)
	}
}

func TestHealthCheckTimeoutWrap(t *testing.T) {
	cause := errors.New("context deadline exceeded")
	err := HealthCheckTimeoutWrap("redis", cause)
	if err.Code() != CodeHealthCheckTimeout {
		t.Errorf("Code() = %v, want %v", err.Code(), CodeHealthCheckTimeout)
	}
	if err.Unwrap() != cause {
		t.Errorf("Unwrap() = %v, want %v", err.Unwrap(), cause)
	}
}

func TestInvalidConfig(t *testing.T) {
	err := InvalidConfig("missing required field")
	if err.Code() != CodeInvalidConfig {
		t.Errorf("Code() = %v, want %v", err.Code(), CodeInvalidConfig)
	}
}

func TestInvalidConfigf(t *testing.T) {
	err := InvalidConfigf("field %q must be positive", "sample_rate")
	if err.Code() != CodeInvalidConfig {
		t.Errorf("Code() = %v, want %v", err.Code(), CodeInvalidConfig)
	}
	expectedMsg := `field "sample_rate" must be positive`
	if err.Message() != expectedMsg {
		t.Errorf("Message() = %v, want %v", err.Message(), expectedMsg)
	}
}

func TestRegistrationFailed(t *testing.T) {
	err := RegistrationFailed("http_requests_total")
	if err.Code() != CodeRegistrationFailed {
		t.Errorf("Code() = %v, want %v", err.Code(), CodeRegistrationFailed)
	}
	expectedMsg := "failed to register metric: http_requests_total"
	if err.Message() != expectedMsg {
		t.Errorf("Message() = %v, want %v", err.Message(), expectedMsg)
	}
}

func TestRegistrationFailedWrap(t *testing.T) {
	cause := errors.New("duplicate metric")
	err := RegistrationFailedWrap("http_requests_total", cause)
	if err.Code() != CodeRegistrationFailed {
		t.Errorf("Code() = %v, want %v", err.Code(), CodeRegistrationFailed)
	}
	if err.Unwrap() != cause {
		t.Errorf("Unwrap() = %v, want %v", err.Unwrap(), cause)
	}
}

func TestLabelCardinalityExceeded(t *testing.T) {
	err := LabelCardinalityExceeded("user_id", 150, 100)
	if err.Code() != CodeLabelCardinalityExceeded {
		t.Errorf("Code() = %v, want %v", err.Code(), CodeLabelCardinalityExceeded)
	}
	expectedMsg := `label "user_id" has 150 unique values, exceeds limit of 100`
	if err.Message() != expectedMsg {
		t.Errorf("Message() = %v, want %v", err.Message(), expectedMsg)
	}
}

func TestSpanCreationFailed(t *testing.T) {
	err := SpanCreationFailed("invalid span name")
	if err.Code() != CodeSpanCreationFailed {
		t.Errorf("Code() = %v, want %v", err.Code(), CodeSpanCreationFailed)
	}
}

func TestSpanCreationFailedWrap(t *testing.T) {
	cause := errors.New("tracer not initialized")
	err := SpanCreationFailedWrap("failed to start span", cause)
	if err.Code() != CodeSpanCreationFailed {
		t.Errorf("Code() = %v, want %v", err.Code(), CodeSpanCreationFailed)
	}
	if err.Unwrap() != cause {
		t.Errorf("Unwrap() = %v, want %v", err.Unwrap(), cause)
	}
}

func TestContextPropagationFailed(t *testing.T) {
	err := ContextPropagationFailed("invalid traceparent header")
	if err.Code() != CodeContextPropagationFailed {
		t.Errorf("Code() = %v, want %v", err.Code(), CodeContextPropagationFailed)
	}
}

func TestContextPropagationFailedWrap(t *testing.T) {
	cause := errors.New("malformed header")
	err := ContextPropagationFailedWrap("failed to extract trace context", cause)
	if err.Code() != CodeContextPropagationFailed {
		t.Errorf("Code() = %v, want %v", err.Code(), CodeContextPropagationFailed)
	}
	if err.Unwrap() != cause {
		t.Errorf("Unwrap() = %v, want %v", err.Unwrap(), cause)
	}
}

func TestShutdownFailed(t *testing.T) {
	err := ShutdownFailed("tracer")
	if err.Code() != CodeShutdownFailed {
		t.Errorf("Code() = %v, want %v", err.Code(), CodeShutdownFailed)
	}
	expectedMsg := "failed to shutdown component: tracer"
	if err.Message() != expectedMsg {
		t.Errorf("Message() = %v, want %v", err.Message(), expectedMsg)
	}
}

func TestShutdownFailedWrap(t *testing.T) {
	cause := errors.New("context canceled")
	err := ShutdownFailedWrap("tracer", cause)
	if err.Code() != CodeShutdownFailed {
		t.Errorf("Code() = %v, want %v", err.Code(), CodeShutdownFailed)
	}
	if err.Unwrap() != cause {
		t.Errorf("Unwrap() = %v, want %v", err.Unwrap(), cause)
	}
}

func TestIsError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"observability error", New(CodeInvalidConfig, "test"), true},
		{"wrapped observability error", fmt.Errorf("wrapper: %w", New(CodeInvalidConfig, "test")), true},
		{"standard error", errors.New("standard"), false},
		{"nil error", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsError(tt.err); got != tt.want {
				t.Errorf("IsError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAsError(t *testing.T) {
	obsErr := New(CodeInvalidConfig, "test")

	tests := []struct {
		name    string
		err     error
		wantNil bool
	}{
		{"observability error", obsErr, false},
		{"wrapped observability error", fmt.Errorf("wrapper: %w", obsErr), false},
		{"standard error", errors.New("standard"), true},
		{"nil error", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AsError(tt.err)
			if tt.wantNil && got != nil {
				t.Errorf("AsError() = %v, want nil", got)
			}
			if !tt.wantNil && got == nil {
				t.Error("AsError() = nil, want non-nil")
			}
		})
	}
}

func TestGetCode(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want Code
	}{
		{"observability error", New(CodeInvalidConfig, "test"), CodeInvalidConfig},
		{"wrapped observability error", fmt.Errorf("wrapper: %w", New(CodeExporterError, "test")), CodeExporterError},
		{"standard error", errors.New("standard"), ""},
		{"nil error", nil, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCode(tt.err); got != tt.want {
				t.Errorf("GetCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsCode(t *testing.T) {
	err := New(CodeInvalidConfig, "test")

	tests := []struct {
		name string
		err  error
		code Code
		want bool
	}{
		{"matching code", err, CodeInvalidConfig, true},
		{"non-matching code", err, CodeExporterError, false},
		{"wrapped matching", fmt.Errorf("w: %w", err), CodeInvalidConfig, true},
		{"standard error", errors.New("std"), CodeInvalidConfig, false},
		{"nil error", nil, CodeInvalidConfig, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCode(tt.err, tt.code); got != tt.want {
				t.Errorf("IsCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsCodeHelpers(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		checkFn  func(error) bool
		expected bool
	}{
		{"IsCollectorInitFailed true", CollectorInitFailed("test"), IsCollectorInitFailed, true},
		{"IsCollectorInitFailed false", TracerInitFailed("test"), IsCollectorInitFailed, false},
		{"IsTracerInitFailed true", TracerInitFailed("test"), IsTracerInitFailed, true},
		{"IsTracerInitFailed false", CollectorInitFailed("test"), IsTracerInitFailed, false},
		{"IsExporterError true", ExporterError("test"), IsExporterError, true},
		{"IsExporterError false", InvalidConfig("test"), IsExporterError, false},
		{"IsHealthCheckFailed true", HealthCheckFailed("db"), IsHealthCheckFailed, true},
		{"IsHealthCheckFailed false", HealthCheckTimeout("db"), IsHealthCheckFailed, false},
		{"IsHealthCheckTimeout true", HealthCheckTimeout("db"), IsHealthCheckTimeout, true},
		{"IsHealthCheckTimeout false", HealthCheckFailed("db"), IsHealthCheckTimeout, false},
		{"IsInvalidConfig true", InvalidConfig("test"), IsInvalidConfig, true},
		{"IsInvalidConfig false", ExporterError("test"), IsInvalidConfig, false},
		{"IsRegistrationFailed true", RegistrationFailed("metric"), IsRegistrationFailed, true},
		{"IsRegistrationFailed false", InvalidConfig("test"), IsRegistrationFailed, false},
		{"IsLabelCardinalityExceeded true", LabelCardinalityExceeded("label", 200, 100), IsLabelCardinalityExceeded, true},
		{"IsLabelCardinalityExceeded false", InvalidConfig("test"), IsLabelCardinalityExceeded, false},
		{"IsSpanCreationFailed true", SpanCreationFailed("test"), IsSpanCreationFailed, true},
		{"IsSpanCreationFailed false", InvalidConfig("test"), IsSpanCreationFailed, false},
		{"IsContextPropagationFailed true", ContextPropagationFailed("test"), IsContextPropagationFailed, true},
		{"IsContextPropagationFailed false", InvalidConfig("test"), IsContextPropagationFailed, false},
		{"IsShutdownFailed true", ShutdownFailed("tracer"), IsShutdownFailed, true},
		{"IsShutdownFailed false", InvalidConfig("test"), IsShutdownFailed, false},
		{"standard error", errors.New("test"), IsInvalidConfig, false},
		{"nil error", nil, IsInvalidConfig, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.checkFn(tt.err); got != tt.expected {
				t.Errorf("check function returned %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestErrorUnwrapChain(t *testing.T) {
	root := errors.New("root cause")
	middle := Wrap(CodeExporterError, "middle", root)
	outer := fmt.Errorf("outer: %w", middle)

	// Test that errors.Is works through the chain.
	if !errors.Is(outer, middle) {
		t.Error("errors.Is should find middle in chain")
	}

	// Test that we can extract the observability error.
	extracted := AsError(outer)
	if extracted == nil {
		t.Fatal("AsError returned nil for wrapped error")
	}
	if extracted.Code() != CodeExporterError {
		t.Errorf("Code() = %v, want %v", extracted.Code(), CodeExporterError)
	}

	// Test that we can still get to the root cause.
	if !errors.Is(outer, root) {
		t.Error("errors.Is should find root in chain")
	}
}
