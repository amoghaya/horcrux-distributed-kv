// LRUPolicy implements EvictionPolicy using LRU strategy
// Most recently used keys are moved to front
// Least recently used keys are evicted from tail
package kv

import "sync"

type LRUNode struct {
	key  string
	prev *LRUNode
	next *LRUNode
}

type LRUPolicy struct {
	capacity int

	items map[string]*LRUNode

	head *LRUNode
	tail *LRUNode

	mu sync.Mutex
}

// create LRU policy
func NewLRUPolicy(cap int) *LRUPolicy {
	return &LRUPolicy{
		capacity: cap,
		items:    make(map[string]*LRUNode),
	}
}

// remove node from linked list
func (l *LRUPolicy) remove(node *LRUNode) {
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		l.head = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	} else {
		l.tail = node.prev
	}
}

// move node to front
func (l *LRUPolicy) addFront(node *LRUNode) {
	node.prev = nil
	node.next = l.head

	if l.head != nil {
		l.head.prev = node
	}

	l.head = node

	if l.tail == nil {
		l.tail = node
	}
}

// move existing node to front
func (l *LRUPolicy) moveToFront(node *LRUNode) {
	l.remove(node)
	l.addFront(node)
}

// called on read access
func (l *LRUPolicy) Access(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	node, ok := l.items[key]
	if !ok {
		return
	}

	l.moveToFront(node)
}

// called on insert
func (l *LRUPolicy) Insert(key string) string {
	l.mu.Lock()
	defer l.mu.Unlock()

	// already exists
	if node, ok := l.items[key]; ok {
		l.moveToFront(node)
		return ""
	}

	node := &LRUNode{
		key: key,
	}

	l.items[key] = node
	l.addFront(node)

	// eviction
	if len(l.items) > l.capacity {
		evicted := l.tail

		l.remove(evicted)
		delete(l.items, evicted.key)

		return evicted.key
	}

	return ""
}

// remove key completely
func (l *LRUPolicy) Delete(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	node, ok := l.items[key]
	if !ok {
		return
	}

	l.remove(node)
	delete(l.items, key)
}
