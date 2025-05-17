package geef

import (
	"fmt"
)

type BunqProvider struct {
	username string
}

func (p *BunqProvider) getRedirectURL(amount int64, currency Currency, description *string) (string, error) {
	amountAsFloat := float64(amount) / 100
	if description != nil {
		return fmt.Sprintf("https://bunq.me/%s/%.2f/%s", p.username, amountAsFloat, *description), nil
	}
	return fmt.Sprintf("https://bunq.me/%s/%.2f", p.username, amountAsFloat), nil
}
