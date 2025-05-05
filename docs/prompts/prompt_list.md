# Go SMS Module Implementation Prompts

This document lists the sequential prompts to guide the implementation of the Go SMS module according to the business analysis requirements.

## Prompts Sequence

1. **[Module Structure](prompt01_module_structure.md)** - Setting up the basic module structure, interfaces, and core components
2. **[Configuration Management](prompt02_configuration.md)** - Implementing Viper-based configuration system
3. **[Message Model](prompt03_message_model.md)** - Creating message models and template rendering system
4. **[Retry Logic](prompt04_retry_logic.md)** - Implementing retry mechanism with exponential backoff
5. **[Core SMS Module](prompt05_sms_module.md)** - Building the main SMS module with provider management
6. **[Twilio Adapter](prompt06_twilio_adapter.md)** - Implementing Twilio provider adapter
7. **[eSMS Adapter](prompt07_esms_adapter.md)** - Implementing eSMS provider adapter 
8. **[SpeedSMS Adapter](prompt08_speedsms_adapter.md)** - Implementing SpeedSMS provider adapter
9. **[Unit Tests](prompt09_unit_tests.md)** - Writing comprehensive tests for all components
10. **[Examples & Documentation](prompt10_examples_docs.md)** - Creating usage examples and documentation
11. **[Project Files](prompt11_project_files.md)** - Setting up README, LICENSE, and other project files

Each prompt file contains specific instructions and requirements for implementing that particular component of the SMS module.
