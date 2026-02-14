package omap

import (
	"cmp"
	"slices"
	"testing"
)

func TestSort(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{
			name:     "Random Order",
			input:    []int{3, 1, 4, 1, 5, 9, 2, 6},
			expected: []int{1, 2, 3, 4, 5, 6, 9}, // Map keys are unique, duplicate '1' overwritten
		},
		{
			name:     "Already Sorted",
			input:    []int{1, 2, 3},
			expected: []int{1, 2, 3},
		},
		{
			name:     "Reverse Sorted",
			input:    []int{3, 2, 1},
			expected: []int{1, 2, 3},
		},
		{
			name:     "Empty",
			input:    []int{},
			expected: []int{},
		},
		{
			name:     "Single Element",
			input:    []int{1},
			expected: []int{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New[int, int]()
			for _, v := range tt.input {
				m.Set(v, v)
			}
			Sort(m)
			if !slices.Equal(m.Keys(), tt.expected) {
				t.Errorf("Sort() keys = %v, want %v", m.Keys(), tt.expected)
			}
		})
	}

	// Test nil map safety
	var nilMap *Map[int, int]
	Sort(nilMap) // Should not panic
}

func TestSortDesc(t *testing.T) {
	m := New[string, int]()
	keys := []string{"foo", "bar", "baz", "qux"}
	for _, k := range keys {
		m.Set(k, 1)
	}

	SortDesc(m)

	expected := []string{"qux", "foo", "baz", "bar"} // Descending order
	if !slices.Equal(m.Keys(), expected) {
		t.Errorf("SortDesc() keys = %v, want %v", m.Keys(), expected)
	}
}

func TestSortFunc(t *testing.T) {
	m := New[string, int]()
	// Insert by length: one(3), three(5), four(4), six(3)
	// "one" is inserted first, "six" is inserted last.
	input := []string{"one", "three", "four", "six"}
	for _, k := range input {
		m.Set(k, len(k))
	}

	// Sort by key length ascending
	SortFunc(m, func(a, b string) int {
		return cmp.Compare(len(a), len(b))
	})

	keys := m.Keys()

	// Check lengths are sorted: 3, 3, 4, 5
	expectedLengths := []int{3, 3, 4, 5}
	for i, k := range keys {
		if len(k) != expectedLengths[i] {
			t.Errorf("Index %d: key %q has length %d, want %d", i, k, len(k), expectedLengths[i])
		}
	}

	// Check stability: "one" should come before "six"
	if keys[0] != "one" || keys[1] != "six" {
		t.Errorf("Stability check failed: expected [one, six, ...], got %v", keys)
	}
}
