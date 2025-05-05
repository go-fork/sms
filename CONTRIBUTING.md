# Contributing to go-sms

Thank you for your interest in contributing to go-sms! This document provides guidelines and instructions for contributing to this project.

## Table of Contents

- [Setting Up Development Environment](#setting-up-development-environment)
- [Coding Style and Conventions](#coding-style-and-conventions)
- [Pull Request Process](#pull-request-process)
- [Testing Requirements](#testing-requirements)
- [Adding a New Provider](#adding-a-new-provider)
- [Documentation Guidelines](#documentation-guidelines)
- [Release Process](#release-process)

## Setting Up Development Environment

### Prerequisites

- Go 1.18 or later
- Git

### Local Development Setup

1. Fork the repository on GitHub
2. Clone your fork locally:
```bash
git clone https://github.com/yourusername/go-sms.git
cd go-sms
```

3. Add the upstream repository as a remote:
```bash
git remote add upstream https://github.com/zinzinday/go-sms.git
```

4. Create a branch for your work:
```bash
git checkout -b feature/your-feature-name
```

5. Install development dependencies:
```bash
go get -u github.com/stretchr/testify/assert
go get -u github.com/stretchr/testify/require
go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
```

## Coding Style and Conventions

### Code Formatting

- Run `go fmt` before committing changes.
- Use `goimports` to organize imports (standard library first, then third-party, then local imports).
- Consider using `golangci-lint` for additional code quality checks.

### Naming Conventions

- Use meaningful variable and function names.
- Use camelCase for private variables and functions, PascalCase for exported ones.
- Package names should be short, lowercase, and avoid underscores or mixedCaps.
- Interface names that describe behavior should end with "-er" (e.g., `Provider`).

### Error Handling

- Always check errors and return them to the caller unless there's a good reason not to.
- Use `fmt.Errorf()` with the `%w` verb to wrap errors with additional context.
- Define custom error types for specific error conditions.

### Documentation

- Document all exported functions, types, and methods.
- Follow the Go documentation conventions (https://blog.golang.org/godoc).
- Keep comments up-to-date when changing code.
- Add examples when appropriate.

## Pull Request Process

1. Ensure your code follows the style guidelines.
2. Update documentation if necessary.
3. Add or update tests as needed.
4. Run all tests and ensure they pass.
5. Commit your changes with a descriptive commit message.
6. Push your branch to your fork.
7. Submit a pull request to the `main` branch.
8. Respond to any feedback or requested changes.

### Commit Message Format

Follow the [Conventional Commits](https://www.conventionalcommits.org/) format:

```
type(scope): short summary

Detailed description if necessary
```

Types include:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code changes that neither fix bugs nor add features
- `test`: Adding or updating tests
- `chore`: Changes to the build process or auxiliary tools

## Testing Requirements

- All new code should include appropriate tests.
- Aim for at least 80% test coverage for new code.
- Use table-driven tests where appropriate.
- Mock external services and APIs for tests.
- Run tests with the race detector periodically: `go test -race ./...`

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Adding a New Provider

Adding a new SMS provider requires creating a new adapter. Here's how:

1. Create a new directory in `/adapters` for your provider (e.g., `/adapters/newprovider`).
2. Initialize a new Go module:
   ```bash
   cd adapters/newprovider
   go mod init github.com/zinzinday/go-sms/adapters/newprovider
   ```

3. Create the following files:
   - `config.go` - Configuration handling
   - `adapter.go` - Provider implementation
   - `adapter_test.go` - Tests
   - `README.md` - Documentation

4. Implement the `model.Provider` interface:
   ```go
   type Provider struct {
       // Provider-specific fields
   }

   func NewProvider(configFile string) (*Provider, error) {
       // Initialize provider with configuration
   }

   func (p *Provider) Name() string {
       return "newprovider"
   }

   func (p *Provider) SendSMS(ctx context.Context, req model.SendSMSRequest) (model.SendSMSResponse, error) {
       // Implement SMS sending
   }

   func (p *Provider) SendVoiceCall(ctx context.Context, req model.SendVoiceRequest) (model.SendVoiceResponse, error) {
       // Implement voice calling
   }
   ```

5. Add configuration handling in `config.go`:
   ```go
   type NewProviderConfig struct {
       // Provider-specific config fields
   }

   func LoadConfig(v *viper.Viper) (*NewProviderConfig, error) {
       // Load and validate config
   }
   ```

6. Write tests in `adapter_test.go`.
7. Create documentation in `README.md`.

## Documentation Guidelines

- Keep documentation simple, clear, and concise.
- Use code examples to illustrate usage.
- Document all available options and their defaults.
- Include information about error handling and troubleshooting.
- Add diagrams when they help explain complex interactions.

## Release Process

The release process is handled by the project maintainers, but here's how it works:

1. Version numbers follow [Semantic Versioning](https://semver.org/).
2. Changes are documented in CHANGELOG.md.
3. Releases are tagged in Git and published to GitHub.
4. The Go module system handles versioning for users.

Thank you for contributing to go-sms!
