package db

import (
	"fmt"
	"os"
	"testing"
)

func setupBenchDB(b *testing.B) *Database {
	b.Helper()
	dir, err := os.MkdirTemp("", "meerkat-bench-*")
	if err != nil {
		b.Fatalf("failed to create temp dir: %v", err)
	}

	db, err := NewDatabase(dir)
	if err != nil {
		os.RemoveAll(dir)
		b.Fatalf("failed to create database: %v", err)
	}

	b.Cleanup(func() {
		db.Cleanup()
	})

	return db
}

func BenchmarkSetKey(b *testing.B) {
	db := setupBenchDB(b)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		db.SetKey(fmt.Sprintf("key-%d", i), fmt.Sprintf("value-%d", i))
	}
}

func BenchmarkGetKey(b *testing.B) {
	db := setupBenchDB(b)

	for i := 0; i < 10000; i++ {
		db.SetKey(fmt.Sprintf("key-%d", i), fmt.Sprintf("value-%d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.GetKey(fmt.Sprintf("key-%d", i%10000))
	}
}

func BenchmarkDeleteKey(b *testing.B) {
	db := setupBenchDB(b)

	for i := 0; i < b.N; i++ {
		db.SetKey(fmt.Sprintf("key-%d", i), fmt.Sprintf("value-%d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.DeleteKey(fmt.Sprintf("key-%d", i))
	}
}
