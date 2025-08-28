package middleware

import (
	"net/http"
	"strconv"
	"time"

	"supplx-gateway-marketplace/internal/obs"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

func Metrics(routeLabel string, upstreamGetter func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w, status: 200}
			obs.RequestsInFlight.Inc()
			next.ServeHTTP(rec, r)
			obs.RequestsInFlight.Dec()
			up := upstreamGetter()
			dur := time.Since(start).Seconds()
			obs.RequestDuration.WithLabelValues(routeLabel, up).Observe(dur)
			statusClass := strconv.Itoa(rec.status/100) + "xx"
			obs.RequestsTotal.WithLabelValues(routeLabel, up, statusClass).Inc()
		})
	}
}


