package ratelimit

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

const luaTokenBucket = `
-- KEYS[1]=bucket key
-- ARGV[1]=capacity, [2]=refill_per_sec, [3]=now_ms, [4]=cost, [5]=burst
-- Stored: tokens, last_ms
local key = KEYS[1]
local capacity = tonumber(ARGV[1])
local refill_per_sec = tonumber(ARGV[2])
local now_ms = tonumber(ARGV[3])
local cost = tonumber(ARGV[4])
local burst = tonumber(ARGV[5])

local data = redis.call('HMGET', key, 'tokens', 'last_ms')
local tokens = tonumber(data[1])
local last_ms = tonumber(data[2])

if tokens == nil then
  tokens = capacity + burst
  last_ms = now_ms
else
  local delta = (now_ms - last_ms) / 1000.0
  tokens = math.min(capacity + burst, tokens + delta * refill_per_sec)
  last_ms = now_ms
end

local allowed = 0
if tokens >= cost then
  tokens = tokens - cost
  allowed = 1
else
  allowed = 0
end

redis.call('HMSET', key, 'tokens', tokens, 'last_ms', last_ms)
redis.call('PEXPIRE', key, math.max(1000, math.floor((capacity + burst) / math.max(1, refill_per_sec) * 1000)))

if allowed == 1 then
  return {1, math.floor(tokens)}
else
  local need = cost - tokens
  local retry = math.ceil(need / math.max(1, refill_per_sec))
  return {0, retry}
end
`

type TokenBucketConfig struct {
	Capacity     int
	RefillPerSec int
	Burst        int
}

// Allow executes the token bucket script. Returns (allowed, remainingOrRetry, error)
func Allow(ctx context.Context, rdb *redis.Client, key string, cfg TokenBucketConfig, cost int, nowMs int64) (bool, int64, error) {
	if rdb == nil {
		return true, int64(cfg.Capacity), nil
	}
	res, err := rdb.Eval(ctx, luaTokenBucket, []string{key}, cfg.Capacity, cfg.RefillPerSec, nowMs, cost, cfg.Burst).Result()
	if err != nil {
		return false, 0, err
	}
	arr, ok := res.([]interface{})
	if !ok || len(arr) != 2 {
		return false, 0, errors.New("unexpected lua result")
	}
	allowed := arr[0].(int64) == 1
	value, _ := arr[1].(int64)
	return allowed, value, nil
}


