package web

import (
	"fmt"
	"net/http"

	"github.com/Aditya-Bhalerao/cohereDB/db"
	"github.com/dgraph-io/badger/v4"
)

type Server struct {
	db   *db.Database
	addr string
}

func NewServer(db *db.Database, addr string) *Server {
	server := Server{
		db:   db,
		addr: addr,
	}

	return &server
}

// Method to register handlers
func (s *Server) RegisterHandlers() {
	http.HandleFunc("/get", s.GetHandler)
	http.HandleFunc("/set", s.SetHandler)
	// Register more handlers here as needed

}

// Method to start the server
func (s *Server) Start() error {
	return http.ListenAndServe(s.addr, nil)
}

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

	response := fmt.Sprintf(`{"key": %q, "value": %q}`, key, value)
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
