# Prompt 3: Message Model

## Objective
Implement the message model and template rendering system for the SMS module.

## Required Files to Create/Complete

1. `/model/message.go` - Message model with template rendering capability
2. `/model/request.go` - Request structures for SMS and voice calls
3. `/model/response.go` - Response structures for API results

## Implementation Requirements

### Message Model
- Create a `Message` struct in `model/message.go` with:
  - `From string` - Sender identifier (phone number, brandname)
  - `To string` - Recipient phone number
  - `By string` - Application identifier (e.g., "MyApp")
  
- Implement a template rendering method:
  - `Render(template string, data map[string]interface{}) string`
  - This method should replace placeholders like `{message}`, `{app_name}` with values from the data map
  - Handle cases where a placeholder is missing from the data map

### Request Structures
- Define in `model/request.go`:
  - `SendSMSRequest` struct including:
    - `Message` embedded struct
    - `Template string` - Optional custom template
    - `Data map[string]interface{}` - Data for template rendering
  
  - `SendVoiceRequest` struct with similar fields as `SendSMSRequest`

### Response Structures
- Define in `model/response.go`:
  - `SendSMSResponse` struct with:
    - `MessageID string` - Provider-specific message identifier
    - `Status string` - Delivery status
    - `Provider string` - Provider used to send the message
    - `Cost float64` - Optional cost information
  
  - `SendVoiceResponse` struct with similar fields as `SendSMSResponse`

## Deliverables
- Complete message model with template rendering capability
- Request structures for SMS and voice calls
- Response structures for API results
