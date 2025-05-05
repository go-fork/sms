# eSMS Adapter for go-sms

This adapter enables sending SMS messages and voice calls through the eSMS API using the go-sms module. eSMS is a popular SMS provider in Vietnam.

## Installation

```bash
go get github.com/zinzinday/go-sms
go get github.com/zinzinday/go-sms/adapters/esms
```

## Configuration

Add eSMS configuration to your config file:

```yaml
# Default provider to use (must be registered through AddProvider)
default_provider: esms

# HTTP client timeout
http_timeout: 10s

# Retry configuration
retry_attempts: 3
retry_delay: 500ms

# Default templates
sms_template: "Your message from {app_name}: {message}"
voice_template: "Your verification code is {code}"

# Provider configurations
providers:
  esms:
    api_key: your_api_key       # Required
    secret: your_secret_key     # Required
    brandname: YourBrandname    # Required for SMS type 2 (branded messages)
    sms_type: 2                 # Optional - defaults to 2 (2=branded, 4=OTP, 8=8xx)
    base_url: http://rest.esms.vn/api  # Optional - defaults to standard eSMS API URL
```

## SMS Types

eSMS supports several types of SMS messages:

- **Type 2**: Brandname SMS (Requires registered brandname)
- **Type 4**: OTP SMS (One-Time Password)
- **Type 8**: 8xx SMS (Using registered 8xx numbers)

## Usage

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/zinzinday/go-sms"
	"github.com/zinzinday/go-sms/adapters/esms"
	"github.com/zinzinday/go-sms/model"
)

func main() {
	// Initialize the SMS module
	configFile := "./config.yaml"
	module, err := sms.NewModule(configFile)
	if err != nil {
		log.Fatalf("Failed to initialize SMS module: %v", err)
	}

	// Initialize the eSMS provider
	esmsProvider, err := esms.NewProvider(configFile)
	if err != nil {
		log.Fatalf("Failed to initialize eSMS provider: %v", err)
	}

	// Add the eSMS provider to the module
	if err := module.AddProvider(esmsProvider); err != nil {
		log.Fatalf("Failed to add eSMS provider: %v", err)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Send an SMS
	smsReq := model.SendSMSRequest{
		Message: model.Message{
			From: "YourBrandname",  // Use your registered brandname
			To:   "+84123456789",   // Phone number in international format
			By:   "MyApp",
		},
		Data: map[string]interface{}{
			"app_name": "MyApp",
			"message":  "Your verification code is 123456",
		},
		Options: map[string]interface{}{
			"is_unicode": 1,  // 1 for Unicode, 0 for non-Unicode
		},
	}

	smsResp, err := module.SendSMS(ctx, smsReq)
	if err != nil {
		log.Fatalf("Failed to send SMS: %v", err)
	}

	fmt.Printf("SMS sent! Message ID: %s, Status: %s\n", smsResp.MessageID, smsResp.Status)

	// Make a voice call with OTP
	voiceReq := model.SendVoiceRequest{
		Message: model.Message{
			To: "+84123456789",
			By: "MyApp",
		},
		Template: "Your verification code is 123456",
		Options: map[string]interface{}{
			"speed": 1.0,         // Voice speed (0.5 to 2.0)
			"retry_times": 2,     // Number of retry attempts
			"otp": "123456",      // Explicitly specify OTP (optional)
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

- Send SMS messages through eSMS API
- Make voice calls for OTP delivery
- Support for different SMS types (branded, OTP, 8xx)
- Unicode support for Vietnamese and other languages
- Automatic OTP extraction from message content
- Configurable options like message scheduling and voice speed

## Options

### SMS Options

| Option | Description | Example |
|--------|-------------|---------|
| `is_unicode` | Unicode support (1=enable, 0=disable) | `1` |
| `schedule_time` | Schedule message for later delivery | `"2023-12-31 12:00:00"` |

### Voice Call Options

| Option | Description | Example |
|--------|-------------|---------|
| `speed` | Voice speed (0.5 to 2.0) | `1.0` |
| `retry_times` | Number of retry attempts | `2` |
| `otp` | Explicitly specify OTP code | `"123456"` |

## Status Codes

This adapter maps eSMS status codes to go-sms status enums:

| eSMS Code | Description | go-sms Status |
|-----------|-------------|---------------|
| 100 | Success | StatusSent |
| 99 | Invalid parameter | StatusFailed |
| 101-105 | Authentication errors | StatusFailed |
| 106-114 | Account/balance errors | StatusFailed |
| 118-122 | Recipient errors | StatusFailed |
| Other | Unknown status | StatusUnknown |

## Limitations

- Voice calling is only supported for OTP delivery in the eSMS API
- The API will automatically extract the numeric OTP code from your voice message
- Alternatively, you can explicitly provide the OTP code using the `otp` option

## License

MIT
