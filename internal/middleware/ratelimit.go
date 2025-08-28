package middleware

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"supplx-gateway-marketplace/internal/middleware/userctx"
	"supplx-gateway-marketplace/pkg/ratelimit"
)

type RateLimiter struct {
	rdb *redis.Client
}

func NewRateLimiter(rdb *redis.Client) *RateLimiter {
	return &RateLimiter{rdb: rdb}
}

type RLPolicy struct {
	Capacity     int
	RefillPerSec int
	Burst        int
	KeyType      string // apiKey | userSub | ip | route
	RouteMatch   string
}

func (m *RateLimiter) WithPolicy(pol RLPolicy, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := m.buildKey(pol, r)
		allowed, val, err := ratelimit.Allow(r.Context(), m.rdb, key, ratelimit.TokenBucketConfig{Capacity: pol.Capacity, RefillPerSec: pol.RefillPerSec, Burst: pol.Burst}, 1, time.Now().UnixMilli())
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		if !allowed {
			w.Header().Set("Retry-After", strconv.FormatInt(val, 10))
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte("rate limited"))
			return
		}
		w.Header().Set("X-RateLimit-Remaining", strconv.FormatInt(val, 10))
		next.ServeHTTP(w, r)
	})
}

func (m *RateLimiter) buildKey(pol RLPolicy, r *http.Request) string {
	parts := []string{"rl"}
	switch pol.KeyType {
	case "apiKey":
		ak := r.Header.Get("X-API-Key")
		if ak == "" {
			ak = "anon"
		}
		parts = append(parts, "api", ak)
	case "userSub":
		sub := userctx.Sub(r.Context())
		if sub == "" {
			sub = "anon"
		}
		parts = append(parts, "sub", sub)
	case "ip":
		ip := clientIP(r)
		parts = append(parts, "ip", ip)
	}
	// per route
	parts = append(parts, "route", pol.RouteMatch)
	return strings.Join(parts, ":")
}

func clientIP(r *http.Request) string {
	if xf := r.Header.Get("X-Forwarded-For"); xf != "" {
		parts := strings.Split(xf, ",")
		return strings.TrimSpace(parts[0])
	}
	h, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return h
}


