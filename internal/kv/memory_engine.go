package kv

import (
	"sync"
	"time"
)

// InMemoryEngine stores key-value pairs in memory
// with WAL, TTL, metrics, and eviction support.
type InMemoryEngine struct {
	store map[string]*Entry

	// global lock for map safety
	mu sync.RWMutex

	eviction EvictionPolicy
	metrics  *Metrics
	wal      *WAL
}

// NewInMemoryEngine creates engine with eviction policy injected
func NewInMemoryEngine(eviction EvictionPolicy, wal *WAL) *InMemoryEngine {
	return &InMemoryEngine{
		store:    make(map[string]*Entry),
		eviction: eviction,
		wal:      wal,
		metrics:  NewMetrics(),
	}
}

// Put inserts or updates a key
func (m *InMemoryEngine) Put(key string, value interface{}) {
	start := time.Now()

	m.mu.Lock()
	defer m.mu.Unlock()

	// WAL first for durability
	if m.wal != nil {
		m.wal.LogWrite("PUT", key, value)
	}

	// eviction tracking
	evicted := m.eviction.Insert(key)

	// remove evicted key
	if evicted != "" {
		delete(m.store, evicted)
		m.metrics.IncEviction()
	}

	// update existing
	if entry, exists := m.store[key]; exists {
		entry.Update(value)
	} else {
		m.store[key] = NewEntry(value)
	}

	m.metrics.IncPut()
	m.metrics.RecordPutLatency(time.Since(start))
}

// Get retrieves a value
func (m *InMemoryEngine) Get(key string) (interface{}, bool) {
	start := time.Now()

	m.mu.RLock()
	entry, ok := m.store[key]
	m.mu.RUnlock()

	m.metrics.IncGet()

	if !ok {
		return nil, false
	}

	// TTL check
	if entry.IsExpired() {
		m.metrics.IncExpired()
		return nil, false
	}

	// update LRU/LFU/etc metadata
	m.eviction.Access(key)

	m.metrics.RecordGetLatency(time.Since(start))

	return entry.Value, true
}

// Delete removes a key
func (m *InMemoryEngine) Delete(key string) bool {
	start := time.Now()

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.wal != nil {
		m.wal.LogWrite("DEL", key, "")
	}

	m.eviction.Delete(key)
	delete(m.store, key)

	m.metrics.IncDel()
	m.metrics.RecordDelLatency(time.Since(start))

	return true
}

// cleanupExpired removes expired keys
func (m *InMemoryEngine) cleanupExpired() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for key, entry := range m.store {
		if entry.IsExpired() {
			delete(m.store, key)
			m.eviction.Delete(key)
			m.metrics.IncExpired()
		}
	}
}

// PutWithTTL inserts key with expiration
func (m *InMemoryEngine) PutWithTTL(
	key string,
	value interface{},
	ttlSeconds int,
) {
	start := time.Now()

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.wal != nil {
		m.wal.LogWrite("PUT_TTL", key, value)
	}

	evicted := m.eviction.Insert(key)

	if evicted != "" {
		delete(m.store, evicted)
		m.metrics.IncEviction()
	}

	m.store[key] = NewEntryWithTTL(value, ttlSeconds)

	m.metrics.IncPut()
	m.metrics.RecordPutLatency(time.Since(start))
}

// Keys returns all keys
func (m *InMemoryEngine) Keys() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	keys := make([]string, 0, len(m.store))

	for key := range m.store {
		keys = append(keys, key)
	}

	return keys
}

// GetEntry returns full entry object
func (m *InMemoryEngine) GetEntry(key string) (*Entry, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entry, ok := m.store[key]

	if !ok {
		return nil, false
	}

	if entry.IsExpired() {
		return nil, false
	}

	return entry, true
}
