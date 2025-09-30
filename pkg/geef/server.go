package geef

import (
	"embed"
	"html/template"
	"log"
	"net/http"
)

//go:embed templates/*.html
var templatesFS embed.FS

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

func (s *Server) serveIndex(w http.ResponseWriter, r *http.Request, description string) {
	tmpl, err := template.ParseFS(templatesFS, "templates/index.html")
	if err != nil {
		log.Printf("Failed to parse template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := map[string]string{
		"Description": description,
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Failed to execute template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handlePayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Serve index page for root path
	if r.URL.Path == "/" {
		s.serveIndex(w, r, "")
		return
	}

	path, err := ParsePaymentPath(r.URL.Path)
	if err != nil {
		// If parsing fails, check if it's a description-only path
		// (i.e., no valid amount format detected)
		descriptionOnly := r.URL.Path[1:] // Remove leading slash
		if descriptionOnly != "" {
			s.serveIndex(w, r, descriptionOnly)
			return
		}
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
