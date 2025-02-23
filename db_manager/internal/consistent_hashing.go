package internal

import (
	"hash/crc32"
	"sort"
	"sync"
)

type ConsistentHasher struct {
	mu    sync.RWMutex
	nodes map[string]struct{} // Stores the actual nodes
	ring  []uint32            // Sorted list of hash values
	keys  map[uint32]string   // Hash → Node mapping
}

func NewConsistentHasher() *ConsistentHasher {
	return &ConsistentHasher{
		nodes: make(map[string]struct{}),
		keys:  make(map[uint32]string),
		ring:  []uint32{},
	}
}

// hashKey generates a hash value for a given key
func (h *ConsistentHasher) hashKey(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

// AddNode adds a new node to the hash ring
func (h *ConsistentHasher) AddNode(node string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.nodes[node]; exists {
		return
	}

	h.nodes[node] = struct{}{}

	hash := h.hashKey(node)
	h.keys[hash] = node
	h.ring = append(h.ring, hash)

	// Keep the ring sorted for efficient lookups
	sort.Slice(h.ring, func(i, j int) bool { return h.ring[i] < h.ring[j] })
}

// RemoveNode removes a node from the hash ring
func (h *ConsistentHasher) RemoveNode(node string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.nodes[node]; !exists {
		return
	}

	delete(h.nodes, node)

	hash := h.hashKey(node)
	delete(h.keys, hash)

	newRing := []uint32{}
	for _, h := range h.ring {
		if h != hash {
			newRing = append(newRing, h)
		}
	}
	h.ring = newRing
}

// GetNode finds the nearest node responsible for the given key
func (h *ConsistentHasher) GetNode(key string) (string, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.ring) == 0 {
		return "", false
	}

	hash := h.hashKey(key)

	idx := sort.Search(len(h.ring), func(i int) bool { return h.ring[i] >= hash })

	if idx == len(h.ring) {
		idx = 0
	}

	node, exists := h.keys[h.ring[idx]]
	return node, exists
}

// Reconcile ensures the hash ring matches the given list of nodes
func (h *ConsistentHasher) Reconcile(nodes []string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	newNodes := make(map[string]struct{})
	for _, node := range nodes {
		newNodes[node] = struct{}{}
	}

	for node := range h.nodes {
		if _, exists := newNodes[node]; !exists {
			h.RemoveNode(node)
		}
	}

	for node := range newNodes {
		if _, exists := h.nodes[node]; !exists {
			h.AddNode(node)
		}
	}
}
