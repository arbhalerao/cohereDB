package http

import (
	"fmt"
	"net/http"

	"github.com/arbhalerao/cohereDB/utils"
	"github.com/dgraph-io/badger"
)

// GetHandler retrieves the value for a given key from the database
func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := r.ParseForm()
	if err != nil {
		utils.Logger.Error().Msgf("[GET] Error parsing form data")
		http.Error(w, fmt.Sprintf(`{"error": "Failed to forward request: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	key := r.Form.Get("key")

	value, err := s.db.GetKey(key)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			utils.Logger.Error().Msgf("[GET] Key %s not found", key)
			http.Error(w, fmt.Sprintf(`{"error": "Key '%s' not found"}`, key), http.StatusNotFound)
			return
		}
		utils.Logger.Error().Msgf("[GET] Error retrieving key %s: %v", key, err)
		http.Error(w, fmt.Sprintf(`{"error": "Failed to get key '%s': %v"}`, key, err), http.StatusInternalServerError)
		return
	}

	utils.Logger.Info().Msgf("[GET] Successfully retrieved key %s", key)

	response := fmt.Sprintf(`{%q}`, value)
	_, err = w.Write([]byte(response))
	if err != nil {
		utils.Logger.Error().Msgf("Error writing response for key '%s': %v", key, err)
		return
	}
}

// Set SetHandler the provided key-value pair in the database
func (s *Server) SetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := r.ParseForm()
	if err != nil {
		utils.Logger.Error().Msgf("[SET] Error parsing form data")
		http.Error(w, fmt.Sprintf(`{"error": "Failed to forward request: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	key := r.Form.Get("key")
	value := r.Form.Get("value")

	err = s.db.SetKey(key, value)
	if err != nil {
		utils.Logger.Error().Msgf("[SET] Error setting key %s: %v", key, err)
		http.Error(w, fmt.Sprintf(`{"error": "Failed to set key '%s': %v"}`, key, err), http.StatusInternalServerError)
		return
	}

	utils.Logger.Info().Msgf("[SET] Successfully set key %s", key)

	response := fmt.Sprintf(`{"message": "Key '%s' set successfully"}`, key)
	_, err = w.Write([]byte(response))
	if err != nil {
		utils.Logger.Error().Msgf("Error writing response for key '%s': %v", key, err)
		return
	}
}

// DeleteHandler removes the key-value pair for the given key from the database
func (s *Server) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := r.ParseForm()
	if err != nil {
		utils.Logger.Error().Msgf("[DELETE] Error parsing form data")
		http.Error(w, fmt.Sprintf(`{"error": "Failed to forward request: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	key := r.Form.Get("key")

	err = s.db.DeleteKey(key)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			utils.Logger.Error().Msgf("[DELETE] Key %s not found", key)
			http.Error(w, fmt.Sprintf(`{"error": "Key '%s' not found"}`, key), http.StatusNotFound)
			return
		}
		utils.Logger.Error().Msgf("[DELETE] Error deleting key %s: %v", key, err)
		http.Error(w, fmt.Sprintf(`{"error": "Failed to delete key '%s': %v"}`, key, err), http.StatusInternalServerError)
		return
	}

	utils.Logger.Info().Msgf("[DELETE] Successfully deleted key %s", key)

	response := fmt.Sprintf(`{"message": "Key '%s' deleted successfully"}`, key)
	_, err = w.Write([]byte(response))
	if err != nil {
		utils.Logger.Error().Msgf("Error writing response for key '%s': %v", key, err)
		return
	}
}
