package cluster

import (
	"hash/fnv"
	"sort"
	"strconv"
)

type HashRing struct {
	nodes      map[uint32]string
	sortedKeys []uint32
	replicas   int
}

// NewHashRing creates consistent hashing ring
func NewHashRing(replicas int) *HashRing {
	return &HashRing{
		nodes:    make(map[uint32]string),
		replicas: replicas,
	}
}

// hash fn
func hashKey(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}

// add node to ring
func (r *HashRing) AddNode(node string) {
	for i := 0; i < r.replicas; i++ {
		virtualKey := node + ":" + strconv.Itoa(i)
		hash := hashKey(virtualKey)

		r.nodes[hash] = node
		r.sortedKeys = append(r.sortedKeys, hash)
	}

	sort.Slice(r.sortedKeys, func(i, j int) bool {
		return r.sortedKeys[i] < r.sortedKeys[j]
	})
}

// get node for key
func (r *HashRing) GetNode(key string) string {
	if len(r.nodes) == 0 {
		return ""
	}

	hash := hashKey(key)

	for _, k := range r.sortedKeys {
		if hash <= k {
			return r.nodes[k]
		}
	}

	return r.nodes[r.sortedKeys[0]]
}

// get n nodes for replication
func (r *HashRing) GetNodes(key string, count int) []string {
	if len(r.nodes) == 0 {
		return nil
	}

	hash := hashKey(key)

	// find start position
	start := 0
	for i, k := range r.sortedKeys {
		if hash <= k {
			start = i
			break
		}
	}

	result := make([]string, 0, count)
	seen := make(map[string]bool)

	// walk ring and collect unique nodes
	for i := 0; len(result) < count && i < len(r.sortedKeys)*2; i++ {
		idx := (start + i) % len(r.sortedKeys)
		node := r.nodes[r.sortedKeys[idx]]

		if !seen[node] {
			seen[node] = true
			result = append(result, node)
		}
	}

	return result
}
