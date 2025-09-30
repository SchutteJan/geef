# AGENTS.md

This file provides guidance to AI code assistants like when working with code in this repository.

## Project Overview

Geef is a minimal Go web application that redirects users to payment request URLs based on URL path parameters. It uses a pluggable provider system to support multiple payment services (currently bunq.me, with tikkie_provider.go stub existing).

### Architecture

**Core Components:**
- `main.go`: Entry point handling CLI flags (provider type, currency, provider-specific config) and server initialization
- `pkg/geef/server.go`: HTTP server with two main handlers:
  - `serveIndex()`: Renders the web form at `/` (with optional pre-filled description)
  - `handlePayment()`: Parses paths and delegates to providers for payment redirect, or falls back to form if path is description-only
- `pkg/geef/templates/index.html`: Go template with Web Component-based form for creating payment requests
- `pkg/geef/provider.go`: Defines the `PaymentRequestProvider` interface (with amount as integer and optional description) and factory function. Provider and currency are specified as runtime arguments to the webserver.
- `pkg/geef/path.go`: Parses URL paths into `PaymentPathArgs` (amount in cents + optional description)
- Provider implementations (e.g., `bunq_provider.go`): Implement `getRedirectURL()` to construct payment URLs for specific payment services

**Design Pattern:**
The provider system is pluggable - new providers can be added by implementing the `PaymentRequestProvider` interface that maps an amount to a payment URL. Amounts are stored internally as integers (cents) to avoid floating-point precision issues. The `ParsePaymentPath` function handles both `.` and `,` as decimal separators and rejects scientific notation.

## Commands

### Development
```bash
# Run the server locally
go run main.go --hostname localhost --port 8080 --provider bunq --bunq-username YOUR_USERNAME

# Run tests
go test ./... -v

# Run a single test
go test ./pkg/geef -v -run TestParsePath

# Format code
go fmt ./...
```

### Docker
```bash
# Build and run with docker-compose
docker-compose up --build

# Build image manually
docker build -t geef .

# Run container
docker run -p 8080:8080 geef --hostname 0.0.0.0 --port 8080 --bunq-username YOUR_USERNAME
```

### Git Hooks
Lefthook is configured for pre-commit hooks:
```bash
# Install hooks (required after clone)
lefthook install
```

The pre-commit hook runs `go fmt` and `go test ./... -v` on all Go files.

## URL Path Format

- `/` → Web form for creating payment requests
- `/10` → 10.00 EUR (redirects to payment provider)
- `/10.50` or `/10,50` → 10.50 EUR (redirects to payment provider)
- `/10/lunch` → 10.00 EUR with description "lunch" (redirects to payment provider)
- `/lunch` → Shows form with description "lunch" pre-filled (any path without valid amount format)

## Adding New Payment Providers

1. Create a new file `pkg/geef/{provider}_provider.go`
2. Implement the `PaymentRequestProvider` interface with `getRedirectURL(amount int64, currency Currency, description *string) (string, error)`
3. Add the provider constant to `pkg/geef/provider.go`
4. Update the `NewProvider()` factory function in `provider.go`
5. Add necessary CLI flags in `main.go` for provider-specific configuration
