package kv

import (
	"os"
	"testing"
)

func TestWALReplay(t *testing.T) {
	os.Remove("test_wal.log")

	wal, err := NewWAL("test_wal.log")
	if err != nil {
		t.Fatal(err)
	}

	lru := NewLRUPolicy(100)

	engine := NewInMemoryEngine(lru, wal)

	engine.Put("user", "harry")

	wal.Close()

	// simulate restart
	replayWal, err := NewWAL("test_wal.log")
	if err != nil {
		t.Fatal(err)
	}

	recovered := NewInMemoryEngine(lru, replayWal)

	replayWal.Replay(recovered)

	val, ok := recovered.Get("user")

	if !ok {
		t.Fatal("expected recovered key")
	}

	if val != "harry" {
		t.Fatalf("expected harry, got %v", val)
	}

	replayWal.Close()

	os.Remove("test_wal.log")
}
