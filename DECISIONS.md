# Key Decisions

- Token bucket vs sliding window: token bucket chosen for burst tolerance and simplicity.
- Redis + Lua scripting for atomic rate limiting across replicas; avoids per-process counters.
- Least-connections as the initial load balancer; EWMA considered future work.
- Idempotent retries limited to safe methods (GET/HEAD) to avoid retry storms.
- Orders write path guarded by per-user (JWT sub) quotas at the gateway.


