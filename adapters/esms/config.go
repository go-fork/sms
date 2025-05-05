package esms

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

// ESMSConfig holds the configuration for the eSMS provider
type ESMSConfig struct {
	// APIKey is the eSMS API key
	APIKey string `mapstructure:"api_key"`

	// Secret is the eSMS secret key
	Secret string `mapstructure:"secret"`

	// Brandname is the registered brand name (optional)
	Brandname string `mapstructure:"brandname"`

	// SMSType is the type of SMS to send (default is 2 for branded messages, 4 for OTP messages)
	SMSType int `mapstructure:"sms_type"`

	// BaseURL is the eSMS API base URL (optional, defaults to standard eSMS API URL)
	BaseURL string `mapstructure:"base_url"`
}

// LoadConfig loads the eSMS configuration from Viper
func LoadConfig(v *viper.Viper) (*ESMSConfig, error) {
	// Look for providers.esms section
	if !v.IsSet("providers.esms") {
		return nil, errors.New("esms configuration not found in config file")
	}

	// Extract the esms section
	esmsConfig := v.Sub("providers.esms")
	if esmsConfig == nil {
		return nil, errors.New("unable to parse esms configuration")
	}

	// Set defaults
	esmsConfig.SetDefault("sms_type", 2) // Default to brandname messages
	esmsConfig.SetDefault("base_url", "http://rest.esms.vn/api")

	// Unmarshal config into struct
	config := &ESMSConfig{}
	if err := esmsConfig.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal esms config: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// Validate validates the eSMS configuration
func (c *ESMSConfig) Validate() error {
	// Check required fields
	if c.APIKey == "" {
		return errors.New("esms api_key is required")
	}

	if c.Secret == "" {
		return errors.New("esms secret is required")
	}

	// Validate SMS type (2 for branded messages, 4 for OTP messages, 8 for 8xx messages)
	validSMSTypes := map[int]bool{2: true, 4: true, 8: true}
	if _, valid := validSMSTypes[c.SMSType]; !valid {
		return errors.New("invalid sms_type value, must be 2, 4, or 8")
	}

	// If SMS type is 2 (branded messages) and no brandname is provided, ensure a fallback
	if c.SMSType == 2 && c.Brandname == "" {
		return errors.New("brandname is required for SMS type 2 (branded messages)")
	}

	return nil
}
