package web

import (
	"fmt"
	"net/http"

	"github.com/dgraph-io/badger/v4"
)

// GetHandler handles read requests
func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	r.ParseForm()
	key := r.Form.Get("key")

	value, err := s.db.GetKey(key)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			http.Error(w, fmt.Sprintf(`{"error": "Key '%s' not found"}`, key), http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf(`{"error": "Failed to get key '%s': %v"`, key, err), http.StatusInternalServerError)
		return
	}

	response := fmt.Sprintf(`{%q}`, value)
	w.Write([]byte(response))
}

// SetHandler handles write requests
func (s *Server) SetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	r.ParseForm()
	key := r.Form.Get("key")
	value := r.Form.Get("value")

	err := s.db.SetKey(key, value)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to set key '%s': %v"}`, key, err), http.StatusInternalServerError)
		return
	}

	response := fmt.Sprintf(`{"message": "Key '%s' set successfully"}`, key)
	w.Write([]byte(response))
}

// DeleteHandler handles delete key requests
func (s *Server) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	r.ParseForm()
	key := r.Form.Get("key")

	err := s.db.DeleteKey(key)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			http.Error(w, fmt.Sprintf(`{"error": "Key '%s' not found"}`, key), http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf(`{"error": "Failed to delete key '%s': %v"`, key, err), http.StatusInternalServerError)
		return
	}

	response := fmt.Sprintf(`{"message": "Key '%s' deleted successfully"}`, key)
	w.Write([]byte(response))
}
