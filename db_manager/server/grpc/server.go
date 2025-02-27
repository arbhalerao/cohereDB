package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/arbha1erao/cohereDB/db_manager/internal"
	"github.com/arbha1erao/cohereDB/pb/db_manager"
	"github.com/arbha1erao/cohereDB/utils"

	"google.golang.org/grpc"
)

type Server struct {
	db_manager.UnimplementedDBManagerServer
	grpc    *grpc.Server
	addr    string
	manager *internal.DBManager
}

func NewServer(addr string, manager *internal.DBManager) *Server {
	grpcServer := grpc.NewServer()
	s := &Server{
		grpc:    grpcServer,
		addr:    addr,
		manager: manager,
	}
	db_manager.RegisterDBManagerServer(grpcServer, s)
	return s
}

// Start initializes the gRPC listener and starts the server.
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		utils.Logger.Fatal().Err(err).Msg("Failed to start gRPC listener")
		return err
	}

	utils.Logger.Info().Msgf("gRPC server listening on %s", s.addr)
	if err := s.grpc.Serve(listener); err != nil {
		utils.Logger.Fatal().Err(err).Msg("gRPC server failed")
		return err
	}

	return nil
}

// Stop gracefully shuts down the gRPC server.
func (s *Server) Stop() {
	utils.Logger.Info().Msg("Shutting down gRPC server...")
	s.grpc.GracefulStop()
}

// Get retrieves the value for a given key from the database
func (s *Server) Get(ctx context.Context, req *db_manager.GetRequest) (*db_manager.GetResponse, error) {
	val, err := s.manager.GetKey(req.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to get key %q: %w", req.Key, err)
	}
	return &db_manager.GetResponse{Value: val}, nil
}

// Set stores the provided key-value pair in the database
func (s *Server) Set(ctx context.Context, req *db_manager.SetRequest) (*db_manager.SetResponse, error) {
	success, err := s.manager.SetKey(req.Key, req.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to set key %q: %w", req.Key, err)
	}
	if !success {
		return nil, fmt.Errorf("failed to set key %q: operation unsuccessful", req.Key)
	}
	return &db_manager.SetResponse{Success: success}, nil
}

// Delete removes the key-value pair for the given key from the database
func (s *Server) Delete(ctx context.Context, req *db_manager.DeleteRequest) (*db_manager.DeleteResponse, error) {
	success, err := s.manager.DeleteKey(req.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to delete key %q: %w", req.Key, err)
	}
	if !success {
		return nil, fmt.Errorf("failed to delete key %q: operation unsuccessful", req.Key)
	}
	return &db_manager.DeleteResponse{Success: success}, nil
}
