package kv

import (
	"strconv"
	"sync"
	"testing"
)

func TestConcurrentAccess(t *testing.T) {
	lru := NewLRUPolicy(1000)

	engine := NewInMemoryEngine(lru, nil)

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			key := "key-" + strconv.Itoa(i)

			engine.Put(key, i)

			val, ok := engine.Get(key)

			if !ok {
				t.Errorf("missing key %s", key)
			}

			if val != i {
				t.Errorf("expected %d got %v", i, val)
			}
		}(i)
	}

	wg.Wait()
}
