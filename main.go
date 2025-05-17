package main

import (
	"flag"
	"log"
	"strings"

	"jan.tf/geef/pkg/geef"
)

func main() {
	hostname := flag.String("hostname", "localhost", "Server hostname")
	port := flag.String("port", "8080", "Server port")
	currencyFlag := flag.String("currency", "eur", "Currency code")
	providerFlag := flag.String("provider", "bunq", "Payment provider")

	// Provider-specific flags
	bunqUsername := flag.String("bunq-username", "", "Bunq.me username (required when provider is 'bunq')")

	flag.Parse()

	// Validate currency
	currency := geef.Currency(strings.ToLower(*currencyFlag))
	if currency != geef.EUR && currency != geef.USD {
		log.Fatalf("Unsupported currency: %s. Supported currencies are: eur, usd", *currencyFlag)
	}

	// Validate provider
	provider := geef.ProviderType(strings.ToLower(*providerFlag))
	if provider != geef.ProviderBunq {
		log.Fatalf("Unsupported provider: %s. Supported providers are: bunq", *providerFlag)
	}

	// Create provider config based on selected provider
	// TODO: Clean this up
	providerConfig := make(map[string]any)
	providerConfig["username"] = *bunqUsername

	server, err := geef.NewServer(currency, provider, providerConfig)
	if err != nil {
		log.Fatal("Failed to create server:", err)
	}

	if err := server.Start(*hostname + ":" + *port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
