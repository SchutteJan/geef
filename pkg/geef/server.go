package geef

import (
	"log"
	"net/http"
)

// Server represents the Geef redirect server
type Server struct {
	provider PaymentRequestProvider
	currency Currency
}

// NewServer creates a new Geef server instance
func NewServer(currency Currency, providerType ProviderType, providerConfig map[string]any) (*Server, error) {
	provider, err := NewProvider(providerType, providerConfig)
	if err != nil {
		return nil, err
	}

	return &Server{
		provider: provider,
		currency: currency,
	}, nil
}

// Start initializes and starts the HTTP server on the specified port
func (s *Server) Start(addr string) error {
	http.HandleFunc("/", s.handlePayment)

	log.Printf("Server starting on %s with currency %s...", addr, s.currency)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) handlePayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path, err := ParsePaymentPath(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
