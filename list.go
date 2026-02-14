package omap

type elem[K comparable, V any] struct {
	next, prev *elem[K, V]
	key        K
	val        V
}

// list is a doubly linked list. It stores elems in insertion order.
type list[K comparable, V any] struct {
	root elem[K, V]
}

func (l *list[K, V]) init() {
	l.root.next = &l.root
	l.root.prev = &l.root
}

func (l *list[K, V]) append(k K, v V) *elem[K, V] {
	e := &elem[K, V]{key: k, val: v}
	e.prev = l.root.prev
	e.next = &l.root
	l.root.prev.next = e
	l.root.prev = e
	return e
}

func (l *list[K, V]) delete(e *elem[K, V]) {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.prev = nil
	e.next = nil
}
