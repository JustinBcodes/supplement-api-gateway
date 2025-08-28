package obs

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "gateway_requests_total", Help: "Total requests"},
		[]string{"route", "upstream", "status_class"},
	)

	RequestsInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{Name: "gateway_requests_in_flight", Help: "In-flight requests"},
	)

	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{Name: "gateway_request_duration_seconds", Help: "Request latency", Buckets: prometheus.DefBuckets},
		[]string{"route", "upstream"},
	)

	RLAllowed = prometheus.NewCounter(
		prometheus.CounterOpts{Name: "gateway_ratelimit_allowed_total", Help: "Rate limit allowed"},
	)
	RLBlocked = prometheus.NewCounter(
		prometheus.CounterOpts{Name: "gateway_ratelimit_blocked_total", Help: "Rate limit blocked"},
	)

	CBState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{Name: "gateway_cb_state", Help: "Circuit breaker state (0 closed, 1 open, 2 half-open)"},
		[]string{"upstream"},
	)

	UpstreamHealthy = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{Name: "gateway_upstream_healthy", Help: "Upstream health (1 healthy, 0 unhealthy)"},
		[]string{"upstream"},
	)
)

func MustRegister() {
	prometheus.MustRegister(RequestsTotal)
	prometheus.MustRegister(RequestsInFlight)
	prometheus.MustRegister(RequestDuration)
	prometheus.MustRegister(RLAllowed)
	prometheus.MustRegister(RLBlocked)
	prometheus.MustRegister(CBState)
	prometheus.MustRegister(UpstreamHealthy)
}


