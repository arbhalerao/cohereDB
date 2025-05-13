package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/arbhalerao/cohereDB/db_manager/internal"
	"github.com/arbhalerao/cohereDB/utils"
	"github.com/google/uuid"
)

type Server struct {
	manager *internal.DBManager
	addr    string
	server  *http.Server
}

type RegisterRequest struct {
	Region   string `json:"region"`
	GRPCAddr string `json:"grpc_addr"`
}

type RegisterResponse struct {
	Success    bool   `json:"success"`
	ServerUUID string `json:"server_uuid,omitempty"`
	Message    string `json:"message"`
}

type HealthResponse struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}

func NewServer(manager *internal.DBManager, addr string) *Server {
	s := &Server{
		manager: manager,
		addr:    addr,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/register", s.registerHandler)
	mux.HandleFunc("/health", s.healthHandler)
	mux.HandleFunc("/servers", s.serversHandler)

	s.server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return s
}

func (s *Server) registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Logger.Error().Err(err).Msg("Failed to decode registration request")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Region == "" || req.GRPCAddr == "" {
		http.Error(w, "Region and grpc_addr are required", http.StatusBadRequest)
		return
	}

	serverUUID := uuid.New().String()

	success := s.manager.AddServer(serverUUID, req.Region, req.GRPCAddr)
	if !success {
		utils.Logger.Error().Msgf("Failed to add server %s", serverUUID)
		response := RegisterResponse{
			Success: false,
			Message: "Failed to register server",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	utils.Logger.Info().Msgf("Successfully registered server %s from region %s at %s",
		serverUUID, req.Region, req.GRPCAddr)

	response := RegisterResponse{
		Success:    true,
		ServerUUID: serverUUID,
		Message:    "Server registered successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := HealthResponse{
		Status: "healthy",
		Time:   time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) serversHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]interface{}{
		"message": "Server list endpoint - implementation pending",
		"status":  "ok",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) Start() error {
	utils.Logger.Info().Msgf("Starting HTTP server on %s", s.addr)
	return s.server.ListenAndServe()
}

func (s *Server) Stop() error {
	utils.Logger.Info().Msg("Shutting down HTTP server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}
