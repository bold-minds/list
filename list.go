// Package list provides set operations on Go slices.
//
// Unique, Union, Intersect, SymmetricDifference, Minus, and Without treat slices as sets,
// returning new slices with predictable, stable output ordering. Every
// operation is nil-safe, never panics on nil or empty input, and never
// mutates its input. Empty results are returned as non-nil empty slices,
// so callers can range over them without nil checks.
//
// All functions are safe for concurrent use on shared inputs: they never
// mutate their arguments, so multiple goroutines may call any function with
// overlapping or identical input slices without external synchronization.
// Each call allocates its own result slice, so callers may freely mutate
// the returned slice without affecting inputs or other callers' results.
//
// Every function takes a type parameter T constrained to comparable.
// This means element types must support the == operator: any primitive
// type, any pointer, any channel, any interface type, or any struct
// whose fields are themselves comparable.
//
// # NaN semantics
//
// Float64 and float32 values follow Go's map-key semantics: NaN is never
// equal to NaN (including itself), so each NaN is a distinct map key.
// This has consequences for list operations:
//
//   - Unique cannot deduplicate NaN values
//   - Union cannot deduplicate NaN values across slices
//   - Intersect can never match NaN against NaN
//   - SymmetricDifference cannot pair NaN values (every NaN is distinct)
//   - Minus cannot remove NaN values
//   - Without cannot remove NaN values
//
// If you need NaN-aware set operations, pre-process your slice to
// replace NaNs with a canonical marker value.
//
// # Runtime panics from non-comparable interface values
//
// Go's comparable constraint is checked at compile time, but interface
// types can carry dynamic values whose concrete type is not comparable
// (e.g., slices or maps stored in an any). Comparing such values causes
// a runtime panic. list does not recover from these panics — it is the
// caller's responsibility to ensure interface slices contain only
// comparable dynamic values.
//
// For documentation and examples, see https://github.com/bold-minds/list.
package list

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
