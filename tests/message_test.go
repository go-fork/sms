package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zinzinday/go-sms/model"
)

// TestMessageCreation tests creating message models
func TestMessageCreation(t *testing.T) {
	tests := []struct {
		name     string
		from     string
		to       string
		by       string
		expected model.Message
	}{
		{
			name: "Valid message",
			from: "Sender",
			to:   "+1234567890",
			by:   "TestApp",
			expected: model.Message{
				From: "Sender",
				To:   "+1234567890",
				By:   "TestApp",
			},
		},
		{
			name: "Empty fields",
			from: "",
			to:   "",
			by:   "",
			expected: model.Message{
				From: "",
				To:   "",
				By:   "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := model.Message{
				From: tt.from,
				To:   tt.to,
				By:   tt.by,
			}
			assert.Equal(t, tt.expected, msg)
		})
	}
}

// TestTemplateRendering tests the template rendering function
func TestTemplateRendering(t *testing.T) {
	tests := []struct {
		name     string
		message  model.Message
		template string
		data     map[string]interface{}
		expected string
	}{
		{
			name: "Basic template with simple variables",
			message: model.Message{
				From: "Sender",
				To:   "+1234567890",
				By:   "TestApp",
			},
			template: "Your message from {app_name}: {message}",
			data: map[string]interface{}{
				"app_name": "TestApp",
				"message":  "Hello, World!",
			},
			expected: "Your message from TestApp: Hello, World!",
		},
		{
			name: "Template with message fields",
			message: model.Message{
				From: "Sender",
				To:   "+1234567890",
				By:   "TestApp",
			},
			template: "From: {from}, To: {to}, By: {by}, Message: {message}",
			data: map[string]interface{}{
				"message": "Hello, World!",
			},
			expected: "From: Sender, To: +1234567890, By: TestApp, Message: Hello, World!",
		},
		{
			name: "Message fields override data map",
			message: model.Message{
				From: "Sender",
				To:   "+1234567890",
				By:   "TestApp",
			},
			template: "From: {from}, To: {to}, By: {by}",
			data: map[string]interface{}{
				"from": "OtherSender",    // Should not override message.From
				"to":   "OtherRecipient", // Should not override message.To
				"by":   "OtherApp",       // Should not override message.By
			},
			expected: "From: Sender, To: +1234567890, By: TestApp",
		},
		{
			name: "Missing variables in data",
			message: model.Message{
				From: "Sender",
				To:   "+1234567890",
				By:   "TestApp",
			},
			template: "Your message: {message}, Code: {code}",
			data: map[string]interface{}{
				"message": "Hello, World!",
				// 'code' is missing
			},
			expected: "Your message: Hello, World!, Code: {code}",
		},
		{
			name: "Empty template",
			message: model.Message{
				From: "Sender",
				To:   "+1234567890",
				By:   "TestApp",
			},
			template: "",
			data: map[string]interface{}{
				"message": "Hello, World!",
			},
			expected: "",
		},
		{
			name: "Nil data map",
			message: model.Message{
				From: "Sender",
				To:   "+1234567890",
				By:   "TestApp",
			},
			template: "From: {from}, To: {to}, By: {by}",
			data:     nil,
			expected: "From: Sender, To: +1234567890, By: TestApp",
		},
		{
			name: "Multiple occurrences of the same variable",
			message: model.Message{
				From: "Sender",
				To:   "+1234567890",
				By:   "TestApp",
			},
			template: "App: {app_name}, Using: {app_name}, By: {app_name}",
			data: map[string]interface{}{
				"app_name": "TestApp",
			},
			expected: "App: TestApp, Using: TestApp, By: TestApp",
		},
		{
			name: "Different variable types",
			message: model.Message{
				From: "Sender",
				To:   "+1234567890",
				By:   "TestApp",
			},
			template: "String: {string}, Number: {number}, Boolean: {boolean}, Float: {float}",
			data: map[string]interface{}{
				"string":  "text",
				"number":  123,
				"boolean": true,
				"float":   3.14,
			},
			expected: "String: text, Number: 123, Boolean: true, Float: 3.14",
		},
		{
			name: "Nested structures",
			message: model.Message{
				From: "Sender",
				To:   "+1234567890",
				By:   "TestApp",
			},
			template: "Complex structure: {complex}",
			data: map[string]interface{}{
				"complex": map[string]interface{}{
					"nested": "value",
				},
			},
			expected: "Complex structure: map[nested:value]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.message.Render(tt.template, tt.data)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestMessageValidation tests the message validation functions
func TestPhoneNumberValidation(t *testing.T) {
	tests := []struct {
		name        string
		phoneNumber string
		isValid     bool
	}{
		{
			name:        "Valid international format with +",
			phoneNumber: "+1234567890",
			isValid:     true,
		},
		{
			name:        "Valid international format with spaces",
			phoneNumber: "+1 234 567 890",
			isValid:     true,
		},
		{
			name:        "Valid international format with hyphens",
			phoneNumber: "+1-234-567-890",
			isValid:     true,
		},
		{
			name:        "Valid international format with parentheses",
			phoneNumber: "+1 (234) 567-890",
			isValid:     true,
		},
		{
			name:        "Valid national format",
			phoneNumber: "1234567890",
			isValid:     true,
		},
		{
			name:        "Valid national format with spaces",
			phoneNumber: "123 456 7890",
			isValid:     true,
		},
		{
			name:        "Valid national format with hyphens",
			phoneNumber: "123-456-7890",
			isValid:     true,
		},
		{
			name:        "Too short",
			phoneNumber: "+123",
			isValid:     false,
		},
		{
			name:        "Contains letters",
			phoneNumber: "+1234abcdef",
			isValid:     false,
		},
		{
			name:        "Empty string",
			phoneNumber: "",
			isValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := model.ValidatePhoneNumber(tt.phoneNumber)
			assert.Equal(t, tt.isValid, result)
		})
	}
}

// TestRequestValidation tests validating request structures
func TestRequestValidation(t *testing.T) {
	tests := []struct {
		name        string
		request     interface{}
		expectError bool
	}{
		{
			name: "Valid SMS request",
			request: model.SendSMSRequest{
				Message: model.Message{
					From: "Sender",
					To:   "+1234567890",
					By:   "TestApp",
				},
			},
			expectError: false,
		},
		{
			name: "Invalid phone number in SMS request",
			request: model.SendSMSRequest{
				Message: model.Message{
					From: "Sender",
					To:   "invalid",
					By:   "TestApp",
				},
			},
			expectError: true,
		},
		{
			name: "Empty sender in SMS request",
			request: model.SendSMSRequest{
				Message: model.Message{
					From: "",
					To:   "+1234567890",
					By:   "TestApp",
				},
			},
			expectError: true,
		},
		{
			name: "Valid voice request",
			request: model.SendVoiceRequest{
				Message: model.Message{
					From: "Sender",
					To:   "+1234567890",
					By:   "TestApp",
				},
			},
			expectError: false,
		},
		{
			name: "Invalid phone number in voice request",
			request: model.SendVoiceRequest{
				Message: model.Message{
					From: "Sender",
					To:   "invalid",
					By:   "TestApp",
				},
			},
			expectError: true,
		},
		{
			name: "Empty sender in voice request",
			request: model.SendVoiceRequest{
				Message: model.Message{
					From: "",
					To:   "+1234567890",
					By:   "TestApp",
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			switch req := tt.request.(type) {
			case model.SendSMSRequest:
				err = req.Validate()
			case model.SendVoiceRequest:
				err = req.Validate()
			}

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestResponseStringMethods tests the String() methods of response models
func TestResponseStringMethods(t *testing.T) {
	// Test SMS response
	smsResp := model.SendSMSResponse{
		MessageID: "msg_123",
		Status:    model.StatusSent,
		Provider:  "test_provider",
	}
	assert.Contains(t, smsResp.String(), "msg_123")
	assert.Contains(t, smsResp.String(), "test_provider")
	assert.Contains(t, smsResp.String(), string(model.StatusSent))

	// Test voice response
	voiceResp := model.SendVoiceResponse{
		CallID:   "call_123",
		Status:   model.CallStatusInitiated,
		Provider: "test_provider",
		Duration: 30,
	}
	assert.Contains(t, voiceResp.String(), "call_123")
	assert.Contains(t, voiceResp.String(), "test_provider")
	assert.Contains(t, voiceResp.String(), string(model.CallStatusInitiated))
	assert.Contains(t, voiceResp.String(), "30s")
}

// TestValidationErrorString tests the ValidationError.Error() method
func TestValidationErrorString(t *testing.T) {
	err := &model.ValidationError{
		Field:   "to",
		Message: "invalid phone number",
	}
	assert.Contains(t, err.Error(), "to")
	assert.Contains(t, err.Error(), "invalid phone number")
}
