package breaker

import (
	"sync"
	"time"
)

type State int

const (
	Closed State = iota
	Open
	HalfOpen
)

type RollingWindow struct {
	buckets []bucket
	size    int
	mu      sync.Mutex
}

type bucket struct {
	success int
	failure int
}

func NewRollingWindow(size int) *RollingWindow {
	return &RollingWindow{buckets: make([]bucket, size), size: size}
}

func (w *RollingWindow) Add(success bool, idx int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	b := &w.buckets[idx%w.size]
	if success {
		b.success++
	} else {
		b.failure++
	}
}

func (w *RollingWindow) Snapshot() (success, failure int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, b := range w.buckets {
		success += b.success
		failure += b.failure
	}
	return
}

type Breaker struct {
	state       State
	openedAt    time.Time
	window      *RollingWindow
	minRequests int
	errorRate   float64
	openFor     time.Duration
	probeLimit  int
	probeCount  int
	mu          sync.Mutex
}

func New(minRequests int, errorRate float64, openFor time.Duration, probeLimit int, windowSize int) *Breaker {
	return &Breaker{
		state:       Closed,
		window:      NewRollingWindow(windowSize),
		minRequests: minRequests,
		errorRate:   errorRate,
		openFor:     openFor,
		probeLimit:  probeLimit,
	}
}

func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}

func (b *Breaker) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	switch b.state {
	case Closed:
		return true
	case Open:
		if time.Since(b.openedAt) >= b.openFor {
			b.state = HalfOpen
			b.probeCount = 0
			return true
		}
		return false
	case HalfOpen:
		if b.probeCount < b.probeLimit {
			b.probeCount++
			return true
		}
		return false
	}
	return false
}

func (b *Breaker) OnResult(success bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.window.Add(success, int(time.Now().Unix()))
	if b.state == HalfOpen {
		// close if high success ratio across probes
		if b.probeCount >= b.probeLimit {
			// naive: require >=80% success
			if successRatio(b.window) >= 0.8 {
				b.state = Closed
			} else {
				b.state = Open
				b.openedAt = time.Now()
			}
		}
		return
	}

	s, f := b.window.Snapshot()
	total := s + f
	if total >= b.minRequests && float64(f)/float64(total) >= b.errorRate {
		b.state = Open
		b.openedAt = time.Now()
	}
}

func successRatio(w *RollingWindow) float64 {
	s, f := w.Snapshot()
	t := s + f
	if t == 0 {
		return 1
	}
	return float64(s) / float64(t)
}


