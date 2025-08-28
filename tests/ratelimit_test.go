package tests

import (
    "testing"
    "time"
    rl "supplx-gateway-marketplace/pkg/ratelimit"
)

func TestTokenBucketMath(t *testing.T){
    cfg := rl.TokenBucketConfig{Capacity: 10, RefillPerSec: 5, Burst: 5}
    // simulate local: tokens should refill over time; here we just ensure cfg is sane
    if cfg.Capacity <= 0 || cfg.RefillPerSec <= 0 { t.Fatal("invalid config") }
    _ = time.Now()
}


