package kv

import (
	"strconv"
	"testing"
)

func BenchmarkPut(b *testing.B) {
	lru := NewLRUPolicy(1000)
	engine := NewInMemoryEngine(lru, nil)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		engine.Put("key"+strconv.Itoa(i), i)
	}
}

func BenchmarkGet(b *testing.B) {
	lru := NewLRUPolicy(1000)
	engine := NewInMemoryEngine(lru, nil)

	// preload data
	for i := 0; i < 100000; i++ {
		engine.Put("key"+strconv.Itoa(i), i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		engine.Get("key500")
	}
}

func BenchmarkDelete(b *testing.B) {
	lru := NewLRUPolicy(1000)
	engine := NewInMemoryEngine(lru, nil)

	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i)

		engine.Put(key, i)

		b.StartTimer()
		engine.Delete(key)
		b.StopTimer()
	}
}
