# list

[![Go Reference](https://pkg.go.dev/badge/github.com/bold-minds/list.svg)](https://pkg.go.dev/github.com/bold-minds/list)
[![Build](https://img.shields.io/github/actions/workflow/status/bold-minds/list/test.yaml?branch=main&label=tests)](https://github.com/bold-minds/list/actions/workflows/test.yaml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/bold-minds/list)](go.mod)
[![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/bold-admin/dd52b54365ac3d2b99754f68a4bcb30d/raw/coverage.json)](https://github.com/bold-minds/list/actions/workflows/test.yaml)

**The list operations Go's stdlib leaves out — set algebra, positional extraction, sampling, reordering, sorting, and substitution — as plain generic functions on `[]T`.**

Go's `slices` package covers sort, search, and mutation of contiguous ranges. It deliberately omits true deduplication, set algebra across multiple slices, positional range extraction with safe bounds, random sampling, reordering, and value-based substitution — the operations you reach for when *shaping* a slice rather than searching or sorting it in place. `list` provides those operations as one family of outcome-named functions, in one repo, so you don't have to compose three libraries to express a single pipeline.

```go
// Deduplicate while preserving order of first occurrence
unique := list.Unique(tags)

// Take the first 10, then sort them ascending, then drop zero values
top := list.Compact(list.Sort(list.FirstN(entries, 10)))

// Sample 5 random users without replacement
sample := list.SampleN(users, 5)

// Merge two cohorts, drop the banned set, sort by signup date
active := list.SortBy(
    list.Minus(list.Union(admins, editors), banned),
    func(a, b User) int { return a.SignedUp.Compare(b.SignedUp) },
)
```

## ✨ Why list?

**Set operations** (comparable):
- 🧹 **`Unique` is true dedup** — preserves order of first occurrence, unlike `slices.Compact` which only removes consecutive duplicates
- 🔗 **`Union` across N slices** — variadic, returns unique elements in first-seen order
- ⚡ **`Intersect` across N slices** — variadic, elements present in every slice
- ➖ **`Minus` for binary difference** — `list.Minus(allUsers, banned)` reads like English
- ⇄ **`SymmetricDifference`** — elements in either slice but not both
- 🚫 **`Without` removes specific values** — not a set operation, just "drop these elements"

**Positional** (any) — safe on out-of-range without panics:
- 🔝 **`FirstN` / `LastN`** — take the head or tail, clamped to length
- 📐 **`Between(s, start, end)`** — like `s[start:end]` but bounds-clamped, never panics
- 🎯 **`At(s, i)`** — single-element access with negative indexes (`At(s, -1)` is the last), returns `(T, bool)`
- **`First` / `Last`** — aliases for `At(s, 0)` / `At(s, -1)`

**Sampling** (any):
- 🎲 **`Sample` / `SampleN`** — uniformly random selection via `math/rand/v2`

**Reordering** (any):
- ↩️ **`Reverse`** — new reversed slice, input untouched
- 🔀 **`Shuffle`** — new randomly-permuted slice, input untouched

**Sorting:**
- 📊 **`Sort` / `SortDesc`** (`cmp.Ordered`) — new sorted slice, ascending or descending. NaN sorts first ascending, last descending.
- 🧮 **`SortBy`** (any) — new sorted slice with a caller-supplied `func(a, b T) int` comparator for non-ordered types

**Zero stripping** (comparable):
- ✨ **`Compact`** — new slice with every element equal to T's zero value removed. (Distinct from stdlib `slices.Compact`, which removes *consecutive duplicates*.)

**Substitution** (comparable):
- 🔁 **`Replace` / `ReplaceFirst`** — value-based replacement, all occurrences or first only

**Contracts for everything:**
- 🛡️ **Order-preserving** — every operation has a defined, stable output order
- 🪶 **Zero dependencies** — pure Go stdlib; `math/rand/v2` for sampling, `slices` for the sort primitives
- 🔒 **Never panics on valid input** — nil, empty, out-of-range all return empty (non-nil) slices or `(_, false)`
- 🧊 **Immutable** — no function mutates its input; every result is a fresh allocation

## 📦 Installation

```bash
go get github.com/bold-minds/list
```

Requires Go 1.23 or later.

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

### `FirstN` / `LastN` / `Between` — positional ranges, bounds-clamped

Out-of-range indexes clamp rather than panic. Negative counts become empty slices. Over-large counts return the whole (copied) slice.

```go
s := []int{10, 20, 30, 40, 50}

list.FirstN(s, 3)           // [10 20 30]
list.FirstN(s, 100)         // [10 20 30 40 50]  (clamped, not out-of-range)
list.FirstN(s, 0)           // []
list.LastN(s, 2)            // [40 50]
list.Between(s, 1, 4)       // [20 30 40]
list.Between(s, -10, 100)   // [10 20 30 40 50]  (both bounds clamped)
list.Between(s, 4, 2)       // []                (inverted range)
```

`Between` is the non-panicking replacement for `s[start:end]` when indexes come from user input, config, or any source where out-of-range is a possibility rather than a bug.

### `At` / `First` / `Last` — single-element access with negative-index support

Returns `(value, ok)` — `ok` is `false` for empty slices or out-of-range indexes. Negative indexes count from the end, so `At(s, -1)` is the last element.

```go
s := []string{"a", "b", "c"}

v, ok := list.At(s, 0)       // ("a", true)
v, ok  = list.At(s, -1)      // ("c", true)
v, ok  = list.At(s, 99)      // ("", false)
v, ok  = list.First(s)       // ("a", true)
v, ok  = list.Last(s)        // ("c", true)
```

Use `At` when you want to extract a single element defensively without the `if i < len(s) { … }` guard at every call site.

### `Sample` / `SampleN` — random selection

`Sample` returns one uniformly-random element; `SampleN` returns *n* distinct random elements (without replacement) via a Fisher–Yates partial shuffle.

```go
users := []User{ /* ... */ }

u, ok := list.Sample(users)          // one random user
five := list.SampleN(users, 5)       // 5 distinct random users
all  := list.SampleN(users, 1000)    // full shuffle if n >= len
```

Uses `math/rand/v2`'s goroutine-safe top-level source. **Not cryptographically secure** — fine for tests, fixtures, A/B bucket assignment, and UI randomization; use `crypto/rand` directly if you need unpredictability guarantees.

### `Reverse` / `Shuffle` — reordering without a key

Both return a new slice; the input is untouched.

```go
list.Reverse([]int{1, 2, 3})  // [3 2 1]
list.Shuffle(users)           // users in random order, input untouched
```

### `Sort` / `SortDesc` / `SortBy` — non-mutating sort

Unlike stdlib `slices.Sort` which sorts in place, `list.Sort` returns a fresh sorted slice. The input is never touched — the same guarantee as every other `list` function.

```go
// Ordered types: cmp.Ordered constraint covers strings, all integer
// widths, and both float types.
list.Sort([]int{3, 1, 2})            // [1 2 3]
list.Sort([]string{"c", "a", "b"})   // [a b c]
list.SortDesc([]int{3, 1, 2})        // [3 2 1]

// Non-ordered types: supply a comparator.
users := []User{ /* ... */ }
byName := list.SortBy(users, func(a, b User) int {
    return strings.Compare(a.Name, b.Name)
})
```

`NaN` sorts before every non-`NaN` value under `Sort` (and after under `SortDesc`), following `cmp.Compare`'s defined behavior.

### `Compact` — strip zero values

Returns a new slice with every element equal to `T`'s zero value removed. This is distinct from stdlib `slices.Compact`, which removes *consecutive duplicates* regardless of value — the naming collision is unfortunate but irreversible.

```go
list.Compact([]string{"a", "", "b", "", "c"})  // [a b c]
list.Compact([]int{1, 0, 2, 0, 3})             // [1 2 3]
list.Compact([]*User{u1, nil, u2, nil})        // [u1 u2]
```

Requires `T comparable` so the zero value can be detected without reflection.

### `Replace` / `ReplaceFirst` — value-based substitution

```go
list.Replace([]int{1, 2, 3, 2, 1}, 2, 99)         // [1 99 3 99 1]
list.ReplaceFirst([]int{1, 2, 3, 2, 1}, 2, 99)    // [1 99 3 2 1]
```

Positional replacement (`s[i] = v`) is already a one-liner in Go, so `list` does not ship an indexed variant.

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

Measured on Go 1.26 (Intel Ultra 9 275HX; library targets Go 1.23+). All operations are O(n) in total input size.

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
// =========================================================================
// Set operations (T comparable)
// =========================================================================

// Unique returns a new slice with duplicate elements removed, preserving
// order of first occurrence.
func Unique[T comparable](s []T) []T

// Union returns the unique elements across all provided slices,
// preserving order of first occurrence as each slice is walked in turn.
func Union[T comparable](slices ...[]T) []T

// Intersect returns the unique elements present in every provided slice.
// Order is taken from the first slice.
func Intersect[T comparable](slices ...[]T) []T

// SymmetricDifference returns unique elements present in a or b but
// not both. Order is taken from a first, then b.
func SymmetricDifference[T comparable](a, b []T) []T

// Minus returns the unique elements of a that are not present in b.
// Binary — for multi-slice subtraction, compose with Union.
func Minus[T comparable](a, b []T) []T

// Without returns a new slice with all occurrences of the specified
// items removed. Does NOT deduplicate remaining elements.
func Without[T comparable](s []T, items ...T) []T

// =========================================================================
// Positional extraction (T any)
// =========================================================================

// FirstN returns the first n elements, clamped to valid range.
func FirstN[T any](s []T, n int) []T

// LastN returns the last n elements, clamped to valid range.
func LastN[T any](s []T, n int) []T

// Between returns s[start:end] with bounds clamped. Never panics.
func Between[T any](s []T, start, end int) []T

// At returns the element at index i with ok=true, or (zero, false)
// on empty slice or out-of-range. Negative indexes count from the end.
func At[T any](s []T, i int) (T, bool)

// First is an alias for At(s, 0).
func First[T any](s []T) (T, bool)

// Last is an alias for At(s, -1).
func Last[T any](s []T) (T, bool)

// =========================================================================
// Sampling (T any)
// =========================================================================

// Sample returns one uniformly-random element via math/rand/v2.
func Sample[T any](s []T) (T, bool)

// SampleN returns n distinct random elements without replacement.
func SampleN[T any](s []T, n int) []T

// =========================================================================
// Reordering (T any)
// =========================================================================

// Reverse returns a new slice in reverse order. Input not mutated.
func Reverse[T any](s []T) []T

// Shuffle returns a new randomly-permuted slice. Input not mutated.
func Shuffle[T any](s []T) []T

// =========================================================================
// Sorting
// =========================================================================

// Sort returns a new slice sorted ascending. T must be cmp.Ordered.
// Input not mutated.
func Sort[T cmp.Ordered](s []T) []T

// SortDesc returns a new slice sorted descending. T must be cmp.Ordered.
func SortDesc[T cmp.Ordered](s []T) []T

// SortBy returns a new slice sorted by a caller-supplied comparator.
// T is any; less must return negative/zero/positive for a<b / a==b / a>b.
func SortBy[T any](s []T, less func(a, b T) int) []T

// =========================================================================
// Zero stripping (T comparable)
// =========================================================================

// Compact returns a new slice with zero values of T removed.
// Distinct from stdlib slices.Compact, which removes consecutive duplicates.
func Compact[T comparable](s []T) []T

// =========================================================================
// Substitution (T comparable)
// =========================================================================

// Replace returns a new slice with every occurrence of old replaced.
func Replace[T comparable](s []T, old, newVal T) []T

// ReplaceFirst returns a new slice with the first occurrence of old replaced.
func ReplaceFirst[T comparable](s []T, old, newVal T) []T
```

## 🤝 Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Bold Minds Go libraries follow a shared set of design principles; read [PRINCIPLES.md](https://github.com/bold-minds/oss/blob/main/PRINCIPLES.md) before opening a PR.

## 📄 License

MIT. See [LICENSE](LICENSE).

## 🔗 Related Projects

- Go standard library [`slices`](https://pkg.go.dev/slices) — covers sort, reverse, contains, index, delete, concat. **Note:** `slices.Compact` only removes *consecutive* duplicates; for true deduplication, use `list.Unique`. The naming similarity between `Compact` and `Unique` is a known source of confusion.
- [`bold-minds/each`](https://github.com/bold-minds/each) — predicate and key-function operations on Go slices (Find, Filter, GroupBy, KeyBy, Partition, Count, Every). `each` handles per-element predicates on a single slice; `list` handles set semantics across one or more slices.
- [`samber/lo`](https://github.com/samber/lo) — comprehensive Go utility library. Includes `Uniq`, `Union`, `Intersect`, `Difference`, and `Without` among ~200 other helpers. `list` is the scoped version that ships only the five set operations and nothing else.
