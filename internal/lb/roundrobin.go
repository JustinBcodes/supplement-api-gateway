package lb

import "sync/atomic"

type RoundRobin struct {
	endpoints []string
	idx       atomic.Uint64
}

func NewRoundRobin(endpoints []string) *RoundRobin {
	return &RoundRobin{endpoints: endpoints}
}

func (r *RoundRobin) Next() string {
	if len(r.endpoints) == 0 {
		return ""
	}
	i := r.idx.Add(1)
	return r.endpoints[int(i-1)%len(r.endpoints)]
}


