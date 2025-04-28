package db_server

import (
	"net/http"

	"github.com/arbha1erao/cohereDB/db"
)

type Server struct {
	db   *db.Database
	addr string
}

func NewServer(db *db.Database, addr string) *Server {
	return &Server{
		db:   db,
		addr: addr,
	}
}

// Method to register handlers
func (s *Server) RegisterHandlers() {
	http.HandleFunc("/get", s.GetHandler)
	http.HandleFunc("/set", s.SetHandler)
	http.HandleFunc("/delete", s.DeleteHandler)

	// Register more handlers here as needed
}

// Method to start the server
func (s *Server) Start() error {
	return http.ListenAndServe(s.addr, nil)
}
