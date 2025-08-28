package lb

import (
	"math"
	"sync"
)

type LeastConnections struct {
	mu        sync.RWMutex
	conns     map[string]int
	endpoints []string
}

func NewLeastConnections(endpoints []string) *LeastConnections {
	conns := make(map[string]int, len(endpoints))
	for _, e := range endpoints {
		conns[e] = 0
	}
	return &LeastConnections{conns: conns, endpoints: endpoints}
}

func (l *LeastConnections) Acquire() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	minConn := math.MaxInt
	var chosen string
	for _, e := range l.endpoints {
		c := l.conns[e]
		if c < minConn {
			minConn = c
			chosen = e
		}
	}
	if chosen != "" {
		l.conns[chosen]++
	}
	return chosen
}

func (l *LeastConnections) Release(endpoint string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, ok := l.conns[endpoint]; ok {
		l.conns[endpoint]--
		if l.conns[endpoint] < 0 {
			l.conns[endpoint] = 0
		}
	}
}


