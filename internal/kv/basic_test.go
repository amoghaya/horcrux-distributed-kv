package kv

import "testing"

func TestPutGet(t *testing.T) {
	lru := NewLRUPolicy(100)
	engine := NewInMemoryEngine(lru, nil)

	engine.Put("name", "harry")

	val, ok := engine.Get("name")

	if !ok {
		t.Fatal("expected key to exist")
	}

	if val != "harry" {
		t.Fatalf("expected harry, got %v", val)
	}
}

func TestDelete(t *testing.T) {
	lru := NewLRUPolicy(100)
	engine := NewInMemoryEngine(lru, nil)

	engine.Put("key", "value")

	engine.Delete("key")

	_, ok := engine.Get("key")

	if ok {
		t.Fatal("expected key to be deleted")
	}
}
