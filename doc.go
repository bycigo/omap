// Licensed under the MIT license, see LICENSE file for details.

/*
Package omap provides an ordered map implementation for Go.

Unlike Go's standard map, omap maintains the insertion order of keys.
It leverages Go iterators (iter.Seq2) and supports generic key-value types.

# Features

  - Ordered: Maintains insertion order of elements
  - Generic: Supports any comparable key and any value type
  - Efficient: O(1) time complexity for Set, Get, Delete, and Has operations
  - Iteration: Supports range iterator via All() method
  - JSON: Supports JSON marshaling/unmarshaling, preserving key order
  - Sorting: Provides in-place sorting capabilities
*/
package omap
