# go-sms

A flexible and extensible Go module for sending SMS messages and voice calls through multiple providers with a unified API.

## Overview

go-sms is a service component library that enables Go applications to send SMS messages and voice calls through various providers using a standardized interface. Instead of integrating with multiple SMS provider APIs directly, you can use go-sms as an abstraction layer to simplify your code and make it more maintainable.

### Key Features

- **Unified API**: Send SMS and voice calls through any supported provider with the same API
- **Multiple Providers**: Support for various SMS providers including Twilio, eSMS, and SpeedSMS
- **Provider Management**: Easily switch between providers at runtime
- **Message Templates**: Dynamic message content with template variable substitution
- **Configuration Management**: Simple YAML-based configuration with validation
- **Retry Mechanism**: Built-in retry logic with exponential backoff
- **Extensible Architecture**: Easily add new provider adapters

## Installation

### Core Module

```bash
go get github.com/go-fork/sms
```

### Provider Adapters

Install only the provider adapters you need:

```bash
# For Twilio
go get github.com/go-fork/sms/adapters/twilio

# For eSMS (Vietnam)
go get github.com/go-fork/sms/adapters/esms

# For SpeedSMS (Vietnam)
go get github.com/go-fork/sms/adapters/speedsms
```

## Quick Start

### 1. Create a Configuration File

Create a `config.yaml` file with your provider credentials:

```yaml
# Default provider to use
default_provider: twilio

# HTTP client configuration
http_timeout: 10s

# Retry configuration
retry_attempts: 3
retry_delay: 500ms

# Default message templates
sms_template: "Your message from {app_name}: {message}"
voice_template: "Your message from {app_name} is {message}"

# Provider configurations
providers:
  twilio:
    account_sid: your_account_sid
    auth_token: your_auth_token
    from_number: +1234567890
  esms:
    api_key: your_api_key
    secret: your_secret_key
    brandname: your_brandname
  speedsms:
    token: your_access_token
    sender: your_sender_id
```

### 2. Initialize the Module and Send an SMS

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
	// Initialize the SMS module with configuration
	module, err := sms.NewModule("./config.yaml")
	if err != nil {
		log.Fatalf("Failed to initialize SMS module: %v", err)
	}

	// Initialize and add the Twilio provider
	twilioProvider, err := twilio.NewProvider("./config.yaml")
	if err != nil {
		log.Fatalf("Failed to initialize Twilio provider: %v", err)
	}
	err = module.AddProvider(twilioProvider)
	if err != nil {
		log.Fatalf("Failed to add Twilio provider: %v", err)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create an SMS request
	request := model.SendSMSRequest{
		Message: model.Message{
			From: "", // Use default from number in config
			To:   "+1234567890", // Recipient's phone number
			By:   "MyApp", // Your application's name
		},
		Data: map[string]interface{}{
			"app_name": "MyApp",
			"message":  "Hello from go-sms!",
		},
	}

	// Send the SMS
	response, err := module.SendSMS(ctx, request)
	if err != nil {
		log.Fatalf("Failed to send SMS: %v", err)
	}

	fmt.Printf("SMS sent successfully! Message ID: %s\n", response.MessageID)
}
```

## Configuration Options

### Top-Level Configuration

| Option | Description | Default | Example |
|--------|-------------|---------|---------|
| `default_provider` | Name of the default provider to use | | `"twilio"` |
| `http_timeout` | Timeout for HTTP requests | `10s` | `"30s"` |
| `retry_attempts` | Number of retry attempts | `3` | `5` |
| `retry_delay` | Initial delay between retries | `500ms` | `"1s"` |
| `sms_template` | Default template for SMS messages | `"Your message is {message}"` | `"Message from {app_name}: {message}"` |
| `voice_template` | Default template for voice calls | `"Your message is {message}"` | `"Message from {app_name}: {message}"` |

### Provider-Specific Configuration

#### Twilio

```yaml
providers:
  twilio:
    account_sid: your_account_sid  # Required
    auth_token: your_auth_token    # Required
    from_number: +1234567890       # Required - must be in E.164 format
    region: us1                    # Optional - defaults to "us1"
    api_version: 2010-04-01        # Optional - defaults to "2010-04-01"
```

#### eSMS (Vietnam)

```yaml
providers:
  esms:
    api_key: your_api_key          # Required
    secret: your_secret_key        # Required
    brandname: your_brandname      # Required for branded messages
    sms_type: 2                    # Optional - defaults to 2 (2=branded, 4=OTP)
    base_url: http://rest.esms.vn/api  # Optional
```

#### SpeedSMS (Vietnam)

```yaml
providers:
  speedsms:
    token: your_access_token       # Required
    sender: your_sender_id         # Optional
    sms_type: 2                    # Optional - defaults to 2
    base_url: https://api.speedsms.vn/index.php  # Optional
```

## Message Model and Templates

The message model represents the core information needed to send a message:

```go
type Message struct {
	From string // Sender identifier (phone number, brandname)
	To   string // Recipient's phone number
	By   string // Application identifier (e.g., "MyApp")
}
```

### Template Variables

You can use the following variables in your message templates:

- `{from}`: The sender's identifier
- `{to}`: The recipient's phone number
- `{by}`: The application identifier
- Any custom variables provided in the `Data` map

### Example Templates

```yaml
# Basic template
sms_template: "Your message is {message}"

# More detailed template
sms_template: "Hello from {app_name}! Your verification code is {code}."

# Template with transaction information
sms_template: "Your order #{order_id} has been confirmed. Total: ${amount}."
```

## API Reference

### Module Initialization

```go
func NewModule(configFile string) (*Module, error)
```

### Provider Management

```go
func (m *Module) AddProvider(provider model.Provider) error
func (m *Module) SwitchProvider(name string) error
func (m *Module) GetProvider(name string) (model.Provider, error)
func (m *Module) GetActiveProvider() (model.Provider, error)
```

### Sending Messages

```go
func (m *Module) SendSMS(ctx context.Context, req model.SendSMSRequest) (model.SendSMSResponse, error)
func (m *Module) SendVoiceCall(ctx context.Context, req model.SendVoiceRequest) (model.SendVoiceResponse, error)
```

### Request Structures

```go
type SendSMSRequest struct {
	Message  model.Message
	Template string // Optional - overrides config template
	Data     map[string]interface{}
	Options  map[string]interface{} // Provider-specific options
}

type SendVoiceRequest struct {
	Message  model.Message
	Template string // Optional - overrides config template
	Data     map[string]interface{}
	Options  map[string]interface{} // Provider-specific options
}
```

### Response Structures

```go
type SendSMSResponse struct {
	MessageID        string
	Status           MessageStatus
	Provider         string
	SentAt           time.Time
	Cost             float64 // Optional
	Currency         string  // Optional
	ProviderResponse map[string]interface{}
}

type SendVoiceResponse struct {
	CallID           string
	Status           CallStatus
	Provider         string
	StartedAt        time.Time
	EndedAt          *time.Time // Optional
	Duration         int        // In seconds
	Cost             float64    // Optional
	Currency         string     // Optional
	ProviderResponse map[string]interface{}
}
```

## Examples

See the `/examples` directory for more comprehensive examples:

- `examples/simple/main.go`: Basic usage example
- `examples/template/main.go`: Examples of using templates
- `examples/multi-provider/main.go`: Working with multiple providers

## Performance Benchmarks

The go-sms module is designed for high performance and reliability:

| Operation | Average Latency | Throughput (msgs/sec) |
|-----------|----------------|----------------------|
| SMS Sending | ~200-500ms | Up to 100 with rate limiting |
| Voice Call Initiation | ~300-700ms | Up to 20 with rate limiting |
| Template Rendering | <1ms | >10,000 |
| Provider Switching | <1ms | >5,000 |

Benchmarks performed on standard cloud infrastructure. Your results may vary depending on network conditions and provider responsiveness.

## Integration Examples

### With Web Frameworks

#### Using with Gin

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/go-fork/sms"
    "github.com/go-fork/sms/adapters/twilio"
)

func setupRouter(smsModule *sms.Module) *gin.Engine {
    r := gin.Default()
    
    r.POST("/send-sms", func(c *gin.Context) {
        var req struct {
            PhoneNumber string `json:"phone_number" binding:"required"`
            Message     string `json:"message" binding:"required"`
        }
        
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(400, gin.H{"error": err.Error()})
            return
        }
        
        smsRequest := model.SendSMSRequest{
            Message: model.Message{
                To: req.PhoneNumber,
                By: "WebApp",
            },
            Data: map[string]interface{}{
                "message": req.Message,
            },
        }
        
        resp, err := smsModule.SendSMS(c.Request.Context(), smsRequest)
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        
        c.JSON(200, gin.H{
            "message_id": resp.MessageID,
            "status": resp.Status,
        })
    })
    
    return r
}
```

#### Using with Echo

```go
import (
    "github.com/labstack/echo/v4"
    "github.com/go-fork/sms"
)

func setupRoutes(e *echo.Echo, smsModule *sms.Module) {
    e.POST("/send-sms", func(c echo.Context) error {
        // Implementation similar to Gin example
    })
}
```

### With GORM for Notifications

```go
import (
    "gorm.io/gorm"
    "github.com/go-fork/sms"
)

type Notification struct {
    gorm.Model
    UserID      uint
    PhoneNumber string
    Message     string
    Status      string
    MessageID   string
}

func SendPendingNotifications(db *gorm.DB, smsModule *sms.Module) error {
    var notifications []Notification
    
    // Find pending notifications
    result := db.Where("status = ?", "pending").Find(&notifications)
    if result.Error != nil {
        return result.Error
    }
    
    for _, notification := range notifications {
        // Send SMS
        req := model.SendSMSRequest{
            Message: model.Message{
                To: notification.PhoneNumber,
                By: "DatabaseNotifier",
            },
            Data: map[string]interface{}{
                "message": notification.Message,
            },
        }
        
        resp, err := smsModule.SendSMS(context.Background(), req)
        
        // Update notification status
        if err != nil {
            db.Model(&notification).Updates(map[string]interface{}{
                "status": "failed",
            })
        } else {
            db.Model(&notification).Updates(map[string]interface{}{
                "status": string(resp.Status),
                "message_id": resp.MessageID,
            })
        }
    }
    
    return nil
}
```

## Comparison with Other Solutions

| Feature | go-sms | aws-sdk-go (SNS) | messagebird/go-rest-api | twilio-go |
|---------|--------|-----------------|------------------------|-----------|
| Multiple Providers | ✅ | ❌ | ❌ | ❌ |
| Template System | ✅ | ❌ | ❌ | ✅ |
| Voice Calls | ✅ | ❌ | ✅ | ✅ |
| Provider Switching | ✅ | ❌ | ❌ | ❌ |
| Retry Logic | ✅ | ✅ | ❌ | ❌ |
| Configuration Management | ✅ | ❌ | ❌ | ❌ |
| Vietnamese Providers | ✅ | ❌ | ❌ | ❌ |

## Community and Support

- **GitHub Issues**: For bug reports and feature requests
- **Discussions**: For questions and community support
- **Contributing**: See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines

## Roadmap

Future plans for go-sms include:

- **Additional Providers**: Implementing more international and regional SMS providers
- **Advanced Analytics**: Message delivery tracking and reporting
- **Batch Sending**: Optimized bulk messaging capabilities
- **Scheduling**: Time-delayed message delivery
- **Webhooks**: Support for delivery status callbacks
- **Admin Dashboard**: Web interface for managing and monitoring messages

## License

This project is licensed under the MIT License - see the LICENSE file for details.
