package internal

import (
	"fmt"
	"testing"
)

func BenchmarkGetNode(b *testing.B) {
	h := NewConsistentHasher()
	for i := 0; i < 10; i++ {
		h.AddNode(fmt.Sprintf("server-%d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.GetNode(fmt.Sprintf("key-%d", i))
	}
}

func BenchmarkGetReplicaNodes(b *testing.B) {
	h := NewConsistentHasher()
	for i := 0; i < 10; i++ {
		h.AddNode(fmt.Sprintf("server-%d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.GetReplicaNodes(fmt.Sprintf("key-%d", i), 2)
	}
}

func BenchmarkAddNode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		h := NewConsistentHasher()
		for j := 0; j < 100; j++ {
			h.AddNode(fmt.Sprintf("server-%d", j))
		}
	}
}

func BenchmarkRemoveNode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		h := NewConsistentHasher()
		for j := 0; j < 100; j++ {
			h.AddNode(fmt.Sprintf("server-%d", j))
		}
		b.StartTimer()
		for j := 0; j < 100; j++ {
			h.RemoveNode(fmt.Sprintf("server-%d", j))
		}
	}
}
