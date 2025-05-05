# Prompt 1: 
Module Structure

## Objective
Set up the basic structure for the Go SMS module including directory layout, core interfaces, and placeholder files.

## Required Files and Directories to Create

1. Create the main module directory structure:
   - `/model` - For message models and provider interfaces
   - `/config` - For configuration management
   - `/client` - For core client functionality
   - `/retry` - For retry logic
   - `/tests` - For unit tests
   - `/examples` - For usage examples

2. Create the basic Go module files:
   - `go.mod` - Initialize with required dependencies
   - `go.sum` - Will be generated automatically
   - `sms.go` - Main module entry point

3. Create the provider interface and message models:
   - `/model/provider.go` - Interface for SMS/Voice providers
   - `/model/message.go` - Basic message model structure
   - `/model/request.go` - Request structures
   - `/model/response.go` - Response structures

## Implementation Requirements

### Module Configuration
- Initialize Go module with name `github.com/go-fork/sms`
- Add initial dependencies:
  - `github.com/spf13/viper` (v1.15.0+)
  - `github.com/go-resty/resty/v2` (v2.7.0+)

### Provider Interface
- Define the `Provider` interface in `model/provider.go` with:
  - `Name() string` - Returns provider name
  - `SendSMS(context.Context, SendSMSRequest) (SendSMSResponse, error)` 
  - `SendVoiceCall(context.Context, SendVoiceRequest) (SendVoiceResponse, error)`

### Module Entry Point
- Create a basic structure in `sms.go` for the module with:
  - `type Module struct` - Main module container
  - `NewModule(configFile string) (*Module, error)` - Constructor
  - `AddProvider(provider Provider) error` - Method to register providers
  - `SwitchProvider(name string) error` - Method to change active provider
  - `SendSMS` and `SendVoiceCall` methods that delegate to the active provider

## Deliverables
- Basic directory structure
- Initialized Go module with dependencies
- Core interfaces and structures
- Basic module implementation
