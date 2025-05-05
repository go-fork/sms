package speedsms

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// SpeedSMSConfig holds the configuration for the SpeedSMS provider
type SpeedSMSConfig struct {
	// Token is the SpeedSMS access token
	Token string `mapstructure:"token"`

	// Sender is the sender ID (optional)
	Sender string `mapstructure:"sender"`

	// BaseURL is the SpeedSMS API base URL (optional, default provided)
	BaseURL string `mapstructure:"base_url"`

	// SMSType is the type of SMS (2: Advertising, 4: OTP/Transactional, 8: Customer Care)
	SMSType int `mapstructure:"sms_type"`
}

// LoadConfig loads the SpeedSMS configuration from Viper
func LoadConfig(v *viper.Viper) (*SpeedSMSConfig, error) {
	// Look for providers.speedsms section
	if !v.IsSet("providers.speedsms") {
		return nil, errors.New("speedsms configuration not found in config file")
	}

	// Extract the speedsms section
	speedConfig := v.Sub("providers.speedsms")
	if speedConfig == nil {
		return nil, errors.New("unable to parse speedsms configuration")
	}

	// Set defaults
	speedConfig.SetDefault("base_url", "https://api.speedsms.vn/index.php")
	speedConfig.SetDefault("sms_type", 2) // Default to advertising messages

	// Unmarshal config into struct
	config := &SpeedSMSConfig{}
	if err := speedConfig.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal speedsms config: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// Validate validates the SpeedSMS configuration
func (c *SpeedSMSConfig) Validate() error {
	// Check required fields
	if c.Token == "" {
		return errors.New("speedsms token is required")
	}

	// Validate token format (basic check)
	if len(c.Token) < 20 {
		return errors.New("speedsms token appears to be invalid (too short)")
	}

	// Validate SMS type (2: Advertising, 4: OTP/Transaction, 8: Customer Care)
	validSMSTypes := map[int]bool{2: true, 4: true, 8: true}
	if _, valid := validSMSTypes[c.SMSType]; !valid {
		return errors.New("invalid sms_type value, must be 2, 4, or 8")
	}

	// If base URL is provided, make sure it's a valid URL
	if c.BaseURL != "" && !strings.HasPrefix(c.BaseURL, "http") {
		return errors.New("base_url must start with http:// or https://")
	}

	return nil
}
