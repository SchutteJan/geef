package geef

import (
	"strings"
	"testing"
)

func TestParsePaymentPath_WholeNumber(t *testing.T) {
	got, err := ParsePaymentPath("/23")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Amount != 2300 {
		t.Errorf("amount = %d, want 2300", got.Amount)
	}
	if got.Description != nil {
		t.Errorf("description = %v, want <nil>", *got.Description)
	}
}

func TestParsePaymentPath_AmountOnly(t *testing.T) {
	got, err := ParsePaymentPath("/10.50")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Amount != 1050 {
		t.Errorf("amount = %d, want 1050", got.Amount)
	}
	if got.Description != nil {
		t.Errorf("description = %v, want <nil>", *got.Description)
	}
}

func TestParsePaymentPath_AmountComma(t *testing.T) {
	got, err := ParsePaymentPath("/10,50")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Amount != 1050 {
		t.Errorf("amount = %d, want 1050", got.Amount)
	}
	if got.Description != nil {
		t.Errorf("description = %v, want <nil>", *got.Description)
	}
}

func TestParsePaymentPath_WithDescription(t *testing.T) {
	got, err := ParsePaymentPath("/10.50/lunch")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Amount != 1050 {
		t.Errorf("amount = %d, want 1050", got.Amount)
	}
	wantDesc := "lunch"
	if got.Description == nil || *got.Description != wantDesc {
		t.Errorf("description = %v, want %q", got.Description, wantDesc)
	}
}

func TestParsePaymentPath_EmptyPath(t *testing.T) {
	_, err := ParsePaymentPath("/")
	if err == nil || !strings.Contains(err.Error(), "empty path") {
		t.Fatalf("error = %v, want to contain %q", err, "empty path")
	}
}

func TestParsePaymentPath_InvalidAmount(t *testing.T) {
	_, err := ParsePaymentPath("/abc")
	if err == nil || !strings.Contains(err.Error(), "invalid amount format") {
		t.Fatalf("error = %v, want to contain %q", err, "invalid amount format")
	}
}

func TestParsePaymentPath_NegativeAmount(t *testing.T) {
	got, err := ParsePaymentPath("/-10.50")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Amount != -1050 {
		t.Errorf("amount = %d, want -1050", got.Amount)
	}
}

func TestParsePaymentPath_LargeAmount(t *testing.T) {
	got, err := ParsePaymentPath("/999999.99")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Amount != 99999999 {
		t.Errorf("amount = %d, want 99999999", got.Amount)
	}
}

func TestParsePaymentPath_MultipleDecimalPlaces(t *testing.T) {
	_, err := ParsePaymentPath("/10.50.50")
	if err == nil || !strings.Contains(err.Error(), "invalid amount format") {
		t.Fatalf("error = %v, want to contain %q", err, "invalid amount format")
	}
}

func TestParsePaymentPath_LeadingZeros(t *testing.T) {
	got, err := ParsePaymentPath("/0010.50")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Amount != 1050 {
		t.Errorf("amount = %d, want 1050", got.Amount)
	}
}

func TestParsePaymentPath_WithSpaces(t *testing.T) {
	got, err := ParsePaymentPath("/ 10.50 ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Amount != 1050 {
		t.Errorf("amount = %d, want 1050", got.Amount)
	}
}

func TestParsePaymentPath_WithCurrencySymbol(t *testing.T) {
	_, err := ParsePaymentPath("/$10.50")
	if err == nil || !strings.Contains(err.Error(), "invalid amount format") {
		t.Fatalf("error = %v, want to contain %q", err, "invalid amount format")
	}
}

func TestParsePaymentPath_SpecialCharacters(t *testing.T) {
	desc := "lunch!@#$%^&*()_+"
	got, err := ParsePaymentPath("/10.50/" + desc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Description == nil || *got.Description != desc {
		t.Errorf("description = %v, want %q", got.Description, desc)
	}
}

func TestParsePaymentPath_ScientificNotation(t *testing.T) {
	_, err := ParsePaymentPath("/1e2")
	if err == nil || !strings.Contains(err.Error(), "invalid amount format") {
		t.Fatalf("error = %v, want to contain %q", err, "invalid amount format")
	}
}

func TestParsePaymentPath_TrailingZeros(t *testing.T) {
	got, err := ParsePaymentPath("/10.500")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Amount != 1050 {
		t.Errorf("amount = %d, want 1050", got.Amount)
	}
}
