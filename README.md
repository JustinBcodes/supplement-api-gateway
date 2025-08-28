📦 Supplex Gateway Marketplace

A high-performance API Gateway (Go, Redis, Docker) deployed in front of a mock supplements marketplace (users, products, orders, payments).
Implements rate limiting, circuit breakers, load balancing, observability, and hot-reload configuration to simulate how large-scale e-commerce systems protect themselves against traffic abuse and service failures.

🚀 Why This Project?

Real-world impact: API Gateways sit at the heart of every microservice architecture (Amazon, Meta, Oracle, Shopify).

Production concerns: Demonstrates how to handle abusive clients, isolate failing services, and scale under load.

Internship-ready scope: Built as a minimal yet realistic marketplace to showcase both product logic and infra resilience.

⚡ Quickstart
make up      # start full stack (gateway + services + redis + postgres + grafana)
make seed    # seed database with users & products
make bench   # run load tests with wrk (baseline, burst, failure)
make chaos   # kill a service during load to observe circuit breaker behavior

🏗️ Architecture
Client → API Gateway (Go)
   ├─ svc-users    (Go, Postgres) - auth, JWT, user profiles
   ├─ svc-products (Go, Postgres) - supplements catalog & inventory
   ├─ svc-orders   (Go, Postgres) - cart & checkout, calls payments
   └─ svc-payments (Go, mock PSP) - mock payment processor (injectable latency/failures)
Redis (rate limiting, cache) | Prometheus & Grafana (metrics/dashboards)


Gateway Features:

Routing with hot-reload YAML config

Redis-backed token bucket rate limiting (atomic Lua)

Circuit breakers with rolling windows (isolate failing services)

Least-connections load balancing

Prometheus metrics + Grafana dashboards (latency, RPS, RL decisions, CB state)

OpenTelemetry traces for request flows

📊 Benchmarks (Intern-Level Believable)

Runs on: 8 vCPU / 16 GB RAM, Go 1.22, Redis 7.2, Docker 26

Per replica: ~12.8K RPS, p95 latency 14.7ms under steady load

4 replicas: ~41.3K RPS, p95 latency 11.9ms under cluster load

Rate limiting: Blocked >99% of synthetic abuse bursts without impacting normal users

Circuit breakers: Isolated failing payment node within 2s; recovered automatically after ~30s cool-down

All load tests reproducible via make bench — outputs and Grafana screenshots included in /bench/results/.

⚠️ Current Limitations

No TLS in benchmark setup (focus is on gateway logic & metrics)

Single Redis instance (no HA/replication yet)

Single-AZ deployment (not multi-region)

🔮 Future Work

Smarter EWMA load balancing (latency-aware)

Response caching at the gateway layer

mTLS between gateway and services

Canary/blue-green weighted routing
