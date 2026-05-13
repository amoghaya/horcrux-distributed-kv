package cluster

import (
	"fmt"
	"sync"
	"time"
)

type FailureDetector struct {
	timeout time.Duration
	nodes   map[string]*Node
	mu      sync.RWMutex
}

// NewFailureDetector creates detector
func NewFailureDetector(timeout time.Duration) *FailureDetector {
	return &FailureDetector{
		timeout: timeout,
		nodes:   make(map[string]*Node),
	}
}

// AddNode safely registers node with detector
func (f *FailureDetector) AddNode(node *Node) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.nodes[node.ID] = node
}

// Start runs background health monitoring loop
func (f *FailureDetector) Start() {
	go func() {
		for {
			f.mu.RLock()

			for id, node := range f.nodes {
				if !node.IsAlive(f.timeout) {
					fmt.Printf(
						"[FAILURE DETECTOR] Node %s marked as dead\n",
						id,
					)
				}
			}

			f.mu.RUnlock()

			time.Sleep(1 * time.Second)
		}
	}()
}
