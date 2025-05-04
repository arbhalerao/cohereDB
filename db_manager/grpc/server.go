package grpc

import (
	"context"
	"net"

	"github.com/arbha1erao/cohereDB/db_manager/grpc/pb"
	"github.com/arbha1erao/cohereDB/db_manager/internal"
	"github.com/arbha1erao/cohereDB/utils"

	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedDBManagerServer
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
	pb.RegisterDBManagerServer(grpcServer, s)
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
// ToDo
func (s *Server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	return &pb.GetResponse{Value: string("")}, nil
}

// Set stores the provided key-value pair in the database
// ToDo
func (s *Server) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
	return &pb.SetResponse{Success: true}, nil
}

// Delete removes the key-value pair for the given key from the database
// ToDo
func (s *Server) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	return &pb.DeleteResponse{Success: true}, nil
}
