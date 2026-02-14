package omap

import "cmp"

// Sort sorts the map in ascending order of keys.
func Sort[K cmp.Ordered, V any](m *Map[K, V]) {
	SortFunc(m, cmp.Compare)
}

// SortDesc sorts the map in descending order of keys.
func SortDesc[K cmp.Ordered, V any](m *Map[K, V]) {
	SortFunc(m, func(k1, k2 K) int {
		return cmp.Compare(k2, k1)
	})
}

// SortFunc sorts the map using a custom comparison function for keys.
func SortFunc[K cmp.Ordered, V any](m *Map[K, V], compare func(k1, k2 K) int) {
	if m == nil || m.Len() < 2 {
		return
	}

	// detach the list from root for sorting
	head := m.kl.root.next
	tail := m.kl.root.prev

	// disconnect the circular links
	head.prev = nil
	tail.next = nil

	// perform merge sort
	sorted := mergeSortList(head, compare)

	// reconnect the sorted list back to root
	m.kl.root.next = sorted
	sorted.prev = &m.kl.root

	// find the new tail and reconnect
	last := sorted
	for last.next != nil {
		last = last.next
	}
	last.next = &m.kl.root
	m.kl.root.prev = last
}

func mergeSortList[K cmp.Ordered, V any](head *elem[K, V], compare func(k1, k2 K) int) *elem[K, V] {
	if head == nil || head.next == nil {
		return head
	}

	// split the list into two halves
	mid, fast := head, head
	for fast.next != nil && fast.next.next != nil {
		mid = mid.next
		fast = fast.next.next
	}
	right := mid.next
	mid.next = nil
	if right != nil {
		right.prev = nil
	}

	// recursively sort both halves
	left := mergeSortList[K, V](head, compare)
	right = mergeSortList[K, V](right, compare)

	// merge the sorted halves
	return mergeList[K, V](left, right, compare)
}

func mergeList[K cmp.Ordered, V any](left, right *elem[K, V], compare func(k1, k2 K) int) *elem[K, V] {
	// create a dummy head for the result list
	var dummy elem[K, V]
	tail := &dummy

	for left != nil && right != nil {
		if compare(left.key, right.key) <= 0 {
			tail.next = left
			left.prev = tail
			left = left.next
		} else {
			tail.next = right
			right.prev = tail
			right = right.next
		}
		tail = tail.next
	}

	// append the remaining nodes
	if left != nil {
		tail.next = left
		left.prev = tail
	}
	if right != nil {
		tail.next = right
		right.prev = tail
	}

	// fix the head's prev pointer
	result := dummy.next
	if result != nil {
		result.prev = nil
	}
	return result
}
