package tests

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zinzinday/go-sms/config"
)

// TestLoadConfig tests loading configurations from various formats
func TestLoadConfig(t *testing.T) {
	// Test loading from YAML
	yamlConfig := `
default_provider: test_provider
http_timeout: 10s
retry_attempts: 3
retry_delay: 500ms
sms_template: "Your message is {message}"
voice_template: "Your message is {message}"

providers:
  test_provider:
    api_key: test_key
    secret: test_secret
  another_provider:
    token: test_token
`
	yamlFile, err := createTempConfig(yamlConfig)
	require.NoError(t, err)
	defer os.Remove(yamlFile)

	cfg, err := config.LoadConfig(yamlFile)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "test_provider", cfg.DefaultProvider)
	assert.Equal(t, 10*time.Second, cfg.HTTPTimeout)
	assert.Equal(t, 3, cfg.RetryAttempts)
	assert.Equal(t, 500*time.Millisecond, cfg.RetryDelay)
	assert.Equal(t, "Your message is {message}", cfg.SMSTemplate)
	assert.Equal(t, "Your message is {message}", cfg.VoiceTemplate)
	assert.NotNil(t, cfg.Providers)
	assert.Contains(t, cfg.Providers, "test_provider")
	assert.Contains(t, cfg.Providers, "another_provider")

	// Test loading from non-existent file
	_, err = config.LoadConfig("non-existent-file.yaml")
	assert.Error(t, err)

	// Test loading from invalid YAML format
	invalidYaml := `
default_provider: test_provider
http_timeout: 10s
retry_attempts: "not a number" # This should be a number
`
	invalidFile, err := createTempConfig(invalidYaml)
	require.NoError(t, err)
	defer os.Remove(invalidFile)

	_, err = config.LoadConfig(invalidFile)
	assert.Error(t, err)

	// Test loading from JSON
	jsonConfig := `{
		"default_provider": "test_provider",
		"http_timeout": "10s",
		"retry_attempts": 3,
		"retry_delay": "500ms",
		"sms_template": "Your message is {message}",
		"voice_template": "Your message is {message}",
		"providers": {
			"test_provider": {
				"api_key": "test_key",
				"secret": "test_secret"
			}
		}
	}`
	jsonFile, err := os.CreateTemp("", "sms-config-*.json")
	require.NoError(t, err)
	defer os.Remove(jsonFile.Name())
	_, err = jsonFile.Write([]byte(jsonConfig))
	require.NoError(t, err)
	jsonFile.Close()

	cfg, err = config.LoadConfig(jsonFile.Name())
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "test_provider", cfg.DefaultProvider)
}

// TestValidateConfig tests the configuration validation
func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      *config.Config
		expectError bool
	}{
		{
			name: "Valid configuration",
			config: &config.Config{
				DefaultProvider: "test_provider",
				HTTPTimeout:     10 * time.Second,
				RetryAttempts:   3,
				RetryDelay:      500 * time.Millisecond,
				SMSTemplate:     "Your message is {message}",
				VoiceTemplate:   "Your message is {message}",
				Providers: map[string]interface{}{
					"test_provider": map[string]interface{}{
						"api_key": "test_key",
					},
				},
			},
			expectError: false,
		},
		{
			name: "Missing default provider",
			config: &config.Config{
				DefaultProvider: "",
				HTTPTimeout:     10 * time.Second,
				RetryAttempts:   3,
				RetryDelay:      500 * time.Millisecond,
				SMSTemplate:     "Your message is {message}",
				VoiceTemplate:   "Your message is {message}",
				Providers: map[string]interface{}{
					"test_provider": map[string]interface{}{
						"api_key": "test_key",
					},
				},
			},
			expectError: true,
		},
		{
			name: "Default provider not in providers",
			config: &config.Config{
				DefaultProvider: "non_existent",
				HTTPTimeout:     10 * time.Second,
				RetryAttempts:   3,
				RetryDelay:      500 * time.Millisecond,
				SMSTemplate:     "Your message is {message}",
				VoiceTemplate:   "Your message is {message}",
				Providers: map[string]interface{}{
					"test_provider": map[string]interface{}{
						"api_key": "test_key",
					},
				},
			},
			expectError: true,
		},
		{
			name: "Invalid HTTP timeout",
			config: &config.Config{
				DefaultProvider: "test_provider",
				HTTPTimeout:     0, // Invalid timeout
				RetryAttempts:   3,
				RetryDelay:      500 * time.Millisecond,
				SMSTemplate:     "Your message is {message}",
				VoiceTemplate:   "Your message is {message}",
				Providers: map[string]interface{}{
					"test_provider": map[string]interface{}{
						"api_key": "test_key",
					},
				},
			},
			expectError: true,
		},
		{
			name: "Invalid retry attempts",
			config: &config.Config{
				DefaultProvider: "test_provider",
				HTTPTimeout:     10 * time.Second,
				RetryAttempts:   -1, // Invalid retry attempts
				RetryDelay:      500 * time.Millisecond,
				SMSTemplate:     "Your message is {message}",
				VoiceTemplate:   "Your message is {message}",
				Providers: map[string]interface{}{
					"test_provider": map[string]interface{}{
						"api_key": "test_key",
					},
				},
			},
			expectError: true,
		},
		{
			name: "Invalid retry delay with retry attempts > 0",
			config: &config.Config{
				DefaultProvider: "test_provider",
				HTTPTimeout:     10 * time.Second,
				RetryAttempts:   3,
				RetryDelay:      0, // Invalid retry delay when retry attempts > 0
				SMSTemplate:     "Your message is {message}",
				VoiceTemplate:   "Your message is {message}",
				Providers: map[string]interface{}{
					"test_provider": map[string]interface{}{
						"api_key": "test_key",
					},
				},
			},
			expectError: true,
		},
		{
			name: "Valid zero retry attempts with zero retry delay",
			config: &config.Config{
				DefaultProvider: "test_provider",
				HTTPTimeout:     10 * time.Second,
				RetryAttempts:   0, // No retries
				RetryDelay:      0, // Valid when no retries
				SMSTemplate:     "Your message is {message}",
				VoiceTemplate:   "Your message is {message}",
				Providers: map[string]interface{}{
					"test_provider": map[string]interface{}{
						"api_key": "test_key",
					},
				},
			},
			expectError: false,
		},
		{
			name: "Missing SMS template",
			config: &config.Config{
				DefaultProvider: "test_provider",
				HTTPTimeout:     10 * time.Second,
				RetryAttempts:   3,
				RetryDelay:      500 * time.Millisecond,
				SMSTemplate:     "", // Missing SMS template
				VoiceTemplate:   "Your message is {message}",
				Providers: map[string]interface{}{
					"test_provider": map[string]interface{}{
						"api_key": "test_key",
					},
				},
			},
			expectError: true,
		},
		{
			name: "Missing voice template",
			config: &config.Config{
				DefaultProvider: "test_provider",
				HTTPTimeout:     10 * time.Second,
				RetryAttempts:   3,
				RetryDelay:      500 * time.Millisecond,
				SMSTemplate:     "Your message is {message}",
				VoiceTemplate:   "", // Missing voice template
				Providers: map[string]interface{}{
					"test_provider": map[string]interface{}{
						"api_key": "test_key",
					},
				},
			},
			expectError: true,
		},
		{
			name: "No providers",
			config: &config.Config{
				DefaultProvider: "test_provider",
				HTTPTimeout:     10 * time.Second,
				RetryAttempts:   3,
				RetryDelay:      500 * time.Millisecond,
				SMSTemplate:     "Your message is {message}",
				VoiceTemplate:   "Your message is {message}",
				Providers:       nil, // No providers
			},
			expectError: true,
		},
		{
			name: "Empty providers map",
			config: &config.Config{
				DefaultProvider: "test_provider",
				HTTPTimeout:     10 * time.Second,
				RetryAttempts:   3,
				RetryDelay:      500 * time.Millisecond,
				SMSTemplate:     "Your message is {message}",
				VoiceTemplate:   "Your message is {message}",
				Providers:       map[string]interface{}{}, // Empty providers map
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestGetProviderConfig tests retrieving provider-specific configuration
func TestGetProviderConfig(t *testing.T) {
	cfg := &config.Config{
		DefaultProvider: "test_provider",
		Providers: map[string]interface{}{
			"test_provider": map[string]interface{}{
				"api_key": "test_key",
				"secret":  "test_secret",
			},
		},
	}

	// Test getting an existing provider's config
	providerConfig, err := cfg.GetProviderConfig("test_provider")
	assert.NoError(t, err)
	assert.NotNil(t, providerConfig)
	assert.Equal(t, "test_key", providerConfig["api_key"])
	assert.Equal(t, "test_secret", providerConfig["secret"])

	// Test getting a non-existent provider's config
	_, err = cfg.GetProviderConfig("non_existent")
	assert.Error(t, err)

	// Test with nil providers map
	nilProvidersCfg := &config.Config{
		DefaultProvider: "test_provider",
		Providers:       nil,
	}
	_, err = nilProvidersCfg.GetProviderConfig("test_provider")
	assert.Error(t, err)

	// Test with invalid provider config format
	invalidCfg := &config.Config{
		DefaultProvider: "test_provider",
		Providers: map[string]interface{}{
			"test_provider": "not a map", // Invalid format, should be a map
		},
	}
	_, err = invalidCfg.GetProviderConfig("test_provider")
	assert.Error(t, err)
}

// TestDefaultValues tests that default values are set correctly
func TestDefaultValues(t *testing.T) {
	// Create a minimal config with just the required fields
	minimalConfig := `
default_provider: test_provider
providers:
  test_provider:
    api_key: test_key
`
	minimalFile, err := createTempConfig(minimalConfig)
	require.NoError(t, err)
	defer os.Remove(minimalFile)

	cfg, err := config.LoadConfig(minimalFile)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// Check that default values were set
	assert.Equal(t, "test_provider", cfg.DefaultProvider)
	assert.Equal(t, config.DefaultHTTPTimeout, cfg.HTTPTimeout)
	assert.Equal(t, config.DefaultRetryAttempts, cfg.RetryAttempts)
	assert.Equal(t, config.DefaultRetryDelay, cfg.RetryDelay)
	assert.Equal(t, config.DefaultSMSTemplate, cfg.SMSTemplate)
	assert.Equal(t, config.DefaultVoiceTemplate, cfg.VoiceTemplate)
}
