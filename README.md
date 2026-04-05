# list

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/bold-minds/list.svg)](https://pkg.go.dev/github.com/bold-minds/list)
[![Go Version](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/bold-minds/list/main/.github/badges/go-version.json)](https://golang.org/doc/go1.21)
[![Latest Release](https://img.shields.io/github/v/release/bold-minds/list?logo=github&color=blueviolet)](https://github.com/bold-minds/list/releases)
[![Last Updated](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/bold-minds/list/main/.github/badges/last-updated.json)](https://github.com/bold-minds/list/commits)
[![golangci-lint](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/bold-minds/list/main/.github/badges/golangci-lint.json)](https://github.com/bold-minds/list/actions/workflows/test.yaml)
[![Coverage](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/bold-minds/list/main/.github/badges/coverage.json)](https://github.com/bold-minds/list/actions/workflows/test.yaml)
[![Dependabot](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/bold-minds/list/main/.github/badges/dependabot.json)](https://github.com/bold-minds/list/security/dependabot)

**Set operations on Go slices.**

Go's `slices` package covers sorting, searching, and mutation, but it deliberately omits the operations that treat a slice as a set: true deduplication, union, intersection, difference, and removal of specific values. `list` provides those five operations as outcome-named standalone functions.

```go
// Before — stdlib slices.Compact only removes CONSECUTIVE duplicates
//          (you have to sort first, losing original order)
sorted := slices.Clone(tags)
slices.Sort(sorted)
unique := slices.Compact(sorted)

// After — preserves order of first occurrence, no sort required
unique := list.Unique(tags)
```

## ✨ Why list?

- 🧹 **`Unique` is true dedup** — preserves order of first occurrence, unlike `slices.Compact` which only removes consecutive duplicates
- 🔗 **`Union` across N slices** — variadic, returns unique elements in first-seen order
- ⚡ **`Intersect` across N slices** — variadic, elements present in every slice
- ➖ **`Minus` for binary difference** — `list.Minus(allUsers, banned)` reads like English
- 🚫 **`Without` removes specific values** — not a set operation, just "drop these elements"
- 🎯 **Order-preserving** — every operation returns elements in a stable, predictable order
- 🪶 **Five functions, one file, zero dependencies** — only what stdlib genuinely skipped

## 📦 Installation

```bash
go get github.com/bold-minds/list
```

Requires Go 1.21 or later.

## 🚀 Quick Start

```go
package main

import (
    "fmt"

    "github.com/bold-minds/list"
)

func main() {
    tags := []string{"go", "web", "api", "go", "web", "auth"}
    admins := []int{1, 2, 3, 4}
    editors := []int{3, 4, 5, 6}
    banned := []int{2, 5}

    // Deduplicate, preserve order of first occurrence
    unique := list.Unique(tags)
    fmt.Println(unique) // [go web api auth]

    // Combine slices, unique across all
    all := list.Union(admins, editors)
    fmt.Println(all) // [1 2 3 4 5 6]

    // Elements present in all input slices
    both := list.Intersect(admins, editors)
    fmt.Println(both) // [3 4]

    // Elements in the first slice but not in the second
    adminsOnly := list.Minus(admins, editors)
    fmt.Println(adminsOnly) // [1 2]

    // Remove specific values
    allowed := list.Without(admins, 2, 4)
    fmt.Println(allowed) // [1 3]

    // Compose them when needed
    activeAdmins := list.Minus(list.Union(admins, editors), banned)
    fmt.Println(activeAdmins) // [1 3 4 6]
}
```

## 🔧 Core Features

### `Unique` — true deduplication

Returns a new slice containing each unique value from the input, in order of first occurrence. Unlike `slices.Compact` (which only removes *consecutive* duplicates), `Unique` handles duplicates anywhere in the slice without requiring a sort first.

```go
list.Unique([]int{1, 2, 2, 3, 1, 4, 2})  // [1 2 3 4]
list.Unique([]string{"go", "web", "go"}) // [go web]
list.Unique([]int{})                     // []
```

Works with any `comparable` type: primitives, pointers, interfaces, and structs whose fields are comparable.

### `Union` — unique across N slices

Returns the unique elements across all provided slices, preserving order of first occurrence as you walk through the inputs.

```go
list.Union([]int{1, 2, 3}, []int{3, 4, 5})             // [1 2 3 4 5]
list.Union([]int{1}, []int{2}, []int{1, 3})            // [1 2 3]
list.Union([]string{"a", "b"}, []string{"b", "c"})     // [a b c]
list.Union[int]()                                      // []
```

Variadic — works with any number of input slices (including zero, which returns an empty slice).

### `Intersect` — elements in every slice

Returns the unique elements that appear in *every* provided slice. Order is taken from the first input slice.

```go
list.Intersect([]int{1, 2, 3, 4}, []int{2, 3, 5})          // [2 3]
list.Intersect([]int{1, 2, 3}, []int{2, 3}, []int{3, 4})   // [3]
list.Intersect[int]([]int{1, 2, 3})                        // [1 2 3]
list.Intersect[int]()                                      // []
```

Variadic, matching `Union`. A single-slice call returns the slice's unique elements (equivalent to `Unique`).

### `Minus` — binary set difference

Returns the unique elements of the first slice that are **not** present in the second slice. Order is taken from the first slice.

```go
list.Minus([]int{1, 2, 3, 4}, []int{2, 4})        // [1 3]
list.Minus([]string{"a", "b", "c"}, []string{"b"}) // [a c]
list.Minus([]int{1, 2, 3}, []int{})               // [1 2 3]
list.Minus([]int{}, []int{1, 2})                  // []
```

Binary (always exactly two arguments) because "subtract multiple slices" is ambiguous — is it `a − (b ∪ c)` or `((a − b) − c)`? If you need to subtract multiple, compose with `Union`: `list.Minus(a, list.Union(b, c, d))`.

### `Without` — remove specific values

Returns a new slice with all occurrences of the specified values removed, preserving the order of remaining elements.

```go
list.Without([]int{1, 2, 3, 2, 1}, 2)             // [1 3 1]
list.Without([]string{"a", "b", "c"}, "b", "c")   // [a]
list.Without([]int{1, 2, 3}, 5)                   // [1 2 3]
```

`Without` does **not** deduplicate remaining elements — if `1` appears three times in the input and you don't remove `1`, it appears three times in the output. For dedup + removal, chain: `list.Unique(list.Without(s, v))`.

## 🛡️ Safety guarantees

- **Never panics on valid input.** Nil slices, empty slices, and zero-variadic calls all return empty (non-nil) slices.
- **Immutable.** `list` never modifies input slices. Every function returns a new slice.
- **Order-preserving.** Every operation has a defined, stable output order.
- **Non-nil results.** Every function returns an empty (non-nil) slice rather than `nil` when there are no results. Safe to range over without nil checks.
- **Zero dependencies.** Pure stdlib underneath.

### NaN semantics

Floating-point `NaN` values follow Go's map-key semantics: `NaN` is never equal to `NaN` (including itself). This means `list` operations do not dedupe, intersect, or exclude NaNs:

```go
nan := math.NaN()
list.Unique([]float64{nan, nan, nan})       // [NaN NaN NaN] — all three remain
list.Minus([]float64{nan, 1.0}, []float64{nan}) // [NaN 1.0] — NaN not removed
```

If you need NaN-aware set operations, pre-process the slice to replace NaNs with a canonical marker.

### Non-comparable interface values

Go's `comparable` constraint is checked at compile time, but interface types can carry dynamic values whose concrete type is not comparable (e.g., a `[]int` stored in an `any`). Comparing such values causes a runtime panic. `list` does not recover from these panics — it is the caller's responsibility to ensure interface slices contain only comparable dynamic values.

## 🏎️ Performance

Measured on Go 1.26 (Intel Ultra 9 275HX; library targets Go 1.21+). All operations are O(n) in total input size.

```
BenchmarkUnique_Small-24             2087532    279.4 ns/op    408 B/op    4 allocs/op
BenchmarkUnique_WithDupes-24         2009478    287.8 ns/op    744 B/op    4 allocs/op
BenchmarkUnion_Two-24                1300902    477.7 ns/op    520 B/op    5 allocs/op
BenchmarkUnion_Three-24               593526    892.6 ns/op   1384 B/op    8 allocs/op
BenchmarkIntersect_Two-24             562664   1012 ns/op     1048 B/op   10 allocs/op
BenchmarkIntersect_Three-24           615122   1056 ns/op     1008 B/op   10 allocs/op
BenchmarkMinus_Basic-24              1402407    428.3 ns/op    736 B/op    7 allocs/op
BenchmarkWithout_SingleItem-24       9485799     61.2 ns/op     80 B/op    1 allocs/op
BenchmarkWithout_MultipleItems-24    6966316     91.9 ns/op     80 B/op    1 allocs/op
```

Allocations come from the map-backed lookups (unavoidable for O(1) set membership) and the result slice. `Without` is the cheapest because its "exclude these items" semantic does not require a dedup pass.

## 🧪 Testing

```bash
go test ./...                      # unit tests
go test -race ./...                # race detection
go test -bench=. -benchmem ./...   # benchmarks
```

Current coverage: 100%.

## 📚 API Reference

```go
// Unique returns a new slice with duplicate elements removed.
// Preserves order of first occurrence. Returns an empty (non-nil)
// slice for nil or empty input.
func Unique[T comparable](s []T) []T

// Union returns the unique elements across all provided slices,
// preserving order of first occurrence as each slice is walked in turn.
// Variadic — accepts zero or more slices.
func Union[T comparable](slices ...[]T) []T

// Intersect returns the unique elements present in every provided slice.
// Order is taken from the first slice. Variadic — a single-slice call
// is equivalent to Unique; a zero-slice call returns an empty slice.
func Intersect[T comparable](slices ...[]T) []T

// Minus returns the unique elements of a that are not present in b.
// Order is taken from a. Binary — always exactly two arguments.
func Minus[T comparable](a, b []T) []T

// Without returns a new slice with all occurrences of the specified items
// removed, preserving the order of remaining elements. Does NOT deduplicate
// remaining elements.
func Without[T comparable](s []T, items ...T) []T
```

## 🤝 Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Bold Minds Go libraries follow a shared set of design principles; read [PRINCIPLES.md](https://github.com/bold-minds/oss/blob/main/PRINCIPLES.md) before opening a PR.

## 📄 License

MIT. See [LICENSE](LICENSE).

## 🔗 Related Projects

- Go standard library [`slices`](https://pkg.go.dev/slices) — covers sort, reverse, contains, index, delete, concat. **Note:** `slices.Compact` only removes *consecutive* duplicates; for true deduplication, use `list.Unique`. The naming similarity between `Compact` and `Unique` is a known source of confusion.
- [`bold-minds/each`](https://github.com/bold-minds/each) — predicate and key-function operations on Go slices (Find, Filter, GroupBy, KeyBy, Partition, Count, Every). `each` handles per-element predicates on a single slice; `list` handles set semantics across one or more slices.
- [`samber/lo`](https://github.com/samber/lo) — comprehensive Go utility library. Includes `Uniq`, `Union`, `Intersect`, `Difference`, and `Without` among ~200 other helpers. `list` is the scoped version that ships only the five set operations and nothing else.
