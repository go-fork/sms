# Prompt 6: Twilio Adapter

## Objective
Implement the Twilio provider adapter as a separate Go module.

## Required Files to Create

1. Create a new directory: `/adapters/twilio/`
2. Inside the Twilio adapter directory:
   - `go.mod` and `go.sum` - Module definition
   - `adapter.go` - Twilio adapter implementation
   - `config.go` - Twilio configuration handling
   - `adapter_test.go` - Unit tests for the adapter
   - `README.md` - Adapter documentation

## Implementation Requirements

### Module Setup
- Initialize a new Go module: `github.com/go-fork/sms/adapters/twilio`
- Add dependencies:
  - Main SMS module (`github.com/go-fork/sms`)
  - HTTP client (`github.com/go-resty/resty/v2`)

### Twilio Configuration
- Implement in `config.go`:
  - `TwilioConfig` struct with fields:
    - `AccountSID string` - Twilio account SID
    - `AuthToken string` - Twilio authentication token
    - `FromNumber string` - Default sender phone number
  - `LoadConfig(config *viper.Viper) (*TwilioConfig, error)` - Load Twilio config from Viper

### Twilio Adapter
- Implement in `adapter.go`:
  - `Provider` struct implementing the `model.Provider` interface
  - `NewProvider(configFile string) (*Provider, error)` - Constructor
  - `Name() string` - Return "twilio" as provider name
  - `SendSMS(ctx context.Context, req model.SendSMSRequest) (model.SendSMSResponse, error)`
  - `SendVoiceCall(ctx context.Context, req model.SendVoiceRequest) (model.SendVoiceResponse, error)`
  
- Implement Twilio API integration:
  - Use Twilio's REST API for sending SMS messages
  - Use Twilio's TwiML API for voice calls
  - Handle authentication and request formatting
  - Parse and map response data to module's response models

### Unit Tests
- Implement tests in `adapter_test.go`:
  - Test configuration loading
  - Test adapter creation
  - Test SMS and voice call sending (with mocked HTTP responses)

## Deliverables
- Complete Twilio adapter as a separate Go module
- Configuration handling specific to Twilio
- Implementation of SMS and voice call sending via Twilio API
- Unit tests for the adapter
- Adapter documentation
