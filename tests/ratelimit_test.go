package tests

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	rl "supplx-gateway-marketplace/pkg/ratelimit"
)

func TestTokenBucketBlocksAndRefills(t *testing.T) {
	server := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: server.Addr()})
	defer client.Close()
	cfg := rl.TokenBucketConfig{Capacity: 2, RefillPerSec: 1, Burst: 0}
	ctx := context.Background()
	for i := 0; i < 2; i++ {
		allowed, _, err := rl.Allow(ctx, client, "client-a", cfg, 1, 1000)
		if err != nil || !allowed {
			t.Fatalf("request %d should be allowed: %v", i, err)
		}
	}
	allowed, retry, err := rl.Allow(ctx, client, "client-a", cfg, 1, 1000)
	if err != nil || allowed || retry < 1 {
		t.Fatalf("expected blocked request with retry, got allowed=%v retry=%d err=%v", allowed, retry, err)
	}
	allowed, _, err = rl.Allow(ctx, client, "client-a", cfg, 1, 2000)
	if err != nil || !allowed {
		t.Fatalf("expected a token after refill: %v", err)
	}
	allowed, _, err = rl.Allow(ctx, client, "client-b", cfg, 1, 2000)
	if err != nil || !allowed {
		t.Fatalf("client buckets must be independent: %v", err)
	}
}
