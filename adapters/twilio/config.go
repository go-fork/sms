package twilio

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// TwilioConfig holds the configuration for the Twilio provider
type TwilioConfig struct {
	// AccountSID is the Twilio account SID
	AccountSID string `mapstructure:"account_sid"`

	// AuthToken is the Twilio authentication token
	AuthToken string `mapstructure:"auth_token"`

	// FromNumber is the default sender phone number
	FromNumber string `mapstructure:"from_number"`

	// Region is the Twilio region (optional, defaults to "us1")
	Region string `mapstructure:"region"`

	// APIVersion is the Twilio API version (optional, defaults to "2010-04-01")
	APIVersion string `mapstructure:"api_version"`
}

// LoadConfig loads the Twilio configuration from Viper
func LoadConfig(v *viper.Viper) (*TwilioConfig, error) {
	// Look for providers.twilio section
	if !v.IsSet("providers.twilio") {
		return nil, errors.New("twilio configuration not found in config file")
	}

	// Extract the twilio section
	twilioConfig := v.Sub("providers.twilio")
	if twilioConfig == nil {
		return nil, errors.New("unable to parse twilio configuration")
	}

	// Set defaults
	twilioConfig.SetDefault("region", "us1")
	twilioConfig.SetDefault("api_version", "2010-04-01")

	// Unmarshal config into struct
	config := &TwilioConfig{}
	if err := twilioConfig.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal twilio config: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// Validate validates the Twilio configuration
func (c *TwilioConfig) Validate() error {
	// Check required fields
	if c.AccountSID == "" {
		return errors.New("twilio account_sid is required")
	}

	if c.AuthToken == "" {
		return errors.New("twilio auth_token is required")
	}

	if c.FromNumber == "" {
		return errors.New("twilio from_number is required")
	}

	// Validate the format of the AccountSID (should start with "AC")
	if !strings.HasPrefix(c.AccountSID, "AC") {
		return errors.New("invalid twilio account_sid format (should start with 'AC')")
	}

	// Validate the phone number format (basic check)
	// E.164 format: +country code followed by number
	if !strings.HasPrefix(c.FromNumber, "+") {
		return errors.New("from_number must be in E.164 format (e.g., +1234567890)")
	}

	return nil
}
