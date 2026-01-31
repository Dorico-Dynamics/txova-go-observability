package health

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

// PostgresChecker checks the health of a PostgreSQL database.
type PostgresChecker struct {
	name     string
	db       *sql.DB
	required bool
}

// NewPostgresChecker creates a new PostgreSQL health checker.
func NewPostgresChecker(name string, db *sql.DB, required bool) *PostgresChecker {
	return &PostgresChecker{
		name:     name,
		db:       db,
		required: required,
	}
}

// Name returns the name of the checker.
func (c *PostgresChecker) Name() string {
	return c.name
}

// Check performs the health check.
func (c *PostgresChecker) Check(ctx context.Context) Result {
	start := time.Now()

	if err := c.db.PingContext(ctx); err != nil {
		return NewUnhealthyResult(time.Since(start), err)
	}

	return NewHealthyResult(time.Since(start))
}

// Required returns whether this check is required.
func (c *PostgresChecker) Required() bool {
	return c.required
}

// Pinger defines the interface for Redis-like clients that support Ping.
type Pinger interface {
	Ping(ctx context.Context) error
}

// RedisChecker checks the health of a Redis connection.
type RedisChecker struct {
	name     string
	client   Pinger
	required bool
}

// NewRedisChecker creates a new Redis health checker.
func NewRedisChecker(name string, client Pinger, required bool) *RedisChecker {
	return &RedisChecker{
		name:     name,
		client:   client,
		required: required,
	}
}

// Name returns the name of the checker.
func (c *RedisChecker) Name() string {
	return c.name
}

// Check performs the health check.
func (c *RedisChecker) Check(ctx context.Context) Result {
	start := time.Now()

	if err := c.client.Ping(ctx); err != nil {
		return NewUnhealthyResult(time.Since(start), err)
	}

	return NewHealthyResult(time.Since(start))
}

// Required returns whether this check is required.
func (c *RedisChecker) Required() bool {
	return c.required
}

// KafkaMetadataFetcher defines the interface for Kafka clients that can fetch metadata.
type KafkaMetadataFetcher interface {
	// GetMetadata fetches cluster metadata to verify connectivity.
	GetMetadata(ctx context.Context) error
}

// KafkaChecker checks the health of a Kafka connection.
type KafkaChecker struct {
	name     string
	client   KafkaMetadataFetcher
	required bool
}

// NewKafkaChecker creates a new Kafka health checker.
func NewKafkaChecker(name string, client KafkaMetadataFetcher, required bool) *KafkaChecker {
	return &KafkaChecker{
		name:     name,
		client:   client,
		required: required,
	}
}

// Name returns the name of the checker.
func (c *KafkaChecker) Name() string {
	return c.name
}

// Check performs the health check.
func (c *KafkaChecker) Check(ctx context.Context) Result {
	start := time.Now()

	if err := c.client.GetMetadata(ctx); err != nil {
		return NewUnhealthyResult(time.Since(start), err)
	}

	return NewHealthyResult(time.Since(start))
}

// Required returns whether this check is required.
func (c *KafkaChecker) Required() bool {
	return c.required
}

// HTTPChecker checks the health of an external HTTP service.
type HTTPChecker struct {
	name       string
	url        string
	client     *http.Client
	required   bool
	expectCode int
}

// NewHTTPChecker creates a new HTTP health checker.
func NewHTTPChecker(name, url string, client *http.Client, required bool) *HTTPChecker {
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}
	return &HTTPChecker{
		name:       name,
		url:        url,
		client:     client,
		required:   required,
		expectCode: http.StatusOK,
	}
}

// WithExpectedStatusCode sets the expected HTTP status code.
func (c *HTTPChecker) WithExpectedStatusCode(code int) *HTTPChecker {
	c.expectCode = code
	return c
}

// Name returns the name of the checker.
func (c *HTTPChecker) Name() string {
	return c.name
}

// Check performs the health check.
func (c *HTTPChecker) Check(ctx context.Context) Result {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url, nil)
	if err != nil {
		return NewUnhealthyResult(time.Since(start), err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return NewUnhealthyResult(time.Since(start), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != c.expectCode {
		return NewUnhealthyResult(time.Since(start),
			fmt.Errorf("unexpected status code: got %d, want %d", resp.StatusCode, c.expectCode))
	}

	return NewHealthyResult(time.Since(start)).WithDetails(map[string]any{
		"status_code": resp.StatusCode,
	})
}

// Required returns whether this check is required.
func (c *HTTPChecker) Required() bool {
	return c.required
}

// FuncChecker wraps a function as a health checker.
type FuncChecker struct {
	name     string
	checkFn  func(ctx context.Context) error
	required bool
}

// NewFuncChecker creates a new function-based health checker.
func NewFuncChecker(name string, checkFn func(ctx context.Context) error, required bool) *FuncChecker {
	return &FuncChecker{
		name:     name,
		checkFn:  checkFn,
		required: required,
	}
}

// Name returns the name of the checker.
func (c *FuncChecker) Name() string {
	return c.name
}

// Check performs the health check.
func (c *FuncChecker) Check(ctx context.Context) Result {
	start := time.Now()

	if err := c.checkFn(ctx); err != nil {
		return NewUnhealthyResult(time.Since(start), err)
	}

	return NewHealthyResult(time.Since(start))
}

// Required returns whether this check is required.
func (c *FuncChecker) Required() bool {
	return c.required
}
