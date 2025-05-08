package internal

import (
	"context"
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
}

func NewDBManager() *DBManager {
	return &DBManager{
		servers: make(map[string]dbServer),
	}
}

// AddServer registers a new DB server and creates a gRPC connection
func (m *DBManager) AddServer(uuid, region, addr string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.servers[uuid]; exists {
		return false
	}

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

	return true
}

// RemoveServer unregisters a DB server and closes its connection
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

	return true
}

// HealthCheckServers verifies if registered DB servers are alive.
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
		} else {
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
}

// GetKey retrieves a value from the appropriate DB server
func (m *DBManager) GetKey(key string) (string, error) {
	uuid := ""
	m.mu.Lock()
	server, exists := m.servers[uuid]
	m.mu.Unlock()

	if !exists {
		return "", nil
	}

	resp, err := server.client.Get(context.Background(), &db_server.GetRequest{Key: key})
	if err != nil {
		return "", err
	}

	return resp.Value, nil
}

// SetKey stores a key-value pair on a specific DB server
func (m *DBManager) SetKey(key, value string) (bool, error) {
	uuid := ""
	m.mu.Lock()
	server, exists := m.servers[uuid]
	m.mu.Unlock()

	if !exists {
		return false, nil
	}

	_, err := server.client.Set(context.Background(), &db_server.SetRequest{Key: key, Value: value})
	if err != nil {
		return false, err
	}

	return true, nil
}

// DeleteKey removes a key-value pair from a specific DB server
func (m *DBManager) DeleteKey(key string) (bool, error) {
	uuid := ""
	m.mu.Lock()
	server, exists := m.servers[uuid]
	m.mu.Unlock()

	if !exists {
		return false, nil
	}

	_, err := server.client.Delete(context.Background(), &db_server.DeleteRequest{Key: key})
	if err != nil {
		return false, err
	}

	return true, nil
}
