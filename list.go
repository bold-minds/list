// Package list provides the list operations Go's stdlib leaves out —
// set algebra, positional extraction, sampling, reordering, sorting,
// and substitution — as plain generic functions on []T.
//
// # Categories
//
// Set operations (T comparable): Unique, Union, Intersect,
// SymmetricDifference, Minus, Without. Treat slices as sets with
// predictable, stable output ordering.
//
// Positional extraction (T any): FirstN, LastN, Between, At, First,
// Last. Extract by index or range. Negative indexes count from the end
// where meaningful; out-of-range indexes clamp or return (_, false)
// rather than panicking.
//
// Sampling (T any): Sample, SampleN. Random element selection from
// math/rand/v2. Use crypto/rand directly if you need cryptographic
// entropy.
//
// Reordering (T any): Reverse, Shuffle. Return new slices — input is
// never mutated.
//
// Sorting: Sort and SortDesc (T cmp.Ordered) for strings, numbers, and
// named types with ordered underlying kinds. SortBy (T any) for
// anything else via a caller-supplied comparator.
//
// Zero stripping (T comparable): Compact. Removes every element equal
// to T's zero value.
//
// Substitution (T comparable): Replace, ReplaceFirst. Value-based
// replacement; replace-by-index is intentionally a one-liner (s[i] = v)
// and not exposed.
//
// # Contracts
//
// Every operation is nil-safe, never panics on nil or empty input, and
// never mutates its input. Empty results are returned as non-nil empty
// slices, so callers can range over them without nil checks.
//
// All functions are safe for concurrent use on shared inputs: they
// never mutate their arguments, so multiple goroutines may call any
// function with overlapping or identical input slices without external
// synchronization. Each call allocates its own result slice, so
// callers may freely mutate the returned slice without affecting
// inputs or other callers' results.
//
// # NaN semantics for set operations
//
// Float64 and float32 values follow Go's map-key semantics: NaN is
// never equal to NaN (including itself), so each NaN is a distinct
// map key. This has consequences for the set operations that use a
// map internally (Unique, Union, Intersect, SymmetricDifference,
// Minus, Without): none of them can match NaN against NaN, and none
// can deduplicate NaN values. The sort and positional operations are
// unaffected — they do not use equality-based deduplication.
//
// If you need NaN-aware set operations, pre-process your slice to
// replace NaNs with a canonical marker value.
//
// # Runtime panics from non-comparable interface values
//
// Go's comparable constraint is checked at compile time, but interface
// types can carry dynamic values whose concrete type is not comparable
// (e.g., slices or maps stored in an any). Comparing such values
// causes a runtime panic. list does not recover from these panics —
// it is the caller's responsibility to ensure interface slices
// contain only comparable dynamic values.
//
// For documentation and examples, see https://github.com/bold-minds/list.
package list

import (
	"cmp"
	mrand "math/rand/v2"
	"slices"
)

// Unique returns a new slice with duplicate elements removed, preserving
// the order of first occurrence. Returns an empty (non-nil) slice for
// nil or empty input.
func Unique[T comparable](s []T) []T {
	if len(s) == 0 {
		return []T{}
	}
	seen := make(map[T]struct{}, len(s))
	// Start result with zero capacity; append grows it on demand. This avoids
	// holding a backing array sized for the worst case when most elements
	// are duplicates.
	result := make([]T, 0)
	for _, v := range s {
		if _, exists := seen[v]; !exists {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

// Union returns the unique elements across all provided slices, preserving
// the order of first occurrence as each slice is walked in turn. Variadic —
// accepts zero or more slices.
func Union[T comparable](slices ...[]T) []T {
	total := 0
	for _, s := range slices {
		total += len(s)
	}
	seen := make(map[T]struct{}, total)
	// Zero initial cap: sum-of-lengths overallocates heavily when slices
	// overlap. Append grows on demand.
	result := make([]T, 0)
	for _, s := range slices {
		for _, v := range s {
			if _, exists := seen[v]; !exists {
				seen[v] = struct{}{}
				result = append(result, v)
			}
		}
	}
	return result
}

// Intersect returns the unique elements present in every provided slice.
// Order is taken from the first slice. Variadic — a single-slice call
// is equivalent to Unique; a zero-slice call returns an empty slice.
// A nil or empty slice anywhere in the input causes the result to be
// empty, since no element can appear in every slice.
func Intersect[T comparable](slices ...[]T) []T {
	if len(slices) == 0 {
		return []T{}
	}
	if len(slices) == 1 {
		return Unique(slices[0])
	}

	// Count how many slices each element appears in, counting each
	// element at most once per slice. A single `seen` map is allocated
	// once and reused across iterations via clear() to avoid per-slice
	// map allocation.
	counts := make(map[T]int, len(slices[0]))
	var seen map[T]struct{}
	for i, s := range slices {
		if seen == nil {
			seen = make(map[T]struct{}, len(s))
		} else {
			clear(seen)
		}
		for _, v := range s {
			if _, exists := seen[v]; exists {
				continue
			}
			seen[v] = struct{}{}
			if i == 0 {
				counts[v] = 1
			} else if counts[v] == i {
				// counts[v] == i means v has appeared in every previous
				// slice; only then do we count this slice's occurrence.
				counts[v]++
			}
		}
	}

	// Emit elements present in all slices, in order from the first slice.
	// We consume the counts map itself as the "already emitted" marker by
	// deleting entries as we emit them, avoiding a second map allocation
	// and a second lookup per element.
	n := len(slices)
	result := make([]T, 0, len(counts))
	for _, v := range slices[0] {
		if counts[v] == n {
			result = append(result, v)
			delete(counts, v)
		}
	}
	return result
}

// SymmetricDifference returns the unique elements present in a or b but not
// in both — i.e., (a ∪ b) − (a ∩ b). Order is taken from a first, then b.
// Binary — always exactly two arguments.
func SymmetricDifference[T comparable](a, b []T) []T {
	if len(a) == 0 {
		return Unique(b)
	}
	if len(b) == 0 {
		return Unique(a)
	}
	inA := make(map[T]struct{}, len(a))
	for _, v := range a {
		inA[v] = struct{}{}
	}
	inB := make(map[T]struct{}, len(b))
	for _, v := range b {
		inB[v] = struct{}{}
	}
	// Two maps suffice: inB doubles as the "skip set" while iterating a
	// (skip if v is in b, otherwise emit and add to inB to dedupe further
	// occurrences in a). Symmetrically, inA is reused while iterating b.
	result := make([]T, 0)
	for _, v := range a {
		if _, skip := inB[v]; skip {
			continue
		}
		inB[v] = struct{}{} // mark emitted so later duplicates in a are skipped
		result = append(result, v)
	}
	for _, v := range b {
		if _, skip := inA[v]; skip {
			continue
		}
		inA[v] = struct{}{} // mark emitted so later duplicates in b are skipped
		result = append(result, v)
	}
	return result
}

// Minus returns the unique elements of a that are not present in b.
// Order is taken from a. Binary — always exactly two arguments.
// To subtract multiple slices, compose with Union:
//
//	list.Minus(a, list.Union(b, c, d))
func Minus[T comparable](a, b []T) []T {
	if len(a) == 0 {
		return []T{}
	}
	exclude := make(map[T]struct{}, len(b))
	for _, v := range b {
		exclude[v] = struct{}{}
	}
	seen := make(map[T]struct{}, len(a))
	result := make([]T, 0, len(a))
	for _, v := range a {
		if _, excluded := exclude[v]; excluded {
			continue
		}
		if _, already := seen[v]; already {
			continue
		}
		seen[v] = struct{}{}
		result = append(result, v)
	}
	return result
}

// Without returns a new slice with all occurrences of the specified items
// removed, preserving the order of remaining elements.
//
// Unlike Unique/Union/Intersect/Minus, Without does NOT deduplicate
// remaining elements. If s contains duplicates of values that are not in
// items, those duplicates remain in the output. For dedup + removal,
// chain: list.Unique(list.Without(s, v)).
func Without[T comparable](s []T, items ...T) []T {
	if len(s) == 0 {
		return []T{}
	}
	if len(items) == 0 {
		// len(s) > 0 from the check above, so append returns a fresh,
		// non-nil slice (preserving the no-alias contract).
		return append([]T(nil), s...)
	}
	exclude := make(map[T]struct{}, len(items))
	for _, v := range items {
		exclude[v] = struct{}{}
	}
	result := make([]T, 0, len(s))
	for _, v := range s {
		if _, excluded := exclude[v]; !excluded {
			result = append(result, v)
		}
	}
	return result
}

// =============================================================================
// Positional extraction — take a subset by index or range
// =============================================================================

// FirstN returns a new slice containing the first n elements of s.
// If n <= 0 returns an empty slice. If n >= len(s) returns a copy of s.
func FirstN[T any](s []T, n int) []T {
	if n <= 0 || len(s) == 0 {
		return []T{}
	}
	if n >= len(s) {
		return append([]T(nil), s...)
	}
	return append([]T(nil), s[:n]...)
}

// LastN returns a new slice containing the last n elements of s.
// If n <= 0 returns an empty slice. If n >= len(s) returns a copy of s.
func LastN[T any](s []T, n int) []T {
	if n <= 0 || len(s) == 0 {
		return []T{}
	}
	if n >= len(s) {
		return append([]T(nil), s...)
	}
	return append([]T(nil), s[len(s)-n:]...)
}

// DropL returns a new slice with the first n elements removed. It is
// the complement of FirstN: FirstN keeps the head, DropL discards it.
// If n <= 0, DropL returns a full copy of s. If n >= len(s), it
// returns an empty slice. Input is never mutated.
//
//	list.DropL([]int{1, 2, 3, 4, 5}, 2) // [3 4 5]
//	list.DropL([]int{1, 2, 3}, 10)      // []
//	list.DropL([]int{1, 2, 3}, 0)       // [1 2 3]
func DropL[T any](s []T, n int) []T {
	if n <= 0 {
		return append([]T(nil), s...)
	}
	if n >= len(s) {
		return []T{}
	}
	return append([]T(nil), s[n:]...)
}

// DropR returns a new slice with the last n elements removed. It is
// the complement of LastN: LastN keeps the tail, DropR discards it.
// If n <= 0, DropR returns a full copy of s. If n >= len(s), it
// returns an empty slice. Input is never mutated.
//
//	list.DropR([]int{1, 2, 3, 4, 5}, 2) // [1 2 3]
//	list.DropR([]int{1, 2, 3}, 10)      // []
//	list.DropR([]int{1, 2, 3}, 0)       // [1 2 3]
func DropR[T any](s []T, n int) []T {
	if n <= 0 {
		return append([]T(nil), s...)
	}
	if n >= len(s) {
		return []T{}
	}
	return append([]T(nil), s[:len(s)-n]...)
}

// Between returns a new slice containing s[start:end], clamped to valid
// bounds. Negative start is clamped to 0; end greater than len(s) is
// clamped to len(s). If the resulting range is empty or inverted
// (start >= end after clamping), an empty slice is returned.
//
// Unlike Go's native slice expression, Between never panics on
// out-of-range indexes.
func Between[T any](s []T, start, end int) []T {
	if len(s) == 0 {
		return []T{}
	}
	if start < 0 {
		start = 0
	}
	if end > len(s) {
		end = len(s)
	}
	if start >= end {
		return []T{}
	}
	return append([]T(nil), s[start:end]...)
}

// At returns the element at index i and whether the index was in range.
// Negative indexes count from the end: At(s, -1) returns the last element.
// Returns (zero, false) for an empty slice or out-of-range index.
func At[T any](s []T, i int) (T, bool) {
	var zero T
	if len(s) == 0 {
		return zero, false
	}
	if i < 0 {
		i += len(s)
	}
	if i < 0 || i >= len(s) {
		return zero, false
	}
	return s[i], true
}

// First returns the first element of s and ok=true, or (zero, false)
// for an empty slice. Equivalent to At(s, 0).
func First[T any](s []T) (T, bool) {
	return At(s, 0)
}

// Last returns the last element of s and ok=true, or (zero, false)
// for an empty slice. Equivalent to At(s, -1).
func Last[T any](s []T) (T, bool) {
	return At(s, -1)
}

// =============================================================================
// Sampling — random selection via math/rand/v2
// =============================================================================

// Sample returns one uniformly-random element of s and ok=true, or
// (zero, false) for an empty slice. Uses math/rand/v2's top-level
// goroutine-safe source — adequate for tests, fixtures, and UI
// randomization, but NOT for cryptographic purposes. Use crypto/rand
// directly if you need unpredictability guarantees.
func Sample[T any](s []T) (T, bool) {
	var zero T
	if len(s) == 0 {
		return zero, false
	}
	return s[mrand.IntN(len(s))], true //nolint:gosec // G404: documented non-crypto use
}

// SampleN returns n distinct elements chosen uniformly at random from
// s, in random order. If n <= 0 returns empty; if n >= len(s) returns
// a Shuffle'd copy of s (implicitly capping n at len(s)).
//
// The algorithm is a Fisher–Yates partial shuffle: O(n) time, O(len(s))
// space for the working copy.
func SampleN[T any](s []T, n int) []T {
	if n <= 0 || len(s) == 0 {
		return []T{}
	}
	if n >= len(s) {
		return Shuffle(s)
	}
	working := append([]T(nil), s...)
	for i := 0; i < n; i++ {
		j := i + mrand.IntN(len(working)-i) //nolint:gosec // G404: documented non-crypto use
		working[i], working[j] = working[j], working[i]
	}
	return working[:n:n]
}

// =============================================================================
// Reordering — reverse and shuffle without a key
// =============================================================================

// Reverse returns a new slice containing the elements of s in reverse
// order. Input is not mutated.
func Reverse[T any](s []T) []T {
	if len(s) == 0 {
		return []T{}
	}
	result := make([]T, len(s))
	for i, v := range s {
		result[len(s)-1-i] = v
	}
	return result
}

// Shuffle returns a new slice containing the elements of s in a
// uniformly-random order. Input is not mutated. Uses the same
// math/rand/v2 source as Sample — adequate for tests and UI but not
// for cryptographic purposes.
func Shuffle[T any](s []T) []T {
	if len(s) == 0 {
		return []T{}
	}
	result := append([]T(nil), s...)
	mrand.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})
	return result
}

// =============================================================================
// Sorting — ordered by default, custom comparator escape hatch
// =============================================================================

// Sort returns a new slice containing the elements of s in ascending
// order. T must satisfy cmp.Ordered (strings, all integer widths, both
// float types, and named types with those underlying kinds). Input is
// not mutated.
//
// NaN handling follows cmp.Compare: NaN sorts before every non-NaN
// value.
func Sort[T cmp.Ordered](s []T) []T {
	if len(s) == 0 {
		return []T{}
	}
	result := slices.Clone(s)
	slices.Sort(result)
	return result
}

// SortDesc returns a new slice containing the elements of s in
// descending order. T must satisfy cmp.Ordered. Input is not mutated.
//
// NaN handling is the inverse of Sort: NaN sorts after every non-NaN
// value (because the comparator is reversed).
func SortDesc[T cmp.Ordered](s []T) []T {
	if len(s) == 0 {
		return []T{}
	}
	result := slices.Clone(s)
	slices.SortFunc(result, func(a, b T) int { return cmp.Compare(b, a) })
	return result
}

// SortBy returns a new slice containing the elements of s sorted by
// the caller-supplied comparator. less must return negative / zero /
// positive for a<b / a==b / a>b, matching slices.SortFunc. Input is
// not mutated.
//
// Use SortBy when T is not cmp.Ordered — structs, interfaces, pointers,
// custom types with non-ordered underlying kinds — or when the sort
// key is a derived field rather than the element itself.
func SortBy[T any](s []T, less func(a, b T) int) []T {
	if len(s) == 0 {
		return []T{}
	}
	result := slices.Clone(s)
	slices.SortFunc(result, less)
	return result
}

// =============================================================================
// Zero stripping
// =============================================================================

// Compact returns a new slice with every element equal to T's zero
// value removed. The zero value is whatever Go produces for `var v T`:
// `""` for strings, `0` for numerics, `false` for bool, nil for
// pointers/maps/slices/interfaces, and the all-zero struct for
// structs with comparable fields.
//
// Compact requires T comparable so zero detection is expressed in pure
// Go (no reflection). Input is not mutated.
//
// This is distinct from stdlib slices.Compact, which removes
// CONSECUTIVE duplicates. Use list.Unique for order-preserving full
// deduplication and list.Compact for zero-value stripping.
func Compact[T comparable](s []T) []T {
	if len(s) == 0 {
		return []T{}
	}
	var zero T
	result := make([]T, 0, len(s))
	for _, v := range s {
		if v != zero {
			result = append(result, v)
		}
	}
	return result
}

// =============================================================================
// Substitution — value-based positional replacement
// =============================================================================

// Replace returns a new slice in which every element equal to old is
// replaced with newVal. Input is not mutated.
//
// Use stdlib `s[i] = v` directly for single-index replacement — list
// does not ship a positional Replace because Go already has it in the
// language.
func Replace[T comparable](s []T, old, newVal T) []T {
	if len(s) == 0 {
		return []T{}
	}
	result := make([]T, len(s))
	for i, v := range s {
		if v == old {
			result[i] = newVal
		} else {
			result[i] = v
		}
	}
	return result
}

// ReplaceFirst returns a new slice in which the first element equal
// to old is replaced with newVal. Subsequent occurrences are left
// unchanged. If old is not found, the returned slice is a copy of s
// with no changes. Input is not mutated.
func ReplaceFirst[T comparable](s []T, old, newVal T) []T {
	if len(s) == 0 {
		return []T{}
	}
	result := make([]T, len(s))
	copy(result, s)
	for i, v := range s {
		if v == old {
			result[i] = newVal
			break
		}
	}
	return result
}
