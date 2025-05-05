package esms

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/zinzinday/go-sms/client"
	"github.com/zinzinday/go-sms/config"
	"github.com/zinzinday/go-sms/model"
)

const (
	// ProviderName is the name of this provider
	ProviderName = "esms"

	// ESMSSendSMSEndpoint is the endpoint for sending SMS
	ESMSSendSMSEndpoint = "/sms/send"

	// ESMSCheckBalanceEndpoint is the endpoint for checking account balance
	ESMSCheckBalanceEndpoint = "/user/balance"

	// ESMSVoiceOTPEndpoint is the endpoint for sending voice OTP
	ESMSVoiceOTPEndpoint = "/voice/otp"
)

// ESMS API response structures
type esmsSMSResponse struct {
	CodeResult      string `json:"CodeResult"`
	CountRegenerate int    `json:"CountRegenerate"`
	ErrorMessage    string `json:"ErrorMessage"`
	SMSID           string `json:"SMSID"`
}

type esmsVoiceResponse struct {
	CodeResult   string `json:"CodeResult"`
	ErrorMessage string `json:"ErrorMessage"`
	CallID       string `json:"CallId"`
}

// Provider implements the model.Provider interface for eSMS
type Provider struct {
	// client is the HTTP client for making API requests
	client *client.Client

	// config holds the eSMS provider configuration
	config *ESMSConfig
}

// NewProvider creates a new eSMS provider instance
func NewProvider(configFile string) (*Provider, error) {
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

	// Load eSMS-specific configuration
	esmsConfig, err := LoadConfig(v)
	if err != nil {
		return nil, fmt.Errorf("failed to load eSMS configuration: %w", err)
	}

	// Create HTTP client
	httpClient := client.NewClient(cfg)

	return &Provider{
		client: httpClient,
		config: esmsConfig,
	}, nil
}

// Name returns the provider name
func (p *Provider) Name() string {
	return ProviderName
}

// SendSMS sends an SMS message using eSMS
func (p *Provider) SendSMS(ctx context.Context, req model.SendSMSRequest) (model.SendSMSResponse, error) {
	// Get the message body from template
	template := req.Template
	if template == "" {
		// If no template provided, use the Message.Body field directly or an empty string
		template = "{message}"
	}

	// Render the message template with provided data
	messageBody := req.Message.Render(template, req.Data)
	if messageBody == "" {
		return model.SendSMSResponse{}, fmt.Errorf("empty message body after rendering template")
	}

	// Determine the sender (from brandname in config or from request)
	sender := p.config.Brandname
	if req.Message.From != "" {
		sender = req.Message.From
	}

	// If SMS type is not brandname (2) and sender is not a phone number, use default sender
	if p.config.SMSType != 2 && !strings.HasPrefix(sender, "+") {
		sender = "" // eSMS will use the default phone number registered with the account
	}

	// Prepare the API request
	endpoint := p.config.BaseURL + ESMSSendSMSEndpoint

	// Build the form parameters
	params := map[string]string{
		"ApiKey":    p.config.APIKey,
		"SecretKey": p.config.Secret,
		"Phone":     req.Message.To,
		"Content":   messageBody,
		"SmsType":   strconv.Itoa(p.config.SMSType),
	}

	// Add sender if available
	if sender != "" {
		params["Brandname"] = sender
	}

	// Add any custom options from the request
	if req.Options != nil {
		// Schedule time if provided (in format YYYY-MM-DD HH:mm:ss)
		if scheduleTime, ok := req.Options["schedule_time"].(string); ok {
			params["TimeSend"] = scheduleTime
		}

		// IsUnicode if provided (0 or 1)
		if isUnicode, ok := req.Options["is_unicode"].(int); ok {
			params["IsUnicode"] = strconv.Itoa(isUnicode)
		}
	}

	// Make the API request to eSMS
	resp, err := p.client.PostForm(ctx, endpoint, params)
	if err != nil {
		return model.SendSMSResponse{}, fmt.Errorf("eSMS API request failed: %w", err)
	}

	// Handle error responses
	if resp.StatusCode() >= 400 {
		return model.SendSMSResponse{}, fmt.Errorf("eSMS API error: %s", resp.String())
	}

	// Parse the response
	var esmsResp esmsSMSResponse
	if err := json.Unmarshal(resp.Body(), &esmsResp); err != nil {
		return model.SendSMSResponse{}, fmt.Errorf("failed to parse eSMS response: %w", err)
	}

	// Check for eSMS error codes
	if esmsResp.CodeResult != "100" {
		return model.SendSMSResponse{}, fmt.Errorf("eSMS error: %s - %s", esmsResp.CodeResult, esmsResp.ErrorMessage)
	}

	// Map eSMS status to our status
	status := mapESMSStatusCode(esmsResp.CodeResult)

	// Convert eSMS response to our response model
	return model.SendSMSResponse{
		MessageID: esmsResp.SMSID,
		Status:    status,
		Provider:  ProviderName,
		SentAt:    time.Now(),
		ProviderResponse: map[string]interface{}{
			"code_result":      esmsResp.CodeResult,
			"sms_id":           esmsResp.SMSID,
			"regenerate_count": esmsResp.CountRegenerate,
		},
	}, nil
}

// SendVoiceCall initiates a voice call using eSMS's OTP voice service
func (p *Provider) SendVoiceCall(ctx context.Context, req model.SendVoiceRequest) (model.SendVoiceResponse, error) {
	// Get the message body from template
	template := req.Template
	if template == "" {
		// If no template provided, use default
		template = "{message}"
	}

	// Render the message template with provided data
	messageText := req.Message.Render(template, req.Data)
	if messageText == "" {
		return model.SendVoiceResponse{}, fmt.Errorf("empty message text after rendering template")
	}

	// Extract the OTP code from the message
	// By default, we'll assume the OTP is 6 digits and try to extract it
	// Otherwise, use the message as is or get it from options
	otp := extractOTPFromMessage(messageText)

	// If OTP is provided in options, use that instead
	if req.Options != nil {
		if otpFromOptions, ok := req.Options["otp"].(string); ok {
			otp = otpFromOptions
		}
	}

	// Prepare the API request
	endpoint := p.config.BaseURL + ESMSVoiceOTPEndpoint

	// Build the form parameters
	params := map[string]string{
		"ApiKey":    p.config.APIKey,
		"SecretKey": p.config.Secret,
		"Phone":     req.Message.To,
		"Code":      otp,
	}

	// Add any custom options from the request
	if req.Options != nil {
		// Speed (if provided) - values 0.5 to 2.0
		if speed, ok := req.Options["speed"].(float64); ok {
			params["Speed"] = fmt.Sprintf("%.1f", speed)
		}

		// Retry times (if provided)
		if retryTimes, ok := req.Options["retry_times"].(int); ok {
			params["Repeat"] = strconv.Itoa(retryTimes)
		}
	}

	// Make the API request to eSMS
	resp, err := p.client.PostForm(ctx, endpoint, params)
	if err != nil {
		return model.SendVoiceResponse{}, fmt.Errorf("eSMS voice API request failed: %w", err)
	}

	// Handle error responses
	if resp.StatusCode() >= 400 {
		return model.SendVoiceResponse{}, fmt.Errorf("eSMS voice API error: %s", resp.String())
	}

	// Parse the response
	var esmsResp esmsVoiceResponse
	if err := json.Unmarshal(resp.Body(), &esmsResp); err != nil {
		return model.SendVoiceResponse{}, fmt.Errorf("failed to parse eSMS voice response: %w", err)
	}

	// Check for eSMS error codes
	if esmsResp.CodeResult != "100" {
		return model.SendVoiceResponse{}, fmt.Errorf("eSMS voice error: %s - %s", esmsResp.CodeResult, esmsResp.ErrorMessage)
	}

	// Map eSMS status to our status
	callStatus := model.CallStatusInitiated
	if esmsResp.CodeResult == "100" {
		callStatus = model.CallStatusInitiated
	} else {
		callStatus = model.CallStatusFailed
	}

	// Convert eSMS response to our response model
	return model.SendVoiceResponse{
		CallID:    esmsResp.CallID,
		Status:    callStatus,
		Provider:  ProviderName,
		StartedAt: time.Now(),
		ProviderResponse: map[string]interface{}{
			"code_result": esmsResp.CodeResult,
			"call_id":     esmsResp.CallID,
		},
	}, nil
}

// mapESMSStatusCode maps eSMS status codes to our status
func mapESMSStatusCode(statusCode string) model.MessageStatus {
	switch statusCode {
	case "100":
		return model.StatusSent
	case "99":
		return model.StatusFailed // Invalid parameter
	case "101", "102", "103", "104", "105":
		return model.StatusFailed // Authentication errors
	case "106", "107", "108", "109", "110", "111", "112", "113", "114":
		return model.StatusFailed // Account or balance errors
	case "118", "119", "120", "121", "122":
		return model.StatusFailed // Recipient errors
	default:
		return model.StatusUnknown
	}
}

// extractOTPFromMessage attempts to extract an OTP code from a message
// It looks for sequences of digits (typically 4-8 digits long for OTPs)
func extractOTPFromMessage(message string) string {
	// Try to find digits of length 6 (common OTP length)
	for i := 0; i <= len(message)-6; i++ {
		if allDigits(message[i : i+6]) {
			return message[i : i+6]
		}
	}

	// Try other common OTP lengths (4, 5, 8 digits)
	for length := 4; length <= 8; length++ {
		if length == 6 {
			continue // Already checked above
		}
		for i := 0; i <= len(message)-length; i++ {
			if allDigits(message[i : i+length]) {
				return message[i : i+length]
			}
		}
	}

	// If no OTP found, return the message itself (the API may handle it differently)
	return message
}

// allDigits checks if a string consists of all digits
func allDigits(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}
