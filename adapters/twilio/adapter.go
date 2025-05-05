package twilio

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/go-fork/sms/client"
	"github.com/go-fork/sms/config"
	"github.com/go-fork/sms/model"
	"github.com/spf13/viper"
)

const (
	// ProviderName is the name of this provider
	ProviderName = "twilio"

	// TwilioBaseURLTemplate is the base URL template for Twilio API
	// Region and AccountSID will be injected into this template
	TwilioBaseURLTemplate = "https://api.%s.twilio.com/%s/Accounts/%s"

	// TwilioSMSEndpoint is the endpoint for sending SMS
	TwilioSMSEndpoint = "/Messages.json"

	// TwilioCallEndpoint is the endpoint for making calls
	TwilioCallEndpoint = "/Calls.json"
)

// Twilio API response structures
type twilioSMSResponse struct {
	SID          string `json:"sid"`
	Status       string `json:"status"`
	DateCreated  string `json:"date_created"`
	DateSent     string `json:"date_sent"`
	Direction    string `json:"direction"`
	Price        string `json:"price"`
	PriceUnit    string `json:"price_unit"`
	ErrorCode    string `json:"error_code,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

type twilioCallResponse struct {
	SID          string `json:"sid"`
	Status       string `json:"status"`
	DateCreated  string `json:"date_created"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
	Duration     string `json:"duration"`
	Price        string `json:"price"`
	PriceUnit    string `json:"price_unit"`
	ErrorCode    string `json:"error_code,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// Provider implements the model.Provider interface for Twilio
type Provider struct {
	// client is the HTTP client for making API requests
	client *client.Client

	// config holds the Twilio provider configuration
	config *TwilioConfig

	// baseURL is the base URL for Twilio API requests
	baseURL string
}

// NewProvider creates a new Twilio provider instance
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

	// Load Twilio-specific configuration
	twilioConfig, err := LoadConfig(v)
	if err != nil {
		return nil, fmt.Errorf("failed to load Twilio configuration: %w", err)
	}

	// Create HTTP client
	httpClient := client.NewClient(cfg)

	// Set basic authentication with Twilio credentials
	httpClient.SetBasicAuth(twilioConfig.AccountSID, twilioConfig.AuthToken)

	// Construct base URL
	baseURL := fmt.Sprintf(
		TwilioBaseURLTemplate,
		twilioConfig.Region,
		twilioConfig.APIVersion,
		twilioConfig.AccountSID,
	)

	return &Provider{
		client:  httpClient,
		config:  twilioConfig,
		baseURL: baseURL,
	}, nil
}

// Name returns the provider name
func (p *Provider) Name() string {
	return ProviderName
}

// SendSMS sends an SMS message using Twilio
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

	// Determine the from number (use the one from request if provided, otherwise use the default)
	from := req.Message.From
	if from == "" {
		from = p.config.FromNumber
	}

	// Build the form data for Twilio API
	formData := map[string]string{
		"To":   req.Message.To,
		"From": from,
		"Body": messageBody,
	}

	// Add any custom options from the request
	if req.Options != nil {
		// Example: Add StatusCallback URL if provided
		if callbackURL, ok := req.Options["status_callback"].(string); ok {
			formData["StatusCallback"] = callbackURL
		}
	}

	// Make the API request to Twilio
	endpoint := p.baseURL + TwilioSMSEndpoint
	resp, err := p.client.PostForm(ctx, endpoint, formData)
	if err != nil {
		return model.SendSMSResponse{}, fmt.Errorf("twilio API request failed: %w", err)
	}

	// Handle error responses
	if resp.StatusCode() >= 400 {
		return model.SendSMSResponse{}, fmt.Errorf("twilio API error: %s", resp.String())
	}

	// Parse the response
	var twilioResp twilioSMSResponse
	if err := json.Unmarshal(resp.Body(), &twilioResp); err != nil {
		return model.SendSMSResponse{}, fmt.Errorf("failed to parse Twilio response: %w", err)
	}

	// If Twilio returned an error
	if twilioResp.ErrorCode != "" {
		return model.SendSMSResponse{}, fmt.Errorf("twilio error: %s - %s", twilioResp.ErrorCode, twilioResp.ErrorMessage)
	}

	// Map Twilio status to our status
	status := mapTwilioSMSStatus(twilioResp.Status)

	// Parse cost if available
	var cost float64 = 0
	if twilioResp.Price != "" {
		fmt.Sscanf(twilioResp.Price, "%f", &cost)
	}

	// Parse sent time
	sentAt := time.Now()
	if twilioResp.DateSent != "" {
		if t, err := time.Parse(time.RFC3339, twilioResp.DateSent); err == nil {
			sentAt = t
		}
	}

	// Convert Twilio response to our response model
	return model.SendSMSResponse{
		MessageID: twilioResp.SID,
		Status:    status,
		Provider:  ProviderName,
		SentAt:    sentAt,
		Cost:      cost,
		Currency:  twilioResp.PriceUnit,
		ProviderResponse: map[string]interface{}{
			"sid":          twilioResp.SID,
			"status":       twilioResp.Status,
			"date_created": twilioResp.DateCreated,
			"date_sent":    twilioResp.DateSent,
			"direction":    twilioResp.Direction,
		},
	}, nil
}

// SendVoiceCall initiates a voice call using Twilio
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

	// Create a TwiML document for the call
	twiml := fmt.Sprintf(`<Response><Say>%s</Say></Response>`,
		url.QueryEscape(messageText))

	// Determine the from number (use the one from request if provided, otherwise use the default)
	from := req.Message.From
	if from == "" {
		from = p.config.FromNumber
	}

	// Build the form data for Twilio API
	formData := map[string]string{
		"To":    req.Message.To,
		"From":  from,
		"Twiml": twiml,
	}

	// Add any custom options from the request
	if req.Options != nil {
		// Voice type
		if voiceType, ok := req.Options["voice"].(string); ok {
			formData["Voice"] = voiceType
		}

		// Language
		if language, ok := req.Options["language"].(string); ok {
			formData["Language"] = language
		}

		// Callback URL
		if callbackURL, ok := req.Options["status_callback"].(string); ok {
			formData["StatusCallback"] = callbackURL
		}
	}

	// Make the API request to Twilio
	endpoint := p.baseURL + TwilioCallEndpoint
	resp, err := p.client.PostForm(ctx, endpoint, formData)
	if err != nil {
		return model.SendVoiceResponse{}, fmt.Errorf("twilio API request failed: %w", err)
	}

	// Handle error responses
	if resp.StatusCode() >= 400 {
		return model.SendVoiceResponse{}, fmt.Errorf("twilio API error: %s", resp.String())
	}

	// Parse the response
	var twilioResp twilioCallResponse
	if err := json.Unmarshal(resp.Body(), &twilioResp); err != nil {
		return model.SendVoiceResponse{}, fmt.Errorf("failed to parse Twilio response: %w", err)
	}

	// If Twilio returned an error
	if twilioResp.ErrorCode != "" {
		return model.SendVoiceResponse{}, fmt.Errorf("twilio error: %s - %s", twilioResp.ErrorCode, twilioResp.ErrorMessage)
	}

	// Map Twilio status to our status
	status := mapTwilioCallStatus(twilioResp.Status)

	// Parse cost if available
	var cost float64 = 0
	if twilioResp.Price != "" {
		fmt.Sscanf(twilioResp.Price, "%f", &cost)
	}

	// Parse start time
	startedAt := time.Now()
	if twilioResp.StartTime != "" {
		if t, err := time.Parse(time.RFC3339, twilioResp.StartTime); err == nil {
			startedAt = t
		}
	}

	// Parse end time if available
	var endedAt *time.Time
	if twilioResp.EndTime != "" {
		if t, err := time.Parse(time.RFC3339, twilioResp.EndTime); err == nil {
			endedAt = &t
		}
	}

	// Parse duration
	var duration int = 0
	if twilioResp.Duration != "" {
		fmt.Sscanf(twilioResp.Duration, "%d", &duration)
	}

	// Convert Twilio response to our response model
	return model.SendVoiceResponse{
		CallID:    twilioResp.SID,
		Status:    status,
		Provider:  ProviderName,
		StartedAt: startedAt,
		EndedAt:   endedAt,
		Duration:  duration,
		Cost:      cost,
		Currency:  twilioResp.PriceUnit,
		ProviderResponse: map[string]interface{}{
			"sid":          twilioResp.SID,
			"status":       twilioResp.Status,
			"date_created": twilioResp.DateCreated,
			"start_time":   twilioResp.StartTime,
			"end_time":     twilioResp.EndTime,
		},
	}, nil
}

// mapTwilioSMSStatus maps Twilio SMS status to our status
func mapTwilioSMSStatus(twilioStatus string) model.MessageStatus {
	switch strings.ToLower(twilioStatus) {
	case "queued":
		return model.StatusPending
	case "sending":
		return model.StatusPending
	case "sent":
		return model.StatusSent
	case "delivered":
		return model.StatusDelivered
	case "undelivered", "failed":
		return model.StatusFailed
	default:
		return model.StatusUnknown
	}
}

// mapTwilioCallStatus maps Twilio call status to our status
func mapTwilioCallStatus(twilioStatus string) model.CallStatus {
	switch strings.ToLower(twilioStatus) {
	case "queued":
		return model.CallStatusQueued
	case "initiated":
		return model.CallStatusInitiated
	case "ringing":
		return model.CallStatusRinging
	case "in-progress":
		return model.CallStatusInProgress
	case "completed":
		return model.CallStatusCompleted
	case "busy":
		return model.CallStatusBusy
	case "no-answer", "failed":
		return model.CallStatusNoAnswer
	case "canceled":
		return model.CallStatusCanceled
	default:
		return model.CallStatusFailed
	}
}
