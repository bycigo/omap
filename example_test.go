package omap

import (
	"encoding/json"
	"fmt"
)

func ExampleNew() {
	m := New[string, int]()
	m.Set("foo", 1)
	fmt.Println(m.Len())
	// Output: 1
}

func ExampleMap_Set() {
	m := New[string, int]()
	m.Set("first", 10)
	m.Set("second", 20)
	// Update existing
	m.Set("first", 30)

	fmt.Println(m.Keys())
	// Output: [first second]
}

func ExampleMap_TrySet() {
	m := New[string, int]()
	m.Set("foo", 1)

	fooAdded := m.TrySet("foo", 2)
	fmt.Println(fooAdded)
	fmt.Println(m.Get("foo"))

	barAdded := m.TrySet("bar", 3)
	fmt.Println(barAdded)
	fmt.Println(m.Get("bar"))
	// Output:
	// false
	// 1
	// true
	// 3
}

func ExampleMap_Get() {
	m := New[string, string]()
	m.Set("greet", "hello")

	val := m.Get("greet")
	fmt.Println(val)
	// Output:
	// hello
}

func ExampleMap_TryGet() {
	m := New[string, string]()
	m.Set("key", "value")

	if val, ok := m.TryGet("key"); ok {
		fmt.Println(val)
	}
	if _, ok := m.TryGet("missing"); !ok {
		fmt.Println("missing key")
	}
	// Output:
	// value
	// missing key
}

func ExampleMap_Delete() {
	m := New[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	m.Delete("b")
	fmt.Println(m.Keys())
	// Output: [a c]
}

func ExampleMap_Has() {
	m := New[string, int]()
	m.Set("x", 42)

	fmt.Println(m.Has("x"))
	fmt.Println(m.Has("y"))
	// Output:
	// true
	// false
}

func ExampleMap_Clear() {
	m := New[string, int]()
	m.Set("a", 1)
	m.Clear()
	fmt.Println(m.Len())
	// Output: 0
}

func ExampleMap_All() {
	m := New[string, int]()
	m.Set("one", 1)
	m.Set("two", 2)
	m.Set("three", 3)

	for k, v := range m.All() {
		fmt.Printf("%s: %d\n", k, v)
	}
	// Output:
	// one: 1
	// two: 2
	// three: 3
}

func ExampleMap_Len() {
	m := New[string, int]()
	fmt.Println(m.Len())
	m.Set("a", 1)
	fmt.Println(m.Len())
	// Output:
	// 0
	// 1
}

func ExampleMap_Keys() {
	m := New[string, int]()
	m.Set("foo", 1)
	m.Set("bar", 2)
	fmt.Println(m.Keys())
	// Output: [foo bar]
}

func ExampleMap_Values() {
	m := New[string, int]()
	m.Set("foo", 1)
	m.Set("bar", 2)
	fmt.Println(m.Values())
	// Output: [1 2]
}

func ExampleMap_Merge() {
	m1 := New[string, int]()
	m1.Set("a", 1)

	m2 := New[string, int]()
	m2.Set("b", 2)
	m2.Set("c", 3)

	m1.Merge(m2)
	fmt.Println(m1.Keys())
	// Output: [a b c]
}

func ExampleMap_MarshalJSON() {
	m := New[string, int]()
	m.Set("foo", 1)
	m.Set("bar", 2)

	data, _ := json.Marshal(m)
	fmt.Println(string(data))
	// Output: {"foo":1,"bar":2}
}

func ExampleMap_Reverse() {
	m := New[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	m.Reverse()
	fmt.Println(m.Keys())
	fmt.Println(m.Values())
	// Output:
	// [c b a]
	// [3 2 1]
}

func ExampleMap_UnmarshalJSON() {
	jsonStr := `{"foo":1,"bar":2,"baz":3}`
	m := New[string, int]()

	_ = json.Unmarshal([]byte(jsonStr), m)
	fmt.Println(m.Keys())
	// Output: [foo bar baz]
}

func ExampleSort() {
	m := New[string, int]()
	m.Set("c", 3)
	m.Set("a", 1)
	m.Set("b", 2)

	Sort(m)
	fmt.Println(m.Keys())
	// Output: [a b c]
}

func ExampleSortDesc() {
	m := New[string, int]()
	m.Set("a", 1)
	m.Set("c", 3)
	m.Set("b", 2)

	SortDesc(m)
	fmt.Println(m.Keys())
	// Output: [c b a]
}

func ExampleSortFunc() {
	m := New[string, int]()
	m.Set("apple", 1)
	m.Set("banana", 2)
	m.Set("cherry", 3)

	// Sort by length of keys descending
	SortFunc(m, func(k1, k2 string) int {
		return len(k2) - len(k1)
	})
	fmt.Println(m.Keys())
	// Output: [banana cherry apple]
}
