package web

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"

	"github.com/Aditya-Bhalerao/cohereDB/db"
	"github.com/Aditya-Bhalerao/cohereDB/utils"
	"github.com/dgraph-io/badger/v4"
)

type Server struct {
	db          *db.Database
	addr        string
	shardIdx    int
	shardCount  int
	serverAddrs *map[int]string
}

func NewServer(db *db.Database, addr string, shardIdx int, shardCount int, serverAddrs *map[int]string) *Server {
	return &Server{
		db:          db,
		addr:        addr,
		shardIdx:    shardIdx,
		shardCount:  shardCount,
		serverAddrs: serverAddrs,
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

// GetHandler handles read requests
func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	r.ParseForm()
	key := r.Form.Get("key")

	targetShard := s.getShard(key)
	if targetShard != s.shardIdx {
		targetAddr := (*s.serverAddrs)[targetShard]
		utils.Logger.Info().Msgf("[GET] Key %s belongs to shard %d, forwarding to %s", key, targetShard, targetAddr)

		resp, err := s.ForwardRequest(targetAddr, r)
		if err != nil {
			utils.Logger.Error().Msgf("[GET] Error forwarding request for key %s to %s: %v", key, targetAddr, err)
			http.Error(w, `{"error": "Failed to forward request"}`, http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		io.Copy(w, resp.Body)
		return
	}

	value, err := s.db.GetKey(key)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			utils.Logger.Error().Msgf("[GET] Key %s not found in shard %d", key, s.shardIdx)
			http.Error(w, fmt.Sprintf(`{"error": "Key '%s' not found"}`, key), http.StatusNotFound)
			return
		}
		utils.Logger.Error().Msgf("[GET] Error retrieving key %s from shard %d: %v", key, s.shardIdx, err)
		http.Error(w, fmt.Sprintf(`{"error": "Failed to get key '%s': %v"}`, key, err), http.StatusInternalServerError)
		return
	}

	utils.Logger.Info().Msgf("[GET] Successfully retrieved key %s from shard %d", key, s.shardIdx)

	response := fmt.Sprintf(`{%q}`, value)
	w.Write([]byte(response))
}

func (s *Server) SetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	r.ParseForm()
	key := r.Form.Get("key")
	value := r.Form.Get("value")

	targetShard := s.getShard(key)
	if targetShard != s.shardIdx {
		targetAddr := (*s.serverAddrs)[targetShard]

		utils.Logger.Info().Msgf("[SET] Key %s belongs to shard %d, forwarding to %s", key, targetShard, targetAddr)

		resp, err := s.ForwardRequest(targetAddr, r)
		if err != nil {
			utils.Logger.Error().Msgf("[SET] Error forwarding request for key %s to %s: %v", key, targetAddr, err)
			http.Error(w, `{"error": "Failed to forward request"}`, http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		io.Copy(w, resp.Body)
		return
	}

	err := s.db.SetKey(key, value)
	if err != nil {
		utils.Logger.Error().Msgf("[SET] Error setting key %s on shard %d: %v", key, s.shardIdx, err)
		http.Error(w, fmt.Sprintf(`{"error": "Failed to set key '%s': %v"}`, key, err), http.StatusInternalServerError)
		return
	}

	utils.Logger.Info().Msgf("[SET] Successfully set key %s on shard %d", key, s.shardIdx)

	response := fmt.Sprintf(`{"message": "Key '%s' set successfully"}`, key)
	w.Write([]byte(response))
}

func (s *Server) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	r.ParseForm()
	key := r.Form.Get("key")

	targetShard := s.getShard(key)
	if targetShard != s.shardIdx {
		targetAddr := (*s.serverAddrs)[targetShard]

		utils.Logger.Info().Msgf("[DELETE] Key %s belongs to shard %d, forwarding to %s", key, targetShard, targetAddr)

		resp, err := s.ForwardRequest(targetAddr, r)
		if err != nil {
			utils.Logger.Error().Msgf("[DELETE] Error forwarding delete request for key %s to %s: %v", key, targetAddr, err)
			http.Error(w, `{"error": "Failed to forward request"}`, http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		io.Copy(w, resp.Body)
		return
	}

	err := s.db.DeleteKey(key)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			utils.Logger.Error().Msgf("[DELETE] Key %s not found in shard %d", key, s.shardIdx)
			http.Error(w, fmt.Sprintf(`{"error": "Key '%s' not found"}`, key), http.StatusNotFound)
			return
		}
		utils.Logger.Error().Msgf("[DELETE] Error deleting key %s from shard %d: %v", key, s.shardIdx, err)
		http.Error(w, fmt.Sprintf(`{"error": "Failed to delete key '%s': %v"}`, key, err), http.StatusInternalServerError)
		return
	}

	utils.Logger.Info().Msgf("[DELETE] Successfully deleted key %s from shard %d", key, s.shardIdx)

	response := fmt.Sprintf(`{"message": "Key '%s' deleted successfully"}`, key)
	w.Write([]byte(response))
}
