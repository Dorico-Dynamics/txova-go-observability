package tracing

import (
	"go.opentelemetry.io/otel/attribute"
)

// Standard attribute keys for spans.
const (
	// Service attributes.
	AttrServiceName    = "service.name"
	AttrServiceVersion = "service.version"

	// User attributes.
	AttrUserID = "user.id"

	// Request attributes.
	AttrRequestID     = "request.id"
	AttrCorrelationID = "correlation.id"

	// HTTP attributes.
	AttrHTTPMethod     = "http.method"
	AttrHTTPRoute      = "http.route"
	AttrHTTPStatusCode = "http.status_code"
	AttrHTTPURL        = "http.url"
	AttrHTTPScheme     = "http.scheme"
	AttrHTTPHost       = "http.host"
	AttrHTTPUserAgent  = "http.user_agent"
	AttrHTTPClientIP   = "http.client_ip"

	// Database attributes.
	AttrDBSystem    = "db.system"
	AttrDBOperation = "db.operation"
	AttrDBName      = "db.name"
	AttrDBStatement = "db.statement"

	// Messaging attributes.
	AttrMessagingSystem      = "messaging.system"
	AttrMessagingDestination = "messaging.destination"
	AttrMessagingOperation   = "messaging.operation"
	AttrMessagingMessageID   = "messaging.message_id"
	AttrMessagingConsumer    = "messaging.consumer.group"

	// Error attributes.
	AttrErrorType    = "error.type"
	AttrErrorMessage = "error.message"

	// Business attributes.
	AttrRideID      = "ride.id"
	AttrDriverID    = "driver.id"
	AttrRiderID     = "rider.id"
	AttrPaymentID   = "payment.id"
	AttrServiceType = "service.type"
	AttrCity        = "city"
)

// ServiceName creates a service name attribute.
func ServiceName(name string) attribute.KeyValue {
	return attribute.String(AttrServiceName, name)
}

// ServiceVersion creates a service version attribute.
func ServiceVersion(version string) attribute.KeyValue {
	return attribute.String(AttrServiceVersion, version)
}

// UserID creates a user ID attribute.
func UserID(id string) attribute.KeyValue {
	return attribute.String(AttrUserID, id)
}

// RequestID creates a request ID attribute.
func RequestID(id string) attribute.KeyValue {
	return attribute.String(AttrRequestID, id)
}

// CorrelationID creates a correlation ID attribute.
func CorrelationID(id string) attribute.KeyValue {
	return attribute.String(AttrCorrelationID, id)
}

// HTTPMethod creates an HTTP method attribute.
func HTTPMethod(method string) attribute.KeyValue {
	return attribute.String(AttrHTTPMethod, method)
}

// HTTPRoute creates an HTTP route attribute.
func HTTPRoute(route string) attribute.KeyValue {
	return attribute.String(AttrHTTPRoute, route)
}

// HTTPStatusCode creates an HTTP status code attribute.
func HTTPStatusCode(code int) attribute.KeyValue {
	return attribute.Int(AttrHTTPStatusCode, code)
}

// HTTPURL creates an HTTP URL attribute.
func HTTPURL(url string) attribute.KeyValue {
	return attribute.String(AttrHTTPURL, url)
}

// HTTPScheme creates an HTTP scheme attribute.
func HTTPScheme(scheme string) attribute.KeyValue {
	return attribute.String(AttrHTTPScheme, scheme)
}

// HTTPHost creates an HTTP host attribute.
func HTTPHost(host string) attribute.KeyValue {
	return attribute.String(AttrHTTPHost, host)
}

// HTTPUserAgent creates an HTTP user agent attribute.
func HTTPUserAgent(userAgent string) attribute.KeyValue {
	return attribute.String(AttrHTTPUserAgent, userAgent)
}

// HTTPClientIP creates an HTTP client IP attribute.
func HTTPClientIP(ip string) attribute.KeyValue {
	return attribute.String(AttrHTTPClientIP, ip)
}

// DBSystem creates a database system attribute.
func DBSystem(system string) attribute.KeyValue {
	return attribute.String(AttrDBSystem, system)
}

// DBOperation creates a database operation attribute.
func DBOperation(operation string) attribute.KeyValue {
	return attribute.String(AttrDBOperation, operation)
}

// DBName creates a database name attribute.
func DBName(name string) attribute.KeyValue {
	return attribute.String(AttrDBName, name)
}

// DBStatement creates a database statement attribute.
func DBStatement(statement string) attribute.KeyValue {
	return attribute.String(AttrDBStatement, statement)
}

// MessagingSystem creates a messaging system attribute.
func MessagingSystem(system string) attribute.KeyValue {
	return attribute.String(AttrMessagingSystem, system)
}

// MessagingDestination creates a messaging destination attribute.
func MessagingDestination(destination string) attribute.KeyValue {
	return attribute.String(AttrMessagingDestination, destination)
}

// MessagingOperation creates a messaging operation attribute.
func MessagingOperation(operation string) attribute.KeyValue {
	return attribute.String(AttrMessagingOperation, operation)
}

// MessagingMessageID creates a messaging message ID attribute.
func MessagingMessageID(id string) attribute.KeyValue {
	return attribute.String(AttrMessagingMessageID, id)
}

// MessagingConsumer creates a messaging consumer group attribute.
func MessagingConsumer(group string) attribute.KeyValue {
	return attribute.String(AttrMessagingConsumer, group)
}

// ErrorType creates an error type attribute.
func ErrorType(errType string) attribute.KeyValue {
	return attribute.String(AttrErrorType, errType)
}

// ErrorMessage creates an error message attribute.
func ErrorMessage(message string) attribute.KeyValue {
	return attribute.String(AttrErrorMessage, message)
}

// RideID creates a ride ID attribute.
func RideID(id string) attribute.KeyValue {
	return attribute.String(AttrRideID, id)
}

// DriverID creates a driver ID attribute.
func DriverID(id string) attribute.KeyValue {
	return attribute.String(AttrDriverID, id)
}

// RiderID creates a rider ID attribute.
func RiderID(id string) attribute.KeyValue {
	return attribute.String(AttrRiderID, id)
}

// PaymentID creates a payment ID attribute.
func PaymentID(id string) attribute.KeyValue {
	return attribute.String(AttrPaymentID, id)
}

// ServiceType creates a service type attribute.
func ServiceType(serviceType string) attribute.KeyValue {
	return attribute.String(AttrServiceType, serviceType)
}

// City creates a city attribute.
func City(city string) attribute.KeyValue {
	return attribute.String(AttrCity, city)
}
