package internal

import (
	"sync"
)

type dbServer struct {
	uuid   string
	region string
	addr   string
}

type DBManager struct {
	mu      sync.Mutex
	servers map[string]dbServer
}

func NewDBManager() *DBManager {
	return &DBManager{
		servers: make(map[string]dbServer),
	}
}

// AddServer handles adding a new DB server to the cluster.
func (m *DBManager) AddServer(uuid, region, addr string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.servers[uuid]; exists {
		return false
	}

	m.servers[uuid] = dbServer{
		uuid:   uuid,
		region: region,
		addr:   addr,
	}

	return true
}

// RemoveServer removes a database server from the cluster.
func (m *DBManager) RemoveServer(uuid string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.servers[uuid]; !exists {
		return false
	}

	delete(m.servers, uuid)

	return true
}

// HealthCheckServers verifies if registered DB servers are alive.
func (m *DBManager) HealthCheckServers() {
	// Ping each server to ensure it's alive
}

// GetKey retrieves the value associated with a given key from the appropriate DB server
func (m *DBManager) GetKey(key string) (string, error) {
	return "", nil
}

// SetKey stores a key-value pair in the appropriate DB server
func (m *DBManager) SetKey(key, value string) (bool, error) {
	return true, nil
}

// DeleteKey removes a key-value pair from the appropriate DB server
func (m *DBManager) DeleteKey(key string) (bool, error) {
	return true, nil
}
