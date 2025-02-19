package http

import (
	"context"
	"net/http"
	"time"

	"github.com/arbha1erao/cohereDB/db"
	"github.com/arbha1erao/cohereDB/utils"
)

type Server struct {
	db     *db.Database
	addr   string
	server *http.Server
}

func NewServer(db *db.Database, addr string) *Server {
	return &Server{
		db:   db,
		addr: addr,
		server: &http.Server{
			Addr: addr,
		},
	}
}

// RegisterHandlers sets up all HTTP routes
func (s *Server) RegisterHandlers() {
	mux := http.NewServeMux()
	mux.HandleFunc("/get", s.GetHandler)
	mux.HandleFunc("/set", s.SetHandler)
	mux.HandleFunc("/delete", s.DeleteHandler)
	// Register more handlers here as needed
	s.server.Handler = mux
}

// Start launches the HTTP server
func (s *Server) Start() error {
	utils.Logger.Info().Msgf("Starting HTTP server on %s", s.addr)

	return s.server.ListenAndServe()
}

// Shutdown gracefully stops the HTTP server
func (s *Server) Shutdown() error {
	utils.Logger.Info().Msg("Shutting down HTTP server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.server.Shutdown(ctx)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error shutting down HTTP server")
		return err
	}

	utils.Logger.Info().Msg("HTTP server shut down successfully")

	return nil
}
