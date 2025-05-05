package retry

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ErrMaxAttemptsReached is returned when all retry attempts have failed
var ErrMaxAttemptsReached = errors.New("maximum retry attempts reached")

// Config holds retry configuration settings
type Config struct {
	// MaxAttempts is the maximum number of retry attempts (including the initial attempt)
	MaxAttempts int

	// InitialDelay is the initial delay before the first retry
	InitialDelay time.Duration

	// MaxDelay is the maximum delay between retries
	MaxDelay time.Duration

	// Multiplier is the factor by which the delay increases after each attempt
	Multiplier float64

	// RetriableErrors is an optional custom function to determine if an error is retriable
	// If nil, the default IsRetriable function will be used
	RetriableErrors func(error) bool
}

// DefaultConfig returns a default retry configuration
func DefaultConfig() Config {
	return Config{
		MaxAttempts:     3,
		InitialDelay:    500 * time.Millisecond,
		MaxDelay:        30 * time.Second,
		Multiplier:      2.0,
		RetriableErrors: nil,
	}
}

// Do executes the given function with exponential backoff retry logic
// It respects context cancellation and deadlines
func Do(ctx context.Context, config Config, fn func() error) error {
	if config.MaxAttempts <= 0 {
		return errors.New("retry attempts must be greater than 0")
	}

	var err error
	delay := config.InitialDelay

	// Use default IsRetriable if no custom function is provided
	isRetriable := IsRetriable
	if config.RetriableErrors != nil {
		isRetriable = config.RetriableErrors
	}

	// Retry loop
	for attempt := 0; attempt < config.MaxAttempts; attempt++ {
		// Execute the function
		err = fn()

		// If successful or error is not retriable, return immediately
		if err == nil {
			return nil
		}

		// If this is the last attempt or the error is not retriable, return the error
		if attempt == config.MaxAttempts-1 || !isRetriable(err) {
			if attempt == config.MaxAttempts-1 {
				return fmt.Errorf("%w: %v", ErrMaxAttemptsReached, err)
			}
			return err
		}

		// Calculate next delay with exponential backoff
		nextDelay := time.Duration(float64(delay) * config.Multiplier)
		if nextDelay > config.MaxDelay {
			nextDelay = config.MaxDelay
		}
		delay = nextDelay

		// Wait for the delay or until context is canceled
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
			// Continue to next attempt
		}
	}

	// This should not be reached with proper implementation
	return err
}

// DoWithOptions is a simpler version of Do that takes individual parameters instead of a Config
func DoWithOptions(ctx context.Context, attempts int, initialDelay time.Duration, fn func() error) error {
	config := Config{
		MaxAttempts:  attempts,
		InitialDelay: initialDelay,
		MaxDelay:     30 * time.Second,
		Multiplier:   2.0,
	}
	return Do(ctx, config, fn)
}

// IsRetriable determines if an error should be retried
// Returns true for network errors, timeouts, and 5xx status codes
func IsRetriable(err error) bool {
	if err == nil {
		return false
	}

	// Check for network errors
	var netErr net.Error
	if errors.As(err, &netErr) {
		// Retry network errors, especially timeouts
		return true
	}

	// Check for specific error types
	if errors.Is(err, context.DeadlineExceeded) ||
		errors.Is(err, context.Canceled) {
		// Don't retry context cancellation or deadline errors
		return false
	}

	// Check for URL errors (often network-related)
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		// URL errors are often temporary and retriable
		return true
	}

	// Check for HTTP status code errors
	// This assumes that HTTP errors contain the status code in the error message
	// Adapt this logic based on how your application represents HTTP errors
	errStr := err.Error()

	// Check for 5xx errors (server errors)
	if strings.Contains(errStr, "status code: 5") ||
		strings.Contains(errStr, "500") ||
		strings.Contains(errStr, "502") ||
		strings.Contains(errStr, "503") ||
		strings.Contains(errStr, "504") {
		return true
	}

	// Check for specific error messages that typically indicate temporary issues
	if strings.Contains(strings.ToLower(errStr), "timeout") ||
		strings.Contains(strings.ToLower(errStr), "connection refused") ||
		strings.Contains(strings.ToLower(errStr), "connection reset") ||
		strings.Contains(strings.ToLower(errStr), "temporary") ||
		strings.Contains(strings.ToLower(errStr), "too many requests") ||
		strings.Contains(strings.ToLower(errStr), "service unavailable") {
		return true
	}

	// For HTTP response errors, we can check more specifically
	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		return httpErr.StatusCode >= 500 || httpErr.StatusCode == 429
	}

	// Default to not retrying unknown errors
	return false
}

// HTTPError represents an HTTP error with status code
// This is an example struct that can be used to pass HTTP errors with status codes
type HTTPError struct {
	StatusCode int
	Message    string
}

// Error implements the error interface
func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP error: status code: %d, message: %s", e.StatusCode, e.Message)
}

// NewHTTPError creates a new HTTPError
func NewHTTPError(statusCode int, message string) *HTTPError {
	return &HTTPError{
		StatusCode: statusCode,
		Message:    message,
	}
}

// RetryableHTTPCodes returns a slice of HTTP status codes that are considered retriable
func RetryableHTTPCodes() []int {
	return []int{
		http.StatusInternalServerError, // 500
		http.StatusBadGateway,          // 502
		http.StatusServiceUnavailable,  // 503
		http.StatusGatewayTimeout,      // 504
		http.StatusTooManyRequests,     // 429
	}
}

// IsRetriableHTTPCode checks if an HTTP status code is retriable
func IsRetriableHTTPCode(statusCode int) bool {
	for _, code := range RetryableHTTPCodes() {
		if statusCode == code {
			return true
		}
	}
	return false
}
