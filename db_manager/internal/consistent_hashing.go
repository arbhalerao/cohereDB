package internal

import (
	"hash/crc32"
	"sort"
	"sync"
)

type ConsistentHasher struct {
	mu    sync.RWMutex
	nodes map[string]struct{}
	ring  []uint32
	keys  map[uint32]string
}

func NewConsistentHasher() *ConsistentHasher {
	return &ConsistentHasher{
		nodes: make(map[string]struct{}),
		keys:  make(map[uint32]string),
		ring:  []uint32{},
	}
}

func (h *ConsistentHasher) hashKey(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

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

	sort.Slice(h.ring, func(i, j int) bool { return h.ring[i] < h.ring[j] })
}

func (h *ConsistentHasher) RemoveNode(node string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.nodes[node]; !exists {
		return
	}

	delete(h.nodes, node)

	hash := h.hashKey(node)
	delete(h.keys, hash)

	newRing := make([]uint32, 0, len(h.ring)-1)
	for _, hashVal := range h.ring {
		if hashVal != hash {
			newRing = append(newRing, hashVal)
		}
	}
	h.ring = newRing
}

func (h *ConsistentHasher) GetNode(key string) (string, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.ring) == 0 {
		return "", false
	}

	hash := h.hashKey(key)

	idx := sort.Search(len(h.ring), func(i int) bool {
		return h.ring[i] >= hash
	})

	if idx == len(h.ring) {
		idx = 0
	}

	node, exists := h.keys[h.ring[idx]]
	return node, exists
}

func (h *ConsistentHasher) Reconcile(nodes []string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	newNodes := make(map[string]struct{})
	for _, node := range nodes {
		newNodes[node] = struct{}{}
	}

	for node := range h.nodes {
		if _, exists := newNodes[node]; !exists {
			delete(h.nodes, node)
			hash := h.hashKey(node)
			delete(h.keys, hash)

			newRing := make([]uint32, 0, len(h.ring)-1)
			for _, hashVal := range h.ring {
				if hashVal != hash {
					newRing = append(newRing, hashVal)
				}
			}
			h.ring = newRing
		}
	}

	for node := range newNodes {
		if _, exists := h.nodes[node]; !exists {
			h.nodes[node] = struct{}{}
			hash := h.hashKey(node)
			h.keys[hash] = node
			h.ring = append(h.ring, hash)
		}
	}

	sort.Slice(h.ring, func(i, j int) bool { return h.ring[i] < h.ring[j] })
}

func (h *ConsistentHasher) GetNodes() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	nodes := make([]string, 0, len(h.nodes))
	for node := range h.nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

func (h *ConsistentHasher) Size() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.nodes)
}
