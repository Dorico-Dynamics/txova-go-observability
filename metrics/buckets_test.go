package metrics

import (
	"testing"
)

func TestHTTPLatencyBuckets(t *testing.T) {
	t.Parallel()

	if len(HTTPLatencyBuckets) == 0 {
		t.Error("HTTPLatencyBuckets is empty")
	}

	// Verify buckets are in ascending order
	for i := 1; i < len(HTTPLatencyBuckets); i++ {
		if HTTPLatencyBuckets[i] <= HTTPLatencyBuckets[i-1] {
			t.Errorf("HTTPLatencyBuckets not in ascending order at index %d: %v <= %v",
				i, HTTPLatencyBuckets[i], HTTPLatencyBuckets[i-1])
		}
	}

	// Verify expected bucket values from PRD
	expected := []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10}
	if len(HTTPLatencyBuckets) != len(expected) {
		t.Errorf("HTTPLatencyBuckets length = %d, want %d", len(HTTPLatencyBuckets), len(expected))
	}
	for i, v := range expected {
		if HTTPLatencyBuckets[i] != v {
			t.Errorf("HTTPLatencyBuckets[%d] = %v, want %v", i, HTTPLatencyBuckets[i], v)
		}
	}
}

func TestDBLatencyBuckets(t *testing.T) {
	t.Parallel()

	if len(DBLatencyBuckets) == 0 {
		t.Error("DBLatencyBuckets is empty")
	}

	// Verify buckets are in ascending order
	for i := 1; i < len(DBLatencyBuckets); i++ {
		if DBLatencyBuckets[i] <= DBLatencyBuckets[i-1] {
			t.Errorf("DBLatencyBuckets not in ascending order at index %d: %v <= %v",
				i, DBLatencyBuckets[i], DBLatencyBuckets[i-1])
		}
	}

	// Verify expected bucket values from PRD
	expected := []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1}
	if len(DBLatencyBuckets) != len(expected) {
		t.Errorf("DBLatencyBuckets length = %d, want %d", len(DBLatencyBuckets), len(expected))
	}
	for i, v := range expected {
		if DBLatencyBuckets[i] != v {
			t.Errorf("DBLatencyBuckets[%d] = %v, want %v", i, DBLatencyBuckets[i], v)
		}
	}
}

func TestDurationBuckets(t *testing.T) {
	t.Parallel()

	if len(DurationBuckets) == 0 {
		t.Error("DurationBuckets is empty")
	}

	// Verify buckets are in ascending order
	for i := 1; i < len(DurationBuckets); i++ {
		if DurationBuckets[i] <= DurationBuckets[i-1] {
			t.Errorf("DurationBuckets not in ascending order at index %d: %v <= %v",
				i, DurationBuckets[i], DurationBuckets[i-1])
		}
	}

	// Verify expected bucket values from PRD (in seconds)
	expected := []float64{60, 300, 600, 900, 1800, 3600}
	if len(DurationBuckets) != len(expected) {
		t.Errorf("DurationBuckets length = %d, want %d", len(DurationBuckets), len(expected))
	}
	for i, v := range expected {
		if DurationBuckets[i] != v {
			t.Errorf("DurationBuckets[%d] = %v, want %v", i, DurationBuckets[i], v)
		}
	}
}

func TestFareBuckets(t *testing.T) {
	t.Parallel()

	if len(FareBuckets) == 0 {
		t.Error("FareBuckets is empty")
	}

	// Verify buckets are in ascending order
	for i := 1; i < len(FareBuckets); i++ {
		if FareBuckets[i] <= FareBuckets[i-1] {
			t.Errorf("FareBuckets not in ascending order at index %d: %v <= %v",
				i, FareBuckets[i], FareBuckets[i-1])
		}
	}

	// Verify expected bucket values from PRD (fare in smallest currency unit)
	expected := []float64{50, 100, 250, 500, 1000, 2500, 5000, 10000, 25000}
	if len(FareBuckets) != len(expected) {
		t.Errorf("FareBuckets length = %d, want %d", len(FareBuckets), len(expected))
	}
	for i, v := range expected {
		if FareBuckets[i] != v {
			t.Errorf("FareBuckets[%d] = %v, want %v", i, FareBuckets[i], v)
		}
	}
}

func TestRequestSizeBuckets(t *testing.T) {
	t.Parallel()

	if len(RequestSizeBuckets) == 0 {
		t.Error("RequestSizeBuckets is empty")
	}

	// Verify buckets are in ascending order
	for i := 1; i < len(RequestSizeBuckets); i++ {
		if RequestSizeBuckets[i] <= RequestSizeBuckets[i-1] {
			t.Errorf("RequestSizeBuckets not in ascending order at index %d: %v <= %v",
				i, RequestSizeBuckets[i], RequestSizeBuckets[i-1])
		}
	}
}

func TestDistanceBuckets(t *testing.T) {
	t.Parallel()

	if len(DistanceBuckets) == 0 {
		t.Error("DistanceBuckets is empty")
	}

	// Verify buckets are in ascending order
	for i := 1; i < len(DistanceBuckets); i++ {
		if DistanceBuckets[i] <= DistanceBuckets[i-1] {
			t.Errorf("DistanceBuckets not in ascending order at index %d: %v <= %v",
				i, DistanceBuckets[i], DistanceBuckets[i-1])
		}
	}
}

func TestPaymentAmountBuckets(t *testing.T) {
	t.Parallel()

	if len(PaymentAmountBuckets) == 0 {
		t.Error("PaymentAmountBuckets is empty")
	}

	// Verify buckets are in ascending order
	for i := 1; i < len(PaymentAmountBuckets); i++ {
		if PaymentAmountBuckets[i] <= PaymentAmountBuckets[i-1] {
			t.Errorf("PaymentAmountBuckets not in ascending order at index %d: %v <= %v",
				i, PaymentAmountBuckets[i], PaymentAmountBuckets[i-1])
		}
	}
}

func TestBucketsHaveReasonableValues(t *testing.T) {
	t.Parallel()

	// HTTP latency should cover sub-second responses
	if HTTPLatencyBuckets[0] >= 1 {
		t.Error("HTTPLatencyBuckets should start below 1 second")
	}

	// DB latency should cover millisecond responses
	if DBLatencyBuckets[0] >= 0.1 {
		t.Error("DBLatencyBuckets should start below 100ms")
	}

	// Duration buckets should cover minutes to hours (in seconds)
	if DurationBuckets[0] < 60 {
		t.Error("DurationBuckets should start at least 1 minute (60s)")
	}
	if DurationBuckets[len(DurationBuckets)-1] < 3600 {
		t.Error("DurationBuckets should cover at least 1 hour (3600s)")
	}

	// Fare buckets should cover realistic fare ranges
	if FareBuckets[0] > 100 {
		t.Error("FareBuckets should start with small fares")
	}
	if FareBuckets[len(FareBuckets)-1] < 10000 {
		t.Error("FareBuckets should cover high fares")
	}
}

func TestAllBucketsPositive(t *testing.T) {
	t.Parallel()

	bucketSets := map[string][]float64{
		"HTTPLatencyBuckets":   HTTPLatencyBuckets,
		"DBLatencyBuckets":     DBLatencyBuckets,
		"DurationBuckets":      DurationBuckets,
		"FareBuckets":          FareBuckets,
		"RequestSizeBuckets":   RequestSizeBuckets,
		"DistanceBuckets":      DistanceBuckets,
		"PaymentAmountBuckets": PaymentAmountBuckets,
	}

	for name, buckets := range bucketSets {
		for i, v := range buckets {
			if v <= 0 {
				t.Errorf("%s[%d] = %v, want positive value", name, i, v)
			}
		}
	}
}
