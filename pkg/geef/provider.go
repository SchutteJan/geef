package geef

import (
	"fmt"
)

type Currency string

const (
	EUR Currency = "eur"
	USD Currency = "usd"
)

type PaymentRequestProvider interface {
	getRedirectURL(amount int64, currency Currency, description *string) (string, error)
}

type ProviderType string

const (
	ProviderBunq ProviderType = "bunq"
)

func NewProvider(providerType ProviderType, config map[string]any) (PaymentRequestProvider, error) {
	switch providerType {
	case ProviderBunq:
		username, ok := config["username"].(string)
		if !ok {
			return nil, fmt.Errorf("bunq provider requires 'username' configuration")
		}
		return &BunqProvider{
			username: username,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported provider type: %s", providerType)
	}
}
