package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/zinzinday/go-sms"
	"github.com/zinzinday/go-sms/adapters/esms"
	"github.com/zinzinday/go-sms/adapters/speedsms"
	"github.com/zinzinday/go-sms/adapters/twilio"
	"github.com/zinzinday/go-sms/model"
)

func main() {
	// Check if config file path is provided as command-line argument
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go /path/to/config.yaml")
	}
	configFile := os.Args[1]

	fmt.Println("=== Multi-Provider SMS Example ===")
	fmt.Printf("Using config file: %s\n\n", configFile)

	// Initialize the SMS module with configuration
	module, err := sms.NewModule(configFile)
	if err != nil {
		log.Fatalf("Failed to initialize SMS module: %v", err)
	}
	fmt.Println("SMS module initialized successfully")

	// Initialize and register multiple providers
	providers := registerProviders(module, configFile)
	if len(providers) == 0 {
		log.Fatal("No providers were successfully registered. Check your configuration.")
	}

	// Get recipient phone number
	phoneNumber := getPhoneNumber()

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Demonstrate different provider usage examples
	sendWithDefaultProvider(ctx, module, phoneNumber)

	// Only try specific providers that were successfully registered
	for _, p := range providers {
		sendWithSpecificProvider(ctx, module, p, phoneNumber)
	}

	if len(providers) >= 2 {
		demonstrateProviderSelection(ctx, module, phoneNumber)
		demonstrateFailoverBetweenProviders(ctx, module, providers, phoneNumber)
	}
}

// registerProviders initializes and registers all available providers
// Returns a list of successfully registered provider names
func registerProviders(module *sms.Module, configFile string) []string {
	var registeredProviders []string

	// Try to initialize each provider
	providers := []struct {
		name     string
		initFunc func(string) (model.Provider, error)
	}{
		{"Twilio", func(cf string) (model.Provider, error) { return twilio.NewProvider(cf) }},
		{"eSMS", func(cf string) (model.Provider, error) { return esms.NewProvider(cf) }},
		{"SpeedSMS", func(cf string) (model.Provider, error) { return speedsms.NewProvider(cf) }},
	}

	for _, p := range providers {
		fmt.Printf("Initializing %s provider...\n", p.name)
		provider, err := p.initFunc(configFile)

		if err != nil {
			fmt.Printf("Failed to initialize %s provider: %v\n", p.name, err)
			continue
		}

		err = module.AddProvider(provider)
		if err != nil {
			fmt.Printf("Failed to add %s provider: %v\n", p.name, err)
			continue
		}

		fmt.Printf("%s provider registered successfully as '%s'\n", p.name, provider.Name())
		registeredProviders = append(registeredProviders, provider.Name())
	}

	fmt.Printf("\nSuccessfully registered %d providers: %v\n\n", len(registeredProviders), registeredProviders)
	return registeredProviders
}

// sendWithDefaultProvider sends an SMS using the default provider
func sendWithDefaultProvider(ctx context.Context, module *sms.Module, phoneNumber string) {
	fmt.Println("\n=== Example 1: Using Default Provider ===")

	// Get the current active provider
	provider, err := module.GetActiveProvider()
	if err != nil {
		fmt.Printf("Error getting active provider: %v\n", err)
		return
	}

	fmt.Printf("Using default provider: %s\n", provider.Name())

	request := model.SendSMSRequest{
		Message: model.Message{
			To: phoneNumber,
			By: "MultiProviderExample",
		},
		Data: map[string]interface{}{
			"app_name": "MultiProviderExample",
			"message":  "This message is sent using the default provider.",
		},
	}

	fmt.Println("Sending SMS...")
	response, err := module.SendSMS(ctx, request)
	if err != nil {
		fmt.Printf("Failed to send SMS: %v\n", err)
		return
	}

	fmt.Println("SMS sent successfully!")
	fmt.Printf("Message ID: %s\n", response.MessageID)
	fmt.Printf("Provider: %s\n", response.Provider)
	fmt.Printf("Status: %s\n", response.Status)
}

// sendWithSpecificProvider sends an SMS using a specific provider
func sendWithSpecificProvider(ctx context.Context, module *sms.Module, providerName string, phoneNumber string) {
	fmt.Printf("\n=== Example 2: Using Specific Provider (%s) ===\n", providerName)

	// Store current provider to restore later
	currentProvider, err := module.GetActiveProvider()
	if err != nil {
		fmt.Printf("Error getting current provider: %v\n", err)
		return
	}

	// Switch to the specified provider
	err = module.SwitchProvider(providerName)
	if err != nil {
		fmt.Printf("Failed to switch to %s provider: %v\n", providerName, err)
		return
	}

	fmt.Printf("Switched to provider: %s\n", providerName)

	request := model.SendSMSRequest{
		Message: model.Message{
			To: phoneNumber,
			By: "MultiProviderExample",
		},
		Data: map[string]interface{}{
			"app_name": "MultiProviderExample",
			"message":  fmt.Sprintf("This message is explicitly sent using the %s provider.", providerName),
		},
	}

	fmt.Println("Sending SMS...")
	response, err := module.SendSMS(ctx, request)
	if err != nil {
		fmt.Printf("Failed to send SMS: %v\n", err)

		// Restore original provider
		module.SwitchProvider(currentProvider.Name())
		return
	}

	fmt.Println("SMS sent successfully!")
	fmt.Printf("Message ID: %s\n", response.MessageID)
	fmt.Printf("Provider: %s\n", response.Provider)
	fmt.Printf("Status: %s\n", response.Status)

	// Restore original provider
	module.SwitchProvider(currentProvider.Name())
	fmt.Printf("Restored original provider: %s\n", currentProvider.Name())
}

// demonstrateProviderSelection shows how to select providers based on message type
func demonstrateProviderSelection(ctx context.Context, module *sms.Module, phoneNumber string) {
	fmt.Println("\n=== Example 3: Provider Selection Based on Message Type ===")

	// Store current provider to restore later
	currentProvider, err := module.GetActiveProvider()
	if err != nil {
		fmt.Printf("Error getting current provider: %v\n", err)
		return
	}

	// Simulate different message types and select appropriate providers
	messageTypes := []struct {
		name              string
		preferredProvider string
		reason            string
		message           string
	}{
		{
			name:              "Promotional",
			preferredProvider: "speedsms",
			reason:            "cost-effective for marketing messages",
			message:           "Check out our summer sale! 20% off all products until June 30.",
		},
		{
			name:              "OTP",
			preferredProvider: "esms",
			reason:            "reliable delivery for critical authentication messages",
			message:           "Your verification code is 123456. Valid for 5 minutes.",
		},
		{
			name:              "Transactional",
			preferredProvider: "twilio",
			reason:            "global reach and delivery receipts for important notices",
			message:           "Your order #12345 has been shipped and will arrive on June 15.",
		},
	}

	for _, mt := range messageTypes {
		fmt.Printf("\nMessage Type: %s\n", mt.name)
		fmt.Printf("Preferred Provider: %s (%s)\n", mt.preferredProvider, mt.reason)

		// Try to switch to preferred provider
		err := module.SwitchProvider(mt.preferredProvider)
		if err != nil {
			// If preferred provider not available, use default
			fmt.Printf("Preferred provider %s not available, using default provider\n", mt.preferredProvider)
		} else {
			fmt.Printf("Using preferred provider: %s\n", mt.preferredProvider)
		}

		// Get current provider (in case we couldn't switch)
		provider, _ := module.GetActiveProvider()

		request := model.SendSMSRequest{
			Message: model.Message{
				To: phoneNumber,
				By: "MultiProviderExample",
			},
			Data: map[string]interface{}{
				"app_name": "MultiProviderExample",
				"message":  mt.message,
			},
		}

		fmt.Printf("Sending %s message...\n", mt.name)
		// Note: In a real application, you would actually send the message
		// For this example, we just show the selection logic
		fmt.Printf("Would send using %s: \"%s\"\n", provider.Name(), mt.message)
	}

	// Restore original provider
	module.SwitchProvider(currentProvider.Name())
	fmt.Printf("\nRestored original provider: %s\n", currentProvider.Name())
}

// demonstrateFailoverBetweenProviders shows how to implement failover between providers
func demonstrateFailoverBetweenProviders(ctx context.Context, module *sms.Module, providers []string, phoneNumber string) {
	fmt.Println("\n=== Example 4: Provider Failover ===")
	fmt.Println("This example demonstrates how to implement failover between providers")

	// Store current provider to restore later
	currentProvider, err := module.GetActiveProvider()
	if err != nil {
		fmt.Printf("Error getting current provider: %v\n", err)
		return
	}

	// Create a message that we'll try to send with multiple providers
	request := model.SendSMSRequest{
		Message: model.Message{
			To: phoneNumber,
			By: "MultiProviderExample",
		},
		Data: map[string]interface{}{
			"app_name": "MultiProviderExample",
			"message":  "This message demonstrates failover between providers.",
		},
	}

	fmt.Println("Simulating provider failover scenario...")
	fmt.Println("In a real application, you would try the next provider when one fails.")
	fmt.Println("Provider preference order:")

	for i, provider := range providers {
		fmt.Printf("%d. %s\n", i+1, provider)
	}

	// Simulate a failover scenario (for demonstration purposes)
	fmt.Println("\nSimulation:")
	fmt.Printf("1. Try to send with %s... (Simulating failure)\n", providers[0])
	fmt.Printf("2. Falling back to %s... (Simulating success)\n", providers[len(providers)-1])

	// Switch to the last provider for demonstration
	err = module.SwitchProvider(providers[len(providers)-1])
	if err != nil {
		fmt.Printf("Error switching provider: %v\n", err)
		return
	}

	// In a real implementation, you would use code like this:
	/*
		var lastError error

		for _, providerName := range providers {
			err := module.SwitchProvider(providerName)
			if err != nil {
				continue // Try next provider
			}

			response, err := module.SendSMS(ctx, request)
			if err != nil {
				lastError = err
				continue // Try next provider
			}

			// Success! We can break out of the loop
			fmt.Printf("SMS sent successfully with %s\n", providerName)
			break
		}

		if lastError != nil {
			fmt.Printf("All providers failed. Last error: %v\n", lastError)
		}
	*/

	// Show what it would look like when successful
	fmt.Println("\nImplementation example for failover logic:")
	fmt.Println("```go")
	fmt.Println("// Try each provider in order until one succeeds")
	fmt.Println("var lastError error")
	fmt.Println("for _, providerName := range providers {")
	fmt.Println("    err := module.SwitchProvider(providerName)")
	fmt.Println("    if err != nil {")
	fmt.Println("        continue // Try next provider")
	fmt.Println("    }")
	fmt.Println("    ")
	fmt.Println("    response, err := module.SendSMS(ctx, request)")
	fmt.Println("    if err != nil {")
	fmt.Println("        lastError = err")
	fmt.Println("        continue // Try next provider")
	fmt.Println("    }")
	fmt.Println("    ")
	fmt.Println("    // Success! We can break out of the loop")
	fmt.Println("    fmt.Printf(\"SMS sent successfully with %s\\n\", providerName)")
	fmt.Println("    break")
	fmt.Println("}")
	fmt.Println("```")

	// Restore original provider
	module.SwitchProvider(currentProvider.Name())
	fmt.Printf("\nRestored original provider: %s\n", currentProvider.Name())
}

// getPhoneNumber prompts the user for a phone number or uses a default one
func getPhoneNumber() string {
	// In a real application, you would get this from user input or your database
	// For this example, let's use a default value that you should replace
	return "+1234567890" // Replace with an actual phone number for testing
}
