package twilio

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
  twilio:
    account_sid: AC1234567890abcdef1234567890abcdef
    auth_token: abcdef1234567890abcdef1234567890
    from_number: +1234567890
    region: us1
    api_version: 2010-04-01
`
	err := v.ReadConfig(strings.NewReader(testConfig))
	assert.NoError(t, err)

	// Test valid configuration
	config, err := LoadConfig(v)
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "AC1234567890abcdef1234567890abcdef", config.AccountSID)
	assert.Equal(t, "abcdef1234567890abcdef1234567890", config.AuthToken)
	assert.Equal(t, "+1234567890", config.FromNumber)
	assert.Equal(t, "us1", config.Region)
	assert.Equal(t, "2010-04-01", config.APIVersion)

	// Test invalid configuration (missing required fields)
	v = viper.New()
	v.SetConfigType("yaml")
	invalidConfig := `
providers:
  twilio:
    account_sid: 
    auth_token: 
    from_number: 
`
	err = v.ReadConfig(strings.NewReader(invalidConfig))
	assert.NoError(t, err)

	config, err = LoadConfig(v)
	assert.Error(t, err)
	assert.Nil(t, config)
}

func TestValidateConfig(t *testing.T) {
	// Test valid config
	validConfig := &TwilioConfig{
		AccountSID: "AC1234567890abcdef1234567890abcdef",
		AuthToken:  "abcdef1234567890abcdef1234567890",
		FromNumber: "+1234567890",
		Region:     "us1",
		APIVersion: "2010-04-01",
	}
	assert.NoError(t, validConfig.Validate())

	// Test invalid AccountSID format
	invalidSIDConfig := &TwilioConfig{
		AccountSID: "1234567890abcdef1234567890abcdef", // Should start with AC
		AuthToken:  "abcdef1234567890abcdef1234567890",
		FromNumber: "+1234567890",
	}
	assert.Error(t, invalidSIDConfig.Validate())

	// Test invalid FromNumber format
	invalidFromConfig := &TwilioConfig{
		AccountSID: "AC1234567890abcdef1234567890abcdef",
		AuthToken:  "abcdef1234567890abcdef1234567890",
		FromNumber: "1234567890", // Should start with +
	}
	assert.Error(t, invalidFromConfig.Validate())

	// Test missing required fields
	for _, config := range []*TwilioConfig{
		&TwilioConfig{AuthToken: "token", FromNumber: "+1234567890"},  // Missing AccountSID
		&TwilioConfig{AccountSID: "AC123", FromNumber: "+1234567890"}, // Missing AuthToken
		&TwilioConfig{AccountSID: "AC123", AuthToken: "token"},        // Missing FromNumber
	} {
		assert.Error(t, config.Validate())
	}
}

func TestSendSMS(t *testing.T) {
	// Create a test server that mimics Twilio API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if this is a request to the Messages endpoint
		if r.URL.Path == "/2010-04-01/Accounts/AC123/Messages.json" {
			// Check method
			assert.Equal(t, http.MethodPost, r.Method)

			// Check auth
			username, password, ok := r.BasicAuth()
			assert.True(t, ok)
			assert.Equal(t, "AC123", username)
			assert.Equal(t, "auth123", password)

			// Parse form data
			err := r.ParseForm()
			assert.NoError(t, err)

			// Check form values
			assert.Equal(t, "+1234567890", r.FormValue("To"))
			assert.Equal(t, "+0987654321", r.FormValue("From"))
			assert.Equal(t, "Hello from MyApp: This is a test message", r.FormValue("Body"))

			// Send response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			resp := twilioSMSResponse{
				SID:         "SM123456789",
				Status:      "sent",
				DateCreated: "2023-05-01T12:30:00Z",
				DateSent:    "2023-05-01T12:30:05Z",
				Direction:   "outbound-api",
				Price:       "-0.0075",
				PriceUnit:   "USD",
			}

			json.NewEncoder(w).Encode(resp)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Create a provider for testing
	provider := &Provider{
		config: &TwilioConfig{
			AccountSID: "AC123",
			AuthToken:  "auth123",
			FromNumber: "+0987654321",
			Region:     "us1",
			APIVersion: "2010-04-01",
		},
		baseURL: server.URL + "/2010-04-01/Accounts/AC123",
	}

	// Create a test client
	client := client.NewClient(&config.Config{
		HTTPTimeout: 10 * time.Second,
	})
	client.SetBasicAuth(provider.config.AccountSID, provider.config.AuthToken)
	provider.client = client

	// Create request
	req := model.SendSMSRequest{
		Message: model.Message{
			From: "+0987654321",
			To:   "+1234567890",
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
	assert.Equal(t, "SM123456789", resp.MessageID)
	assert.Equal(t, model.StatusSent, resp.Status)
	assert.Equal(t, ProviderName, resp.Provider)
	assert.Equal(t, 0.0075, resp.Cost)
	assert.Equal(t, "USD", resp.Currency)
}

func TestSendVoiceCall(t *testing.T) {
	// Create a test server that mimics Twilio API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if this is a request to the Calls endpoint
		if r.URL.Path == "/2010-04-01/Accounts/AC123/Calls.json" {
			// Check method
			assert.Equal(t, http.MethodPost, r.Method)

			// Check auth
			username, password, ok := r.BasicAuth()
			assert.True(t, ok)
			assert.Equal(t, "AC123", username)
			assert.Equal(t, "auth123", password)

			// Parse form data
			err := r.ParseForm()
			assert.NoError(t, err)

			// Check form values
			assert.Equal(t, "+1234567890", r.FormValue("To"))
			assert.Equal(t, "+0987654321", r.FormValue("From"))
			assert.Contains(t, r.FormValue("Twiml"), "This+is+a+test+voice+call")

			// Send response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			resp := twilioCallResponse{
				SID:         "CA123456789",
				Status:      "initiated",
				DateCreated: "2023-05-01T12:30:00Z",
				StartTime:   "2023-05-01T12:30:05Z",
				EndTime:     "",
				Duration:    "0",
				Price:       "-0.0150",
				PriceUnit:   "USD",
			}

			json.NewEncoder(w).Encode(resp)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Create a provider for testing
	provider := &Provider{
		config: &TwilioConfig{
			AccountSID: "AC123",
			AuthToken:  "auth123",
			FromNumber: "+0987654321",
			Region:     "us1",
			APIVersion: "2010-04-01",
		},
		baseURL: server.URL + "/2010-04-01/Accounts/AC123",
	}

	// Create a test client
	client := client.NewClient(&config.Config{
		HTTPTimeout: 10 * time.Second,
	})
	client.SetBasicAuth(provider.config.AccountSID, provider.config.AuthToken)
	provider.client = client

	// Create request
	req := model.SendVoiceRequest{
		Message: model.Message{
			From: "+0987654321",
			To:   "+1234567890",
			By:   "MyApp",
		},
		Template: "This is a test voice call from {by}",
		Data: map[string]interface{}{
			"by": "MyApp",
		},
		Options: map[string]interface{}{
			"voice":    "woman",
			"language": "en-US",
		},
	}

	// Send call
	resp, err := provider.SendVoiceCall(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// Check response
	assert.Equal(t, "CA123456789", resp.CallID)
	assert.Equal(t, model.CallStatusInitiated, resp.Status)
	assert.Equal(t, ProviderName, resp.Provider)
	assert.Equal(t, 0.0150, resp.Cost)
	assert.Equal(t, "USD", resp.Currency)
	assert.Equal(t, 0, resp.Duration)
}
