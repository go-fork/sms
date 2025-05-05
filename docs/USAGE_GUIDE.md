# go-sms Usage Guide

This comprehensive guide explains how to use the go-sms module effectively in your Go applications.

## Table of Contents

- [Installation](#installation)
- [Configuration](#configuration)
- [Basic Usage](#basic-usage)
- [Template System](#template-system)
- [Working with Multiple Providers](#working-with-multiple-providers)
- [Error Handling](#error-handling)
- [Provider-Specific Features](#provider-specific-features)
- [Best Practices](#best-practices)
- [Advanced Usage](#advanced-usage)
- [Troubleshooting](#troubleshooting)

## Installation

### Requirements
- Go 1.18 or higher
- A valid configuration file
- API credentials for at least one SMS provider

### Installing the Module

```bash
# Install the core module
go get github.com/zinzinday/go-sms

# Install only the provider adapters you need
go get github.com/zinzinday/go-sms/adapters/twilio
go get github.com/zinzinday/go-sms/adapters/esms
go get github.com/zinzinday/go-sms/adapters/speedsms
```

## Configuration

The go-sms module requires a configuration file that defines both global settings and provider-specific configurations.

### Configuration File Structure

```yaml
# Global configuration
default_provider: twilio    # The provider to use by default
http_timeout: 10s           # Timeout for HTTP requests
retry_attempts: 3           # Number of retry attempts for failed requests
retry_delay: 500ms          # Delay between retry attempts (increases exponentially)
sms_template: "Your message from {app_name}: {message}"    # Default SMS template
voice_template: "Your message from {app_name} is {message}"    # Default voice call template

# Provider configurations
providers:
  # Twilio configuration
  twilio:
    account_sid: your_account_sid
    auth_token: your_auth_token
    from_number: +1234567890
    region: us1              # Optional
    api_version: 2010-04-01  # Optional
  
  # eSMS configuration
  esms:
    api_key: your_api_key
    secret: your_secret_key
    brandname: your_brandname  # For branded messages
    sms_type: 2               # 2=Brandname, 4=OTP
  
  # SpeedSMS configuration
  speedsms:
    token: your_access_token
    sender: your_sender_id
    sms_type: 2               # SMS type
```

### Configuration File Location

You'll need to provide the path to your configuration file when initializing the module:

```go
module, err := sms.NewModule("/path/to/your/config.yaml")
```

## Basic Usage

### Initializing the Module

```go
package main

import (
    "github.com/zinzinday/go-sms"
    "github.com/zinzinday/go-sms/adapters/twilio"
)

func main() {
    // Initialize the module with your configuration file
    module, err := sms.NewModule("./config.yaml")
    if err != nil {
        panic(err)
    }
    
    // Initialize a provider
    twilioProvider, err := twilio.NewProvider("./config.yaml")
    if err != nil {
        panic(err)
    }
    
    // Add the provider to the module
    err = module.AddProvider(twilioProvider)
    if err != nil {
        panic(err)
    }
    
    // The module is now ready to use
}
```

### Sending an SMS

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/zinzinday/go-sms"
    "github.com/zinzinday/go-sms/adapters/twilio"
    "github.com/zinzinday/go-sms/model"
)

func main() {
    // Initialize module and provider (see above)
    
    // Create a context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    // Create an SMS request
    request := model.SendSMSRequest{
        Message: model.Message{
            // You can leave From empty to use the default from_number in your config
            From: "",
            To:   "+1234567890", // The recipient's phone number
            By:   "MyApp",       // Your application identifier
        },
        Data: map[string]interface{}{
            "app_name": "MyApp",
            "message":  "Hello from go-sms!",
        },
        // Template is optional - if empty, the default from config will be used
    }
    
    // Send the SMS
    response, err := module.SendSMS(ctx, request)
    if err != nil {
        fmt.Printf("Failed to send SMS: %v\n", err)
        return
    }
    
    // Print the response
    fmt.Printf("SMS sent successfully!\n")
    fmt.Printf("Message ID: %s\n", response.MessageID)
    fmt.Printf("Status: %s\n", response.Status)
    fmt.Printf("Provider: %s\n", response.Provider)
    fmt.Printf("Sent at: %v\n", response.SentAt)
}
```

### Making a Voice Call

```go
// Create a voice call request
voiceRequest := model.SendVoiceRequest{
    Message: model.Message{
        From: "",
        To:   "+1234567890",
        By:   "MyApp",
    },
    Data: map[string]interface{}{
        "app_name": "MyApp",
        "message":  "This is your verification call with code 123456",
    },
    // Voice-specific options can be provided in the Options map
    Options: map[string]interface{}{
        "voice": "female",
        "language": "en-US",
    },
}

// Send the voice call
response, err := module.SendVoiceCall(ctx, voiceRequest)
if err != nil {
    fmt.Printf("Failed to make voice call: %v\n", err)
    return
}

// Print the response
fmt.Printf("Voice call initiated!\n")
fmt.Printf("Call ID: %s\n", response.CallID)
fmt.Printf("Status: %s\n", response.Status)
```

## Template System

The template system allows you to create dynamic message content by replacing placeholders with values.

### Basic Templating

```go
// Define a template with placeholders in {curly_braces}
template := "Hello {name}! Your verification code is {code}."

// Provide the values for the placeholders
data := map[string]interface{}{
    "name": "John",
    "code": "123456",
}

// Create a message
message := model.Message{
    To: "+1234567890",
    By: "MyApp",
}

// Render the template
renderedText := message.Render(template, data)
// Result: "Hello John! Your verification code is 123456."
```

### Using Templates in Requests

```go
// Create a request with a custom template
request := model.SendSMSRequest{
    Message: model.Message{
        To: "+1234567890",
        By: "MyApp",
    },
    // Custom template - overrides the default template in config
    Template: "Hello {name}! Your order #{order_id} has been {status}.",
    Data: map[string]interface{}{
        "name": "Jane",
        "order_id": "12345",
        "status": "shipped",
    },
}

// Send the SMS
response, err := module.SendSMS(ctx, request)
```

### Template Special Variables

The following special variables are always available in templates:

- `{from}` - The sender's identifier (from Message.From)
- `{to}` - The recipient's phone number (from Message.To)
- `{by}` - The application identifier (from Message.By)

## Working with Multiple Providers

The go-sms module allows you to work with multiple SMS providers and switch between them as needed.

### Registering Multiple Providers

```go
// Initialize the module
module, err := sms.NewModule("./config.yaml")
if err != nil {
    panic(err)
}

// Initialize Twilio provider
twilioProvider, err := twilio.NewProvider("./config.yaml")
if err != nil {
    fmt.Printf("Failed to initialize Twilio: %v\n", err)
} else {
    module.AddProvider(twilioProvider)
}

// Initialize eSMS provider
esmsProvider, err := esms.NewProvider("./config.yaml")
if err != nil {
    fmt.Printf("Failed to initialize eSMS: %v\n", err)
} else {
    module.AddProvider(esmsProvider)
}
```

### Switching Between Providers

```go
// Switch to the eSMS provider
err = module.SwitchProvider("esms")
if err != nil {
    fmt.Printf("Failed to switch provider: %v\n", err)
    return
}

// Now the SMS will be sent using eSMS
response, err := module.SendSMS(ctx, request)
```

### Provider Selection Logic

You can implement smart provider selection logic based on message type:

```go
// Select provider based on message type
func selectProvider(module *sms.Module, messageType string) error {
    switch messageType {
    case "otp":
        // For OTP messages, use eSMS for reliability
        return module.SwitchProvider("esms")
    case "marketing":
        // For marketing messages, use SpeedSMS for cost efficiency
        return module.SwitchProvider("speedsms")
    default:
        // For other messages, use Twilio for global reach
        return module.SwitchProvider("twilio")
    }
}
```

### Implementing Failover

You can implement failover between providers:

```go
// Try each provider in order until one succeeds
func sendWithFailover(ctx context.Context, module *sms.Module, request model.SendSMSRequest) (model.SendSMSResponse, error) {
    providers := []string{"twilio", "esms", "speedsms"}
    var lastError error
    
    for _, providerName := range providers {
        err := module.SwitchProvider(providerName)
        if err != nil {
            continue // Provider not available, try next one
        }
        
        response, err := module.SendSMS(ctx, request)
        if err != nil {
            lastError = err
            continue // Failed to send, try next provider
        }
        
        // Success!
        return response, nil
    }
    
    // All providers failed
    return model.SendSMSResponse{}, fmt.Errorf("all providers failed, last error: %w", lastError)
}
```

## Error Handling

The go-sms module uses Go's error handling patterns to report failures.

### Common Errors

```go
response, err := module.SendSMS(ctx, request)
if err != nil {
    // Handle different error types
    switch {
    case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
        fmt.Println("Operation timed out or was canceled")
    case strings.Contains(err.Error(), "no active provider"):
        fmt.Println("No provider is configured or active")
    case strings.Contains(err.Error(), "validation"):
        fmt.Println("Invalid request data:", err)
    case strings.Contains(err.Error(), "rate limit"):
        fmt.Println("Provider rate limit exceeded")
    default:
        fmt.Println("Unexpected error:", err)
    }
    return
}
```

### Retry Logic

The module includes built-in retry logic for transient errors:

```go
// Configure retry settings in your config.yaml
// retry_attempts: 3     # Number of attempts
// retry_delay: 500ms    # Initial delay (increases exponentially)

// These settings apply automatically when sending messages
```

## Best Practices

### Configuration Management

- Keep your configuration file secure (don't commit it to public repositories)
- Use environment variables or a secrets manager for sensitive values
- Create different configuration files for development, testing, and production

### Message Content

- Keep messages concise (SMS has a 160 character limit)
- Include an opt-out option for marketing messages
- Follow regulatory requirements (include sender identification)

### Error Handling

- Always check for errors when sending messages
- Implement graceful degradation when a provider fails
- Log message IDs for troubleshooting

### Performance

- Reuse the module instance for multiple messages
- Set appropriate timeouts for your context
- Consider batch processing for high-volume sending

## Advanced Usage

### Custom Provider Adapters

You can create your own provider adapters by implementing the `model.Provider` interface:

```go
type MyProvider struct {
    // Provider-specific fields
}

func NewProvider(configFile string) (*MyProvider, error) {
    // Initialize provider with configuration
}

func (p *MyProvider) Name() string {
    return "myprovider"
}

func (p *MyProvider) SendSMS(ctx context.Context, req model.SendSMSRequest) (model.SendSMSResponse, error) {
    // Implement SMS sending
}

func (p *MyProvider) SendVoiceCall(ctx context.Context, req model.SendVoiceRequest) (model.SendVoiceResponse, error) {
    // Implement voice calling
}
```

### Custom Template Functions

You can extend the template system by adding custom rendering logic:

```go
// Custom template rendering
func customRender(template string, data map[string]interface{}) string {
    // Your custom rendering logic
    // ...
    return renderedText
}

// Use your custom renderer
request.Template = customRender(myTemplate, myData)
```

## Troubleshooting

### Common Issues

1. **"No active provider set"**
   - Make sure you've successfully added at least one provider
   - Check if you've switched to a provider that doesn't exist

2. **"Failed to load configuration"**
   - Verify the config file path is correct
   - Check if the config file has valid YAML/JSON format
   - Ensure all required fields are present

3. **"Provider X not found"**
   - Make sure you've added the provider
   - Check that the provider name matches exactly
   - Verify the provider adapter is properly initialized

4. **HTTP Errors**
   - Check your network connection
   - Verify your provider credentials
   - Check if provider service is operational
   - Check if you've exceeded rate limits

5. **Message Not Delivered**
   - Verify recipient number format (should be E.164)
   - Check provider dashboard for delivery status
   - Look for error messages in provider responses

### Generating Test Messages

For testing, you can use:

```go
// Test SMS sending with a randomly generated code
func sendTestSMS(module *sms.Module, phoneNumber string) error {
    // Generate random code
    code := fmt.Sprintf("%06d", rand.Intn(1000000))
    
    req := model.SendSMSRequest{
        Message: model.Message{
            To: phoneNumber,
            By: "TestApp",
        },
        Data: map[string]interface{}{
            "code": code,
            "app_name": "TestApp",
        },
    }
    
    _, err := module.SendSMS(context.Background(), req)
    return err
}
```

### Getting Help

If you encounter issues:

1. Check the documentation and examples
2. Look for similar issues on GitHub
3. Create a new issue with detailed information
4. Reach out to the community

## Need More Help?

Visit the [GitHub repository](https://github.com/zinzinday/go-sms) for more examples, documentation, and community support.
