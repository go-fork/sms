package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/zinzinday/go-sms"
	"github.com/zinzinday/go-sms/model"
)

func main() {
	// Initialize the SMS module with configuration
	smsModule, err := sms.NewModule("config.yaml")
	if err != nil {
		log.Fatalf("Failed to initialize SMS module: %v", err)
	}

	// In a real application, you would register providers here
	// smsModule.AddProvider(twilio.NewProvider())

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create an SMS request
	request := model.SendSMSRequest{
		Message: model.Message{
			From: "SenderName",
			To:   "+1234567890",
			By:   "ExampleApp",
		},
		Data: map[string]interface{}{
			"app_name": "ExampleApp",
			"message":  "Hello from the SMS module!",
		},
	}

	// Send the SMS
	response, err := smsModule.SendSMS(ctx, request)
	if err != nil {
		log.Fatalf("Failed to send SMS: %v", err)
	}

	fmt.Printf("SMS sent successfully! Message ID: %s, Status: %s\n",
		response.MessageID, response.Status)
}
