package metrics

// Histogram bucket presets for different metric types.
// These are based on expected distributions for Txova services.

// HTTPLatencyBuckets defines histogram buckets for HTTP request latency in seconds.
// Covers range from 5ms to 10s, suitable for API response times.
var HTTPLatencyBuckets = []float64{
	0.005, // 5ms
	0.01,  // 10ms
	0.025, // 25ms
	0.05,  // 50ms
	0.1,   // 100ms
	0.25,  // 250ms
	0.5,   // 500ms
	1,     // 1s
	2.5,   // 2.5s
	5,     // 5s
	10,    // 10s
}

// DBLatencyBuckets defines histogram buckets for database query latency in seconds.
// Covers range from 1ms to 1s, suitable for database operations.
var DBLatencyBuckets = []float64{
	0.001, // 1ms
	0.005, // 5ms
	0.01,  // 10ms
	0.025, // 25ms
	0.05,  // 50ms
	0.1,   // 100ms
	0.25,  // 250ms
	0.5,   // 500ms
	1,     // 1s
}

// DurationBuckets defines histogram buckets for longer durations in seconds.
// Covers range from 1 minute to 1 hour, suitable for ride durations.
var DurationBuckets = []float64{
	60,   // 1 minute
	300,  // 5 minutes
	600,  // 10 minutes
	900,  // 15 minutes
	1800, // 30 minutes
	3600, // 1 hour
}

// FareBuckets defines histogram buckets for fare amounts in MZN (Mozambican Metical).
// Covers range from 50 MZN to 25,000 MZN, suitable for ride fares.
var FareBuckets = []float64{
	50,    // 50 MZN
	100,   // 100 MZN
	250,   // 250 MZN
	500,   // 500 MZN
	1000,  // 1,000 MZN
	2500,  // 2,500 MZN
	5000,  // 5,000 MZN
	10000, // 10,000 MZN
	25000, // 25,000 MZN
}

// RequestSizeBuckets defines histogram buckets for HTTP request/response body sizes in bytes.
// Covers range from 100 bytes to 10MB.
var RequestSizeBuckets = []float64{
	100,      // 100 B
	1000,     // 1 KB
	10000,    // 10 KB
	100000,   // 100 KB
	1000000,  // 1 MB
	10000000, // 10 MB
}

// DistanceBuckets defines histogram buckets for ride distances in kilometers.
// Covers range from 1km to 50km, suitable for urban ride distances.
var DistanceBuckets = []float64{
	1,  // 1 km
	2,  // 2 km
	5,  // 5 km
	10, // 10 km
	15, // 15 km
	20, // 20 km
	30, // 30 km
	50, // 50 km
}

// PaymentAmountBuckets defines histogram buckets for payment amounts in MZN centavos.
// Covers range from 50 MZN to 25,000 MZN (as centavos: 5000 to 2,500,000).
var PaymentAmountBuckets = []float64{
	5000,    // 50 MZN
	10000,   // 100 MZN
	25000,   // 250 MZN
	50000,   // 500 MZN
	100000,  // 1,000 MZN
	250000,  // 2,500 MZN
	500000,  // 5,000 MZN
	1000000, // 10,000 MZN
	2500000, // 25,000 MZN
}
