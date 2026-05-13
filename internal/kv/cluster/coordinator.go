package cluster

import (
	"fmt"
	"horcrux/internal/kv"
	"sync"
	"time"
)

type Coordinator struct {
	ring              *HashRing
	nodes             map[string]*Node
	replicationFactor int
	writeQuorum       int
	readQuorum        int
	timeout           time.Duration
	mu                sync.RWMutex
}

// NewCoordinator creates a cluster coordinator
func NewCoordinator(replication, w, r int) *Coordinator {
	return &Coordinator{
		ring:              NewHashRing(3),
		nodes:             make(map[string]*Node),
		replicationFactor: replication,
		writeQuorum:       w,
		readQuorum:        r,
		timeout:           2 * time.Second,
	}
}

// AddNode registers a node into the cluster
func (c *Coordinator) AddNode(id string, store kv.Storage) {
	c.mu.Lock()
	defer c.mu.Unlock()

	node := &Node{
		ID:    id,
		Store: store,
	}

	node.MarkAlive() //initial mark alive

	c.nodes[id] = node
	c.ring.AddNode(id)

	rebalancer := NewRebalancer(c)
	rebalancer.Rebalance() //not good as whole cluster is frozen during rebalance, but simple for demo
}

// Put writes data with replication + quorum control
func (c *Coordinator) Put(key string, value interface{}) {
	nodes := c.ring.GetNodes(key, c.replicationFactor)

	success := 0

	for _, id := range nodes {
		node, ok := c.nodes[id]
		if !ok {
			continue
		}

		if node.IsAlive(c.timeout) {
			node.Store.Put(key, value)
			node.MarkAlive()
			success++
		}

		if success >= c.writeQuorum {
			break
		}
	}
}

// Get reads data using quorum-aware replica lookup
func (c *Coordinator) Get(key string) (interface{}, bool) {
	nodes := c.ring.GetNodes(key, c.replicationFactor)

	var latestEntry *kv.Entry
	var staleNodes []*Node

	success := 0

	for _, id := range nodes {
		node := c.nodes[id]

		if !node.IsAlive(c.timeout) {
			continue
		}

		entry, ok := node.Store.GetEntry(key)

		if !ok {
			continue
		}

		success++

		// newest version wins
		if latestEntry == nil || entry.Version > latestEntry.Version {

			// previous latest becomes stale
			if latestEntry != nil {
				for _, sid := range nodes {
					snode := c.nodes[sid]

					e, ok := snode.Store.GetEntry(key)

					if ok && e.Version < entry.Version {
						staleNodes = append(staleNodes, snode)
					}
				}
			}

			latestEntry = entry
		}

		if success >= c.readQuorum {
			break
		}
	}

	if latestEntry == nil {
		return nil, false
	}

	// background read repair
	go c.repairReplicas(key, latestEntry, staleNodes)

	return latestEntry.Value, true
}

// Delete removes key across replicas with quorum enforcement
func (c *Coordinator) Delete(key string) {
	nodes := c.ring.GetNodes(key, c.replicationFactor)

	success := 0

	for _, id := range nodes {
		node, ok := c.nodes[id]
		if !ok {
			continue
		}

		if node.IsAlive(c.timeout) {
			node.Store.Delete(key)
			node.MarkAlive()
			success++
		}

		if success >= c.writeQuorum {
			break
		}
	}
}

func (c *Coordinator) repairReplicas(
	key string,
	latest *kv.Entry,
	staleNodes []*Node,
) {
	for _, node := range staleNodes {
		node.Store.Put(key, latest.Value)

		fmt.Printf(
			"[READ REPAIR] repaired key=%s on node=%s\n",
			key,
			node.ID,
		)
	}
}
