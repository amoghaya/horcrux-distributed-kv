package kv

import (
	"sync"
	"time"
)

// EntryPool reduces GC pressure by reusing Entry objects.
var EntryPool = sync.Pool{
	New: func() interface{} {
		return &Entry{}
	},
}

// GetEntry fetches reusable Entry from pool.
func GetEntry(value interface{}) *Entry {
	e := EntryPool.Get().(*Entry)

	now := time.Now()

	e.Value = value
	e.CreatedAt = now
	e.UpdatedAt = now
	e.ExpiresAt = time.Time{}

	// version for replication/read repair
	e.Version = now.UnixNano()

	return e
}

// PutEntry resets Entry and returns it to pool.
func PutEntry(e *Entry) {
	if e == nil {
		return
	}

	// clear all fields to avoid stale metadata bugs
	e.Value = nil
	e.CreatedAt = time.Time{}
	e.UpdatedAt = time.Time{}
	e.ExpiresAt = time.Time{}
	e.Version = 0

	EntryPool.Put(e)
}
