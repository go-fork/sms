package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-fork/sms/client"
	"github.com/go-fork/sms/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestClientCreation tests client initialization with different configurations
func TestClientCreation(t *testing.T) {
	tests := []struct {
		name        string
		config      *config.Config
		expectError bool
	}{
		{
			name: "Valid configuration",
			config: &config.Config{
				HTTPTimeout: 10 * time.Second,
			},
			expectError: false,
		},
		{
			name:        "Nil configuration",
			config:      nil,
			expectError: false, // Should use default config
		},
		{
			name: "Zero timeout",
			config: &config.Config{
				HTTPTimeout: 0,
			},
			expectError: false, // Should use default timeout
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := client.NewClient(tt.config)
			assert.NotNil(t, c, "Client should not be nil")
		})
	}
}

// TestClientHeaders tests setting headers on the client
func TestClientHeaders(t *testing.T) {
	c := client.NewClient(&config.Config{
		HTTPTimeout: 5 * time.Second,
	})

	// Test setting a single header
	c.SetHeader("X-Test-Header", "test-value")

	// Test setting auth token
	c.SetAuthToken("test-token")

	// Test setting basic auth
	c.SetBasicAuth("username", "password")

	// Since we can't directly access the resty client's headers,
	// we'll test them through an HTTP request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "test-value", r.Header.Get("X-Test-Header"))
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	resp, err := c.Get(context.Background(), server.URL)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
}

// TestClientMethods tests the different HTTP methods
func TestClientMethods(t *testing.T) {
	c := client.NewClient(&config.Config{
		HTTPTimeout: 5 * time.Second,
	})

	// Create a test server that validates the HTTP method
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{"method": r.Method}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Test GET
	resp, err := c.Get(context.Background(), server.URL)
	require.NoError(t, err)
	var result map[string]string
	err = json.Unmarshal(resp.Body(), &result)
	require.NoError(t, err)
	assert.Equal(t, "GET", result["method"])

	// Test POST with JSON body
	body := map[string]string{"key": "value"}
	resp, err = c.Post(context.Background(), server.URL, body)
	require.NoError(t, err)
	err = json.Unmarshal(resp.Body(), &result)
	require.NoError(t, err)
	assert.Equal(t, "POST", result["method"])

	// Test POST with form data
	formData := map[string]string{"form_key": "form_value"}
	resp, err = c.PostForm(context.Background(), server.URL, formData)
	require.NoError(t, err)
	err = json.Unmarshal(resp.Body(), &result)
	require.NoError(t, err)
	assert.Equal(t, "POST", result["method"])
}

// TestClientTimeout tests that the client respects timeout settings
func TestClientTimeout(t *testing.T) {
	c := client.NewClient(&config.Config{
		HTTPTimeout: 100 * time.Millisecond, // Very short timeout
	})

	// Create a test server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond) // Longer than timeout
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// The request should timeout
	_, err := c.Get(context.Background(), server.URL)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "timeout")
}

// TestClientProcessResponse tests the ProcessResponse helper
func TestClientProcessResponse(t *testing.T) {
	c := client.NewClient(&config.Config{
		HTTPTimeout: 5 * time.Second,
	})

	tests := []struct {
		name         string
		statusCode   int
		responseBody string
		expectError  bool
	}{
		{
			name:         "Success response",
			statusCode:   http.StatusOK,
			responseBody: `{"success": true}`,
			expectError:  false,
		},
		{
			name:         "Error response",
			statusCode:   http.StatusBadRequest,
			responseBody: `{"error": "Bad request"}`,
			expectError:  true,
		},
		{
			name:         "Server error",
			statusCode:   http.StatusInternalServerError,
			responseBody: `{"error": "Internal server error"}`,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			resp, err := c.Get(context.Background(), server.URL)
			require.NoError(t, err)

			body, err := c.ProcessResponse(resp, nil)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, []byte(tt.responseBody), body)
			}
		})
	}
}

// TestClientBaseURL tests setting the base URL
func TestClientBaseURL(t *testing.T) {
	c := client.NewClient(&config.Config{
		HTTPTimeout: 5 * time.Second,
	})

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test", r.URL.Path)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Set base URL and make request to a path
	c.SetBaseURL(server.URL)
	resp, err := c.Get(context.Background(), "/test")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
}
