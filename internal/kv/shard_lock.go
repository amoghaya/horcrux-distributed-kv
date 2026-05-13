package kv

import "sync"

// ShardLock provides lock striping to reduce contention
// Each shard has its own RWMutex
type ShardLock struct {
	locks []sync.RWMutex
	size  int
}

// NewShardLock initializes N locks
func NewShardLock(size int) *ShardLock {
	return &ShardLock{
		locks: make([]sync.RWMutex, size),
		size:  size,
	}
}

// getShard returns which lock to use
func (s *ShardLock) getShard(key string) *sync.RWMutex {
	hash := fnv32(key)
	return &s.locks[hash%uint32(s.size)]
}

// FNV hash (fast + stable)
func fnv32(key string) uint32 {
	var hash uint32 = 2166136261
	for i := 0; i < len(key); i++ {
		hash ^= uint32(key[i])
		hash *= 16777619
	}
	return hash
}
