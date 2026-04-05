# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- **Scope expanded** from "set operations on Go slices" to "list operations on Go slices" — the package now covers positional extraction, sampling, reordering, sorting, zero stripping, and value substitution in addition to the original set algebra. This is the architecturally honest home for these operations (Ruby, Python, and JS all treat positional and set ops as methods on the same list type), and it avoids forcing callers to import a separate repo to operate on the same data shape. Package doc comment and README rewritten to reflect the new scope. This is purely additive — every v0.1.x call site continues to work unchanged.

### Added

**Positional extraction (`T any`):**
- `FirstN[T](s, n)` — first n elements, safe on out-of-range.
- `LastN[T](s, n)` — last n elements, safe on out-of-range.
- `Between[T](s, start, end)` — `s[start:end]` with bounds clamped to valid range. Never panics on out-of-range indexes.
- `At[T](s, i) (T, bool)` — single element access with negative-index support (`At(s, -1)` returns the last element) and ok-on-in-range.
- `First[T](s) (T, bool)` — alias for `At(s, 0)`.
- `Last[T](s) (T, bool)` — alias for `At(s, -1)`.

**Sampling (`T any`):**
- `Sample[T](s) (T, bool)` — one uniformly-random element via `math/rand/v2`.
- `SampleN[T](s, n)` — n distinct random elements via Fisher–Yates partial shuffle.

**Reordering (`T any`):**
- `Reverse[T](s)` — new reversed slice, input untouched.
- `Shuffle[T](s)` — new randomly-permuted slice, input untouched.

**Sorting:**
- `Sort[T cmp.Ordered](s)` — new ascending-sorted slice for strings, numerics, and named ordered types. NaN sorts first (cmp.Compare semantics).
- `SortDesc[T cmp.Ordered](s)` — new descending-sorted slice.
- `SortBy[T any](s, less)` — new sorted slice with caller-supplied `func(a, b T) int` comparator for non-ordered types.

**Zero stripping (`T comparable`):**
- `Compact[T](s)` — new slice with every element equal to T's zero value removed. Distinct from stdlib `slices.Compact`, which removes consecutive duplicates.

**Substitution (`T comparable`):**
- `Replace[T](s, old, newVal)` — new slice with every occurrence of `old` replaced.
- `ReplaceFirst[T](s, old, newVal)` — new slice with the first occurrence replaced.

### Added — existing set operations
- `SymmetricDifference[T comparable](a, b []T) []T` — unique elements present in `a` or `b` but not both, i.e. `(a ∪ b) − (a ∩ b)` (from earlier unreleased work).

## [0.1.0] — Initial release

### Added
- `Unique[T comparable](s []T) []T` — deduplicate a slice, preserving order of first occurrence
- `Union[T comparable](slices ...[]T) []T` — unique elements across N slices, variadic
- `Intersect[T comparable](slices ...[]T) []T` — elements present in every provided slice, variadic
- `Minus[T comparable](a, b []T) []T` — elements of `a` not in `b`, binary
- `Without[T comparable](s []T, items ...T) []T` — remove specific values, preserves remaining duplicates
- Full support for custom comparable types (`type UserID string`), struct types with comparable fields, and pointer types
- Documented NaN semantics for floating-point slices (follows Go's map-key rules)
- 100% statement coverage (enforced in CI via `COVERAGE_THRESHOLD=100`) including adversarial edge cases: nil/empty handling, immutability, result-aliasing checks, NaN behavior, struct key correctness, pointer-identity semantics, and runtime panics on non-comparable interface values
- Zero external dependencies — pure stdlib

### Deliberate non-goals
- No operations on maps (use stdlib `maps`)
- No sorting (use stdlib `slices.Sort`)
- No predicate-based operations on a single slice (those live in `bold-minds/each`)
- No NaN-aware float operations — caller must pre-process slices if NaN matching is required

### Requires
- Go 1.23 or later
