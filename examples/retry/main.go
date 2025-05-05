package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/zinzinday/go-sms/retry"
)

// simulateAPICall simulates an API call that might fail
func simulateAPICall() error {
	// Simulate random failures
	rand.Seed(time.Now().UnixNano())

	// 60% chance of failure
	if rand.Float32() < 0.6 {
		// Simulate different types of errors
		errorTypes := []error{
			retry.NewHTTPError(500, "Internal Server Error"),
			retry.NewHTTPError(502, "Bad Gateway"),
			retry.NewHTTPError(503, "Service Unavailable"),
			errors.New("connection reset by peer"),
			errors.New("timeout occurred"),
		}

		return errorTypes[rand.Intn(len(errorTypes))]
	}

	// 40% chance of success
	return nil
}

func main() {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Configure retry settings
	config := retry.Config{
		MaxAttempts:  5,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     2 * time.Second,
		Multiplier:   1.5,
	}

	// Track attempts
	attempts := 0

	// Execute with retry
	err := retry.Do(ctx, config, func() error {
		attempts++
		fmt.Printf("Attempt %d: Sending SMS... ", attempts)

		// Simulate API call
		err := simulateAPICall()
		if err != nil {
			fmt.Printf("Failed: %v\n", err)
			return err
		}

		fmt.Println("Success!")
		return nil
	})

	// Check final result
	if err != nil {
		if errors.Is(err, retry.ErrMaxAttemptsReached) {
			log.Fatalf("Failed after %d attempts: %v", attempts, err)
		} else {
			log.Fatalf("Failed with non-retriable error: %v", err)
		}
	} else {
		fmt.Printf("Operation completed successfully after %d attempts\n", attempts)
	}
}
