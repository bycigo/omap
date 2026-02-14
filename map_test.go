package omap

import (
	"slices"
	"testing"
)

func TestMap_New(t *testing.T) {
	m := New[string, int]()
	if m == nil {
		t.Fatal("New() returned nil")
	}
	if m.Len() != 0 {
		t.Errorf("expected empty map, got len %d", m.Len())
	}
}

func TestMap_Set(t *testing.T) {
	m := New[string, int]()
	m.Set("one", 1)
	m.Set("two", 2)

	if m.Len() != 2 {
		t.Errorf("Len() = %d, want 2", m.Len())
	}

	if val := m.Get("one"); val != 1 {
		t.Errorf("Get(one) = %v, want 1", val)
	}

	if val := m.Get("two"); val != 2 {
		t.Errorf("Get(two) = %v, want 2", val)
	}
}

func TestMap_TrySet(t *testing.T) {
	m := New[string, int]()
	m.Set("foo", 1)

	if m.TrySet("foo", 2) {
		t.Error("TrySet(foo) should fail on existing key")
	}

	if !m.TrySet("bar", 3) {
		t.Error("TrySet(bar) should succeed on new key")
	}
	if val := m.Get("bar"); val != 3 {
		t.Errorf("Get(bar) after TrySet = %v, want 3", val)
	}
}

func TestMap_Get(t *testing.T) {
	m := New[string, int]()
	m.Set("one", 1)
	m.Set("two", 2)

	val := m.Get("one")
	if val != 1 {
		t.Errorf("Get(one) = %v, want 1", val)
	}

	val = m.Get("two")
	if val != 2 {
		t.Errorf("Get(two) = %v, want 2", val)
	}

	val = m.Get("three")
	if val != 0 {
		t.Errorf("Get(three) = %v, want 0", val)
	}
}

func TestMap_TryGet(t *testing.T) {
	m := New[string, int]()
	m.Set("one", 1)
	m.Set("two", 2)

	val, ok := m.TryGet("one")
	if val != 1 {
		t.Errorf("TryGet(one) = (%v, %v), want (1, true)", val, ok)
	}

	val, ok = m.TryGet("two")
	if val != 2 {
		t.Errorf("TryGet(two) = (%v, %v), want (2, true)", val, ok)
	}

	val, ok = m.TryGet("three")
	if ok {
		t.Errorf("TryGet(three) = (%v, %v), want (0, false)", val, ok)
	}
}

func TestMap_Delete(t *testing.T) {
	m := New[string, int]()
	m.Set("one", 1)
	m.Set("two", 2)
	m.Set("three", 3)

	m.Delete("two")

	if m.Len() != 2 {
		t.Errorf("Len() = %d, want 2", m.Len())
	}

	ok := m.Has("two")
	if ok {
		t.Error("Get(two) found deleted key")
	}

	keys := m.Keys()
	wantKeys := []string{"one", "three"}
	if !slices.Equal(keys, wantKeys) {
		t.Errorf("Keys() = %v, want %v", keys, wantKeys)
	}

	// Delete non-existent
	m.Delete("four")
	if m.Len() != 2 {
		t.Errorf("Len() after delete non-existent = %d, want 2", m.Len())
	}

	// Delete multiple
	m.Delete("one", "three")
	if m.Len() != 0 {
		t.Errorf("Len() after delete multiple = %d, want 0", m.Len())
	}
}

func TestMap_Has(t *testing.T) {
	m := New[string, int]()
	m.Set("one", 1)

	if !m.Has("one") {
		t.Error("Has(one) = false, want true")
	}
	if m.Has("two") {
		t.Error("Has(two) = true, want false")
	}
}

func TestMap_Clear(t *testing.T) {
	m := New[string, int]()
	m.Set("one", 1)
	m.Clear()
	if m.Len() != 0 {
		t.Errorf("Len() after clear = %d, want 0", m.Len())
	}
	ok := m.Has("one")
	if ok {
		t.Error("Has(one) found key after clear")
	}
}

func TestMap_All(t *testing.T) {
	m := New[string, int]()
	data := []struct {
		k string
		v int
	}{
		{"one", 1},
		{"two", 2},
		{"three", 3},
	}

	for _, d := range data {
		m.Set(d.k, d.v)
	}

	var i int
	for k, v := range m.All() {
		if k != data[i].k || v != data[i].v {
			t.Errorf("Range() index %d = (%v, %v), want (%v, %v)", i, k, v, data[i].k, data[i].v)
		}
		i++
	}
	if i != len(data) {
		t.Errorf("Range() count = %d, want %d", i, len(data))
	}
}

func TestMap_Len(t *testing.T) {
	m := New[string, int]()
	if m.Len() != 0 {
		t.Errorf("Len() initially = %d, want 0", m.Len())
	}
	m.Set("one", 1)
	if m.Len() != 1 {
		t.Errorf("Len() after set = %d, want 1", m.Len())
	}
}

func TestMap_Keys_Values(t *testing.T) {
	m := New[string, int]()
	m.Set("one", 1)
	m.Set("two", 2)
	m.Set("three", 3)

	keys := m.Keys()
	wantKeys := []string{"one", "two", "three"}
	if !slices.Equal(keys, wantKeys) {
		t.Errorf("Keys() = %v, want %v", keys, wantKeys)
	}

	vals := m.Values()
	wantVals := []int{1, 2, 3}
	if !slices.Equal(vals, wantVals) {
		t.Errorf("Values() = %v, want %v", vals, wantVals)
	}
}

func TestMap_Merge(t *testing.T) {
	m1 := New[string, int]()
	m1.Set("one", 1)

	m2 := New[string, int]()
	m2.Set("two", 2)
	m2.Set("three", 3)

	m1.Merge(m2)

	if m1.Len() != 3 {
		t.Errorf("Len() after merge = %d, want 3", m1.Len())
	}

	// Verify order
	wantKeys := []string{"one", "two", "three"}
	if !slices.Equal(m1.Keys(), wantKeys) {
		t.Errorf("Keys() after merge = %v, want %v", m1.Keys(), wantKeys)
	}
}

func TestMap_Reverse(t *testing.T) {
	// Test empty map
	m := New[string, int]()
	m.Reverse()
	if m.Len() != 0 {
		t.Errorf("Len() after reverse empty = %d, want 0", m.Len())
	}

	// Test single element
	m.Set("one", 1)
	m.Reverse()
	if !slices.Equal(m.Keys(), []string{"one"}) {
		t.Errorf("Keys() after reverse single = %v, want [one]", m.Keys())
	}

	// Test multiple elements
	m.Set("two", 2)
	m.Set("three", 3)
	m.Reverse()
	wantKeys := []string{"three", "two", "one"}
	if !slices.Equal(m.Keys(), wantKeys) {
		t.Errorf("Keys() after reverse = %v, want %v", m.Keys(), wantKeys)
	}
	wantVals := []int{3, 2, 1}
	if !slices.Equal(m.Values(), wantVals) {
		t.Errorf("Values() after reverse = %v, want %v", m.Values(), wantVals)
	}
}
