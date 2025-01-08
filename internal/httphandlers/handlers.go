package httphandlers

import (
	"encoding/json"
	"net/http"

	"github.com/buildwithme/ethparser/internal/parser"
)

// or wherever your Parser interface is
// for storage.Transaction

// Handlers wraps a Parser instance to serve HTTP requests.
type Handlers struct {
	Parser parser.Parser
}

// New returns a struct with all route handlers bound to a Parser.
func New(p parser.Parser) *Handlers {
	return &Handlers{Parser: p}
}

func (s *Handlers) RegisterHandlers() {
	// GET /current-block
	http.HandleFunc("/current-block", s.HandleCurrentBlock)

	// POST /subscribe?address=0x123...
	http.HandleFunc("/subscribe", s.HandleSubscribe)

	// GET /transactions?address=0x123...
	http.HandleFunc("/transactions", s.HandleTransactions)
}

// HandleCurrentBlock responds with the last parsed block.
//   - Only supports GET, otherwise 405 Method Not Allowed
func (h *Handlers) HandleCurrentBlock(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		currentBlock := h.Parser.GetCurrentBlock()
		resp := map[string]int{"currentBlock": currentBlock}
		writeJSON(w, resp)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleSubscribe adds an address to the observer list.
//   - Expects POST with `address` query param
//   - Returns {"subscribed":true|false}
//   - Responds 400 if `address` is missing, or 405 for non-POST
func (h *Handlers) HandleSubscribe(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		address := r.URL.Query().Get("address")
		if address == "" {
			http.Error(w, "Missing 'address' query parameter", http.StatusBadRequest)
			return
		}
		subscribed := h.Parser.Subscribe(address)
		writeJSON(w, map[string]bool{"subscribed": subscribed})
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleTransactions returns inbound/outbound transactions for a given address.
//   - Expects GET with `address` query param
//   - Returns []storage.Transaction in JSON
//   - Responds 400 if `address` is missing, or 405 for non-GET
func (h *Handlers) HandleTransactions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		address := r.URL.Query().Get("address")
		if address == "" {
			http.Error(w, "Missing 'address' query parameter", http.StatusBadRequest)
			return
		}
		txs := h.Parser.GetTransactions(address)
		writeJSON(w, txs)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// writeJSON is a small helper to consistently write JSON responses.
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(data)
}
