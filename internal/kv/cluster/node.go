package cluster

import (
	"horcrux/internal/kv"
	"sync"
	"time"
)

type Node struct {
	ID    string
	Store kv.Storage

	mu       sync.RWMutex
	lastSeen time.Time
}

func (n *Node) MarkAlive() {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.lastSeen = time.Now()
}

func (n *Node) IsAlive(timeout time.Duration) bool {
	n.mu.RLock()
	defer n.mu.RUnlock()

	if n.lastSeen.IsZero() {
		return false
	}

	return time.Since(n.lastSeen) < timeout
}
