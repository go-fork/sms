package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/zinzinday/go-sms"
	"github.com/zinzinday/go-sms/model"
)

// SimpleProvider is a mock provider that implements the Provider interface
type SimpleProvider struct {
	name string
}

// NewSimpleProvider creates a new simple provider
func NewSimpleProvider(name string) *SimpleProvider {
	return &SimpleProvider{name: name}
}

// Name returns the provider name
func (p *SimpleProvider) Name() string {
	return p.name
}

// SendSMS simulates sending an SMS
func (p *SimpleProvider) SendSMS(ctx context.Context, req model.SendSMSRequest) (model.SendSMSResponse, error) {
	fmt.Printf("[%s] Sending SMS to %s: %s\n",
		p.name, req.Message.To, req.Message.Render(req.Template, req.Data))

	// Simulate processing time
	time.Sleep(200 * time.Millisecond)

	return model.SendSMSResponse{
		MessageID: "msg_" + time.Now().Format("20060102150405"),
		Status:    model.StatusSent,
		Provider:  p.name,
		SentAt:    time.Now(),
	}, nil
}

// SendVoiceCall simulates making a voice call
func (p *SimpleProvider) SendVoiceCall(ctx context.Context, req model.SendVoiceRequest) (model.SendVoiceResponse, error) {
	fmt.Printf("[%s] Making voice call to %s: %s\n",
		p.name, req.Message.To, req.Message.Render(req.Template, req.Data))

	// Simulate processing time
	time.Sleep(500 * time.Millisecond)

	return model.SendVoiceResponse{
		CallID:    "call_" + time.Now().Format("20060102150405"),
		Status:    model.CallStatusInitiated,
		Provider:  p.name,
		StartedAt: time.Now(),
		Duration:  0,
	}, nil
}

func main() {
	// Create a temporary configuration file
	configContent := `
default_provider: provider1
http_timeout: 5s
retry_attempts: 3
retry_delay: 500ms
sms_template: "Your message from {app_name}: {message}"
voice_template: "Your message from {app_name} is {message}"

providers:
  provider1:
    api_key: test_key1
  provider2:
    api_key: test_key2
`
	configFile := "temp_config.yaml"

	// Write to temp file (in a real application, you would use an existing config file)
	// This is just for demonstration
	// ... write configContent to configFile ...

	// Initialize the SMS module
	module, err := sms.NewModule(configFile)
	if err != nil {
		log.Fatalf("Failed to initialize SMS module: %v", err)
	}

	// Create and register providers
	provider1 := NewSimpleProvider("provider1")
	provider2 := NewSimpleProvider("provider2")

	// Register providers
	if err := module.AddProvider(provider1); err != nil {
		log.Fatalf("Failed to add provider1: %v", err)
	}

	if err := module.AddProvider(provider2); err != nil {
		log.Fatalf("Failed to add provider2: %v", err)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Send an SMS using the default provider (provider1)
	smsReq := model.SendSMSRequest{
		Message: model.Message{
			From: "Test-Sender",
			To:   "+1234567890",
			By:   "ExampleApp",
		},
		Data: map[string]interface{}{
			"app_name": "ExampleApp",
			"message":  "Hello from the SMS module!",
		},
	}

	smsResp, err := module.SendSMS(ctx, smsReq)
	if err != nil {
		log.Fatalf("Failed to send SMS: %v", err)
	}

	fmt.Printf("SMS sent successfully! Message ID: %s, Status: %s, Provider: %s\n",
		smsResp.MessageID, smsResp.Status, smsResp.Provider)

	// Switch to provider2
	if err := module.SwitchProvider("provider2"); err != nil {
		log.Fatalf("Failed to switch provider: %v", err)
	}

	// Send a voice call using provider2
	voiceReq := model.SendVoiceRequest{
		Message: model.Message{
			From: "VoiceBot",
			To:   "+1234567890",
			By:   "ExampleApp",
		},
		Data: map[string]interface{}{
			"app_name": "ExampleApp",
			"message":  "This is a voice call from the SMS module!",
		},
	}

	voiceResp, err := module.SendVoiceCall(ctx, voiceReq)
	if err != nil {
		log.Fatalf("Failed to send voice call: %v", err)
	}

	fmt.Printf("Voice call initiated! Call ID: %s, Status: %s, Provider: %s\n",
		voiceResp.CallID, voiceResp.Status, voiceResp.Provider)
}
