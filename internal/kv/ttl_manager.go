package kv

import (
	"sync"
	"time"
)

// TTLManager handles active expiration of keys
type TTLManager struct {
	engine *InMemoryEngine

	stop chan struct{}
	wg   sync.WaitGroup
}

// NewTTLManager creates background cleaner
func NewTTLManager(engine *InMemoryEngine) *TTLManager {
	return &TTLManager{
		engine: engine,
		stop:   make(chan struct{}),
	}
}

// Start runs cleanup loop
func (t *TTLManager) Start() {
	t.wg.Add(1)

	go func() {
		defer t.wg.Done()

		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				t.engine.cleanupExpired()
			case <-t.stop:
				return
			}
		}
	}()
}

// Stop shuts down cleaner
func (t *TTLManager) Stop() {
	close(t.stop)
	t.wg.Wait()
}
