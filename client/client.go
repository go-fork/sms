package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-fork/sms/config"
	"github.com/go-resty/resty/v2"
)

// Client represents an HTTP client for making API requests to providers
type Client struct {
	// restyClient is the underlying HTTP client from the resty library
	restyClient *resty.Client

	// config holds the client configuration
	config *config.Config
}

// NewClient creates a new HTTP client with the provided configuration
func NewClient(config *config.Config) *Client {
	if config == nil {
		// Create a default configuration if none is provided
		// This should not happen in normal operation but prevents nil pointer panics
		defaultTimeout := 10 * time.Second
		config = &config.Config{
			HTTPTimeout: defaultTimeout,
		}
	}

	// Create a new resty client with configuration
	restyClient := resty.New()

	// Configure timeouts
	restyClient.SetTimeout(config.HTTPTimeout)

	// Set reasonable defaults for connections
	restyClient.SetTransport(&http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     90 * time.Second,
	})

	// Disable resty's built-in retry to use our custom retry logic
	restyClient.SetRetryCount(0)

	return &Client{
		restyClient: restyClient,
		config:      config,
	}
}

// R returns a new request object for building and executing requests
func (c *Client) R() *resty.Request {
	return c.restyClient.R()
}

// SetBaseURL sets the base URL for all requests
func (c *Client) SetBaseURL(url string) *Client {
	c.restyClient.SetBaseURL(url)
	return c
}

// SetHeader sets a header for all requests
func (c *Client) SetHeader(key, value string) *Client {
	c.restyClient.SetHeader(key, value)
	return c
}

// SetAuthToken sets the Authorization header with a Bearer token for all requests
func (c *Client) SetAuthToken(token string) *Client {
	c.restyClient.SetAuthToken(token)
	return c
}

// SetBasicAuth sets basic authentication for all requests
func (c *Client) SetBasicAuth(username, password string) *Client {
	c.restyClient.SetBasicAuth(username, password)
	return c
}

// Get performs a GET request with context
func (c *Client) Get(ctx context.Context, url string) (*resty.Response, error) {
	return c.R().SetContext(ctx).Get(url)
}

// Post performs a POST request with context and JSON body
func (c *Client) Post(ctx context.Context, url string, body interface{}) (*resty.Response, error) {
	return c.R().SetContext(ctx).SetBody(body).Post(url)
}

// PostForm performs a POST request with context and form data
func (c *Client) PostForm(ctx context.Context, url string, formData map[string]string) (*resty.Response, error) {
	return c.R().SetContext(ctx).SetFormData(formData).Post(url)
}

// ProcessResponse is a helper to process API responses and handle errors
func (c *Client) ProcessResponse(resp *resty.Response, err error) ([]byte, error) {
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Check for HTTP errors
	if !resp.IsSuccess() {
		return nil, fmt.Errorf("HTTP error: status code: %d, body: %s",
			resp.StatusCode(), resp.String())
	}

	return resp.Body(), nil
}
