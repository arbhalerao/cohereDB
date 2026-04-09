package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/arbhalerao/meerkat/db_manager/internal"
	"github.com/arbhalerao/meerkat/pb/db_manager"
	"github.com/arbhalerao/meerkat/utils"

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

func (s *Server) Stop() {
	utils.Logger.Info().Msg("Shutting down gRPC server...")
	s.grpc.GracefulStop()
}

func (s *Server) Get(ctx context.Context, req *db_manager.GetRequest) (*db_manager.GetResponse, error) {
	val, err := s.manager.GetKey(req.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to get key %q: %w", req.Key, err)
	}
	return &db_manager.GetResponse{Value: val}, nil
}

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
