package db_server

import (
	"context"
	"fmt"

	"github.com/arbha1erao/cohereDB/db"
	"github.com/arbha1erao/cohereDB/db_server/grpc/pb"
)

type server struct {
	pb.UnimplementedDBServiceServer
	db *db.Database
}

func NewServer(db *db.Database) *server {
	return &server{db: db}
}

// Get retrieves the value for a given key from the database
func (s *server) Get(ctx context.Context, req *pb.KeyRequest) (*pb.KeyResponse, error) {
	val, err := s.db.GetKey(req.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to get key '%s': %v", req.Key, err)
	}

	return &pb.KeyResponse{Value: string(val)}, nil
}

// Set stores the provided key-value pair in the database
func (s *server) Set(ctx context.Context, req *pb.KeyValueRequest) (*pb.SetResponse, error) {
	err := s.db.SetKey(req.Key, req.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to set key '%s' with value '%s': %v", req.Key, req.Value, err)
	}

	return &pb.SetResponse{Success: true}, nil
}

// Delete removes the key-value pair for the given key from the database
func (s *server) Delete(ctx context.Context, req *pb.KeyRequest) (*pb.DeleteResponse, error) {
	err := s.db.DeleteKey(req.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to delete key '%s': %v", req.Key, err)
	}

	return &pb.DeleteResponse{Success: true}, nil
}
