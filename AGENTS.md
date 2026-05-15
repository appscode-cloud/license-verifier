# AGENTS.md

This file provides guidance to coding agents (e.g. Claude Code, claude.ai/code) when working with code in this repository.

## Repository purpose

Go module `go.bytebuilders.dev/license-verifier` — a **library** that every paid AppsCode operator imports to enforce its license at runtime. Imported via blank-import (`_ "go.bytebuilders.dev/license-verifier/info"`) so a missing/invalid license refuses to start the binary. No standalone binary of its own.

## Architecture

- `lib.go` — top-level verifier entry points (`Verify`, etc.).
- `apis/` — API types describing license content.
- `client/` — generated client used to talk to the license server.
- `info/` — the package consumed via blank-import; carries license metadata baked into the binary at build time.
- `kubernetes/` — Kubernetes-aware helpers (look up a license Secret, refresh it, etc.).
- `bin/` — pre-built helper binaries used by CI.
- `hack/`, `Makefile` — codegen + lint harness.

## Common commands

- `make ci` — full CI pipeline.
- `make gen` — regenerate clients / generated files after API type changes.
- `make fmt`, `make lint`, `make unit-tests` / `make test` — standard.
- `make verify` — codegen + module-tidy verification.
- `make add-license` / `make check-license` — manage license headers.

Run a single Go test:

```
go test . -run TestName -v
```

## Conventions

- Module path is `go.bytebuilders.dev/license-verifier` (vanity URL); imports must use that.
- **Every exported symbol is API.** Downstream operators pin this dep; breaking changes ripple across the whole ACE / KubeDB / KubeVault / Stash ecosystem.
- The `info/` package is consumed via blank-import. Don't move or rename its public symbols.
- The Kubernetes-aware helpers in `kubernetes/` expect a Secret containing a JWT-style license — don't change that contract without coordinating with the binaries that mount it.
- License: `LICENSE`. Sign off commits (`git commit -s`).
