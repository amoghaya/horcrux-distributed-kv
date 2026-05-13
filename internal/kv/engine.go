package kv

// StorageEngine defines the contract for all storage backends.
// This abstraction is CRITICAL for future distributed system design.
type StorageEngine interface {
	Put(key string, value interface{})
	Get(key string) (interface{}, bool)
	Delete(key string) bool
	PutWithTTL(key string, value interface{}, ttlSeconds int)
}
