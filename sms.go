package sms

import (
	"context"
	"fmt"
	"time"

	"github.com/zinzinday/go-sms/config"
	"github.com/zinzinday/go-sms/model"
	"github.com/zinzinday/go-sms/retry"
)

// Module represents the main SMS module that manages providers and handles message sending
type Module struct {
	// config holds the module configuration
	config *config.Config

	// providers is a map of registered providers by name
	providers map[string]model.Provider

	// activeProvider is the currently active provider
	activeProvider model.Provider
}

// NewModule creates a new SMS module instance with the given configuration file
func NewModule(configFile string) (*Module, error) {
	// Load configuration from file
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Create a new module with empty providers map
	module := &Module{
		config:    cfg,
		providers: make(map[string]model.Provider),
	}

	return module, nil
}

// AddProvider registers a provider with the module
func (m *Module) AddProvider(provider model.Provider) error {
	providerName := provider.Name()

	// Check if a provider with the same name already exists
	if _, exists := m.providers[providerName]; exists {
		return fmt.Errorf("provider with name '%s' is already registered", providerName)
	}

	// Add the provider to the map
	m.providers[providerName] = provider

	// If this is the first provider or matches the default provider in config, set it as active
	if m.activeProvider == nil || m.config.DefaultProvider == providerName {
		m.activeProvider = provider
	}

	return nil
}

// SwitchProvider changes the active provider to the one with the specified name
func (m *Module) SwitchProvider(name string) error {
	provider, exists := m.providers[name]
	if !exists {
		return fmt.Errorf("provider '%s' not found", name)
	}

	m.activeProvider = provider
	return nil
}

// GetProvider returns a provider by name
func (m *Module) GetProvider(name string) (model.Provider, error) {
	provider, exists := m.providers[name]
	if !exists {
		return nil, fmt.Errorf("provider '%s' not found", name)
	}

	return provider, nil
}

// GetActiveProvider returns the currently active provider
func (m *Module) GetActiveProvider() (model.Provider, error) {
	if m.activeProvider == nil {
		return nil, fmt.Errorf("no active provider set")
	}

	return m.activeProvider, nil
}

// SendSMS sends an SMS message using the active provider with retry logic
func (m *Module) SendSMS(ctx context.Context, req model.SendSMSRequest) (model.SendSMSResponse, error) {
	if m.activeProvider == nil {
		return model.SendSMSResponse{}, fmt.Errorf("no active provider set")
	}

	// Validate the request
	if err := req.Validate(); err != nil {
		return model.SendSMSResponse{}, err
	}

	// Create retry configuration
	retryConfig := retry.Config{
		MaxAttempts:  m.config.RetryAttempts,
		InitialDelay: m.config.RetryDelay,
		MaxDelay:     30 * time.Second, // Maximum delay between retries
		Multiplier:   2.0,              // Exponential backoff multiplier
	}

	// Initialize response variable
	var response model.SendSMSResponse

	// Execute with retry
	err := retry.Do(ctx, retryConfig, func() error {
		var err error
		response, err = m.activeProvider.SendSMS(ctx, req)
		return err
	})

	if err != nil {
		return model.SendSMSResponse{}, fmt.Errorf("failed to send SMS after %d attempts: %w",
			m.config.RetryAttempts, err)
	}

	// Ensure the provider field is set
	if response.Provider == "" {
		response.Provider = m.activeProvider.Name()
	}

	return response, nil
}

// SendVoiceCall initiates a voice call using the active provider with retry logic
func (m *Module) SendVoiceCall(ctx context.Context, req model.SendVoiceRequest) (model.SendVoiceResponse, error) {
	if m.activeProvider == nil {
		return model.SendVoiceResponse{}, fmt.Errorf("no active provider set")
	}

	// Validate the request
	if err := req.Validate(); err != nil {
		return model.SendVoiceResponse{}, err
	}

	// Create retry configuration
	retryConfig := retry.Config{
		MaxAttempts:  m.config.RetryAttempts,
		InitialDelay: m.config.RetryDelay,
		MaxDelay:     30 * time.Second, // Maximum delay between retries
		Multiplier:   2.0,              // Exponential backoff multiplier
	}

	// Initialize response variable
	var response model.SendVoiceResponse

	// Execute with retry
	err := retry.Do(ctx, retryConfig, func() error {
		var err error
		response, err = m.activeProvider.SendVoiceCall(ctx, req)
		return err
	})

	if err != nil {
		return model.SendVoiceResponse{}, fmt.Errorf("failed to send voice call after %d attempts: %w",
			m.config.RetryAttempts, err)
	}

	// Ensure the provider field is set
	if response.Provider == "" {
		response.Provider = m.activeProvider.Name()
	}

	return response, nil
}
