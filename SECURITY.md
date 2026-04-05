# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 0.x.x   | :white_check_mark: |

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security vulnerability, please follow these steps:

### 1. **Do Not** Create a Public Issue

Please do not report security vulnerabilities through public GitHub issues, discussions, or pull requests.

### 2. Report Privately

Send an email to **security@boldminds.tech** with:

- **Subject**: Security Vulnerability in bold-minds/list
- **Description**: Detailed description of the vulnerability
- **Steps to Reproduce**: Clear steps to reproduce the issue
- **Impact**: Potential impact and severity assessment
- **Suggested Fix**: If you have ideas for a fix (optional)

### 3. Response Timeline

- **Initial Response**: Within 48 hours
- **Status Update**: Within 7 days
- **Resolution**: Varies based on complexity, typically within 30 days

## Security Considerations

`list` is a pure-computation library with a very small attack surface:

- **No network I/O.** `list` does not make network calls.
- **No file I/O.** `list` does not read or write files.
- **No reflection.** All operations use Go's `comparable` type constraint and concrete map lookups.
- **No external dependencies.** Pure Go stdlib.
- **Immutable.** `list` never modifies input slices.
- **Nil-safe.** All functions handle nil inputs without panicking.

### Known runtime-panic sources from caller misuse

`list` does not panic on any documented input. However, there are two
situations where caller mistakes can cause Go runtime panics that propagate
through `list`:

1. **Non-comparable interface values.** If a caller passes a slice of `any`
   containing dynamic values whose concrete type is not comparable (e.g., a
   `[]int` stored inside an `any`), comparing those values causes a Go
   runtime panic in the `map[any]struct{}` used for deduplication. `list`
   does not recover these panics. Callers must ensure that interface-typed
   slices contain only comparable dynamic values.

2. **Extremely large inputs.** Map allocation failure (out of memory) on
   multi-billion-element slices would produce a Go runtime panic. `list`
   does not preallocate based on caller-controlled size hints, so this is
   only a concern for genuinely huge inputs.

### NaN handling

Floating-point NaN values follow Go's map-key semantics (NaN is not equal
to NaN). This means `list` cannot deduplicate, intersect, or match NaN
values — each NaN is treated as a distinct element. This is documented
behavior, not a security issue, but callers processing untrusted float
data should be aware.

## Security Updates

Security updates will be released as patch versions (e.g., 0.1.1),
documented in CHANGELOG.md, and announced through GitHub releases.

## Acknowledgments

We appreciate responsible disclosure and will acknowledge security
researchers who help improve the security of this project.
