package geef

import (
	"fmt"
	"strconv"
	"strings"
)

// PaymentPathArgs represents the parsed components of a payment URL path
type PaymentPathArgs struct {
	Amount      int64
	Description *string
}

// ParsePaymentPath parses the payment URL path into amount and description
func ParsePaymentPath(path string) (*PaymentPathArgs, error) {
	// Remove leading slash and split path
	path = strings.TrimPrefix(path, "/")
	if path == "" {
		return nil, fmt.Errorf("empty path")
	}

	// Split path into amount and description
	parts := strings.SplitN(path, "/", 2)

	// Parse amount, handling both . and , as decimal separators
	amountStr := strings.ReplaceAll(parts[0], ",", ".")
	amountStr = strings.TrimSpace(amountStr)

	// Check for scientific notation
	if strings.ContainsAny(amountStr, "eE") {
		return nil, fmt.Errorf("invalid amount format: scientific notation not allowed")
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid amount format: %w", err)
	}

	// Convert to smallest currency unit (e.g., cents)
	amountInCents := int64(amount * 100)

	// Get description if present
	var description *string
	if len(parts) > 1 {
		desc := parts[1]
		description = &desc
	}

	return &PaymentPathArgs{
		Amount:      amountInCents,
		Description: description,
	}, nil
}
