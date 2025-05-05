# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- Support for additional providers (Plivo, Stringee)
- Rate limiting capabilities
- Message delivery status tracking

### Changed
- Improved error handling for timeout scenarios
- Enhanced template rendering performance

## [1.0.0] - 2023-07-01
### Added
- Initial release of the go-sms module
- Core functionality for sending SMS and voice messages
- Support for multiple providers through adapter pattern:
  - Twilio provider adapter
  - eSMS provider adapter (Vietnam)
  - SpeedSMS provider adapter (Vietnam)
- Provider management with easy switching capability
- Message model with template rendering capabilities
- Configuration management using Viper
- Retry mechanism with exponential backoff
- Comprehensive error handling and validation
- Example applications for simple usage, templates, and multi-provider scenarios
- Full documentation and user guides

### Security
- Secure credential handling through configuration
- Support for different authentication methods based on provider requirements

## Version Numbering Explanation

Given a version number MAJOR.MINOR.PATCH, increments are determined as follows:

1. MAJOR version when incompatible API changes are made
2. MINOR version when functionality is added in a backward compatible manner
3. PATCH version when backward compatible bug fixes are implemented

## How to Upgrade

### Upgrading from beta versions to 1.0.0
- Update provider initialization to use the new `NewProvider(configFile)` pattern
- Replace any direct client usage with the module-level Send methods
- Update configuration files to match the new schema (see example.yaml)

## Contributors
- Initial development team
- Community contributors (see GitHub contributors page)
