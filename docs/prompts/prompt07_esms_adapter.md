# Prompt 7: eSMS Adapter

## Objective
Implement the eSMS provider adapter as a separate Go module.

## Required Files to Create

1. Create a new directory: `/adapters/esms/`
2. Inside the eSMS adapter directory:
   - `go.mod` and `go.sum` - Module definition
   - `adapter.go` - eSMS adapter implementation
   - `config.go` - eSMS configuration handling
   - `adapter_test.go` - Unit tests for the adapter
   - `README.md` - Adapter documentation

## Implementation Requirements

### Module Setup
- Initialize a new Go module: `github.com/zinzinday/go-sms/adapters/esms`
- Add dependencies:
  - Main SMS module (`github.com/zinzinday/go-sms`)
  - HTTP client (`github.com/go-resty/resty/v2`)

### eSMS Configuration
- Implement in `config.go`:
  - `ESMSConfig` struct with fields:
    - `APIKey string` - eSMS API key
    - `Secret string` - eSMS secret key
    - Additional eSMS-specific configuration fields
  - `LoadConfig(config *viper.Viper) (*ESMSConfig, error)` - Load eSMS config from Viper

### eSMS Adapter
- Implement in `adapter.go`:
  - `Provider` struct implementing the `model.Provider` interface
  - `NewProvider(configFile string) (*Provider, error)` - Constructor
  - `Name() string` - Return "esms" as provider name
  - `SendSMS(ctx context.Context, req model.SendSMSRequest) (model.SendSMSResponse, error)`
  - `SendVoiceCall(ctx context.Context, req model.SendVoiceRequest) (model.SendVoiceResponse, error)`
  
- Implement eSMS API integration:
  - Use eSMS's REST API for sending SMS messages
  - Implement voice call functionality if supported
  - Handle authentication and request formatting
  - Parse and map response data to module's response models

### Unit Tests
- Implement tests in `adapter_test.go`:
  - Test configuration loading
  - Test adapter creation
  - Test SMS sending (with mocked HTTP responses)
  - Test voice call sending if supported

## Deliverables
- Complete eSMS adapter as a separate Go module
- Configuration handling specific to eSMS
- Implementation of SMS sending via eSMS API
- Voice call implementation if supported by eSMS
- Unit tests for the adapter
- Adapter documentation
