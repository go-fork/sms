package model

import "fmt"

// SendSMSRequest represents a request to send an SMS
type SendSMSRequest struct {
	// Message contains the core message information (From, To, By)
	Message Message `json:"message"`

	// Template is an optional message template to use
	// If empty, the default template from configuration will be used
	Template string `json:"template,omitempty"`

	// Data contains variables to bind into the template
	Data map[string]interface{} `json:"data,omitempty"`

	// Options contains provider-specific options
	Options map[string]interface{} `json:"options,omitempty"`
}

// SendVoiceRequest represents a request to make a voice call
type SendVoiceRequest struct {
	// Message contains the core message information (From, To, By)
	Message Message `json:"message"`

	// Template is an optional voice script template to use
	// If empty, the default template from configuration will be used
	Template string `json:"template,omitempty"`

	// Data contains variables to bind into the template
	Data map[string]interface{} `json:"data,omitempty"`

	// Options contains provider-specific options such as:
	// - voice_type: The type of voice to use (male/female)
	// - language: The language code (en-US, vi-VN, etc.)
	// - speed: The speech rate (0.8 to 1.2)
	Options map[string]interface{} `json:"options,omitempty"`
}

// Validate performs basic validation on a SendSMSRequest
func (r *SendSMSRequest) Validate() error {
	// Validate phone numbers
	if !ValidatePhoneNumber(r.Message.To) {
		return &ValidationError{Field: "to", Message: "invalid recipient phone number"}
	}

	// Ensure From is not empty
	if r.Message.From == "" {
		return &ValidationError{Field: "from", Message: "sender identifier cannot be empty"}
	}

	return nil
}

// Validate performs basic validation on a SendVoiceRequest
func (r *SendVoiceRequest) Validate() error {
	// Validate phone numbers
	if !ValidatePhoneNumber(r.Message.To) {
		return &ValidationError{Field: "to", Message: "invalid recipient phone number"}
	}

	// Ensure From is not empty
	if r.Message.From == "" {
		return &ValidationError{Field: "from", Message: "sender identifier cannot be empty"}
	}

	return nil
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

// Error returns the error message
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s - %s", e.Field, e.Message)
}
