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
	revolutUsername := flag.String("revolut-username", "", "Revolut.me username (required when provider is 'revolut')")

	flag.Parse()

	// Validate currency
	currency := geef.Currency(strings.ToLower(*currencyFlag))
	if currency != geef.EUR && currency != geef.USD {
		log.Fatalf("Unsupported currency: %s. Supported currencies are: eur, usd", *currencyFlag)
	}

	// Validate provider
	provider := geef.ProviderType(strings.ToLower(*providerFlag))
	if provider != geef.ProviderBunq && provider != geef.ProviderRevolut {
		log.Fatalf("Unsupported provider: %s. Supported providers are: bunq, revolut", *providerFlag)
	}

	// Validate config (for provider)
	if provider == geef.ProviderBunq && *bunqUsername == "" {
		log.Fatal("Argument --bunq-username is required when provider is 'bunq'")
	}

	// Create provider config based on selected provider
	// TODO: Clean this up
	providerConfig := make(map[string]any)
	switch provider {
	case geef.ProviderBunq:
		if *bunqUsername == "" {
			log.Fatal("--bunq-username is required when provider is 'bunq'")
		}
		providerConfig["username"] = *bunqUsername
	case geef.ProviderRevolut:
		if *revolutUsername == "" {
			log.Fatal("--revolut-username is required when provider is 'revolut'")
		}
		providerConfig["username"] = *revolutUsername
	}

	server, err := geef.NewServer(currency, provider, providerConfig)
	if err != nil {
		log.Fatal("Failed to create server:", err)
	}

	if err := server.Start(*hostname + ":" + *port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
