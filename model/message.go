package model

import (
	"fmt"
	"regexp"
	"strings"
)

// Message represents the core message structure used for both SMS and Voice calls
type Message struct {
	// From represents the sender identifier (phone number, brandname, etc.)
	From string `json:"from"`

	// To represents the recipient's phone number
	To string `json:"to"`

	// By represents the application or service that initiated the message
	By string `json:"by"`
}

// Render processes a template string with provided data to generate message content
// It replaces placeholders in the format {key} with corresponding values from data
func (m *Message) Render(template string, data map[string]interface{}) string {
	// If no template is provided, return empty string
	if template == "" {
		return ""
	}

	// Add message fields to the data map if not already present
	if data == nil {
		data = make(map[string]interface{})
	}

	// Add message fields to data if not explicitly provided
	if _, exists := data["from"]; !exists {
		data["from"] = m.From
	}
	if _, exists := data["to"]; !exists {
		data["to"] = m.To
	}
	if _, exists := data["by"]; !exists {
		data["by"] = m.By
	}

	// Regular expression to find placeholders like {key}
	re := regexp.MustCompile(`{([^{}]+)}`)

	// Replace all placeholders with their values
	result := re.ReplaceAllStringFunc(template, func(placeholder string) string {
		// Extract the key from the placeholder (removing { and })
		key := placeholder[1 : len(placeholder)-1]

		// Check if the key exists in the data map
		if value, exists := data[key]; exists {
			return fmt.Sprintf("%v", value)
		}

		// Return the original placeholder if the key doesn't exist
		return placeholder
	})

	return result
}

// ValidatePhoneNumber performs basic validation on a phone number
// Returns true if the phone number appears to be valid
func ValidatePhoneNumber(phoneNumber string) bool {
	// Basic validation for demonstration
	// In a production environment, consider using a dedicated phone number validation library

	// Remove common formatting characters
	cleaned := strings.ReplaceAll(phoneNumber, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")

	// Check if it starts with + and has at least 8 digits
	if strings.HasPrefix(cleaned, "+") {
		// International format
		digits := cleaned[1:]
		if len(digits) >= 8 && regexp.MustCompile(`^\d+$`).MatchString(digits) {
			return true
		}
	} else {
		// National format
		if len(cleaned) >= 8 && regexp.MustCompile(`^\d+$`).MatchString(cleaned) {
			return true
		}
	}

	return false
}
