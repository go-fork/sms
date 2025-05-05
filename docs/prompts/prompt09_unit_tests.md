# Prompt 9: Unit Tests

## Objective
Implement comprehensive unit tests for all components of the SMS module.

## Required Files to Create

1. In the `/tests/` directory:
   - `client_test.go` - Tests for the HTTP client
   - `sms_test.go` - Tests for the main module
   - `config_test.go` - Tests for configuration handling
   - `message_test.go` - Tests for message model and template

## Implementation Requirements

### General Testing Guidelines
- Use the `github.com/stretchr/testify` package for assertions
- Implement table-driven tests where appropriate
- Mock external dependencies (HTTP requests, file system)
- Aim for at least 80% test coverage

### Client Tests
- In `client_test.go`:
  - Test client creation with various configurations
  - Test HTTP request handling (using httptest.Server)
  - Test retry logic
  - Test timeout handling

### SMS Module Tests
- In `sms_test.go`:
  - Test module creation with valid and invalid configurations
  - Test provider management (adding, switching, retrieving)
  - Test sending SMS and voice calls (with mocked providers)
  - Test error scenarios and edge cases

### Configuration Tests
- In `config_test.go`:
  - Test loading configurations from YAML, JSON, and other formats
  - Test validation of required fields
  - Test handling of invalid configurations
  - Test default values

### Message Model Tests
- In `message_test.go`:
  - Test message model creation
  - Test template rendering with various data inputs
  - Test handling of missing template variables
  - Test edge cases (empty templates, nil data maps)

## Test Data
- Create test fixtures for:
  - Sample configuration files
  - Sample requests and responses
  - Mock provider implementations

## Deliverables
- Comprehensive unit tests for all components
- Test helpers and mock implementations
- Documentation of test scenarios and coverage
