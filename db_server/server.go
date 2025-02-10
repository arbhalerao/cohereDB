package db_server

import (
	"crypto/sha256"
	"fmt"
	"net/http"

	"github.com/arbha1erao/cohereDB/cohere/config"
	"github.com/arbha1erao/cohereDB/cohere/db"
	"github.com/arbha1erao/cohereDB/cohere/utils"
)

type Server struct {
	db          *db.Database
	addr        string
	shardIdx    int
	shardCount  int
	serverAddrs *map[int]string
}

func NewServer(db *db.Database, addr string, config config.Config) *Server {
	peerServers := make(map[int]string)
	for _, srv := range config.PeerServers {
		peerServers[srv.Shard] = srv.Addr
	}

	return &Server{
		db:          db,
		addr:        addr,
		shardIdx:    config.Server.Shard,
		shardCount:  config.Database.ShardCount,
		serverAddrs: &peerServers,
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

// Hash the key and return the shard number
func (s *Server) getShard(key string) int {
	hash := sha256.Sum256([]byte(key))
	return int(hash[0]) % s.shardCount
}

// ForwardRequest forwards an HTTP request to the specified target address and returns the response
func (s *Server) ForwardRequest(targetAddr string, originalReq *http.Request) (*http.Response, error) {
	forwardURL := fmt.Sprintf("http://%s%s", targetAddr, originalReq.URL.Path)

	forwardReq, err := http.NewRequest(originalReq.Method, forwardURL, nil)
	if err != nil {
		utils.Logger.Error().Msgf("[FORWARD] Error creating forward request: %v", err)
		return nil, fmt.Errorf("failed to create forward request: %v", err)
	}

	// Copy all headers from original request
	for key, values := range originalReq.Header {
		for _, value := range values {
			forwardReq.Header.Add(key, value)
		}
	}

	// Copy query parameters and form values
	q := originalReq.URL.Query()
	for key, values := range originalReq.Form {
		for _, value := range values {
			q.Add(key, value)
		}
	}
	forwardReq.URL.RawQuery = q.Encode()

	client := &http.Client{}

	resp, err := client.Do(forwardReq)
	if err != nil {
		utils.Logger.Error().Msgf("[FORWARD] Error forwarding request to %s: %v", targetAddr, err)
		return nil, fmt.Errorf("failed to forward request to %s: %v", targetAddr, err)
	}

	utils.Logger.Info().Msgf("[FORWARD] Successfully forwarded request to %s, status: %d", targetAddr, resp.StatusCode)

	return resp, nil
}
