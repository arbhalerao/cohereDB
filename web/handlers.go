package web

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Aditya-Bhalerao/cohereDB/utils"
	"github.com/dgraph-io/badger/v4"
)

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
