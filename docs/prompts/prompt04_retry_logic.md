# Prompt 4: Retry Logic

## Objective
Implement retry logic with exponential backoff for handling temporary failures when sending messages.

## Required Files to Create

1. `/retry/retry.go` - Implementation of retry mechanism

## Implementation Requirements

### Retry Function
- Create a generic retry function:
  - `func Do(ctx context.Context, attempts int, delay time.Duration, fn func() error) error`
  - The function should:
    - Execute the provided function `fn`
    - If `fn` returns an error, wait for the specified delay
    - Retry up to the specified number of attempts
    - Use exponential backoff (double the delay after each attempt)
    - Respect context cancellation
    - Return the last error if all attempts fail

### Error Classification
- Implement a helper function to determine if an error is retriable:
  - `func IsRetriable(err error) bool`
  - Consider network errors, timeouts, and server errors (5xx) as retriable
  - Authentication errors, client errors (4xx), and validation errors should not be retried

### Context Integration
- Ensure the retry mechanism respects context deadlines and cancellation
- Allow for early termination of retry loop if context is cancelled

## Deliverables
- Complete retry implementation with exponential backoff
- Helper functions for error classification
- Context-aware retry mechanism
