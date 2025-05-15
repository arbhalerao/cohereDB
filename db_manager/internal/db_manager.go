package internal

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/arbhalerao/cohereDB/pb/db_server"
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

const ReplicationFactor = 2

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
	ActiveServers.Inc()

	return true
}

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
	ActiveServers.Dec()

	return true
}

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

func (m *DBManager) getReplicaServers(key string) ([]dbServer, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	uuids := m.hasher.GetReplicaNodes(key, ReplicationFactor)
	if len(uuids) == 0 {
		return nil, fmt.Errorf("no available database servers")
	}

	var servers []dbServer
	for _, uuid := range uuids {
		if server, exists := m.servers[uuid]; exists {
			servers = append(servers, server)
		}
	}

	if len(servers) == 0 {
		return nil, fmt.Errorf("no reachable servers for key %q", key)
	}

	return servers, nil
}

func (m *DBManager) GetKey(key string) (string, error) {
	start := time.Now()
	defer func() {
		RequestDuration.WithLabelValues("get").Observe(time.Since(start).Seconds())
	}()

	servers, err := m.getReplicaServers(key)
	if err != nil {
		RequestsTotal.WithLabelValues("get", "error").Inc()
		return "", err
	}

	var lastErr error
	for _, server := range servers {
		resp, err := server.client.Get(context.Background(), &db_server.GetRequest{Key: key})
		if err == nil {
			RequestsTotal.WithLabelValues("get", "success").Inc()
			return resp.Value, nil
		}
		lastErr = err
	}

	RequestsTotal.WithLabelValues("get", "error").Inc()
	return "", fmt.Errorf("all replicas failed for key %q: %v", key, lastErr)
}

func (m *DBManager) SetKey(key, value string) (bool, error) {
	start := time.Now()
	defer func() {
		RequestDuration.WithLabelValues("set").Observe(time.Since(start).Seconds())
	}()

	servers, err := m.getReplicaServers(key)
	if err != nil {
		RequestsTotal.WithLabelValues("set", "error").Inc()
		return false, err
	}

	successCount := 0
	var lastErr error
	for _, server := range servers {
		_, err := server.client.Set(context.Background(), &db_server.SetRequest{Key: key, Value: value})
		if err != nil {
			lastErr = err
			ReplicationWrites.WithLabelValues("failure").Inc()
			continue
		}
		successCount++
		ReplicationWrites.WithLabelValues("success").Inc()
	}

	if successCount == 0 {
		RequestsTotal.WithLabelValues("set", "error").Inc()
		return false, fmt.Errorf("failed to write to any replica for key %q: %v", key, lastErr)
	}

	RequestsTotal.WithLabelValues("set", "success").Inc()
	return true, nil
}

func (m *DBManager) DeleteKey(key string) (bool, error) {
	start := time.Now()
	defer func() {
		RequestDuration.WithLabelValues("delete").Observe(time.Since(start).Seconds())
	}()

	servers, err := m.getReplicaServers(key)
	if err != nil {
		RequestsTotal.WithLabelValues("delete", "error").Inc()
		return false, err
	}

	successCount := 0
	var lastErr error
	for _, server := range servers {
		_, err := server.client.Delete(context.Background(), &db_server.DeleteRequest{Key: key})
		if err != nil {
			lastErr = err
			continue
		}
		successCount++
	}

	if successCount == 0 {
		RequestsTotal.WithLabelValues("delete", "error").Inc()
		return false, fmt.Errorf("failed to delete from any replica for key %q: %v", key, lastErr)
	}

	RequestsTotal.WithLabelValues("delete", "success").Inc()
	return true, nil
}

func (m *DBManager) ReconcileServers() {
	m.mu.Lock()
	activeNodes := make([]string, 0, len(m.servers))
	for uuid := range m.servers {
		activeNodes = append(activeNodes, uuid)
	}
	m.mu.Unlock()

	m.hasher.Reconcile(activeNodes)
}
