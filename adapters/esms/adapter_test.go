package esms

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/zinzinday/go-sms/client"
	"github.com/zinzinday/go-sms/config"
	"github.com/zinzinday/go-sms/model"
)

func TestLoadConfig(t *testing.T) {
	// Create a viper instance with test configuration
	v := viper.New()
	v.SetConfigType("yaml")
	testConfig := `
providers:
  esms:
    api_key: test_api_key
    secret: test_secret
    brandname: TestBrand
    sms_type: 2
    base_url: http://test.esms.vn/api
`
	err := v.ReadConfig(strings.NewReader(testConfig))
	assert.NoError(t, err)

	// Test valid configuration
	config, err := LoadConfig(v)
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "test_api_key", config.APIKey)
	assert.Equal(t, "test_secret", config.Secret)
	assert.Equal(t, "TestBrand", config.Brandname)
	assert.Equal(t, 2, config.SMSType)
	assert.Equal(t, "http://test.esms.vn/api", config.BaseURL)

	// Test invalid configuration (missing required fields)
	v = viper.New()
	v.SetConfigType("yaml")
	invalidConfig := `
providers:
  esms:
    api_key: 
    secret: 
    sms_type: 999
`
	err = v.ReadConfig(strings.NewReader(invalidConfig))
	assert.NoError(t, err)

	config, err = LoadConfig(v)
	assert.Error(t, err)
	assert.Nil(t, config)
}

func TestValidateConfig(t *testing.T) {
	// Test valid config
	validConfig := &ESMSConfig{
		APIKey:    "test_api_key",
		Secret:    "test_secret",
		Brandname: "TestBrand",
		SMSType:   2,
		BaseURL:   "http://test.esms.vn/api",
	}
	assert.NoError(t, validConfig.Validate())

	// Test invalid SMS type
	invalidTypeConfig := &ESMSConfig{
		APIKey:    "test_api_key",
		Secret:    "test_secret",
		Brandname: "TestBrand",
		SMSType:   999, // Invalid SMS type
		BaseURL:   "http://test.esms.vn/api",
	}
	assert.Error(t, invalidTypeConfig.Validate())

	// Test missing brandname for type 2
	missingBrandConfig := &ESMSConfig{
		APIKey:    "test_api_key",
		Secret:    "test_secret",
		Brandname: "", // Empty brandname
		SMSType:   2,  // Requires brandname
		BaseURL:   "http://test.esms.vn/api",
	}
	assert.Error(t, missingBrandConfig.Validate())

	// Test missing required fields
	for _, config := range []*ESMSConfig{
		&ESMSConfig{Secret: "secret", Brandname: "Brand", SMSType: 2}, // Missing APIKey
		&ESMSConfig{APIKey: "key", Brandname: "Brand", SMSType: 2},    // Missing Secret
	} {
		assert.Error(t, config.Validate())
	}
}

func TestSendSMS(t *testing.T) {
	// Create a test server that mimics eSMS API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if this is a request to the SMS endpoint
		if r.URL.Path == "/api/sms/send" {
			// Check method
			assert.Equal(t, http.MethodPost, r.Method)

			// Parse form data
			err := r.ParseForm()
			assert.NoError(t, err)

			// Check form values
			assert.Equal(t, "test_api_key", r.FormValue("ApiKey"))
			assert.Equal(t, "test_secret", r.FormValue("SecretKey"))
			assert.Equal(t, "+84123456789", r.FormValue("Phone"))
			assert.Equal(t, "TestBrand", r.FormValue("Brandname"))
			assert.Equal(t, "Hello from MyApp: This is a test message", r.FormValue("Content"))
			assert.Equal(t, "2", r.FormValue("SmsType"))

			// Send response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			resp := esmsSMSResponse{
				CodeResult:      "100",
				CountRegenerate: 0,
				ErrorMessage:    "",
				SMSID:           "SMS123456789",
			}

			json.NewEncoder(w).Encode(resp)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Create a provider for testing
	provider := &Provider{
		config: &ESMSConfig{
			APIKey:    "test_api_key",
			Secret:    "test_secret",
			Brandname: "TestBrand",
			SMSType:   2,
			BaseURL:   server.URL + "/api",
		},
	}

	// Create a test client
	provider.client = client.NewClient(&config.Config{
		HTTPTimeout: 10 * time.Second,
	})

	// Create request
	req := model.SendSMSRequest{
		Message: model.Message{
			From: "TestBrand",
			To:   "+84123456789",
			By:   "MyApp",
		},
		Template: "Hello from {by}: {message}",
		Data: map[string]interface{}{
			"message": "This is a test message",
			"by":      "MyApp",
		},
	}

	// Send message
	resp, err := provider.SendSMS(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// Check response
	assert.Equal(t, "SMS123456789", resp.MessageID)
	assert.Equal(t, model.StatusSent, resp.Status)
	assert.Equal(t, ProviderName, resp.Provider)
}

func TestSendVoiceCall(t *testing.T) {
	// Create a test server that mimics eSMS Voice API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if this is a request to the Voice OTP endpoint
		if r.URL.Path == "/api/voice/otp" {
			// Check method
			assert.Equal(t, http.MethodPost, r.Method)

			// Parse form data
			err := r.ParseForm()
			assert.NoError(t, err)

			// Check form values
			assert.Equal(t, "test_api_key", r.FormValue("ApiKey"))
			assert.Equal(t, "test_secret", r.FormValue("SecretKey"))
			assert.Equal(t, "+84123456789", r.FormValue("Phone"))
			assert.Equal(t, "123456", r.FormValue("Code"))

			// Send response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			resp := esmsVoiceResponse{
				CodeResult:   "100",
				ErrorMessage: "",
				CallID:       "CALL123456789",
			}

			json.NewEncoder(w).Encode(resp)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Create a provider for testing
	provider := &Provider{
		config: &ESMSConfig{
			APIKey:    "test_api_key",
			Secret:    "test_secret",
			Brandname: "TestBrand",
			SMSType:   2,
			BaseURL:   server.URL + "/api",
		},
	}

	// Create a test client
	provider.client = client.NewClient(&config.Config{
		HTTPTimeout: 10 * time.Second,
	})

	// Create request
	req := model.SendVoiceRequest{
		Message: model.Message{
			From: "TestBrand",
			To:   "+84123456789",
			By:   "MyApp",
		},
		Template: "Your OTP code is 123456",
		Data:     map[string]interface{}{},
		Options: map[string]interface{}{
			"speed":       1.0,
			"retry_times": 2,
		},
	}

	// Send voice call
	resp, err := provider.SendVoiceCall(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// Check response
	assert.Equal(t, "CALL123456789", resp.CallID)
	assert.Equal(t, model.CallStatusInitiated, resp.Status)
	assert.Equal(t, ProviderName, resp.Provider)
}

func TestExtractOTPFromMessage(t *testing.T) {
	testCases := []struct {
		message  string
		expected string
	}{
		{"Your OTP code is 123456", "123456"},
		{"Please use 9876 as your verification code", "9876"},
		{"Code: 54321", "54321"},
		{"12345678 is your OTP", "12345678"},
		{"No OTP here", "No OTP here"},                    // No digits found, returns the message
		{"Random text with some digits 123 in it", "123"}, // Found digits but not enough
		{"Use code 1234 for login", "1234"},
	}

	for _, tc := range testCases {
		t.Run(tc.message, func(t *testing.T) {
			result := extractOTPFromMessage(tc.message)
			assert.Equal(t, tc.expected, result)
		})
	}
}
