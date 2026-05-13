// Storage is a simplified interface for cluster nodes.
// This abstraction allows the Coordinator to interact with any KVStore implementation
//behavior contract
// This is a critical abstraction for the cluster design, as it decouples the coordinator from specific KVStore implementations.

package kv

// Storage defines cluster-facing storage behavior.
type Storage interface {
	Put(key string, value interface{})
	Get(key string) (interface{}, bool)
	Delete(key string) bool

	PutWithTTL(key string, value interface{}, ttlSeconds int)

	// needed for replication + read repair
	GetEntry(key string) (*Entry, bool)

	// useful for rebalancing
	Keys() []string
}
