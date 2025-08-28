package gateway

import (
	"sync"
	"time"

	"supplx-gateway-marketplace/internal/lb"
    "supplx-gateway-marketplace/pkg/breaker"
)

type routeState struct {
	mu       sync.Mutex
	lc       *lb.LeastConnections
	lastUp   string
    breakers map[string]*breaker.Breaker
}

func newRouteState(upstreams []string) *routeState {
	rs := &routeState{lc: lb.NewLeastConnections(upstreams), breakers: make(map[string]*breaker.Breaker)}
	for _, u := range upstreams {
		rs.breakers[u] = breaker.New(50, 0.5, 30*time.Second, 5, 20)
	}
	return rs
}

func (r *routeState) acquire() string {
	r.mu.Lock()
	up := r.lc.Acquire()
	r.lastUp = up
	r.mu.Unlock()
	return up
}

func (r *routeState) release(up string) {
	r.mu.Lock()
	r.lc.Release(up)
	r.mu.Unlock()
}

func (r *routeState) getBreaker(up string) *breaker.Breaker {
	r.mu.Lock()
	defer r.mu.Unlock()
	if b, ok := r.breakers[up]; ok { return b }
	b := breaker.New(50, 0.5, 30*time.Second, 5, 20)
	r.breakers[up] = b
	return b
}


