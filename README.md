# Geef

Simple webserver that redirects to a payment request provider based on path parameters.

I.e this URL `https://mydomain.com/10` gets redirected to a payment request of 10,00 EUR.

## Features

- Whole amounts: `/10` -> 10.00 EUR
- Decimal amounts: `/10.50` -> 10.50 EUR
- Descriptions: `/10/lunch` -> 10.00 EUR with message "lunch"

## Usage

```bash
go run main.go --help
```

## Payment Request Providers

- bunq.me

## Requirements

Go (1.24+): See https://go.dev/

Lefthook (for pre-commit hook):

```bash
go install github.com/evilmartians/lefthook@latest
lefthook install
```
