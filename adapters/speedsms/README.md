# SpeedSMS Adapter for go-sms

This adapter enables sending SMS messages through the SpeedSMS API using the go-sms module. SpeedSMS is a popular SMS provider in Vietnam.

## Installation

```bash
go get github.com/go-fork/sms
go get github.com/go-fork/sms/adapters/speedsms
```

## Configuration

Add SpeedSMS configuration to your config file:

```yaml
# Default provider to use (must be registered through AddProvider)
default_provider: speedsms

# HTTP client timeout
http_timeout: 10s

# Retry configuration
retry_attempts: 3
retry_delay: 500ms

# Default templates
sms_template: "Your message from {app_name}: {message}"

# Provider configurations
providers:
  speedsms:
    token: your_speedsms_access_token  # Required
    sender: YourBrandname              # Optional - sender name/ID
    sms_type: 2                        # Optional - defaults to 2 (2=Advertising, 4=OTP, 8=CustomerCare)
    base_url: https://api.speedsms.vn/index.php  # Optional - defaults to standard SpeedSMS API URL
```

## SMS Types

SpeedSMS supports several types of SMS messages:

- **Type 2**: Advertising SMS
- **Type 4**: OTP/Transaction SMS
- **Type 8**: Customer Care SMS

Note: The correct SMS type is important for message delivery as it affects routing and legal compliance.

## Usage

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-fork/sms"
	"github.com/go-fork/sms/adapters/speedsms"
	"github.com/go-fork/sms/model"
)

func main() {
	// Initialize the SMS module
	configFile := "./config.yaml"
	module, err := sms.NewModule(configFile)
	if err != nil {
		log.Fatalf("Failed to initialize SMS module: %v", err)
	}

	// Initialize the SpeedSMS provider
	speedProvider, err := speedsms.NewProvider(configFile)
	if err != nil {
		log.Fatalf("Failed to initialize SpeedSMS provider: %v", err)
	}

	// Add the SpeedSMS provider to the module
	if err := module.AddProvider(speedProvider); err != nil {
		log.Fatalf("Failed to add SpeedSMS provider: %v", err)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Send an SMS
	smsReq := model.SendSMSRequest{
		Message: model.Message{
			From: "YourBrandname",  // Optional sender ID
			To:   "+84123456789",   // Phone number in international format
			By:   "MyApp",
		},
		Data: map[string]interface{}{
			"app_name": "MyApp",
			"message":  "Your verification code is 123456",
		},
		Options: map[string]interface{}{
			"sms_type": 4,  // Override the SMS type for this message (use 4 for OTP)
		},
	}

	smsResp, err := module.SendSMS(ctx, smsReq)
	if err != nil {
		log.Fatalf("Failed to send SMS: %v", err)
	}

	fmt.Printf("SMS sent! Message ID: %s, Status: %s\n", smsResp.MessageID, smsResp.Status)
	
	// Check your balance (SpeedSMS specific feature)
	if speedProvider, ok := module.GetProvider("speedsms").(*speedsms.Provider); ok {
		balance, err := speedProvider.GetBalance(ctx)
		if err != nil {
			log.Printf("Failed to get balance: %v", err)
		} else {
			fmt.Printf("Current balance: %.2f\n", balance)
		}
	}
}
```

## Features

- Send SMS messages through SpeedSMS API
- Support for different SMS types (advertising, OTP, customer care)
- Configurable sender ID
- Check account balance (provider-specific feature)

## Options

### SMS Options

| Option | Description | Example |
|--------|-------------|---------|
| `sms_type` | Type of SMS to send | `2` (advertising), `4` (OTP), `8` (customer care) |

## Limitations

- **Voice Calling**: SpeedSMS does not natively support voice calls, so the `SendVoiceCall` method will return an error.
- **Bulk Messaging**: For bulk messaging, you may need to implement your own batching logic as this adapter sends to one recipient at a time.

## Status Codes

This adapter maps SpeedSMS status codes to go-sms status enums:

| SpeedSMS Status | Description | go-sms Status |
|-----------|-------------|---------------|
| success | Message accepted | StatusSent |
| error | Message failed | StatusFailed |

## Troubleshooting

- **Authentication Errors**: Make sure your token is correct and has sufficient permissions.
- **Message Delivery Issues**: Check that you're using the correct SMS type for your message content.
- **Invalid Phone Numbers**: Ensure phone numbers are in international format (e.g., +84123456789).

## License

MIT
