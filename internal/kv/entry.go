// Entry represents a stored value in Horcrux.
// This is NOT just a value; it carries metadata needed for real DB systems.
package kv

import "time"

// Entry represents stored value + metadata
type Entry struct {
	Value     interface{}
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time

	Version int64
}

// NewEntry creates normal entry
func NewEntry(value interface{}) *Entry {
	now := time.Now()

	return &Entry{
		Value:     value,
		CreatedAt: now,
		UpdatedAt: now,
		Version:   now.UnixNano(),
	}
}

// NewEntryWithTTL creates expiring entry
func NewEntryWithTTL(value interface{}, ttlSeconds int) *Entry {
	now := time.Now()

	return &Entry{
		Value:     value,
		CreatedAt: now,
		UpdatedAt: now,
		ExpiresAt: now.Add(time.Duration(ttlSeconds) * time.Second),
		Version:   now.UnixNano(),
	}
}

// Update modifies value
func (e *Entry) Update(value interface{}) {
	e.Value = value
	e.UpdatedAt = time.Now()

	// version bump
	e.Version = time.Now().UnixNano()
}

// IsExpired checks expiration
func (e *Entry) IsExpired() bool {
	if e.ExpiresAt.IsZero() {
		return false
	}

	return time.Now().After(e.ExpiresAt)
}
