package kv

// EvictionPolicy defines behavior for cache eviction strategies.
// core abstraction that makes Horcrux extensible
type EvictionPolicy interface {
	// Access is called whenever a key is read or written
	Access(key string)

	// Insert returns evicted key if eviction happened
	Insert(key string) (evicted string)

	// Delete removes tracking of a key
	Delete(key string)
}
