# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] — Initial release

### Added
- `Unique[T comparable](s []T) []T` — deduplicate a slice, preserving order of first occurrence
- `Union[T comparable](slices ...[]T) []T` — unique elements across N slices, variadic
- `Intersect[T comparable](slices ...[]T) []T` — elements present in every provided slice, variadic
- `Minus[T comparable](a, b []T) []T` — elements of `a` not in `b`, binary
- `Without[T comparable](s []T, items ...T) []T` — remove specific values, preserves remaining duplicates
- Full support for custom comparable types (`type UserID string`), struct types with comparable fields, and pointer types
- Documented NaN semantics for floating-point slices (follows Go's map-key rules)
- 100% test coverage including adversarial edge cases: nil/empty handling, immutability, result-aliasing checks, NaN behavior, struct key correctness
- Zero external dependencies — pure stdlib

### Deliberate non-goals
- No operations on maps (use stdlib `maps`)
- No sorting (use stdlib `slices.Sort`)
- No predicate-based operations on a single slice (those live in `bold-minds/each`)
- No NaN-aware float operations — caller must pre-process slices if NaN matching is required

### Requires
- Go 1.21 or later
