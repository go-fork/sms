package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/go-fork/sms"
	"github.com/go-fork/sms/adapters/twilio"
	"github.com/go-fork/sms/model"
)

// This example demonstrates how to implement a rate limiter
// to prevent exceeding provider rate limits when sending
// large volumes of messages

func main() {
	// Check if config file path is provided
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go /path/to/config.yaml")
	}
	configFile := os.Args[1]

	fmt.Println("=== SMS Rate Limiter Example ===")
	fmt.Printf("Using config file: %s\n\n", configFile)

	// Initialize the SMS module
	module, err := sms.NewModule(configFile)
	if err != nil {
		log.Fatalf("Failed to initialize SMS module: %v", err)
	}

	// Initialize the Twilio provider
	twilioProvider, err := twilio.NewProvider(configFile)
	if err != nil {
		log.Fatalf("Failed to initialize Twilio provider: %v", err)
	}

	// Add the provider to the module
	if err := module.AddProvider(twilioProvider); err != nil {
		log.Fatalf("Failed to add Twilio provider: %v", err)
	}

	// Create a rate limiter
	// Most SMS providers have rate limits:
	// - Twilio: 100 messages per second
	// - eSMS: 50 messages per second
	// - SpeedSMS: 30 messages per second
	// We'll use a conservative 10 messages per second here
	limiter := rate.NewLimiter(rate.Limit(10), 1) // 10 messages per second, burst of 1

	// Sample list of phone numbers to send to
	recipients := []string{
		"+1234567890",
		"+2345678901",
		"+3456789012",
		// Add more recipients as needed for testing
	}

	fmt.Printf("Sending messages to %d recipients with rate limiting...\n", len(recipients))

	// Create a wait group to wait for all messages to be sent
	var wg sync.WaitGroup
	wg.Add(len(recipients))

	// Track successes and failures
	var (
		successCount int
		failureCount int
		mu           sync.Mutex // To protect counts during concurrent updates
	)

	// Start time for statistics
	startTime := time.Now()

	// Send messages with rate limiting
	for i, recipient := range recipients {
		// Wait for rate limiter's permission
		if err := limiter.Wait(context.Background()); err != nil {
			log.Printf("Rate limiter error: %v", err)
			continue
		}

		// Create a unique message for each recipient
		go func(index int, phoneNumber string) {
			defer wg.Done()

			// Create a context with timeout for this specific message
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// Create the SMS request
			request := model.SendSMSRequest{
				Message: model.Message{
					To: phoneNumber,
					By: "RateLimiterExample",
				},
				Data: map[string]interface{}{
					"app_name": "RateLimiterExample",
					"message":  fmt.Sprintf("This is test message #%d with rate limiting", index+1),
				},
			}

			// Send the SMS
			fmt.Printf("Sending message #%d to %s...\n", index+1, phoneNumber)
			response, err := module.SendSMS(ctx, request)

			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				fmt.Printf("Failed to send message #%d: %v\n", index+1, err)
				failureCount++
				return
			}

			fmt.Printf("Message #%d sent successfully! ID: %s\n", index+1, response.MessageID)
			successCount++
		}(i, recipient)
	}

	// Wait for all messages to be sent
	wg.Wait()

	// Calculate elapsed time
	elapsed := time.Since(startTime)

	// Display statistics
	fmt.Println("\n=== Sending Statistics ===")
	fmt.Printf("Total messages: %d\n", len(recipients))
	fmt.Printf("Successful: %d\n", successCount)
	fmt.Printf("Failed: %d\n", failureCount)
	fmt.Printf("Time elapsed: %v\n", elapsed)
	fmt.Printf("Average rate: %.2f messages/second\n", float64(len(recipients))/elapsed.Seconds())
}
