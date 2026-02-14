package omap

import (
	"encoding/json"
	"slices"
	"testing"
)

func TestMap_MarshalJSON(t *testing.T) {
	t.Run("Regular", func(t *testing.T) {
		type T string
		m := New[T, any]()
		m.Set("foo", "bar")
		m.Set("baz", 123)

		b, err := json.Marshal(m)
		if err != nil {
			t.Fatalf("MarshalJSON failed: %v", err)
		}

		expected := `{"foo":"bar","baz":123}`
		if string(b) != expected {
			t.Errorf("MarshalJSON = %s, want %s", string(b), expected)
		}
	})

	t.Run("Number keys", func(t *testing.T) {
		type T int
		m := New[T, string]()
		m.Set(1, "one")
		m.Set(2, "two")

		b, err := json.Marshal(m)
		if err != nil {
			t.Fatalf("MarshalJSON failed: %v", err)
		}

		expected := `{"1":"one","2":"two"}`
		if string(b) != expected {
			t.Errorf("MarshalJSON = %s, want %s", string(b), expected)
		}
	})
}

func TestMap_UnmarshalJSON(t *testing.T) {
	t.Run("Regular", func(t *testing.T) {
		type T string
		jsonStr := `{"foo":"bar","baz":123}`
		m := New[T, any]()
		if err := json.Unmarshal([]byte(jsonStr), m); err != nil {
			t.Fatalf("UnmarshalJSON failed: %v", err)
		}

		if m.Len() != 2 {
			t.Errorf("Len() = %d, want 2", m.Len())
		}

		val, ok := m.TryGet("foo")
		if !ok || val != "bar" {
			t.Errorf("TryGet(foo) = %v, want bar", val)
		}

		ok = m.Has("baz")
		if !ok {
			t.Errorf("baz not found")
		}

		keys := m.Keys()
		wantKeys := []T{"foo", "baz"}
		if !slices.Equal(keys, wantKeys) {
			t.Errorf("Keys() = %v, want %v", keys, wantKeys)
		}
	})

	t.Run("Number keys", func(t *testing.T) {
		type T int
		jsonStr := `{"1":"bar","2":123}`
		m := New[T, any]()
		if err := json.Unmarshal([]byte(jsonStr), m); err != nil {
			t.Fatalf("UnmarshalJSON failed: %v", err)
		}

		if m.Len() != 2 {
			t.Errorf("Len() = %d, want 2", m.Len())
		}

		val, ok := m.TryGet(1)
		if !ok || val != "bar" {
			t.Errorf("TryGetGet(foo) = %v, want bar", val)
		}

		keys := m.Keys()
		wantKeys := []T{1, 2}
		if !slices.Equal(keys, wantKeys) {
			t.Errorf("Keys() = %v, want %v", keys, wantKeys)
		}
	})
}
