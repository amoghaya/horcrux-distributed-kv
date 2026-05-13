package kv

import (
	"sync"
	"time"
)

// Metrics tracks runtime stats of the KV engine
type Metrics struct {
	mu sync.RWMutex

	// operation counters
	getCount int64
	putCount int64
	delCount int64

	// eviction + ttl
	evictions int64
	expired   int64

	// latency tracking
	getLatency []time.Duration
	putLatency []time.Duration
}

// NewMetrics creates metrics instance
func NewMetrics() *Metrics {
	return &Metrics{
		getLatency: make([]time.Duration, 0),
		putLatency: make([]time.Duration, 0),
	}
}

//----------counters----------//

func (m *Metrics) IncGet() {
	m.mu.Lock()
	m.getCount++
	m.mu.Unlock()
}

func (m *Metrics) IncPut() {
	m.mu.Lock()
	m.putCount++
	m.mu.Unlock()
}

func (m *Metrics) IncDel() {
	m.mu.Lock()
	m.delCount++
	m.mu.Unlock()
}

func (m *Metrics) IncEviction() {
	m.mu.Lock()
	m.evictions++
	m.mu.Unlock()
}

func (m *Metrics) IncExpired() {
	m.mu.Lock()
	m.expired++
	m.mu.Unlock()
}

//latency tracking

func (m *Metrics) RecordGetLatency(d time.Duration) {
	m.mu.Lock()
	m.getLatency = append(m.getLatency, d)
	m.mu.Unlock()
}

func (m *Metrics) RecordPutLatency(d time.Duration) {
	m.mu.Lock()
	m.putLatency = append(m.putLatency, d)
	m.mu.Unlock()
}

func (m *Metrics) RecordDelLatency(d time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// optionally store separate slice
}

//BASIC REPORT FUNCTION

func (m *Metrics) Snapshot() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return map[string]interface{}{
		"gets":      m.getCount,
		"puts":      m.putCount,
		"deletes":   m.delCount,
		"evictions": m.evictions,
		"expired":   m.expired,
	}
}
