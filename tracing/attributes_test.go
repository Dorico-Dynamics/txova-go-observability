package tracing

import (
	"testing"
)

func TestAttributeConstants(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		constant string
		expected string
	}{
		// Service attributes
		{"AttrServiceName", AttrServiceName, "service.name"},
		{"AttrServiceVersion", AttrServiceVersion, "service.version"},

		// User attributes
		{"AttrUserID", AttrUserID, "user.id"},

		// Request attributes
		{"AttrRequestID", AttrRequestID, "request.id"},
		{"AttrCorrelationID", AttrCorrelationID, "correlation.id"},

		// HTTP attributes
		{"AttrHTTPMethod", AttrHTTPMethod, "http.method"},
		{"AttrHTTPRoute", AttrHTTPRoute, "http.route"},
		{"AttrHTTPStatusCode", AttrHTTPStatusCode, "http.status_code"},
		{"AttrHTTPURL", AttrHTTPURL, "http.url"},
		{"AttrHTTPScheme", AttrHTTPScheme, "http.scheme"},
		{"AttrHTTPHost", AttrHTTPHost, "http.host"},
		{"AttrHTTPUserAgent", AttrHTTPUserAgent, "http.user_agent"},
		{"AttrHTTPClientIP", AttrHTTPClientIP, "http.client_ip"},

		// Database attributes
		{"AttrDBSystem", AttrDBSystem, "db.system"},
		{"AttrDBOperation", AttrDBOperation, "db.operation"},
		{"AttrDBName", AttrDBName, "db.name"},
		{"AttrDBStatement", AttrDBStatement, "db.statement"},

		// Messaging attributes
		{"AttrMessagingSystem", AttrMessagingSystem, "messaging.system"},
		{"AttrMessagingDestination", AttrMessagingDestination, "messaging.destination"},
		{"AttrMessagingOperation", AttrMessagingOperation, "messaging.operation"},
		{"AttrMessagingMessageID", AttrMessagingMessageID, "messaging.message_id"},
		{"AttrMessagingConsumer", AttrMessagingConsumer, "messaging.consumer.group"},

		// Error attributes
		{"AttrErrorType", AttrErrorType, "error.type"},
		{"AttrErrorMessage", AttrErrorMessage, "error.message"},

		// Business attributes
		{"AttrRideID", AttrRideID, "ride.id"},
		{"AttrDriverID", AttrDriverID, "driver.id"},
		{"AttrRiderID", AttrRiderID, "rider.id"},
		{"AttrPaymentID", AttrPaymentID, "payment.id"},
		{"AttrServiceType", AttrServiceType, "service.type"},
		{"AttrCity", AttrCity, "city"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.constant != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, tt.constant, tt.expected)
			}
		})
	}
}

func TestServiceName(t *testing.T) {
	t.Parallel()

	attr := ServiceName("ride-service")
	if string(attr.Key) != AttrServiceName {
		t.Errorf("Key = %v, want %v", attr.Key, AttrServiceName)
	}
	if attr.Value.AsString() != "ride-service" {
		t.Errorf("Value = %v, want ride-service", attr.Value.AsString())
	}
}

func TestServiceVersion(t *testing.T) {
	t.Parallel()

	attr := ServiceVersion("v1.2.3")
	if string(attr.Key) != AttrServiceVersion {
		t.Errorf("Key = %v, want %v", attr.Key, AttrServiceVersion)
	}
	if attr.Value.AsString() != "v1.2.3" {
		t.Errorf("Value = %v, want v1.2.3", attr.Value.AsString())
	}
}

func TestUserID(t *testing.T) {
	t.Parallel()

	attr := UserID("user-123")
	if string(attr.Key) != AttrUserID {
		t.Errorf("Key = %v, want %v", attr.Key, AttrUserID)
	}
	if attr.Value.AsString() != "user-123" {
		t.Errorf("Value = %v, want user-123", attr.Value.AsString())
	}
}

func TestRequestID(t *testing.T) {
	t.Parallel()

	attr := RequestID("req-abc-123")
	if string(attr.Key) != AttrRequestID {
		t.Errorf("Key = %v, want %v", attr.Key, AttrRequestID)
	}
	if attr.Value.AsString() != "req-abc-123" {
		t.Errorf("Value = %v, want req-abc-123", attr.Value.AsString())
	}
}

func TestCorrelationID(t *testing.T) {
	t.Parallel()

	attr := CorrelationID("corr-xyz-789")
	if string(attr.Key) != AttrCorrelationID {
		t.Errorf("Key = %v, want %v", attr.Key, AttrCorrelationID)
	}
	if attr.Value.AsString() != "corr-xyz-789" {
		t.Errorf("Value = %v, want corr-xyz-789", attr.Value.AsString())
	}
}

func TestHTTPMethod(t *testing.T) {
	t.Parallel()

	attr := HTTPMethod("POST")
	if string(attr.Key) != AttrHTTPMethod {
		t.Errorf("Key = %v, want %v", attr.Key, AttrHTTPMethod)
	}
	if attr.Value.AsString() != "POST" {
		t.Errorf("Value = %v, want POST", attr.Value.AsString())
	}
}

func TestHTTPRoute(t *testing.T) {
	t.Parallel()

	attr := HTTPRoute("/api/v1/rides/{id}")
	if string(attr.Key) != AttrHTTPRoute {
		t.Errorf("Key = %v, want %v", attr.Key, AttrHTTPRoute)
	}
	if attr.Value.AsString() != "/api/v1/rides/{id}" {
		t.Errorf("Value = %v, want /api/v1/rides/{id}", attr.Value.AsString())
	}
}

func TestHTTPStatusCode(t *testing.T) {
	t.Parallel()

	attr := HTTPStatusCode(200)
	if string(attr.Key) != AttrHTTPStatusCode {
		t.Errorf("Key = %v, want %v", attr.Key, AttrHTTPStatusCode)
	}
	if attr.Value.AsInt64() != 200 {
		t.Errorf("Value = %v, want 200", attr.Value.AsInt64())
	}
}

func TestHTTPURL(t *testing.T) {
	t.Parallel()

	attr := HTTPURL("https://api.txova.com/v1/rides")
	if string(attr.Key) != AttrHTTPURL {
		t.Errorf("Key = %v, want %v", attr.Key, AttrHTTPURL)
	}
	if attr.Value.AsString() != "https://api.txova.com/v1/rides" {
		t.Errorf("Value = %v, want https://api.txova.com/v1/rides", attr.Value.AsString())
	}
}

func TestDBSystem(t *testing.T) {
	t.Parallel()

	attr := DBSystem("postgresql")
	if string(attr.Key) != AttrDBSystem {
		t.Errorf("Key = %v, want %v", attr.Key, AttrDBSystem)
	}
	if attr.Value.AsString() != "postgresql" {
		t.Errorf("Value = %v, want postgresql", attr.Value.AsString())
	}
}

func TestDBOperation(t *testing.T) {
	t.Parallel()

	attr := DBOperation("SELECT")
	if string(attr.Key) != AttrDBOperation {
		t.Errorf("Key = %v, want %v", attr.Key, AttrDBOperation)
	}
	if attr.Value.AsString() != "SELECT" {
		t.Errorf("Value = %v, want SELECT", attr.Value.AsString())
	}
}

func TestMessagingSystem(t *testing.T) {
	t.Parallel()

	attr := MessagingSystem("kafka")
	if string(attr.Key) != AttrMessagingSystem {
		t.Errorf("Key = %v, want %v", attr.Key, AttrMessagingSystem)
	}
	if attr.Value.AsString() != "kafka" {
		t.Errorf("Value = %v, want kafka", attr.Value.AsString())
	}
}

func TestMessagingDestination(t *testing.T) {
	t.Parallel()

	attr := MessagingDestination("ride_events")
	if string(attr.Key) != AttrMessagingDestination {
		t.Errorf("Key = %v, want %v", attr.Key, AttrMessagingDestination)
	}
	if attr.Value.AsString() != "ride_events" {
		t.Errorf("Value = %v, want ride_events", attr.Value.AsString())
	}
}

func TestErrorType(t *testing.T) {
	t.Parallel()

	attr := ErrorType("ValidationError")
	if string(attr.Key) != AttrErrorType {
		t.Errorf("Key = %v, want %v", attr.Key, AttrErrorType)
	}
	if attr.Value.AsString() != "ValidationError" {
		t.Errorf("Value = %v, want ValidationError", attr.Value.AsString())
	}
}

func TestErrorMessage(t *testing.T) {
	t.Parallel()

	attr := ErrorMessage("invalid input")
	if string(attr.Key) != AttrErrorMessage {
		t.Errorf("Key = %v, want %v", attr.Key, AttrErrorMessage)
	}
	if attr.Value.AsString() != "invalid input" {
		t.Errorf("Value = %v, want invalid input", attr.Value.AsString())
	}
}

func TestRideID(t *testing.T) {
	t.Parallel()

	attr := RideID("ride-456")
	if string(attr.Key) != AttrRideID {
		t.Errorf("Key = %v, want %v", attr.Key, AttrRideID)
	}
	if attr.Value.AsString() != "ride-456" {
		t.Errorf("Value = %v, want ride-456", attr.Value.AsString())
	}
}

func TestDriverID(t *testing.T) {
	t.Parallel()

	attr := DriverID("driver-789")
	if string(attr.Key) != AttrDriverID {
		t.Errorf("Key = %v, want %v", attr.Key, AttrDriverID)
	}
	if attr.Value.AsString() != "driver-789" {
		t.Errorf("Value = %v, want driver-789", attr.Value.AsString())
	}
}

func TestRiderID(t *testing.T) {
	t.Parallel()

	attr := RiderID("rider-101")
	if string(attr.Key) != AttrRiderID {
		t.Errorf("Key = %v, want %v", attr.Key, AttrRiderID)
	}
	if attr.Value.AsString() != "rider-101" {
		t.Errorf("Value = %v, want rider-101", attr.Value.AsString())
	}
}

func TestPaymentID(t *testing.T) {
	t.Parallel()

	attr := PaymentID("pay-202")
	if string(attr.Key) != AttrPaymentID {
		t.Errorf("Key = %v, want %v", attr.Key, AttrPaymentID)
	}
	if attr.Value.AsString() != "pay-202" {
		t.Errorf("Value = %v, want pay-202", attr.Value.AsString())
	}
}

func TestServiceType(t *testing.T) {
	t.Parallel()

	attr := ServiceType("premium")
	if string(attr.Key) != AttrServiceType {
		t.Errorf("Key = %v, want %v", attr.Key, AttrServiceType)
	}
	if attr.Value.AsString() != "premium" {
		t.Errorf("Value = %v, want premium", attr.Value.AsString())
	}
}

func TestCity(t *testing.T) {
	t.Parallel()

	attr := City("maputo")
	if string(attr.Key) != AttrCity {
		t.Errorf("Key = %v, want %v", attr.Key, AttrCity)
	}
	if attr.Value.AsString() != "maputo" {
		t.Errorf("Value = %v, want maputo", attr.Value.AsString())
	}
}
