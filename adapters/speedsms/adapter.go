package speedsms

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-fork/sms/client"
	"github.com/go-fork/sms/config"
	"github.com/go-fork/sms/model"
	"github.com/spf13/viper"
)

const (
	// ProviderName is the name of this provider
	ProviderName = "speedsms"

	// SpeedSMSSendSMSEndpoint is the endpoint for sending SMS
	SpeedSMSSendSMSEndpoint = "/sms/send"

	// SpeedSMSCheckBalanceEndpoint is the endpoint for checking account balance
	SpeedSMSCheckBalanceEndpoint = "/user/balance"
)

// SpeedSMS API request/response structures
type speedSMSSendRequest struct {
	To      []string `json:"to"`
	Content string   `json:"content"`
	Type    int      `json:"sms_type"`
	Sender  string   `json:"sender,omitempty"`
}

type speedSMSResponse struct {
	Status  string   `json:"status"`
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    []string `json:"data,omitempty"`
}

// Provider implements the model.Provider interface for SpeedSMS
type Provider struct {
	// client is the HTTP client for making API requests
	client *client.Client

	// config holds the SpeedSMS provider configuration
	config *SpeedSMSConfig
}

// NewProvider creates a new SpeedSMS provider instance
func NewProvider(configFile string) (model.Provider, error) {
	// Load the main configuration
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Create a new Viper instance
	v := viper.New()
	v.SetConfigFile(configFile)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Load SpeedSMS-specific configuration
	speedConfig, err := LoadConfig(v)
	if err != nil {
		return nil, fmt.Errorf("failed to load SpeedSMS configuration: %w", err)
	}

	// Create HTTP client
	httpClient := client.NewClient(cfg)

	// Set authorization header with token
	httpClient.SetHeader("Authorization", speedConfig.Token)
	httpClient.SetHeader("Content-Type", "application/json")

	return &Provider{
		client: httpClient,
		config: speedConfig,
	}, nil
}

// Name returns the provider name
func (p *Provider) Name() string {
	return ProviderName
}

// SendSMS sends an SMS message using SpeedSMS
func (p *Provider) SendSMS(ctx context.Context, req model.SendSMSRequest) (model.SendSMSResponse, error) {
	// Get the message body from template
	template := req.Template
	if template == "" {
		// If no template provided, use default template from config
		template = "{message}"
	}

	// Render the message template with provided data
	messageBody := req.Message.Render(template, req.Data)
	if messageBody == "" {
		return model.SendSMSResponse{}, fmt.Errorf("empty message body after rendering template")
	}

	// Determine the sender (from brandname in config or from request)
	sender := p.config.Sender
	if req.Message.From != "" {
		sender = req.Message.From
	}

	// Prepare the API request
	endpoint := p.config.BaseURL + SpeedSMSSendSMSEndpoint

	// Get SMS type from config or options
	smsType := p.config.SMSType
	if req.Options != nil {
		if optSmsType, ok := req.Options["sms_type"].(int); ok {
			smsType = optSmsType
		}
	}

	// SpeedSMS requires phone numbers as an array, but we're sending to just one
	phoneNumbers := []string{req.Message.To}

	// Build the request body
	reqBody := speedSMSSendRequest{
		To:      phoneNumbers,
		Content: messageBody,
		Type:    smsType,
	}

	// Add sender if available
	if sender != "" {
		reqBody.Sender = sender
	}

	// Make the API request to SpeedSMS
	resp, err := p.client.R().
		SetContext(ctx).
		SetBody(reqBody).
		Post(endpoint)

	if err != nil {
		return model.SendSMSResponse{}, fmt.Errorf("SpeedSMS API request failed: %w", err)
	}

	// Handle error responses
	if resp.StatusCode() >= 400 {
		return model.SendSMSResponse{}, fmt.Errorf("SpeedSMS API error: %s", resp.String())
	}

	// Parse the response
	var speedResp speedSMSResponse
	if err := json.Unmarshal(resp.Body(), &speedResp); err != nil {
		return model.SendSMSResponse{}, fmt.Errorf("failed to parse SpeedSMS response: %w", err)
	}

	// Check for SpeedSMS error codes
	if speedResp.Status != "success" {
		return model.SendSMSResponse{}, fmt.Errorf("SpeedSMS error: %d - %s", speedResp.Code, speedResp.Message)
	}

	// Generate a message ID if SpeedSMS didn't provide one
	// SpeedSMS returns success data in the 'data' field which contains the phone numbers
	messageID := fmt.Sprintf("sms_%d", time.Now().UnixNano())
	if len(speedResp.Data) > 0 {
		messageID = fmt.Sprintf("sms_%s_%d", strings.Join(speedResp.Data, "_"), time.Now().Unix())
	}

	// Map SpeedSMS status to our status
	status := model.StatusSent
	if speedResp.Status != "success" {
		status = model.StatusFailed
	}

	// Convert SpeedSMS response to our response model
	return model.SendSMSResponse{
		MessageID: messageID,
		Status:    status,
		Provider:  ProviderName,
		SentAt:    time.Now(),
		ProviderResponse: map[string]interface{}{
			"status":  speedResp.Status,
			"code":    speedResp.Code,
			"message": speedResp.Message,
			"data":    speedResp.Data,
		},
	}, nil
}

// SendVoiceCall initiates a voice call using SpeedSMS
// Note: SpeedSMS doesn't support voice calls directly, so this method returns an error
func (p *Provider) SendVoiceCall(ctx context.Context, req model.SendVoiceRequest) (model.SendVoiceResponse, error) {
	// SpeedSMS doesn't support voice calls natively
	// Return an appropriate error
	return model.SendVoiceResponse{}, fmt.Errorf("voice calls are not supported by the SpeedSMS provider")
}

// GetBalance returns the current balance of the SpeedSMS account
// This is a provider-specific method that isn't part of the Provider interface
func (p *Provider) GetBalance(ctx context.Context) (float64, error) {
	endpoint := p.config.BaseURL + SpeedSMSCheckBalanceEndpoint

	resp, err := p.client.R().
		SetContext(ctx).
		Get(endpoint)

	if err != nil {
		return 0, fmt.Errorf("failed to get balance: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	var result struct {
		Status  string  `json:"status"`
		Code    int     `json:"code"`
		Message string  `json:"message"`
		Data    float64 `json:"data"`
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return 0, fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Status != "success" {
		return 0, fmt.Errorf("error getting balance: %s", result.Message)
	}

	return result.Data, nil
}
