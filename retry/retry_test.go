package retry

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"
)

func TestDo(t *testing.T) {
	tests := []struct {
		name           string
		config         Config
		fn             func() error
		expectedCalls  int
		expectedError  bool
		contextTimeout time.Duration
	}{
		{
			name:          "Success on first attempt",
			config:        DefaultConfig(),
			fn:            func() error { return nil },
			expectedCalls: 1,
			expectedError: false,
		},
		{
			name:   "Success after retries",
			config: DefaultConfig(),
			fn: (func() func() error {
				calls := 0
				return func() error {
					calls++
					if calls < 3 {
						return &HTTPError{StatusCode: 500, Message: "Server error"}
					}
					return nil
				}
			})(),
			expectedCalls: 3,
			expectedError: false,
		},
		{
			name:   "Max attempts reached",
			config: DefaultConfig(),
			fn: func() error {
				return &HTTPError{StatusCode: 500, Message: "Server error"}
			},
			expectedCalls: 3, // Default MaxAttempts is 3
			expectedError: true,
		},
		{
			name:   "Non-retriable error",
			config: DefaultConfig(),
			fn: func() error {
				return &HTTPError{StatusCode: 400, Message: "Bad request"}
			},
			expectedCalls: 1, // Should stop after first non-retriable error
			expectedError: true,
		},
		{
			name:   "Context cancellation",
			config: DefaultConfig(),
			fn: func() error {
				return &HTTPError{StatusCode: 500, Message: "Server error"}
			},
			expectedCalls:  1, // Should stop after context is cancelled
			expectedError:  true,
			contextTimeout: 50 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			wrappedFn := func() error {
				calls++
				return tt.fn()
			}

			var ctx context.Context
			var cancel context.CancelFunc

			if tt.contextTimeout > 0 {
				ctx, cancel = context.WithTimeout(context.Background(), tt.contextTimeout)
			} else {
				ctx, cancel = context.WithCancel(context.Background())
			}
			defer cancel()

			// For the context cancellation test, make the initial delay longer
			// than the context timeout to ensure cancellation happens
			if tt.name == "Context cancellation" {
				tt.config.InitialDelay = 100 * time.Millisecond
			}

			err := Do(ctx, tt.config, wrappedFn)

			if (err != nil) != tt.expectedError {
				t.Errorf("Do() error = %v, expectedError %v", err, tt.expectedError)
			}

			if calls != tt.expectedCalls {
				t.Errorf("Do() called function %d times, expected %d", calls, tt.expectedCalls)
			}

			if tt.name == "Context cancellation" && !errors.Is(err, context.DeadlineExceeded) {
				t.Errorf("Expected context deadline exceeded error, got %v", err)
			}
		})
	}
}

func TestIsRetriable(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "Nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "Network error",
			err:      &net.OpError{Op: "dial", Err: errors.New("connection refused")},
			expected: true,
		},
		{
			name:     "Context deadline exceeded",
			err:      context.DeadlineExceeded,
			expected: false,
		},
		{
			name:     "Context canceled",
			err:      context.Canceled,
			expected: false,
		},
		{
			name:     "HTTP 500 error",
			err:      &HTTPError{StatusCode: 500, Message: "Internal server error"},
			expected: true,
		},
		{
			name:     "HTTP 429 error",
			err:      &HTTPError{StatusCode: 429, Message: "Too many requests"},
			expected: true,
		},
		{
			name:     "HTTP 400 error",
			err:      &HTTPError{StatusCode: 400, Message: "Bad request"},
			expected: false,
		},
		{
			name:     "HTTP 401 error",
			err:      &HTTPError{StatusCode: 401, Message: "Unauthorized"},
			expected: false,
		},
		{
			name:     "Error with timeout in message",
			err:      errors.New("operation timed out"),
			expected: true,
		},
		{
			name:     "Error with connection refused in message",
			err:      errors.New("connection refused"),
			expected: true,
		},
		{
			name:     "Generic error",
			err:      errors.New("some error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRetriable(tt.err); got != tt.expected {
				t.Errorf("IsRetriable() = %v, expected %v", got, tt.expected)
			}
		})
	}
}
