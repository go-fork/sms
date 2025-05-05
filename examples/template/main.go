package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-fork/sms"
	"github.com/go-fork/sms/adapters/twilio"
	"github.com/go-fork/sms/model"
)

func main() {
	// Check if config file path is provided as command-line argument
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go /path/to/config.yaml")
	}
	configFile := os.Args[1]

	fmt.Println("=== SMS Template Examples ===")

	// Initialize the SMS module
	module, err := sms.NewModule(configFile)
	if err != nil {
		log.Fatalf("Failed to initialize SMS module: %v", err)
	}

	// Initialize and add the Twilio provider
	twilioProvider, err := twilio.NewProvider(configFile)
	if err != nil {
		log.Fatalf("Failed to initialize Twilio provider: %v", err)
	}
	if err := module.AddProvider(twilioProvider); err != nil {
		log.Fatalf("Failed to add Twilio provider: %v", err)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get recipient phone number
	phoneNumber := getPhoneNumber()

	// Demonstrate different template examples
	demonstrateDefaultTemplate(ctx, module, phoneNumber)
	demonstrateCustomTemplate(ctx, module, phoneNumber)
	demonstrateComplexTemplate(ctx, module, phoneNumber)
	demonstratePreviewTemplate(phoneNumber)
}

// demonstrateDefaultTemplate sends an SMS using the default template from config
func demonstrateDefaultTemplate(ctx context.Context, module *sms.Module, phoneNumber string) {
	fmt.Println("\n=== Example 1: Using Default Template ===")
	fmt.Println("This example uses the default template defined in your config file")

	request := model.SendSMSRequest{
		Message: model.Message{
			To: phoneNumber,
			By: "TemplateExample",
		},
		Data: map[string]interface{}{
			"app_name": "TemplateExample",
			"message":  "This message uses the default template from config",
		},
		// Template is empty, so it will use the default from config
	}

	fmt.Println("Sending SMS with default template...")
	response, err := module.SendSMS(ctx, request)
	if err != nil {
		fmt.Printf("Failed to send SMS: %v\n", err)
		return
	}

	fmt.Println("SMS sent successfully!")
	fmt.Printf("Message ID: %s\n", response.MessageID)
	fmt.Printf("Status: %s\n", response.Status)
}

// demonstrateCustomTemplate sends an SMS using a custom template
func demonstrateCustomTemplate(ctx context.Context, module *sms.Module, phoneNumber string) {
	fmt.Println("\n=== Example 2: Using Custom Template ===")
	fmt.Println("This example overrides the default template with a custom one")

	// Define a custom template
	customTemplate := "Hello {name}! Your verification code is {code}. This SMS was sent by {by}."

	request := model.SendSMSRequest{
		Message: model.Message{
			To: phoneNumber,
			By: "TemplateExample",
		},
		Template: customTemplate, // Override the default template
		Data: map[string]interface{}{
			"name": "John Doe",
			"code": "123456",
		},
	}

	// Show what will be sent
	renderedContent := request.Message.Render(customTemplate, request.Data)
	fmt.Printf("Template: %s\n", customTemplate)
	fmt.Printf("Rendered content: %s\n", renderedContent)

	fmt.Println("Sending SMS with custom template...")
	response, err := module.SendSMS(ctx, request)
	if err != nil {
		fmt.Printf("Failed to send SMS: %v\n", err)
		return
	}

	fmt.Println("SMS sent successfully!")
	fmt.Printf("Message ID: %s\n", response.MessageID)
	fmt.Printf("Status: %s\n", response.Status)
}

// demonstrateComplexTemplate sends an SMS using a more complex template
func demonstrateComplexTemplate(ctx context.Context, module *sms.Module, phoneNumber string) {
	fmt.Println("\n=== Example 3: Using Complex Template ===")
	fmt.Println("This example uses a complex template with multiple variables")

	// Define a more complex template
	complexTemplate := "Dear {customer_name}, your {product_name} order #{order_id} has been {status}. " +
		"Expected delivery: {delivery_date}. " +
		"Track your order at: {tracking_url}. " +
		"Thank you for shopping with {company_name}!"

	request := model.SendSMSRequest{
		Message: model.Message{
			To: phoneNumber,
			By: "TemplateExample",
		},
		Template: complexTemplate,
		Data: map[string]interface{}{
			"customer_name": "Jane Smith",
			"product_name":  "Premium Headphones",
			"order_id":      "ORD-12345",
			"status":        "shipped",
			"delivery_date": "June 15, 2023",
			"tracking_url":  "https://track.example.com/ORD-12345",
			"company_name":  "AudioWorld",
		},
	}

	// Show what will be sent
	renderedContent := request.Message.Render(complexTemplate, request.Data)
	fmt.Printf("Template: %s\n", complexTemplate)
	fmt.Printf("Rendered content: %s\n", renderedContent)

	fmt.Println("Sending SMS with complex template...")
	response, err := module.SendSMS(ctx, request)
	if err != nil {
		fmt.Printf("Failed to send SMS: %v\n", err)
		return
	}

	fmt.Println("SMS sent successfully!")
	fmt.Printf("Message ID: %s\n", response.MessageID)
	fmt.Printf("Status: %s\n", response.Status)
}

// demonstratePreviewTemplate shows template rendering without sending an SMS
func demonstratePreviewTemplate(phoneNumber string) {
	fmt.Println("\n=== Example 4: Template Preview (No SMS Sent) ===")
	fmt.Println("This example demonstrates how to preview templates without sending")

	// Create a message
	message := model.Message{
		From: "TemplateDemo",
		To:   phoneNumber,
		By:   "TemplateExample",
	}

	// Define a template with missing variables
	template := "Hello {name}, welcome to {service}! Your account ID is {account_id}."

	// Data with a missing variable
	data := map[string]interface{}{
		"name":    "Alice",
		"service": "Our Amazing Service",
		// Note: account_id is missing!
	}

	// Render the template
	renderedContent := message.Render(template, data)

	fmt.Printf("Template: %s\n", template)
	fmt.Printf("Data: %+v\n", data)
	fmt.Printf("Rendered content: %s\n", renderedContent)
	fmt.Println("Notice that missing variables remain as placeholders in the text")

	// Show how message fields are available in templates
	template2 := "Message from {from} to {to}, sent by application {by}"
	renderedContent2 := message.Render(template2, nil)

	fmt.Printf("\nTemplate with message fields: %s\n", template2)
	fmt.Printf("Rendered content: %s\n", renderedContent2)
	fmt.Println("Message fields (from, to, by) are always available in templates")
}

// getPhoneNumber prompts the user for a phone number or uses a default one
func getPhoneNumber() string {
	// In a real application, you would get this from user input or your database
	// For this example, let's use a default value that you should replace
	return "+1234567890" // Replace with an actual phone number for testing
}
