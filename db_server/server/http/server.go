package http

import (
	"context"
	"net/http"
	"time"

	"github.com/arbhalerao/cohereDB/db"
	"github.com/arbhalerao/cohereDB/utils"
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

func (s *Server) RegisterHandlers() {
	mux := http.NewServeMux()
	mux.HandleFunc("/get", s.GetHandler)
	mux.HandleFunc("/set", s.SetHandler)
	mux.HandleFunc("/delete", s.DeleteHandler)
	s.server.Handler = mux
}

func (s *Server) Start() error {
	utils.Logger.Info().Msgf("Starting HTTP server on %s", s.addr)

	return s.server.ListenAndServe()
}

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
