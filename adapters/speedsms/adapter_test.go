package speedsms

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-fork/sms/client"
	"github.com/go-fork/sms/config"
	"github.com/go-fork/sms/model"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Create a viper instance with test configuration
	v := viper.New()
	v.SetConfigType("yaml")
	testConfig := `
providers:
  speedsms:
    token: test_token_with_at_least_20_characters
    sender: TestBrand
    sms_type: 2
    base_url: https://test.speedsms.vn/index.php
`
	err := v.ReadConfig(strings.NewReader(testConfig))
	assert.NoError(t, err)

	// Test valid configuration
	config, err := LoadConfig(v)
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "test_token_with_at_least_20_characters", config.Token)
	assert.Equal(t, "TestBrand", config.Sender)
	assert.Equal(t, 2, config.SMSType)
	assert.Equal(t, "https://test.speedsms.vn/index.php", config.BaseURL)

	// Test invalid configuration (missing required fields)
	v = viper.New()
	v.SetConfigType("yaml")
	invalidConfig := `
providers:
  speedsms:
    token: 
    sender: TestBrand
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
	validConfig := &SpeedSMSConfig{
		Token:   "test_token_with_at_least_20_characters",
		Sender:  "TestBrand",
		SMSType: 2,
		BaseURL: "https://api.speedsms.vn/index.php",
	}
	assert.NoError(t, validConfig.Validate())

	// Test invalid SMS type
	invalidTypeConfig := &SpeedSMSConfig{
		Token:   "test_token_with_at_least_20_characters",
		Sender:  "TestBrand",
		SMSType: 999, // Invalid SMS type
		BaseURL: "https://api.speedsms.vn/index.php",
	}
	assert.Error(t, invalidTypeConfig.Validate())

	// Test invalid token (too short)
	invalidTokenConfig := &SpeedSMSConfig{
		Token:   "short_token",
		Sender:  "TestBrand",
		SMSType: 2,
		BaseURL: "https://api.speedsms.vn/index.php",
	}
	assert.Error(t, invalidTokenConfig.Validate())

	// Test invalid base URL
	invalidURLConfig := &SpeedSMSConfig{
		Token:   "test_token_with_at_least_20_characters",
		Sender:  "TestBrand",
		SMSType: 2,
		BaseURL: "invalid-url", // Missing http/https prefix
	}
	assert.Error(t, invalidURLConfig.Validate())

	// Test missing required field (token)
	missingTokenConfig := &SpeedSMSConfig{
		Token:   "", // Empty token
		Sender:  "TestBrand",
		SMSType: 2,
		BaseURL: "https://api.speedsms.vn/index.php",
	}
	assert.Error(t, missingTokenConfig.Validate())
}

func TestSendSMS(t *testing.T) {
	// Create a test server that mimics SpeedSMS API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if this is a request to the SMS endpoint
		if r.URL.Path == "/index.php/sms/send" {
			// Check method
			assert.Equal(t, http.MethodPost, r.Method)

			// Check headers
			assert.Equal(t, "test_token_with_at_least_20_characters", r.Header.Get("Authorization"))
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			// Decode request body
			var reqBody speedSMSSendRequest
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&reqBody)
			assert.NoError(t, err)

			// Check request values
			assert.Equal(t, []string{"+84123456789"}, reqBody.To)
			assert.Equal(t, "Hello from MyApp: This is a test message", reqBody.Content)
			assert.Equal(t, 2, reqBody.Type)
			assert.Equal(t, "TestBrand", reqBody.Sender)

			// Send success response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			resp := speedSMSResponse{
				Status:  "success",
				Code:    0,
				Message: "Request processed successfully",
				Data:    []string{"+84123456789"},
			}

			json.NewEncoder(w).Encode(resp)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Create a provider for testing
	provider := &Provider{
		config: &SpeedSMSConfig{
			Token:   "test_token_with_at_least_20_characters",
			Sender:  "TestBrand",
			SMSType: 2,
			BaseURL: server.URL + "/index.php",
		},
	}

	// Create a test client
	provider.client = client.NewClient(&config.Config{
		HTTPTimeout: 10 * time.Second,
	})
	// Set headers
	provider.client.SetHeader("Authorization", provider.config.Token)
	provider.client.SetHeader("Content-Type", "application/json")

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
	assert.Contains(t, resp.MessageID, "sms_")
	assert.Equal(t, model.StatusSent, resp.Status)
	assert.Equal(t, ProviderName, resp.Provider)

	// Test error response
	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		resp := speedSMSResponse{
			Status:  "error",
			Code:    1001,
			Message: "Invalid phone number",
		}

		json.NewEncoder(w).Encode(resp)
	}))
	defer errorServer.Close()

	// Update provider with error server
	provider.config.BaseURL = errorServer.URL + "/index.php"

	// Send message to error server
	_, err = provider.SendSMS(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SpeedSMS error")
}

func TestSendVoiceCall(t *testing.T) {
	// Create a provider for testing
	provider := &Provider{
		config: &SpeedSMSConfig{
			Token:   "test_token_with_at_least_20_characters",
			Sender:  "TestBrand",
			SMSType: 2,
			BaseURL: "https://api.speedsms.vn/index.php",
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
	}

	// Send voice call (should fail since SpeedSMS doesn't support voice)
	_, err := provider.SendVoiceCall(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not supported")
}
