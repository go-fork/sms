package config

import (
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
