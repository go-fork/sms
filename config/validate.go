package config

import (
	"errors"
	"fmt"
)

// ErrMissingDefaultProvider indicates that the default provider is not specified
var ErrMissingDefaultProvider = errors.New("default provider is required")

// ErrNoProvidersConfigured indicates that no providers are configured
var ErrNoProvidersConfigured = errors.New("at least one provider must be configured")

// ErrInvalidHTTPTimeout indicates an invalid HTTP timeout value
var ErrInvalidHTTPTimeout = errors.New("HTTP timeout must be greater than 0")

// ErrInvalidRetryAttempts indicates an invalid retry attempts value
var ErrInvalidRetryAttempts = errors.New("retry attempts must be non-negative")

// ErrInvalidRetryDelay indicates an invalid retry delay value
var ErrInvalidRetryDelay = errors.New("retry delay must be greater than 0")

// ErrMissingSMSTemplate indicates a missing SMS template
var ErrMissingSMSTemplate = errors.New("SMS template is required")

// ErrMissingVoiceTemplate indicates a missing voice template
var ErrMissingVoiceTemplate = errors.New("voice template is required")

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate default provider
	if c.DefaultProvider == "" {
		return ErrMissingDefaultProvider
	}

	// Validate providers
	if c.Providers == nil || len(c.Providers) == 0 {
		return ErrNoProvidersConfigured
	}

	// Verify that the default provider exists in the configured providers
	if _, ok := c.Providers[c.DefaultProvider]; !ok {
		return fmt.Errorf("default provider '%s' not found in configured providers", c.DefaultProvider)
	}

	// Validate HTTP timeout
	if c.HTTPTimeout <= 0 {
		return ErrInvalidHTTPTimeout
	}

	// Validate retry attempts (0 means no retries, which is valid)
	if c.RetryAttempts < 0 {
		return ErrInvalidRetryAttempts
	}

	// Validate retry delay (only if retry attempts > 0)
	if c.RetryAttempts > 0 && c.RetryDelay <= 0 {
		return ErrInvalidRetryDelay
	}

	// Validate SMS template
	if c.SMSTemplate == "" {
		return ErrMissingSMSTemplate
	}

	// Validate voice template
	if c.VoiceTemplate == "" {
		return ErrMissingVoiceTemplate
	}

	return nil
}

// ValidateProviderConfig validates provider-specific configuration
// This is a helper function that providers can use to validate their configurations
func ValidateProviderConfig(config map[string]interface{}, requiredFields ...string) error {
	for _, field := range requiredFields {
		value, exists := config[field]
		if !exists {
			return fmt.Errorf("missing required field: %s", field)
		}

		// Check if string fields are not empty
		if strValue, isString := value.(string); isString && strValue == "" {
			return fmt.Errorf("field %s cannot be empty", field)
		}
	}

	return nil
}
