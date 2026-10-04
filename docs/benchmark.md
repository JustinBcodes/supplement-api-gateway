# Gateway loopback benchmark

Measured on 2026-10-03 with Go 1.22.0 on Windows/amd64 and an AMD Ryzen 9 9900X3D.

```bash
go test ./internal/gateway -run '^$' -bench BenchmarkGatewayProxy -benchmem -benchtime=2s -count=3
```

| Run | Iterations | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: | ---: |
| 1 | 40,030 | 67,608 | 47,766 | 109 |
| 2 | 34,232 | 69,417 | 48,026 | 110 |
| 3 | 33,049 | 85,641 | 48,366 | 110 |

Median: **69,417 ns/op** (69.4 µs/op). The benchmark uses an `httptest` upstream on the same machine and measures one proxied request through the gateway handler. It excludes Redis, JWT verification, real network conditions, and multi-replica deployment. Use `scripts/bench_wrk.sh` only after starting the stack to measure the full system; results will vary by environment.
