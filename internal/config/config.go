package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type RetryPolicy struct {
	Max      int      `yaml:"max"`
	BackoffMs int      `yaml:"backoff_ms"`
	JitterMs int       `yaml:"jitter_ms"`
	Methods  []string `yaml:"methods"`
}

type RateLimitPolicy struct {
	Key           string `yaml:"key"`
	Capacity      int    `yaml:"capacity"`
	RefillPerSec  int    `yaml:"refill_per_sec"`
	Burst         int    `yaml:"burst"`
}

type CircuitBreakerPolicy struct {
	WindowSeconds int     `yaml:"window_s"`
	ErrorRate     float64 `yaml:"error_rate"`
	MinRequests   int     `yaml:"min_requests"`
	OpenMs        int     `yaml:"open_ms"`
	Probe         int     `yaml:"probe"`
}

type AuthPolicy struct {
	Type    string `yaml:"type"`
	JWKSURL string `yaml:"jwks_url"`
}

type Policies struct {
	Auth            *AuthPolicy            `yaml:"auth"`
	RateLimit       *RateLimitPolicy       `yaml:"rate_limit"`
	CircuitBreaker  *CircuitBreakerPolicy  `yaml:"circuit_breaker"`
	Retries         *RetryPolicy           `yaml:"retries"`
}

type Route struct {
	Match     string   `yaml:"match"`
	Upstreams []string `yaml:"upstreams"`
	Policies  Policies `yaml:"policies"`
}

type GatewayConfig struct {
	Routes []Route `yaml:"routes"`
}

func Load(path string) (*GatewayConfig, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var cfg GatewayConfig
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return &cfg, nil
}


