package kv

// KVStore is the public API layer of Horcrux.
// It delegates all operations to the storage engine.
type KVStore struct {
	engine StorageEngine
}

// NewKVStore initializes the KV system
func NewKVStore(engine StorageEngine) *KVStore {
	return &KVStore{
		engine: engine,
	}
}

// Put stores a value
func (k *KVStore) Put(key string, value interface{}) {
	k.engine.Put(key, value)
}

// Get retrieves a value
func (k *KVStore) Get(key string) (interface{}, bool) {
	return k.engine.Get(key)
}

// Delete removes a key
func (k *KVStore) Delete(key string) bool {
	return k.engine.Delete(key)
}

func (k *KVStore) PutWithTTL(key string, value interface{}, ttl int) {
	k.engine.PutWithTTL(key, value, ttl)
}

func (k *KVStore) GetEntry(key string) (*Entry, bool) {
	engine, ok := k.engine.(*InMemoryEngine)
	if !ok {
		return nil, false
	}

	return engine.GetEntry(key)
}
