package omap

import (
	"testing"
)

func TestList_append(t *testing.T) {
	l := list[string, int]{}
	l.init()

	// Append first element
	e1 := l.append("a", 1)
	if e1.key != "a" {
		t.Errorf("Expected key 'a', got '%v'", e1.key)
	}
	if e1.val != 1 {
		t.Errorf("Expected val 1, got %v", e1.val)
	}
	if l.root.next != e1 {
		t.Errorf("Expected root.next to be e1")
	}
	if l.root.prev != e1 {
		t.Errorf("Expected root.prev to be e1")
	}
	if e1.next != &l.root {
		t.Errorf("Expected e1.next to be root")
	}
	if e1.prev != &l.root {
		t.Errorf("Expected e1.prev to be root")
	}

	// Append second element
	e2 := l.append("b", 2)
	if e2.key != "b" {
		t.Errorf("Expected key 'b', got '%v'", e2.key)
	}
	if e2.val != 2 {
		t.Errorf("Expected val 2, got %v", e2.val)
	}
	if l.root.next != e1 {
		t.Errorf("Expected root.next to be e1")
	}
	if l.root.prev != e2 {
		t.Errorf("Expected root.prev to be e2")
	}
	if e1.next != e2 {
		t.Errorf("Expected e1.next to be e2")
	}
	if e2.prev != e1 {
		t.Errorf("Expected e2.prev to be e1")
	}
	if e2.next != &l.root {
		t.Errorf("Expected e2.next to be root")
	}
}

func TestList_delete(t *testing.T) {
	l := list[string, int]{}
	l.init()

	e1 := l.append("a", 1)
	e2 := l.append("b", 2)
	e3 := l.append("c", 3)

	// Delete middle element
	l.delete(e2)
	if e1.next != e3 {
		t.Errorf("After deleting e2, e1.next should be e3")
	}
	if e3.prev != e1 {
		t.Errorf("After deleting e2, e3.prev should be e1")
	}
	if e2.next != nil || e2.prev != nil {
		t.Errorf("Deleted element's pointers should be nil")
	}

	// Delete tail element
	l.delete(e3)
	if e1.next != &l.root {
		t.Errorf("After deleting e3, e1.next should be root")
	}
	if l.root.prev != e1 {
		t.Errorf("After deleting e3, root.prev should be e1")
	}

	// Delete head element
	l.delete(e1)
	if l.root.next != &l.root {
		t.Errorf("After deleting e1, root.next should be root")
	}
	if l.root.prev != &l.root {
		t.Errorf("After deleting e1, root.prev should be root")
	}
}
