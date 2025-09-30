package geef

import (
	"fmt"
)

type RevolutProvider struct {
	username string
}

func (p *RevolutProvider) getRedirectURL(amount int64, currency Currency, description *string) (string, error) {
	amountAsFloat := float64(amount) / 100
	amountFormatted := fmt.Sprintf("%.0f", amountAsFloat)

	if description != nil {
		return fmt.Sprintf("https://revolut.me/%s/%s%s/%s", p.username, currency, amountFormatted, *description), nil
	}
	return fmt.Sprintf("https://revolut.me/%s/%s%s", p.username, currency, amountFormatted), nil
}
