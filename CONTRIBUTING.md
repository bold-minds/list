# Contributing to `list`

Thanks for your interest in contributing. This guide covers the operational process. For the **why** — the design principles every contribution is tested against — see **[bold-minds/oss/PRINCIPLES.md](https://github.com/bold-minds/oss/blob/main/PRINCIPLES.md)**.

## 🎯 Before You Start

Every contribution is measured against the four Bold Minds principles: **outcome naming**, **one way to do each thing**, **get out of the way**, and **non-goals explicit**. If your proposed change doesn't honor these, it will not be merged.

**Read [PRINCIPLES.md](https://github.com/bold-minds/oss/blob/main/PRINCIPLES.md) first.** It's the load-bearing document.

## 🔧 Development Setup

**Requirements:** Go 1.21 or later, Git, Bash.

```bash
git clone https://github.com/bold-minds/list.git
cd list
go test ./...              # unit tests
go test -race ./...        # race detection
go test -bench=. ./...     # benchmarks
./scripts/validate.sh      # full validation pipeline (local mode)
./scripts/validate.sh ci   # strict CI mode
```

Your contribution must pass `./scripts/validate.sh ci` before submitting.

## 📁 Project Structure

```
list/
├── list.go                # Implementation (single file)
├── list_test.go           # Unit tests
├── bench_test.go          # Benchmarks
├── examples/              # Runnable examples
├── scripts/
│   └── validate.sh        # Validation pipeline
├── README.md
├── CONTRIBUTING.md        # This file
├── CHANGELOG.md
├── CODE_OF_CONDUCT.md
├── SECURITY.md
├── LICENSE
└── go.mod
```

Flat layout. No `internal/` directory.

## 🎨 Code Style

### Naming
- Outcome naming per PRINCIPLES.md. Function names describe the set operation performed (`Unique`, `Union`, `Intersect`, `Minus`, `Without`).

### Error Handling
- Functions **must not panic** on valid input (nil, empty, or otherwise).
- No error returns — set operations either succeed or return empty slices.
- Never `Must*` variants.

### Documentation
- Every exported function has a doc comment.
- Edge cases (nil, empty, NaN, aliasing, immutability) documented in the package doc and README.

### Dependencies
- **Zero external dependencies.** `list` is pure stdlib.

## 🧪 Testing

**Coverage target: 100% of exported functions.**

```bash
go test -v ./...
go test -race ./...
go test -cover ./...
go test -bench=. -benchmem ./...
```

**Every PR must include adversarial tests.** In addition to happy-path coverage, tests must verify:

1. **Non-nil empty returns** — no function returns nil for empty or missing results.
2. **Immutability** — the input slice is byte-identical before and after the call.
3. **Result aliasing** — mutating the returned slice must not affect the input.
4. **NaN semantics** (for float-aware features) — behavior is consistent with Go's map-key rules.
5. **Custom comparable types** — named types (`type UserID string`) and struct keys work correctly.

## 📝 Pull Request Process

### PR Checklist

- [ ] **Outcome naming** — does the function name describe what the caller gets?
- [ ] **One way** — does any existing function already do this?
- [ ] **Get out of the way** — can a Go dev use this from the signature alone?
- [ ] **Non-goals** — does this violate any of the library's stated non-goals?
- [ ] Tests cover 100% of new code
- [ ] Adversarial tests included (nil, immutability, aliasing, NaN where applicable)
- [ ] Benchmarks added for new exported functions
- [ ] README updated
- [ ] CHANGELOG.md updated
- [ ] `./scripts/validate.sh ci` passes locally

## 🆕 Adding a New Function

`list` is deliberately tiny (five functions). New additions must clear a high bar:

1. Read the library's non-goals in [README.md](README.md) and [PRINCIPLES.md](https://github.com/bold-minds/oss/blob/main/PRINCIPLES.md).
2. Prove the stdlib gap. Current Go's `slices` package is more capable than many people realize.
3. Show real-world evidence of the pain.
4. If the function operates on maps or uses a predicate, it belongs in a different library (`each` for predicates, `stdlib maps` for maps).

## 🏷️ Versioning and Releases

- Semantic versioning
- v0.x: API may change between minor versions
- v1.0+: breaking changes require major version bump
- Every release updates CHANGELOG.md

## 🙏 Code of Conduct

See [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).

## 📄 License

By contributing, you agree your contributions are licensed under the MIT License.
