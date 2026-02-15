# omap

`omap` provides an ordered map implementation for Go. Unlike Go's standard `map`, `omap` maintains the insertion order
of keys. It leverages Go iterators and supports generic key-value types.

## Features

- **Ordered**: Maintains insertion order of elements.
- **Generic**: Supports any `comparable` key and any value type.
- **Efficient**: O(1) for Set, Get, Delete.
- **Iteration**: Supports `range` iterator.
- **JSON**: Supports JSON marshaling/unmarshaling, preserving order.
- **Sorting**: Provides in-place sorting capabilities.

## Installation

```bash
go get github.com/bycigo/omap
```

## Usage

### Creating a omap

```go
// creating by omap.New()
m := omap.New[string, any]()

// with initial capacity
m := omap.Make[string, any](10)
```

### Basic Operations

#### Set & Get & Has

```go
m.Set("foo", 1)
m.Set("bar", 2)

added := m.TrySet("foo", 3) // added will be false, because "foo" already exists
added := m.TrySet("baz", 3) // added will be true, because "baz" does not exist

if val := m.Get("foo"); ok {
	fmt.Println(val)
}

if val, ok := m.TryGet("foo"); ok {
	fmt.Println(val)
}
if val, ok := m.TryGet("x"); !ok {
	fmt.Println("Key 'x' does not exist")
}

if m.Has("bar") {
	fmt.Println("Key 'bar' exists")
}
```

#### Delete & Clear

```go
m.Delete("foo") // Removes "foo"
m.Clear()	   // Removes all elements
```

#### Length

```go
fmt.Println(m.Len()) // Returns number of elements
```

### Iteration

Iterate over key-value pairs in insertion order:

```go
for k, v := range m.All() {
	fmt.Printf("%s: %d\n", k, v)
}
```

### Keys & Values

Retrieve all keys or values as a slice, maintaining the current order:

```go
keys := m.Keys()	 // Returns []K
values := m.Values() // Returns []V
```

### Merging Maps

Merge other `omap` instances into the current one:

```go
m1 := omap.New[string, int]()
m1.Set("a", 1)

m2 := omap.New[string, int]()
m2.Set("b", 2)

m1.Merge(m2)
// m1 now contains a:1, b:2
```

### Reversing Maps

Reverse the order of key-value pairs in-place:

```go
m := omap.New[string, int]()
m.Set("a", 1)
m.Set("b", 2)

m.Reverse()

fmt.Println(m.Keys())   // Output: [b a]
fmt.Println(m.Values()) // Output: [2 1]
```

### Sorting

`omap` provides in-place sorting methods, but requires keys to be `cmp.Ordered`.

```go
// Sort in ascending order
omap.Sort(m)

// Sort in descending order
omap.SortDesc(m)

// Sort with a custom comparison function
omap.SortFunc(m, func(k1, k2 string) int {
	// Example: sort by length of keys
	return len(k1) - len(k2)
})
```

### JSON Serialization

`omap` implements `json.Marshaler` and `json.Unmarshaler` interfaces, ensuring JSON objects preserve key order during
both marshaling and unmarshaling.

```go
// Marshal
data, err := json.Marshal(m)
// Output: {"foo":1,"bar":2}

// Unmarshal
err := json.Unmarshal(data, m)
```

