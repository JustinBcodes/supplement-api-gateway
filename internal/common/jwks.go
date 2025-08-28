package common

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type JWKS struct {
	Keys []map[string]any `json:"keys"`
}

type JWKSCache struct {
	mu     sync.RWMutex
	set    map[string]JWKS
	client *http.Client
}

func NewJWKSCache() *JWKSCache {
	return &JWKSCache{set: make(map[string]JWKS), client: &http.Client{Timeout: 5 * time.Second}}
}

func (c *JWKSCache) Get(url string) (JWKS, error) {
	c.mu.RLock()
	if v, ok := c.set[url]; ok {
		c.mu.RUnlock()
		return v, nil
	}
	c.mu.RUnlock()

	resp, err := c.client.Get(url)
	if err != nil {
		return JWKS{}, err
	}
	defer resp.Body.Close()
	var out JWKS
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return JWKS{}, err
	}
	c.mu.Lock()
	c.set[url] = out
	c.mu.Unlock()
	return out, nil
}


