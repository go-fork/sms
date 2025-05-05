package model

import (
	"fmt"
	"time"
)

// MessageStatus represents the status of a message
type MessageStatus string

const (
	// StatusPending indicates the message is queued for delivery
	StatusPending MessageStatus = "pending"

	// StatusSent indicates the message has been sent to the provider
	StatusSent MessageStatus = "sent"

	// StatusDelivered indicates the message has been delivered to the recipient
	StatusDelivered MessageStatus = "delivered"

	// StatusFailed indicates the message delivery has failed
	StatusFailed MessageStatus = "failed"

	// StatusUnknown indicates the message status is unknown
	StatusUnknown MessageStatus = "unknown"
)

// SendSMSResponse represents the response after sending an SMS
type SendSMSResponse struct {
	// MessageID is the unique identifier assigned by the provider
	MessageID string `json:"message_id"`

	// Status represents the delivery status (pending, sent, delivered, failed, unknown)
	Status MessageStatus `json:"status"`

	// Provider is the name of the provider used to send the message
	Provider string `json:"provider"`

	// SentAt is the timestamp when the message was sent
	SentAt time.Time `json:"sent_at"`

	// Cost is the cost of sending the message (if available)
	Cost float64 `json:"cost,omitempty"`

	// Currency is the currency of the cost (if cost is provided)
	Currency string `json:"currency,omitempty"`

	// ProviderResponse contains the raw response from the provider
	ProviderResponse map[string]interface{} `json:"provider_response,omitempty"`
}

// CallStatus represents the status of a voice call
type CallStatus string

const (
	// CallStatusQueued indicates the call is queued
	CallStatusQueued CallStatus = "queued"

	// CallStatusInitiated indicates the call has been initiated
	CallStatusInitiated CallStatus = "initiated"

	// CallStatusRinging indicates the recipient's phone is ringing
	CallStatusRinging CallStatus = "ringing"

	// CallStatusInProgress indicates the call is in progress
	CallStatusInProgress CallStatus = "in-progress"

	// CallStatusCompleted indicates the call has completed successfully
	CallStatusCompleted CallStatus = "completed"

	// CallStatusBusy indicates the recipient's line was busy
	CallStatusBusy CallStatus = "busy"

	// CallStatusNoAnswer indicates the recipient didn't answer
	CallStatusNoAnswer CallStatus = "no-answer"

	// CallStatusFailed indicates the call failed
	CallStatusFailed CallStatus = "failed"

	// CallStatusCanceled indicates the call was canceled
	CallStatusCanceled CallStatus = "canceled"
)

// SendVoiceResponse represents the response after initiating a voice call
type SendVoiceResponse struct {
	// CallID is the unique identifier assigned by the provider
	CallID string `json:"call_id"`

	// Status represents the call status
	Status CallStatus `json:"status"`

	// Provider is the name of the provider used to make the call
	Provider string `json:"provider"`

	// StartedAt is the timestamp when the call was initiated
	StartedAt time.Time `json:"started_at"`

	// EndedAt is the timestamp when the call ended (if completed)
	EndedAt *time.Time `json:"ended_at,omitempty"`

	// Duration is the call duration in seconds (0 if not yet completed)
	Duration int `json:"duration"`

	// Cost is the cost of the call (if available)
	Cost float64 `json:"cost,omitempty"`

	// Currency is the currency of the cost (if cost is provided)
	Currency string `json:"currency,omitempty"`

	// ProviderResponse contains the raw response from the provider
	ProviderResponse map[string]interface{} `json:"provider_response,omitempty"`
}

// String returns a string representation of the SendSMSResponse
func (r *SendSMSResponse) String() string {
	return fmt.Sprintf("SMS [%s] via %s: %s", r.MessageID, r.Provider, r.Status)
}

// String returns a string representation of the SendVoiceResponse
func (r *SendVoiceResponse) String() string {
	return fmt.Sprintf("Voice Call [%s] via %s: %s (Duration: %ds)",
		r.CallID, r.Provider, r.Status, r.Duration)
}
