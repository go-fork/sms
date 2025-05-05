# Prompt 8: SpeedSMS Adapter

## Objective
Implement the SpeedSMS provider adapter as a separate Go module.

## Required Files to Create

1. Create a new directory: `/adapters/speedsms/`
2. Inside the SpeedSMS adapter directory:
   - `go.mod` and `go.sum` - Module definition
   - `adapter.go` - SpeedSMS adapter implementation
   - `config.go` - SpeedSMS configuration handling
   - `adapter_test.go` - Unit tests for the adapter
   - `README.md` - Adapter documentation

## Implementation Requirements

### Module Setup
- Initialize a new Go module: `github.com/go-fork/sms/adapters/speedsms`
- Add dependencies:
  - Main SMS module (`github.com/go-fork/sms`)
  - HTTP client (`github.com/go-resty/resty/v2`)

### SpeedSMS Configuration
- Implement in `config.go`:
  - `SpeedSMSConfig` struct with fields:
    - `Token string` - SpeedSMS authentication token
    - Additional SpeedSMS-specific configuration fields
  - `LoadConfig(config *viper.Viper) (*SpeedSMSConfig, error)` - Load SpeedSMS config from Viper

### SpeedSMS Adapter
- Implement in `adapter.go`:
  - `Provider` struct implementing the `model.Provider` interface
  - `NewProvider(configFile string) (*Provider, error)` - Constructor
  - `Name() string` - Return "speedsms" as provider name
  - `SendSMS(ctx context.Context, req model.SendSMSRequest) (model.SendSMSResponse, error)`
  - `SendVoiceCall(ctx context.Context, req model.SendVoiceRequest) (model.SendVoiceResponse, error)`
  
- Implement SpeedSMS API integration:
  - Use SpeedSMS's REST API for sending SMS messages
  - Implement voice call functionality if supported
  - Handle authentication using the token
  - Parse and map response data to module's response models

### Unit Tests
- Implement tests in `adapter_test.go`:
  - Test configuration loading
  - Test adapter creation
  - Test SMS sending (with mocked HTTP responses)
  - Test voice call sending if supported

## Deliverables
- Complete SpeedSMS adapter as a separate Go module
- Configuration handling specific to SpeedSMS
- Implementation of SMS sending via SpeedSMS API
- Voice call implementation if supported by SpeedSMS
- Unit tests for the adapter
- Adapter documentation
