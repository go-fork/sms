package tests

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zinzinday/go-sms"
	"github.com/zinzinday/go-sms/model"
)

// MockProvider implements model.Provider for testing
type MockProvider struct {
	mock.Mock
}

func (m *MockProvider) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockProvider) SendSMS(ctx context.Context, req model.SendSMSRequest) (model.SendSMSResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(model.SendSMSResponse), args.Error(1)
}

func (m *MockProvider) SendVoiceCall(ctx context.Context, req model.SendVoiceRequest) (model.SendVoiceResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(model.SendVoiceResponse), args.Error(1)
}

// createTempConfig creates a temporary config file for testing
func createTempConfig(content string) (string, error) {
	tmpfile, err := os.CreateTemp("", "sms-config-*.yaml")
	if err != nil {
		return "", err
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
		return "", err
	}

	if err := tmpfile.Close(); err != nil {
		os.Remove(tmpfile.Name())
		return "", err
	}

	return tmpfile.Name(), nil
}

// TestNewModule tests module creation with various configurations
func TestNewModule(t *testing.T) {
	// Valid configuration
	validConfig := `
default_provider: test_provider
http_timeout: 10s
retry_attempts: 3
retry_delay: 500ms
sms_template: "Your message is {message}"
voice_template: "Your message is {message}"

providers:
  test_provider:
    api_key: test_key
`

	// Invalid configuration - missing default_provider
	invalidConfig := `
http_timeout: 10s
retry_attempts: 3
retry_delay: 500ms
`

	tests := []struct {
		name        string
		config      string
		expectError bool
	}{
		{
			name:        "Valid configuration",
			config:      validConfig,
			expectError: false,
		},
		{
			name:        "Invalid configuration",
			config:      invalidConfig,
			expectError: true,
		},
		{
			name:        "Non-existent config file",
			config:      "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var configFile string
			var err error

			if tt.config != "" {
				configFile, err = createTempConfig(tt.config)
				require.NoError(t, err)
				defer os.Remove(configFile)
			} else {
				configFile = "non-existent-file.yaml"
			}

			module, err := sms.NewModule(configFile)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, module)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, module)
			}
		})
	}
}

// TestProviderManagement tests adding, switching, and retrieving providers
func TestProviderManagement(t *testing.T) {
	configFile, err := createTempConfig(`
default_provider: provider1
http_timeout: 10s
retry_attempts: 3
retry_delay: 500ms

providers:
  provider1:
    api_key: key1
  provider2:
    api_key: key2
`)
	require.NoError(t, err)
	defer os.Remove(configFile)

	module, err := sms.NewModule(configFile)
	require.NoError(t, err)

	// Create mock providers
	provider1 := new(MockProvider)
	provider1.On("Name").Return("provider1")

	provider2 := new(MockProvider)
	provider2.On("Name").Return("provider2")

	provider3 := new(MockProvider)
	provider3.On("Name").Return("provider3")

	// Test AddProvider
	err = module.AddProvider(provider1)
	assert.NoError(t, err)

	err = module.AddProvider(provider2)
	assert.NoError(t, err)

	// Test adding a provider with the same name
	err = module.AddProvider(provider1)
	assert.Error(t, err)

	// Test SwitchProvider
	err = module.SwitchProvider("provider2")
	assert.NoError(t, err)

	// Test switching to a non-existent provider
	err = module.SwitchProvider("non-existent")
	assert.Error(t, err)

	// Test GetProvider
	p, err := module.GetProvider("provider1")
	assert.NoError(t, err)
	assert.Equal(t, "provider1", p.Name())

	// Test getting a non-existent provider
	p, err = module.GetProvider("non-existent")
	assert.Error(t, err)
	assert.Nil(t, p)

	// Test GetActiveProvider
	active, err := module.GetActiveProvider()
	assert.NoError(t, err)
	assert.Equal(t, "provider2", active.Name())

	// Test with no active provider (create a new module without adding providers)
	newModule, err := sms.NewModule(configFile)
	require.NoError(t, err)
	_, err = newModule.GetActiveProvider()
	assert.Error(t, err)
}

// TestSendSMS tests sending SMS messages
func TestSendSMS(t *testing.T) {
	configFile, err := createTempConfig(`
default_provider: test_provider
http_timeout: 10s
retry_attempts: 3
retry_delay: 500ms

providers:
  test_provider:
    api_key: test_key
`)
	require.NoError(t, err)
	defer os.Remove(configFile)

	module, err := sms.NewModule(configFile)
	require.NoError(t, err)

	// Create mock provider
	provider := new(MockProvider)
	provider.On("Name").Return("test_provider")

	// Set up expected calls for different scenarios
	successReq := model.SendSMSRequest{
		Message: model.Message{
			From: "Sender",
			To:   "+1234567890",
			By:   "TestApp",
		},
		Data: map[string]interface{}{
			"message": "Test message",
		},
	}

	successResp := model.SendSMSResponse{
		MessageID: "msg_123",
		Status:    model.StatusSent,
		Provider:  "test_provider",
		SentAt:    time.Now(),
	}

	provider.On("SendSMS", mock.Anything, successReq).Return(successResp, nil)

	errorReq := model.SendSMSRequest{
		Message: model.Message{
			From: "Sender",
			To:   "invalid",
			By:   "TestApp",
		},
	}

	provider.On("SendSMS", mock.Anything, errorReq).Return(model.SendSMSResponse{}, errors.New("invalid phone number"))

	// Add the mock provider
	err = module.AddProvider(provider)
	require.NoError(t, err)

	// Test successful SMS sending
	resp, err := module.SendSMS(context.Background(), successReq)
	assert.NoError(t, err)
	assert.Equal(t, "msg_123", resp.MessageID)
	assert.Equal(t, model.StatusSent, resp.Status)
	assert.Equal(t, "test_provider", resp.Provider)

	// Test error case
	_, err = module.SendSMS(context.Background(), errorReq)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid phone number")

	// Test with no active provider
	newModule, err := sms.NewModule(configFile)
	require.NoError(t, err)
	_, err = newModule.SendSMS(context.Background(), successReq)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no active provider")

	// Verify all expectations were met
	provider.AssertExpectations(t)
}

// TestSendVoiceCall tests sending voice calls
func TestSendVoiceCall(t *testing.T) {
	configFile, err := createTempConfig(`
default_provider: test_provider
http_timeout: 10s
retry_attempts: 3
retry_delay: 500ms

providers:
  test_provider:
    api_key: test_key
`)
	require.NoError(t, err)
	defer os.Remove(configFile)

	module, err := sms.NewModule(configFile)
	require.NoError(t, err)

	// Create mock provider
	provider := new(MockProvider)
	provider.On("Name").Return("test_provider")

	// Set up expected calls for different scenarios
	successReq := model.SendVoiceRequest{
		Message: model.Message{
			From: "Sender",
			To:   "+1234567890",
			By:   "TestApp",
		},
		Data: map[string]interface{}{
			"message": "Test message",
		},
	}

	successResp := model.SendVoiceResponse{
		CallID:    "call_123",
		Status:    model.CallStatusInitiated,
		Provider:  "test_provider",
		StartedAt: time.Now(),
		Duration:  0,
	}

	provider.On("SendVoiceCall", mock.Anything, successReq).Return(successResp, nil)

	errorReq := model.SendVoiceRequest{
		Message: model.Message{
			From: "Sender",
			To:   "invalid",
			By:   "TestApp",
		},
	}

	provider.On("SendVoiceCall", mock.Anything, errorReq).Return(model.SendVoiceResponse{}, errors.New("invalid phone number"))

	// Add the mock provider
	err = module.AddProvider(provider)
	require.NoError(t, err)

	// Test successful voice call
	resp, err := module.SendVoiceCall(context.Background(), successReq)
	assert.NoError(t, err)
	assert.Equal(t, "call_123", resp.CallID)
	assert.Equal(t, model.CallStatusInitiated, resp.Status)
	assert.Equal(t, "test_provider", resp.Provider)

	// Test error case
	_, err = module.SendVoiceCall(context.Background(), errorReq)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid phone number")

	// Test with no active provider
	newModule, err := sms.NewModule(configFile)
	require.NoError(t, err)
	_, err = newModule.SendVoiceCall(context.Background(), successReq)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no active provider")

	// Verify all expectations were met
	provider.AssertExpectations(t)
}
