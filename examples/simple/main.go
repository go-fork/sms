package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/zinzinday/go-sms"
	"github.com/zinzinday/go-sms/adapters/twilio" // Import the Twilio adapter
	"github.com/zinzinday/go-sms/model"
)

func main() {
	// Check if config file path is provided as command-line argument
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go /path/to/config.yaml")
	}
	configFile := os.Args[1]

	fmt.Println("=== Simple SMS Example ===")
	fmt.Printf("Using config file: %s\n\n", configFile)

	// Step 1: Initialize the SMS module with configuration
	fmt.Println("1. Initializing SMS module...")
	module, err := sms.NewModule(configFile)
	if err != nil {
		log.Fatalf("Failed to initialize SMS module: %v", err)
	}
	fmt.Println("   SMS module initialized successfully")

	// Step 2: Initialize the Twilio provider
	fmt.Println("2. Initializing Twilio provider...")
	twilioProvider, err := twilio.NewProvider(configFile)
	if err != nil {
		log.Fatalf("Failed to initialize Twilio provider: %v", err)
	}
	fmt.Println("   Twilio provider initialized successfully")

	// Step 3: Add the provider to the module
	fmt.Println("3. Adding Twilio provider to the module...")
	if err := module.AddProvider(twilioProvider); err != nil {
		log.Fatalf("Failed to add Twilio provider: %v", err)
	}
	fmt.Println("   Twilio provider added successfully")

	// Step 4: Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Step 5: Create an SMS request
	fmt.Println("4. Creating SMS request...")
	phoneNumber := getPhoneNumber()
	request := model.SendSMSRequest{
		Message: model.Message{
			// From field can be left empty to use the default from number in config
			From: "",
			To:   phoneNumber,
			By:   "SimpleExample",
		},
		Data: map[string]interface{}{
			"app_name": "SimpleExample",
			"message":  "This is a test message from the go-sms simple example!",
		},
		// Template field is empty, so it will use the default template from config
	}
	fmt.Println("   SMS request created successfully")

	// Step 6: Send the SMS
	fmt.Println("5. Sending SMS...")
	response, err := module.SendSMS(ctx, request)
	if err != nil {
		log.Fatalf("Failed to send SMS: %v", err)
	}

	// Step 7: Display the result
	fmt.Println("\n=== SMS Sent Successfully ===")
	fmt.Printf("Message ID: %s\n", response.MessageID)
	fmt.Printf("Status: %s\n", response.Status)
	fmt.Printf("Provider: %s\n", response.Provider)
	fmt.Printf("Sent at: %v\n", response.SentAt)

	if response.Cost > 0 {
		fmt.Printf("Cost: %v %s\n", response.Cost, response.Currency)
	}
}

// getPhoneNumber prompts the user for a phone number or uses a default one
func getPhoneNumber() string {
	// In a real application, you would get this from user input or your database
	// For this example, let's use a default value that you should replace
	return "+1234567890" // Replace with an actual phone number for testing
}
