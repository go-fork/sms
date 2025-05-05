# Prompt 5: Core SMS Module

## Objective
Implement the core SMS module that manages providers and handles message sending.

## Required Files to Create/Complete

1. `/sms.go` - Main module entry point and provider management
2. `/client/client.go` - Core client implementation

## Implementation Requirements

### SMS Module
- Complete the `Module` struct in `sms.go` with:
  - `config *config.Config` - Module configuration
  - `providers map[string]model.Provider` - Map of registered providers
  - `activeProvider model.Provider` - Currently active provider

- Implement provider management methods:
  - `NewModule(configFile string) (*Module, error)` - Constructor that loads configuration
  - `AddProvider(provider model.Provider) error` - Register a new provider
  - `SwitchProvider(name string) error` - Change the active provider
  - `GetProvider(name string) (model.Provider, error)` - Get a provider by name

- Implement message sending methods:
  - `SendSMS(ctx context.Context, req model.SendSMSRequest) (model.SendSMSResponse, error)`
  - `SendVoiceCall(ctx context.Context, req model.SendVoiceRequest) (model.SendVoiceResponse, error)`
  - These methods should delegate to the active provider and apply retry logic

### Client Implementation
- Implement a core HTTP client in `client/client.go`:
  - `NewClient(config *config.Config) *Client` - Create a new HTTP client with configuration
  - Configure timeouts, retries, and other HTTP client settings
  - Use `github.com/go-resty/resty/v2` for HTTP requests

## Deliverables
- Complete SMS module with provider management
- Message sending implementation with retry logic
- HTTP client configuration for provider adapters
