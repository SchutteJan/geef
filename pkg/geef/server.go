package geef

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

//go:embed templates/*.html
var templatesFS embed.FS

// Server represents the Geef redirect server
type Server struct {
	provider       PaymentRequestProvider
	currency       Currency
	noAutoRedirect bool
}

// NewServer creates a new Geef server instance
func NewServer(currency Currency, providerType ProviderType, providerConfig map[string]any, noAutoRedirect bool) (*Server, error) {
	provider, err := NewProvider(providerType, providerConfig)
	if err != nil {
		return nil, err
	}

	return &Server{
		provider:       provider,
		currency:       currency,
		noAutoRedirect: noAutoRedirect,
	}, nil
}

// Start initializes and starts the HTTP server on the specified port
func (s *Server) Start(addr string) error {
	http.HandleFunc("/", s.handlePayment)

	log.Printf("Server starting on %s with currency %s...", addr, s.currency)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) serveIndex(w http.ResponseWriter, r *http.Request, amount int64, description string) {
	tmpl, err := template.ParseFS(templatesFS, "templates/index.html")
	if err != nil {
		log.Printf("Failed to parse template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Convert amount from cents to decimal string
	amountStr := ""
	if amount > 0 {
		amountFloat := float64(amount) / 100
		amountStr = formatAmount(amountFloat)
	}

	data := map[string]any{
		"Amount":         amountStr,
		"Description":    description,
		"NoAutoRedirect": s.noAutoRedirect,
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Failed to execute template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func formatAmount(amount float64) string {
	// Format with 2 decimal places, using comma as decimal separator (nl-NL style)
	formatted := fmt.Sprintf("%.2f", amount)
	// Replace dot with comma for Dutch formatting
	return strings.Replace(formatted, ".", ",", 1)
}

func (s *Server) handlePayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Serve index page for root path
	if r.URL.Path == "/" {
		s.serveIndex(w, r, 0, "")
		return
	}

	path, err := ParsePaymentPath(r.URL.Path)
	if err != nil {
		// If parsing fails, check if it's a description-only path
		// (i.e., no valid amount format detected)
		descriptionOnly := r.URL.Path[1:] // Remove leading slash
		if descriptionOnly != "" {
			s.serveIndex(w, r, 0, descriptionOnly)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if we should redirect or show the form
	shouldRedirect := !s.noAutoRedirect || r.URL.Query().Get("redirect") == "true"

	if !shouldRedirect {
		// Show form page with amount and description pre-filled
		description := ""
		if path.Description != nil {
			description = *path.Description
		}
		s.serveIndex(w, r, path.Amount, description)
		return
	}

	url, err := s.provider.getRedirectURL(path.Amount, s.currency, path.Description)
	if err != nil {
		log.Printf("Failed to create payment URL: %v", err)
		http.Error(w, "Failed to create payment URL", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
