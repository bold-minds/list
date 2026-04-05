// Package list provides set operations on Go slices.
//
// Unique, Union, Intersect, Minus, and Without treat slices as sets,
// returning new slices with predictable, stable output ordering. Every
// operation is nil-safe, never panics on nil or empty input, and never
// mutates its input. Empty results are returned as non-nil empty slices,
// so callers can range over them without nil checks.
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
	result := make([]T, 0, len(s))
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
	seen := make(map[T]struct{})
	result := []T{}
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
func Intersect[T comparable](slices ...[]T) []T {
	if len(slices) == 0 {
		return []T{}
	}
	if len(slices) == 1 {
		return Unique(slices[0])
	}

	// Count how many slices each element appears in, counting each
	// element at most once per slice.
	counts := make(map[T]int)
	for i, s := range slices {
		seen := make(map[T]struct{})
		for _, v := range s {
			if _, exists := seen[v]; exists {
				continue
			}
			seen[v] = struct{}{}
			if i == 0 {
				counts[v] = 1
			} else if counts[v] == i {
				counts[v]++
			}
		}
	}

	// Emit elements present in all slices, in order from the first slice.
	result := []T{}
	emitted := make(map[T]struct{})
	for _, v := range slices[0] {
		if _, already := emitted[v]; already {
			continue
		}
		if counts[v] == len(slices) {
			emitted[v] = struct{}{}
			result = append(result, v)
		}
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
		result := make([]T, len(s))
		copy(result, s)
		return result
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
