package internal

import (
	"context"
	"fmt"
	"sync"

	"github.com/arbha1erao/cohereDB/pb/db_server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type dbServer struct {
	uuid   string
	region string
	addr   string
	conn   *grpc.ClientConn
	client db_server.DBServerClient
}

type DBManager struct {
	mu      sync.Mutex
	servers map[string]dbServer
	hasher  *ConsistentHasher
}

func NewDBManager() *DBManager {
	return &DBManager{
		servers: make(map[string]dbServer),
		hasher:  NewConsistentHasher(),
	}
}

// AddServer registers a new DB server, creates a gRPC connection, and updates the consistent hash ring
func (m *DBManager) AddServer(uuid, region, addr string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.servers[uuid]; exists {
		return false
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return false
	}

	client := db_server.NewDBServerClient(conn)

	m.servers[uuid] = dbServer{
		uuid:   uuid,
		region: region,
		addr:   addr,
		conn:   conn,
		client: client,
	}

	m.hasher.AddNode(uuid)

	return true
}

// RemoveServer unregisters a DB server, closes its connection, and updates the consistent hash ring.
func (m *DBManager) RemoveServer(uuid string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	server, exists := m.servers[uuid]
	if !exists {
		return false
	}

	if server.conn != nil {
		server.conn.Close()
	}

	delete(m.servers, uuid)

	m.hasher.RemoveNode(uuid)

	return true
}

// HealthCheckServers verifies if registered DB servers are alive and removes unresponsive ones
func (m *DBManager) HealthCheckServers() {
	m.mu.Lock()
	servers := make(map[string]dbServer, len(m.servers))
	for k, v := range m.servers {
		servers[k] = v
	}
	m.mu.Unlock()

	var toRemove []string

	for uuid, server := range servers {
		_, err := server.client.HealthCheck(context.Background(), &db_server.HealthCheckRequest{})
		if err != nil {
			toRemove = append(toRemove, uuid)
		}
	}

	// Remove unresponsive servers
	m.mu.Lock()
	for _, uuid := range toRemove {
		if server, exists := m.servers[uuid]; exists {
			if server.conn != nil {
				server.conn.Close()
			}
			delete(m.servers, uuid)
		}
	}
	m.mu.Unlock()

	m.ReconcileServers()
}

// GetKey retrieves a value from the appropriate DB server
func (m *DBManager) GetKey(key string) (string, error) {
	m.mu.Lock()
	uuid, exists := m.hasher.GetNode(key)
	m.mu.Unlock()

	if !exists {
		return "", fmt.Errorf("no available database servers")
	}

	m.mu.Lock()
	server, serverExists := m.servers[uuid]
	m.mu.Unlock()

	if !serverExists {
		return "", fmt.Errorf("server not found: %s", uuid)
	}

	resp, err := server.client.Get(context.Background(), &db_server.GetRequest{Key: key})
	if err != nil {
		return "", err
	}

	return resp.Value, nil
}

// SetKey stores a key-value pair on a specific DB server
func (m *DBManager) SetKey(key, value string) (bool, error) {
	m.mu.Lock()
	uuid, exists := m.hasher.GetNode(key)
	m.mu.Unlock()

	if !exists {
		return false, fmt.Errorf("no available database servers")
	}

	m.mu.Lock()
	server, serverExists := m.servers[uuid]
	m.mu.Unlock()

	if !serverExists {
		return false, fmt.Errorf("server not found: %s", uuid)
	}

	_, err := server.client.Set(context.Background(), &db_server.SetRequest{Key: key, Value: value})
	if err != nil {
		return false, err
	}

	return true, nil
}

// DeleteKey removes a key-value pair from a specific DB server
func (m *DBManager) DeleteKey(key string) (bool, error) {
	m.mu.Lock()
	uuid, exists := m.hasher.GetNode(key)
	m.mu.Unlock()

	if !exists {
		return false, fmt.Errorf("no available database servers")
	}

	m.mu.Lock()
	server, serverExists := m.servers[uuid]
	m.mu.Unlock()

	if !serverExists {
		return false, fmt.Errorf("server not found: %s", uuid)
	}

	_, err := server.client.Delete(context.Background(), &db_server.DeleteRequest{Key: key})
	if err != nil {
		return false, err
	}

	return true, nil
}

// ReconcileServers ensures the hash ring is updated with active servers
func (m *DBManager) ReconcileServers() {
	m.mu.Lock()
	activeNodes := make([]string, 0, len(m.servers))
	for uuid := range m.servers {
		activeNodes = append(activeNodes, uuid)
	}
	m.mu.Unlock()

	m.hasher.Reconcile(activeNodes)
}
