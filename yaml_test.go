package omap

import (
	"slices"
	"testing"

	"go.yaml.in/yaml/v3"
)

func TestMap_MarshalYAML(t *testing.T) {
	t.Run("StringKeys", func(t *testing.T) {
		m := New[string, any]()
		m.Set("foo", "bar")
		m.Set("baz", 123)

		b, err := yaml.Marshal(m)
		if err != nil {
			t.Fatalf("MarshalYAML failed: %v", err)
		}

		expected := "foo: bar\nbaz: 123\n"
		if string(b) != expected {
			t.Errorf("MarshalYAML = %q, want %q", string(b), expected)
		}
	})

	t.Run("IntKeys", func(t *testing.T) {
		m := New[int, string]()
		m.Set(1, "one")
		m.Set(2, "two")

		b, err := yaml.Marshal(m)
		if err != nil {
			t.Fatalf("MarshalYAML failed: %v", err)
		}

		expected := "1: one\n2: two\n"
		if string(b) != expected {
			t.Errorf("MarshalYAML = %q, want %q", string(b), expected)
		}
	})

	t.Run("Empty", func(t *testing.T) {
		m := New[string, string]()

		b, err := yaml.Marshal(m)
		if err != nil {
			t.Fatalf("MarshalYAML failed: %v", err)
		}

		expected := "{}\n"
		if string(b) != expected {
			t.Errorf("MarshalYAML = %q, want %q", string(b), expected)
		}
	})

	t.Run("NestedMap", func(t *testing.T) {
		inner := New[string, int]()
		inner.Set("x", 1)
		inner.Set("z", 2)

		outer := New[string, *Map[string, int]]()
		outer.Set("nested", inner)

		b, err := yaml.Marshal(outer)
		if err != nil {
			t.Fatalf("MarshalYAML failed: %v", err)
		}

		expected := "nested:\n    x: 1\n    z: 2\n"
		if string(b) != expected {
			t.Errorf("MarshalYAML = %q, want %q", string(b), expected)
		}
	})

	t.Run("PreservesOrder", func(t *testing.T) {
		m := New[string, int]()
		m.Set("c", 3)
		m.Set("a", 1)
		m.Set("b", 2)

		b, err := yaml.Marshal(m)
		if err != nil {
			t.Fatalf("MarshalYAML failed: %v", err)
		}

		expected := "c: 3\na: 1\nb: 2\n"
		if string(b) != expected {
			t.Errorf("MarshalYAML = %q, want %q", string(b), expected)
		}
	})

	t.Run("NilMap", func(t *testing.T) {
		var m Map[string, int]

		b, err := yaml.Marshal(&m)
		if err != nil {
			t.Fatalf("MarshalYAML failed: %v", err)
		}

		expected := "{}\n"
		if string(b) != expected {
			t.Errorf("MarshalYAML = %q, want %q", string(b), expected)
		}
	})
}

func TestMap_UnmarshalYAML(t *testing.T) {
	t.Run("StringKeys", func(t *testing.T) {
		yamlStr := "foo: bar\nbaz: 123\n"
		m := New[string, any]()
		if err := yaml.Unmarshal([]byte(yamlStr), m); err != nil {
			t.Fatalf("UnmarshalYAML failed: %v", err)
		}

		if m.Len() != 2 {
			t.Errorf("Len() = %d, want 2", m.Len())
		}

		val, ok := m.TryGet("foo")
		if !ok || val != "bar" {
			t.Errorf("TryGet(foo) = %v, want bar", val)
		}

		keys := m.Keys()
		wantKeys := []string{"foo", "baz"}
		if !slices.Equal(keys, wantKeys) {
			t.Errorf("Keys() = %v, want %v", keys, wantKeys)
		}
	})

	t.Run("IntKeys", func(t *testing.T) {
		yamlStr := "1: one\n2: two\n"
		m := New[int, string]()
		if err := yaml.Unmarshal([]byte(yamlStr), m); err != nil {
			t.Fatalf("UnmarshalYAML failed: %v", err)
		}

		if m.Len() != 2 {
			t.Errorf("Len() = %d, want 2", m.Len())
		}

		if val := m.Get(1); val != "one" {
			t.Errorf("Get(1) = %v, want one", val)
		}
		if val := m.Get(2); val != "two" {
			t.Errorf("Get(2) = %v, want two", val)
		}

		keys := m.Keys()
		wantKeys := []int{1, 2}
		if !slices.Equal(keys, wantKeys) {
			t.Errorf("Keys() = %v, want %v", keys, wantKeys)
		}
	})

	t.Run("PreservesOrder", func(t *testing.T) {
		yamlStr := "c: 3\na: 1\nb: 2\n"
		m := New[string, int]()
		if err := yaml.Unmarshal([]byte(yamlStr), m); err != nil {
			t.Fatalf("UnmarshalYAML failed: %v", err)
		}

		keys := m.Keys()
		wantKeys := []string{"c", "a", "b"}
		if !slices.Equal(keys, wantKeys) {
			t.Errorf("Keys() = %v, want %v", keys, wantKeys)
		}

		values := m.Values()
		wantValues := []int{3, 1, 2}
		if !slices.Equal(values, wantValues) {
			t.Errorf("Values() = %v, want %v", values, wantValues)
		}
	})

	t.Run("Empty", func(t *testing.T) {
		yamlStr := "{}\n"
		m := New[string, int]()
		if err := yaml.Unmarshal([]byte(yamlStr), m); err != nil {
			t.Fatalf("UnmarshalYAML failed: %v", err)
		}

		if m.Len() != 0 {
			t.Errorf("Len() = %d, want 0", m.Len())
		}
	})

	t.Run("SliceValues", func(t *testing.T) {
		yamlStr := "a:\n  - 1\n  - 2\nb:\n  - 3\n  - 4\n"
		m := New[string, []int]()
		if err := yaml.Unmarshal([]byte(yamlStr), m); err != nil {
			t.Fatalf("UnmarshalYAML failed: %v", err)
		}

		if m.Len() != 2 {
			t.Errorf("Len() = %d, want 2", m.Len())
		}

		a := m.Get("a")
		if !slices.Equal(a, []int{1, 2}) {
			t.Errorf("Get(a) = %v, want [1 2]", a)
		}

		b := m.Get("b")
		if !slices.Equal(b, []int{3, 4}) {
			t.Errorf("Get(b) = %v, want [3 4]", b)
		}
	})

	t.Run("NestedMap", func(t *testing.T) {
		yamlStr := "outer:\n    x: 1\n    y: 2\n"
		m := New[string, *Map[string, int]]()
		if err := yaml.Unmarshal([]byte(yamlStr), m); err != nil {
			t.Fatalf("UnmarshalYAML failed: %v", err)
		}

		inner, ok := m.TryGet("outer")
		if !ok {
			t.Fatal("outer not found")
		}
		if inner.Get("x") != 1 {
			t.Errorf("inner.Get(x) = %v, want 1", inner.Get("x"))
		}
		if inner.Get("y") != 2 {
			t.Errorf("inner.Get(y) = %v, want 2", inner.Get("y"))
		}
	})

	t.Run("InvalidInput_NotMapping", func(t *testing.T) {
		yamlStr := "- a\n- b\n"
		m := New[string, string]()
		err := yaml.Unmarshal([]byte(yamlStr), m)
		if err == nil {
			t.Error("expected error for non-mapping input")
		}
	})

	t.Run("ZeroValueMap", func(t *testing.T) {
		yamlStr := "foo: bar\nbaz: qux\n"
		var m Map[string, string]
		if err := yaml.Unmarshal([]byte(yamlStr), &m); err != nil {
			t.Fatalf("UnmarshalYAML failed: %v", err)
		}

		if m.Len() != 2 {
			t.Errorf("Len() = %d, want 2", m.Len())
		}

		if val := m.Get("foo"); val != "bar" {
			t.Errorf("Get(foo) = %v, want bar", val)
		}
	})
}

func TestMap_YAMLRoundTrip(t *testing.T) {
	t.Run("StringString", func(t *testing.T) {
		m := New[string, string]()
		m.Set("hello", "world")
		m.Set("foo", "bar")

		b, err := yaml.Marshal(m)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		m2 := New[string, string]()
		if err := yaml.Unmarshal(b, m2); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if !slices.Equal(m.Keys(), m2.Keys()) {
			t.Errorf("keys mismatch: %v != %v", m.Keys(), m2.Keys())
		}
		if !slices.Equal(m.Values(), m2.Values()) {
			t.Errorf("values mismatch: %v != %v", m.Values(), m2.Values())
		}
	})

	t.Run("IntString", func(t *testing.T) {
		m := New[int, string]()
		m.Set(10, "ten")
		m.Set(20, "twenty")
		m.Set(30, "thirty")

		b, err := yaml.Marshal(m)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		m2 := New[int, string]()
		if err := yaml.Unmarshal(b, m2); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if !slices.Equal(m.Keys(), m2.Keys()) {
			t.Errorf("keys mismatch: %v != %v", m.Keys(), m2.Keys())
		}
		if !slices.Equal(m.Values(), m2.Values()) {
			t.Errorf("values mismatch: %v != %v", m.Values(), m2.Values())
		}
	})

	t.Run("SliceValues", func(t *testing.T) {
		m := New[string, []int]()
		m.Set("first", []int{1, 2, 3})
		m.Set("second", []int{4, 5, 6})

		b, err := yaml.Marshal(m)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		m2 := New[string, []int]()
		if err := yaml.Unmarshal(b, m2); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if !slices.Equal(m.Keys(), m2.Keys()) {
			t.Errorf("keys mismatch: %v != %v", m.Keys(), m2.Keys())
		}
		for _, k := range m.Keys() {
			if !slices.Equal(m.Get(k), m2.Get(k)) {
				t.Errorf("values for key %q mismatch: %v != %v", k, m.Get(k), m2.Get(k))
			}
		}
	})
}
