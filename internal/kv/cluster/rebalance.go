package cluster

type Rebalancer struct {
	coordinator *Coordinator
}

// constructor
func NewRebalancer(c *Coordinator) *Rebalancer {
	return &Rebalancer{
		coordinator: c,
	}
}

// Rebalance redistributes keys across nodes when cluster topology changes
func (r *Rebalancer) Rebalance() {
	for _, node := range r.coordinator.nodes {

		keys := node.Store.Keys()

		for _, key := range keys {

			owners := r.coordinator.ring.GetNodes(
				key,
				r.coordinator.replicationFactor,
			)

			shouldExist := false

			for _, owner := range owners {
				if owner == node.ID {
					shouldExist = true
					break
				}
			}

			// node no longer owns key
			if !shouldExist {

				val, ok := node.Store.Get(key)
				if !ok {
					continue
				}

				// migrate to correct owners
				for _, ownerID := range owners {

					target, exists := r.coordinator.nodes[ownerID]
					if !exists {
						continue
					}

					target.Store.Put(key, val)
				}

				// optional cleanup
				node.Store.Delete(key)
			}
		}
	}
}

// For every key:

// recompute correct owners
// check if current node still owns it
// if not:
// copy to correct nodes
// remove locally

// COMPLEXITY ANALYSIS
// Current approach:
// O(total_keys × replication_factor)

// Expensive.

// WHY REAL SYSTEMS ARE BETTER

// Cassandra/Dynamo:

// partition metadata
// vnode ownership
// SSTable streaming
// token ranges

// avoid full scans.
