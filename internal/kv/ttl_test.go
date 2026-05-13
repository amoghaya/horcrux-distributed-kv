package kv

import (
	"testing"
	"time"
)

func TestTTLExpiration(t *testing.T) {
	lru := NewLRUPolicy(100)
	engine := NewInMemoryEngine(lru, nil)

	engine.PutWithTTL("session", "active", 1)

	time.Sleep(2 * time.Second)

	_, ok := engine.Get("session")

	if ok {
		t.Fatal("expected key to expire")
	}
}
