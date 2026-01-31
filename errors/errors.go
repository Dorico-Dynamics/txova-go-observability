// Package errors provides observability-specific error types for the Txova platform.
// It defines error codes for metrics collection, tracing, and health check operations.
package errors

import (
	"errors"
	"fmt"
)

// Code represents an observability error code.
type Code string

// Observability error codes.
const (
	// CodeCollectorInitFailed indicates a metrics collector failed to initialize.
	CodeCollectorInitFailed Code = "COLLECTOR_INIT_FAILED"
	// CodeTracerInitFailed indicates the tracer failed to initialize.
	CodeTracerInitFailed Code = "TRACER_INIT_FAILED"
	// CodeExporterError indicates an error exporting metrics or traces.
	CodeExporterError Code = "EXPORTER_ERROR"
	// CodeHealthCheckFailed indicates a health check failed.
	CodeHealthCheckFailed Code = "HEALTH_CHECK_FAILED"
	// CodeHealthCheckTimeout indicates a health check timed out.
	CodeHealthCheckTimeout Code = "HEALTH_CHECK_TIMEOUT"
	// CodeInvalidConfig indicates invalid configuration.
	CodeInvalidConfig Code = "INVALID_CONFIG"
	// CodeRegistrationFailed indicates metric registration failed.
	CodeRegistrationFailed Code = "REGISTRATION_FAILED"
	// CodeLabelCardinalityExceeded indicates too many unique label values.
	CodeLabelCardinalityExceeded Code = "LABEL_CARDINALITY_EXCEEDED"
	// CodeSpanCreationFailed indicates span creation failed.
	CodeSpanCreationFailed Code = "SPAN_CREATION_FAILED"
	// CodeContextPropagationFailed indicates context propagation failed.
	CodeContextPropagationFailed Code = "CONTEXT_PROPAGATION_FAILED"
	// CodeShutdownFailed indicates graceful shutdown failed.
	CodeShutdownFailed Code = "SHUTDOWN_FAILED"
)

// String returns the string representation of the error code.
func (c Code) String() string {
	return string(c)
}

// Error represents an observability error with a code, message, and optional cause.
type Error struct {
	code    Code
	message string
	cause   error
}

// New creates a new Error with the given code and message.
func New(code Code, message string) *Error {
	return &Error{
		code:    code,
		message: message,
	}
}

// Newf creates a new Error with the given code and formatted message.
func Newf(code Code, format string, args ...any) *Error {
	return &Error{
		code:    code,
		message: fmt.Sprintf(format, args...),
	}
}

// Wrap creates a new Error that wraps an existing error.
func Wrap(code Code, message string, cause error) *Error {
	return &Error{
		code:    code,
		message: message,
		cause:   cause,
	}
}

// Wrapf creates a new Error that wraps an existing error with a formatted message.
func Wrapf(code Code, cause error, format string, args ...any) *Error {
	return &Error{
		code:    code,
		message: fmt.Sprintf(format, args...),
		cause:   cause,
	}
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.code, e.message, e.cause)
	}
	return fmt.Sprintf("%s: %s", e.code, e.message)
}

// Code returns the error code.
func (e *Error) Code() Code {
	return e.code
}

// Message returns the error message.
func (e *Error) Message() string {
	return e.message
}

// Unwrap returns the wrapped error, if any.
func (e *Error) Unwrap() error {
	return e.cause
}

// Is reports whether the target error is an Error with the same code.
func (e *Error) Is(target error) bool {
	var obsErr *Error
	if errors.As(target, &obsErr) {
		return e.code == obsErr.code
	}
	return false
}

// WithMessage returns a new Error with the same code but a different message.
func (e *Error) WithMessage(message string) *Error {
	return &Error{
		code:    e.code,
		message: message,
		cause:   e.cause,
	}
}

// WithCause returns a new Error with the same code and message but wrapping a different cause.
func (e *Error) WithCause(cause error) *Error {
	return &Error{
		code:    e.code,
		message: e.message,
		cause:   cause,
	}
}

// Constructors for common error types.

// CollectorInitFailed creates a new collector initialization error.
func CollectorInitFailed(message string) *Error {
	return New(CodeCollectorInitFailed, message)
}

// CollectorInitFailedWrap creates a new collector initialization error wrapping a cause.
func CollectorInitFailedWrap(message string, cause error) *Error {
	return Wrap(CodeCollectorInitFailed, message, cause)
}

// TracerInitFailed creates a new tracer initialization error.
func TracerInitFailed(message string) *Error {
	return New(CodeTracerInitFailed, message)
}

// TracerInitFailedWrap creates a new tracer initialization error wrapping a cause.
func TracerInitFailedWrap(message string, cause error) *Error {
	return Wrap(CodeTracerInitFailed, message, cause)
}

// ExporterError creates a new exporter error.
func ExporterError(message string) *Error {
	return New(CodeExporterError, message)
}

// ExporterErrorWrap creates a new exporter error wrapping a cause.
func ExporterErrorWrap(message string, cause error) *Error {
	return Wrap(CodeExporterError, message, cause)
}

// HealthCheckFailed creates a new health check failure error.
func HealthCheckFailed(component string) *Error {
	return Newf(CodeHealthCheckFailed, "health check failed for component: %s", component)
}

// HealthCheckFailedWrap creates a new health check failure error wrapping a cause.
func HealthCheckFailedWrap(component string, cause error) *Error {
	return Wrapf(CodeHealthCheckFailed, cause, "health check failed for component: %s", component)
}

// HealthCheckTimeout creates a new health check timeout error.
func HealthCheckTimeout(component string) *Error {
	return Newf(CodeHealthCheckTimeout, "health check timed out for component: %s", component)
}

// HealthCheckTimeoutWrap creates a new health check timeout error wrapping a cause.
func HealthCheckTimeoutWrap(component string, cause error) *Error {
	return Wrapf(CodeHealthCheckTimeout, cause, "health check timed out for component: %s", component)
}

// InvalidConfig creates a new invalid configuration error.
func InvalidConfig(message string) *Error {
	return New(CodeInvalidConfig, message)
}

// InvalidConfigf creates a new invalid configuration error with a formatted message.
func InvalidConfigf(format string, args ...any) *Error {
	return Newf(CodeInvalidConfig, format, args...)
}

// RegistrationFailed creates a new registration failure error.
func RegistrationFailed(metric string) *Error {
	return Newf(CodeRegistrationFailed, "failed to register metric: %s", metric)
}

// RegistrationFailedWrap creates a new registration failure error wrapping a cause.
func RegistrationFailedWrap(metric string, cause error) *Error {
	return Wrapf(CodeRegistrationFailed, cause, "failed to register metric: %s", metric)
}

// LabelCardinalityExceeded creates a new label cardinality exceeded error.
func LabelCardinalityExceeded(label string, count int, limit int) *Error {
	return Newf(CodeLabelCardinalityExceeded, "label %q has %d unique values, exceeds limit of %d", label, count, limit)
}

// SpanCreationFailed creates a new span creation failure error.
func SpanCreationFailed(message string) *Error {
	return New(CodeSpanCreationFailed, message)
}

// SpanCreationFailedWrap creates a new span creation failure error wrapping a cause.
func SpanCreationFailedWrap(message string, cause error) *Error {
	return Wrap(CodeSpanCreationFailed, message, cause)
}

// ContextPropagationFailed creates a new context propagation failure error.
func ContextPropagationFailed(message string) *Error {
	return New(CodeContextPropagationFailed, message)
}

// ContextPropagationFailedWrap creates a new context propagation failure error wrapping a cause.
func ContextPropagationFailedWrap(message string, cause error) *Error {
	return Wrap(CodeContextPropagationFailed, message, cause)
}

// ShutdownFailed creates a new shutdown failure error.
func ShutdownFailed(component string) *Error {
	return Newf(CodeShutdownFailed, "failed to shutdown component: %s", component)
}

// ShutdownFailedWrap creates a new shutdown failure error wrapping a cause.
func ShutdownFailedWrap(component string, cause error) *Error {
	return Wrapf(CodeShutdownFailed, cause, "failed to shutdown component: %s", component)
}

// Helper functions for error checking.

// IsError checks if the given error is an observability Error.
func IsError(err error) bool {
	var obsErr *Error
	return errors.As(err, &obsErr)
}

// AsError attempts to extract an Error from the given error.
// Returns nil if the error is not an observability Error.
func AsError(err error) *Error {
	var obsErr *Error
	if errors.As(err, &obsErr) {
		return obsErr
	}
	return nil
}

// GetCode returns the error code from an error if it's an observability Error,
// or an empty string otherwise.
func GetCode(err error) Code {
	if obsErr := AsError(err); obsErr != nil {
		return obsErr.Code()
	}
	return ""
}

// IsCode checks if the given error is an observability Error with the specified code.
func IsCode(err error, code Code) bool {
	if obsErr := AsError(err); obsErr != nil {
		return obsErr.Code() == code
	}
	return false
}

// IsCollectorInitFailed checks if the error is a collector initialization failure.
func IsCollectorInitFailed(err error) bool {
	return IsCode(err, CodeCollectorInitFailed)
}

// IsTracerInitFailed checks if the error is a tracer initialization failure.
func IsTracerInitFailed(err error) bool {
	return IsCode(err, CodeTracerInitFailed)
}

// IsExporterError checks if the error is an exporter error.
func IsExporterError(err error) bool {
	return IsCode(err, CodeExporterError)
}

// IsHealthCheckFailed checks if the error is a health check failure.
func IsHealthCheckFailed(err error) bool {
	return IsCode(err, CodeHealthCheckFailed)
}

// IsHealthCheckTimeout checks if the error is a health check timeout.
func IsHealthCheckTimeout(err error) bool {
	return IsCode(err, CodeHealthCheckTimeout)
}

// IsInvalidConfig checks if the error is an invalid configuration error.
func IsInvalidConfig(err error) bool {
	return IsCode(err, CodeInvalidConfig)
}

// IsRegistrationFailed checks if the error is a registration failure.
func IsRegistrationFailed(err error) bool {
	return IsCode(err, CodeRegistrationFailed)
}

// IsLabelCardinalityExceeded checks if the error is a label cardinality exceeded error.
func IsLabelCardinalityExceeded(err error) bool {
	return IsCode(err, CodeLabelCardinalityExceeded)
}

// IsSpanCreationFailed checks if the error is a span creation failure.
func IsSpanCreationFailed(err error) bool {
	return IsCode(err, CodeSpanCreationFailed)
}

// IsContextPropagationFailed checks if the error is a context propagation failure.
func IsContextPropagationFailed(err error) bool {
	return IsCode(err, CodeContextPropagationFailed)
}

// IsShutdownFailed checks if the error is a shutdown failure.
func IsShutdownFailed(err error) bool {
	return IsCode(err, CodeShutdownFailed)
}
