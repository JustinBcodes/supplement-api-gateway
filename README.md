Supplements Marketplace + API Gateway (supplx-gateway-marketplace)

What/Why

This repository contains a minimal microservices-based supplements marketplace protected by a high-throughput Go API Gateway. The gateway provides routing, Redis-backed token-bucket rate limiting, circuit breakers, least-connections load balancing, Prometheus metrics, OpenTelemetry tracing, and hot-reload configuration.

Quickstart

1. make up
2. make seed
3. make bench
4. make chaos

Hardware used for reference runs and full results will be added after first end-to-end implementation and benchmarking.

Architecture (high level)

Client → api-gateway (Go)
  ├─ svc-users (Go, Postgres)
  ├─ svc-products (Go, Postgres)
  ├─ svc-orders (Go, Postgres, talks to payments)
  └─ svc-payments (Go, mock PSP)

Limitations

- No TLS in the initial bench setup
- Single Redis instance
- Single AZ environment

Future Work

- EWMA load balancing, response caching, mTLS, canary weights


