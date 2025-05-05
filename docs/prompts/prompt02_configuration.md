# Prompt 2: Configuration Management

## Objective
Implement the configuration management system using Viper to handle module configuration from a provided configuration file.

## Required Files to Create

1. `/config/config.go` - Main configuration handling logic
2. `/config/validate.go` - Configuration validation logic
3. `/config/example.yaml` - Example configuration file

## Implementation Requirements

### Configuration Structure
- Create a `Config` struct in `config/config.go` with fields:
  - `DefaultProvider string` - Name of the default provider to use
  - `HTTPTimeout time.Duration` - Timeout for HTTP requests
  - `RetryAttempts int` - Number of retry attempts
  - `RetryDelay time.Duration` - Delay between retries
  - `SMSTemplate string` - Default template for SMS messages
  - `VoiceTemplate string` - Default template for voice calls
  - `Providers map[string]interface{}` - Provider-specific configurations

### Configuration Loading
- Implement a `LoadConfig(configFile string) (*Config, error)` function using Viper that:
  - Accepts a path to a configuration file (YAML, JSON, etc.)
  - Loads the configuration from the file
  - Populates the `Config` struct
  - Validates the configuration
  - Returns the populated Config or an error

### Configuration Validation
- Implement validation in `config/validate.go`:
  - `Validate() error` method for the Config struct
  - Verify required fields (default_provider, etc.)
  - Validate timeout and retry settings
  - Check for at least one provider configuration

### Example Configuration
- Create an example YAML configuration in `config/example.yaml` with:
  - Top-level configuration (default_provider, http_timeout, etc.)
  - Sample provider configurations for Twilio, eSMS, and SpeedSMS
  - Example templates for SMS and voice messages

## Deliverables
- Complete configuration handling system
- Configuration validation logic
- Example configuration file
