package omap

import (
	"iter"
)

// Map represents an ordered map that maintains elements in the order of their insertion.
type Map[K comparable, V any] struct {
	kv map[K]*elem[K, V]
	kl list[K, V]
}

// New creates and returns a new Map instance.
func New[K comparable, V any]() *Map[K, V] {
	return Make[K, V](0)
}

// Make creates and returns a new Map instance with the specified capacity.
func Make[K comparable, V any](capacity int) *Map[K, V] {
	m := Map[K, V]{}
	m.init(capacity)
	return &m
}

func (m *Map[K, V]) lazyInit() {
	if m.kv == nil {
		m.init(0)
	}
}

func (m *Map[K, V]) init(capacity int) {
	m.kv = make(map[K]*elem[K, V], capacity)
	m.kl.init()
}

// Set adds a key-value pair to the map.
// If the key already exists, its value is updated.
// If the key does not exist, it is appended to the end of the insertion order list.
func (m *Map[K, V]) Set(key K, value V) {
	m.lazyInit()
	if e, exists := m.kv[key]; exists {
		e.val = value
	} else {
		e = m.kl.append(key, value)
		m.kv[key] = e
	}
}

// TrySet adds a key-value pair to the map only if the key does not already exist.
// It returns true if the key-value pair was added, and false if the key already exists.
func (m *Map[K, V]) TrySet(key K, value V) bool {
	m.lazyInit()
	if e, exists := m.kv[key]; !exists {
		e = m.kl.append(key, value)
		m.kv[key] = e
		return true
	}
	return false
}

// Get retrieves the value associated with the given key.
func (m *Map[K, V]) Get(key K) (value V) {
	e, ok := m.kv[key]
	if ok {
		value = e.val
	}
	return
}

// TryGet retrieves the value associated with the given key.
// It returns the value and true if the key exists, otherwise the zero value and false.
func (m *Map[K, V]) TryGet(key K) (value V, ok bool) {
	e, ok := m.kv[key]
	if ok {
		value = e.val
	}
	return
}

// Delete removes the key-value pair associated with the given key from the map.
// It is no-op if the key does not exist.
func (m *Map[K, V]) Delete(keys ...K) {
	for _, k := range keys {
		if e, ok := m.kv[k]; ok {
			m.kl.delete(e)
			delete(m.kv, k)
		}
	}
}

// Has checks if the given key exists in the map.
func (m *Map[K, V]) Has(key K) bool {
	_, ok := m.kv[key]
	return ok
}

// Clear removes all key-value pairs from the map.
func (m *Map[K, V]) Clear() {
	m.init(0)
}

// All returns an iterator over the map's entries in insertion order.
func (m *Map[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for e := m.kl.root.next; e != nil && e != &m.kl.root; {
			next := e.next
			if !yield(e.key, e.val) {
				return
			}
			e = next
		}
	}
}

// Len returns the number of key-value pairs in the map.
func (m *Map[K, V]) Len() int {
	return len(m.kv)
}

// Keys returns a slice of all keys in the map, in the order they were inserted.
func (m *Map[K, V]) Keys() []K {
	keys := make([]K, 0, len(m.kv))
	for e := m.kl.root.next; e != nil && e != &m.kl.root; e = e.next {
		keys = append(keys, e.key)
	}
	return keys
}

// Values returns a slice of all values in the map, in the order their keys were inserted.
func (m *Map[K, V]) Values() []V {
	values := make([]V, 0, len(m.kv))
	for e := m.kl.root.next; e != nil && e != &m.kl.root; e = e.next {
		values = append(values, e.val)
	}
	return values
}

// Merge merges the key-value pairs from the target maps into the current map.
func (m *Map[K, V]) Merge(target ...*Map[K, V]) {
	for _, item := range target {
		for k, v := range item.All() {
			m.Set(k, v)
		}
	}
}

// Reverse reverses the order of elements in the map.
func (m *Map[K, V]) Reverse() {
	if m.kv == nil || m.kl.root.next == &m.kl.root {
		return
	}

	curr := m.kl.root.next
	for curr != &m.kl.root {
		curr.next, curr.prev = curr.prev, curr.next
		curr = curr.prev
	}
	m.kl.root.next, m.kl.root.prev = m.kl.root.prev, m.kl.root.next
}
