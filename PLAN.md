# phenotype-go-kit — Implementation Plan

## Overview

Go infrastructure toolkit for the Phenotype ecosystem. Provides small, focused packages for logging, data structures, polling, registries, and cross-cutting concerns extracted from production Phenotype services.

## Phase 1 — Core Primitives (Complete)

| ID | Task | Status |
|----|------|--------|
| P1.1 | `logctx` — context-scoped slog.Logger injection | Done |
| P1.2 | `ringbuffer` — generic fixed-capacity circular buffer | Done |
| P1.3 | `waitfor` — polling with exponential backoff and testable clocks | Done |
| P1.4 | `registry` — generic thread-safe key-value registry with ref counting | Done |

## Phase 2 — Infrastructure Utilities (Complete)

| ID | Task | Status | Depends On |
|----|------|--------|------------|
| P2.1 | `cache` — distributed and in-process caching layer | Done | P1.1 |
| P2.2 | `circuit` — circuit breaker with configurable thresholds | Done | P1.1 |
| P2.3 | `retry` — retry with jitter and backoff strategies | Done | P1.3 |
| P2.4 | `ratelimit` — token bucket and sliding window rate limiting | Done | P1.1 |
| P2.5 | `metrics` — Prometheus-compatible metrics collection | Done | P1.1 |
| P2.6 | `tracing` — OpenTelemetry tracing helpers | Done | P1.1 |
| P2.7 | `logging` — structured log formatters and enrichers | Done | P1.1 |
| P2.8 | `health` — HTTP health check handlers (liveness, readiness) | Done | P1.1 |

## Phase 3 — Application Patterns (Complete)

| ID | Task | Status | Depends On |
|----|------|--------|------------|
| P3.1 | `auth` — authentication middleware and token validation | Done | P2.1 |
| P3.2 | `oauth2` — OAuth2 client and server flows | Done | P3.1 |
| P3.3 | `cors` — CORS middleware with configurable policies | Done | — |
| P3.4 | `config` — structured configuration loading with validation | Done | — |
| P3.5 | `secrets` — secrets management client (Vault, env-based) | Done | P3.4 |
| P3.6 | `db` — database connection pooling and migration helpers | Done | P3.4 |
| P3.7 | `storage` — blob storage abstraction (S3-compatible) | Done | P3.4 |
| P3.8 | `bus` — event bus with pub/sub and dead-letter support | Done | P2.5 |
| P3.9 | `webhook` — inbound webhook handling and signature verification | Done | P3.1 |

## Phase 4 — Platform Integration (In Progress)

| ID | Task | Status | Depends On |
|----|------|--------|------------|
| P4.1 | `discovery` — service discovery client (Consul, DNS) | In Progress | P3.4 |
| P4.2 | `plugins` — plugin host runtime (Extism/WASM) | In Progress | P3.4 |
| P4.3 | `domain` — domain model base types and value objects | In Progress | — |
| P4.4 | `contracts` — interface contracts and validation helpers | In Progress | P4.3 |
| P4.5 | `trigger` — event trigger and scheduler primitives | Planned | P3.8 |
| P4.6 | `chaos` — fault injection for resilience testing | Planned | P2.2, P2.3 |
| P4.7 | `validation` — fluent validation builder | In Progress | — |
| P4.8 | `transform` — data transformation pipeline primitives | Planned | P4.3 |

## Phase 5 — Hardening and Observability (Planned)

| ID | Task | Status | Depends On |
|----|------|--------|------------|
| P5.1 | 80%+ unit test coverage across all packages | Planned | All P1–P4 |
| P5.2 | Integration tests with real backends (Redis, Postgres) | Planned | P2.1, P3.6 |
| P5.3 | Benchmark suite for hot-path packages | Planned | P1–P2 |
| P5.4 | OpenAPI / protobuf contract exports | Planned | P4.4 |
| P5.5 | Publish to pkg.go.dev with full godoc | Planned | P5.1 |
| P5.6 | GitHub Actions CI (lint, vet, test, coverage) | In Progress | — |

## DAG Summary

```
P1.1-P1.4 (core) -> P2.x (infra) -> P3.x (app) -> P4.x (platform) -> P5.x (harden)
```

## References

- PRD: `PRD.md`
- Functional Requirements: `FUNCTIONAL_REQUIREMENTS.md`
- Architecture: `ARCHITECTURE.md`
- ADR: `ADR.md`
