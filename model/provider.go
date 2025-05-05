package model

import (
	"context"
)

// Provider defines the interface for SMS and Voice Call providers
type Provider interface {
	// Name returns the provider's name
	Name() string

	// SendSMS sends an SMS message through the provider
	SendSMS(ctx context.Context, request SendSMSRequest) (SendSMSResponse, error)

	// SendVoiceCall initiates a voice call through the provider
	SendVoiceCall(ctx context.Context, request SendVoiceRequest) (SendVoiceResponse, error)
}
