# Twilio Adapter for go-sms

This adapter enables sending SMS messages and making voice calls through the Twilio API using the go-sms module.

## Installation

```bash
go get github.com/go-fork/sms
go get github.com/go-fork/sms/adapters/twilio
```

## Configuration

Add Twilio configuration to your config file:

```yaml
# Default provider to use (must be registered through AddProvider)
default_provider: twilio

# HTTP client timeout
http_timeout: 10s

# Retry configuration
retry_attempts: 3
retry_delay: 500ms

# Default templates
sms_template: "Your message from {app_name}: {message}"
voice_template: "Your message from {app_name} is {message}"

# Provider configurations
providers:
  twilio:
    account_sid: your_account_sid  # Required
    auth_token: your_auth_token    # Required
    from_number: +1234567890       # Required - must be in E.164 format
    region: us1                    # Optional - defaults to "us1"
    api_version: 2010-04-01        # Optional - defaults to "2010-04-01"
```

## Usage

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-fork/sms"
	"github.com/go-fork/sms/adapters/twilio"
	"github.com/go-fork/sms/model"
)

func main() {
	// Initialize the SMS module
	configFile := "./config.yaml"
	module, err := sms.NewModule(configFile)
	if err != nil {
		log.Fatalf("Failed to initialize SMS module: %v", err)
	}

	// Initialize the Twilio provider
	twilioProvider, err := twilio.NewProvider(configFile)
	if err != nil {
		log.Fatalf("Failed to initialize Twilio provider: %v", err)
	}

	// Add the Twilio provider to the module
	if err := module.AddProvider(twilioProvider); err != nil {
		log.Fatalf("Failed to add Twilio provider: %v", err)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Send an SMS
	smsReq := model.SendSMSRequest{
		Message: model.Message{
			From: "",  // Leave empty to use the default from_number in config
			To:   "+1234567890",
			By:   "MyApp",
		},
		Data: map[string]interface{}{
			"app_name": "MyApp",
			"message":  "Your verification code is 123456",
		},
	}

	smsResp, err := module.SendSMS(ctx, smsReq)
	if err != nil {
		log.Fatalf("Failed to send SMS: %v", err)
	}

	fmt.Printf("SMS sent! Message ID: %s, Status: %s\n", smsResp.MessageID, smsResp.Status)

	// Make a voice call
	voiceReq := model.SendVoiceRequest{
		Message: model.Message{
			From: "",  // Leave empty to use the default from_number in config
			To:   "+1234567890",
			By:   "MyApp",
		},
		Data: map[string]interface{}{
			"app_name": "MyApp",
			"message":  "Your verification code is 123456",
		},
		Options: map[string]interface{}{
			"voice": "woman",      // Optional: "man" or "woman"
			"language": "en-US",   // Optional: language code
		},
	}

	voiceResp, err := module.SendVoiceCall(ctx, voiceReq)
	if err != nil {
		log.Fatalf("Failed to make voice call: %v", err)
	}

	fmt.Printf("Voice call initiated! Call ID: %s, Status: %s\n", voiceResp.CallID, voiceResp.Status)
}
```

## Features

- Send SMS messages through Twilio's API
- Make voice calls using Twilio's TwiML API
- Templates for dynamic message content
- Configurable options like voice type and language
- Full integration with go-sms module retry and configuration systems

## Options

### SMS Options

| Option | Description | Example |
|--------|-------------|---------|
| `status_callback` | URL to receive status updates | `https://example.com/status` |

### Voice Call Options

| Option | Description | Example |
|--------|-------------|---------|
| `voice` | Voice type | `"man"` or `"woman"` |
| `language` | Language code | `"en-US"`, `"fr-FR"`, `"es-ES"` |
| `status_callback` | URL to receive status updates | `https://example.com/status` |

## Status Mapping

This adapter maps Twilio status codes to go-sms status enums:

### SMS Status Mapping

| Twilio Status | go-sms Status |
|---------------|---------------|
| queued | StatusPending |
| sending | StatusPending |
| sent | StatusSent |
| delivered | StatusDelivered |
| undelivered, failed | StatusFailed |
| (others) | StatusUnknown |

### Call Status Mapping

| Twilio Status | go-sms Status |
|---------------|---------------|
| queued | CallStatusQueued |
| initiated | CallStatusInitiated |
| ringing | CallStatusRinging |
| in-progress | CallStatusInProgress |
| completed | CallStatusCompleted |
| busy | CallStatusBusy |
| no-answer, failed | CallStatusNoAnswer |
| canceled | CallStatusCanceled |
| (others) | CallStatusFailed |

## License

MIT
