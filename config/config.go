package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

const (
	// DefaultHTTPTimeout is the default timeout for HTTP requests
	DefaultHTTPTimeout = 10 * time.Second

	// DefaultRetryAttempts is the default number of retry attempts
	DefaultRetryAttempts = 3

	// DefaultRetryDelay is the default delay between retries
	DefaultRetryDelay = 500 * time.Millisecond

	// DefaultSMSTemplate is the default template for SMS messages
	DefaultSMSTemplate = "Your message is {message}"

	// DefaultVoiceTemplate is the default template for voice calls
	DefaultVoiceTemplate = "Your message is {message}"
)

// Error definitions for configuration validation
var (
	// ErrMissingDefaultProvider indicates that the default provider is not specified
	ErrMissingDefaultProvider = errors.New("default provider is required")

	// ErrNoProvidersConfigured indicates that no providers are configured
	ErrNoProvidersConfigured = errors.New("at least one provider must be configured")

	// ErrInvalidHTTPTimeout indicates an invalid HTTP timeout value
	ErrInvalidHTTPTimeout = errors.New("HTTP timeout must be greater than 0")

	// ErrInvalidRetryAttempts indicates an invalid retry attempts value
	ErrInvalidRetryAttempts = errors.New("retry attempts must be non-negative")

	// ErrInvalidRetryDelay indicates an invalid retry delay value
	ErrInvalidRetryDelay = errors.New("retry delay must be greater than 0")

	// ErrMissingSMSTemplate indicates a missing SMS template
	ErrMissingSMSTemplate = errors.New("SMS template is required")

	// ErrMissingVoiceTemplate indicates a missing voice template
	ErrMissingVoiceTemplate = errors.New("voice template is required")
)

// ConfigProvider defines the interface for configuration access
type ConfigProvider interface {
	// Get basic configuration
	GetHTTPTimeout() time.Duration
	GetRetryAttempts() int
	GetRetryDelay() time.Duration
	GetDefaultProvider() string

	// Get template configuration
	GetSMSTemplate() string
	GetVoiceTemplate() string

	// Provider configuration
	GetProviderConfig(providerName string) (map[string]interface{}, error)

	// Validation
	Validate() error
}

// Config represents the module configuration
type Config struct {
	// DefaultProvider is the name of the default provider to use
	DefaultProvider string `mapstructure:"default_provider"`

	// HTTPTimeout is the timeout for HTTP requests
	HTTPTimeout time.Duration `mapstructure:"http_timeout"`

	// RetryAttempts is the number of retry attempts
	RetryAttempts int `mapstructure:"retry_attempts"`

	// RetryDelay is the delay between retries
	RetryDelay time.Duration `mapstructure:"retry_delay"`

	// SMSTemplate is the default template for SMS messages
	SMSTemplate string `mapstructure:"sms_template"`

	// VoiceTemplate is the default template for voice calls
	VoiceTemplate string `mapstructure:"voice_template"`

	// Providers contains provider-specific configurations
	Providers map[string]interface{} `mapstructure:"providers"`
}

// Implement ConfigProvider interface

// GetHTTPTimeout returns the configured HTTP timeout
func (c *Config) GetHTTPTimeout() time.Duration {
	return c.HTTPTimeout
}

// GetRetryAttempts returns the configured retry attempts
func (c *Config) GetRetryAttempts() int {
	return c.RetryAttempts
}

// GetRetryDelay returns the configured retry delay
func (c *Config) GetRetryDelay() time.Duration {
	return c.RetryDelay
}

// GetDefaultProvider returns the name of the default provider
func (c *Config) GetDefaultProvider() string {
	return c.DefaultProvider
}

// GetSMSTemplate returns the configured SMS template
func (c *Config) GetSMSTemplate() string {
	return c.SMSTemplate
}

// GetVoiceTemplate returns the configured voice template
func (c *Config) GetVoiceTemplate() string {
	return c.VoiceTemplate
}

// LoadConfig loads configuration from the specified file path
func LoadConfig(configFile string) (*Config, error) {
	v := viper.New()

	// Set default values
	v.SetDefault("http_timeout", DefaultHTTPTimeout)
	v.SetDefault("retry_attempts", DefaultRetryAttempts)
	v.SetDefault("retry_delay", DefaultRetryDelay)
	v.SetDefault("sms_template", DefaultSMSTemplate)
	v.SetDefault("voice_template", DefaultVoiceTemplate)

	// Set configuration file
	v.SetConfigFile(configFile)

	// Read configuration file
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse duration strings for timeout and retry delay
	httpTimeoutStr := v.GetString("http_timeout")
	if httpTimeout, err := time.ParseDuration(httpTimeoutStr); err == nil {
		v.Set("http_timeout", httpTimeout)
	} else {
		return nil, fmt.Errorf("invalid http_timeout format: %s", httpTimeoutStr)
	}

	retryDelayStr := v.GetString("retry_delay")
	if retryDelay, err := time.ParseDuration(retryDelayStr); err == nil {
		v.Set("retry_delay", retryDelay)
	} else {
		return nil, fmt.Errorf("invalid retry_delay format: %s", retryDelayStr)
	}

	// Unmarshal config into struct
	config := &Config{}
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// GetProviderConfig returns the configuration for a specific provider
func (c *Config) GetProviderConfig(providerName string) (map[string]interface{}, error) {
	if c.Providers == nil {
		return nil, fmt.Errorf("no providers configured")
	}

	providerConfig, ok := c.Providers[providerName]
	if !ok {
		return nil, fmt.Errorf("provider %s not found in configuration", providerName)
	}

	config, ok := providerConfig.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid configuration format for provider %s", providerName)
	}

	return config, nil
}

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
