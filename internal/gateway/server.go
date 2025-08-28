package gateway

import (
	"errors"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"supplx-gateway-marketplace/internal/common"
	"supplx-gateway-marketplace/internal/config"
	"supplx-gateway-marketplace/internal/middleware"
	"supplx-gateway-marketplace/internal/obs"
	"supplx-gateway-marketplace/internal/proxy"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	cfg   atomic.Value // *config.GatewayConfig
	mu    sync.RWMutex
	state map[string]*routeState
	rl    *middleware.RateLimiter
	jwks  *common.JWKSCache
}

func NewServer(cfg *config.GatewayConfig) *Server {
	s := &Server{state: make(map[string]*routeState)}
	s.cfg.Store(cfg)
	s.rebuild(cfg)
	return s
}

func (s *Server) UpdateConfig(cfg *config.GatewayConfig) {
	s.cfg.Store(cfg)
	s.rebuild(cfg)
}

func (s *Server) SetRateLimiter(rl *middleware.RateLimiter) { s.rl = rl }
func (s *Server) SetJWKSCache(j *common.JWKSCache) { s.jwks = j }

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/", s.route)
	// Chain: request id -> mux (per-route middleware applied dynamically)
	obs.MustRegister()
	return middleware.WithRequestID(mux)
}

func (s *Server) route(w http.ResponseWriter, r *http.Request) {
	cfg := s.cfg.Load().(*config.GatewayConfig)
	path := r.URL.Path
	for _, rt := range cfg.Routes {
		if matchPath(rt.Match, path) {
			rs := s.getState(rt)
			up := rs.acquire()
			defer rs.release(up)
			if up == "" {
				w.WriteHeader(http.StatusBadGateway)
				_, _ = w.Write([]byte("no upstreams"))
				return
			}
			// Circuit breaker check
			br := rs.getBreaker(up)
			obs.CBState.WithLabelValues(up).Set(float64(br.State()))
			// In open state, reject immediately
			if !br.Allow() {
				obs.CBState.WithLabelValues(up).Set(float64(br.State()))
				w.WriteHeader(http.StatusServiceUnavailable)
				_, _ = w.Write([]byte("circuit open"))
				return
			}

			base := http.HandlerFunc(func(rw http.ResponseWriter, rq *http.Request) {
				_ = proxy.SimpleReverseProxy(up, rw, rq)
			})
			h := middleware.Metrics(rt.Match, func() string { return up })(base)
			// per-route auth if configured
			if ap := rt.Policies.Auth; ap != nil && s.jwks != nil {
				a := &middleware.AuthMiddleware{JWKS: s.jwks, Cfg: middleware.AuthConfig{Issuer: "", Audience: "", JWKSURL: ap.JWKSURL}}
				h = a.Wrap(h)
			}
			// per-route rate limit if configured
			if rl := rt.Policies.RateLimit; rl != nil && s.rl != nil {
				pol := middleware.RLPolicy{Capacity: rl.Capacity, RefillPerSec: rl.RefillPerSec, Burst: rl.Burst, KeyType: rl.Key, RouteMatch: rt.Match}
				h = s.rl.WithPolicy(pol, h)
			}
			// retries for idempotent methods
			if rp := rt.Policies.Retries; rp != nil && (r.Method == http.MethodGet || r.Method == http.MethodHead) {
				if err := retryServe(h, w, r, rp.Max, rp.BackoffMs, rp.JitterMs); err != nil {
					br.OnResult(false)
					obs.CBState.WithLabelValues(up).Set(float64(br.State()))
					log.Printf("proxy error: %v", err)
					w.WriteHeader(http.StatusBadGateway)
					_, _ = w.Write([]byte("upstream error"))
				} else {
					br.OnResult(true)
					obs.CBState.WithLabelValues(up).Set(float64(br.State()))
				}
				return
			}
			if err := serveOnce(h, w, r); err != nil {
				br.OnResult(false)
				obs.CBState.WithLabelValues(up).Set(float64(br.State()))
				log.Printf("proxy error: %v", err)
				w.WriteHeader(http.StatusBadGateway)
				_, _ = w.Write([]byte("upstream error"))
			} else {
				br.OnResult(true)
				obs.CBState.WithLabelValues(up).Set(float64(br.State()))
			}
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte("route not found"))
}

func serveOnce(h http.Handler, w http.ResponseWriter, r *http.Request) error {
	rr := &errRecorder{ResponseWriter: w}
	h.ServeHTTP(rr, r)
	return rr.err
}

type errRecorder struct{
	http.ResponseWriter
	err error
}

func (e *errRecorder) WriteHeader(status int){
	if status >= 500 && e.err == nil {
		e.err = errors.New("upstream error")
	}
	e.ResponseWriter.WriteHeader(status)
}

func (s *Server) rebuild(cfg *config.GatewayConfig){
	s.mu.Lock()
	defer s.mu.Unlock()
	m := make(map[string]*routeState, len(cfg.Routes))
	for _, rt := range cfg.Routes {
		m[rt.Match] = newRouteState(rt.Upstreams)
	}
	s.state = m
}

func (s *Server) getState(rt config.Route) *routeState {
	s.mu.RLock()
	rs, ok := s.state[rt.Match]
	s.mu.RUnlock()
	if ok { return rs }
	s.mu.Lock()
	defer s.mu.Unlock()
	if rs, ok = s.state[rt.Match]; ok { return rs }
	rs = newRouteState(rt.Upstreams)
	s.state[rt.Match] = rs
	return rs
}

func retryServe(h http.Handler, w http.ResponseWriter, r *http.Request, max int, backoffMs int, jitterMs int) error {
	for i := 0; i <= max; i++ {
		err := serveOnce(h, w, r)
		if err == nil { return nil }
		if i == max { return err }
		sleep := time.Duration(backoffMs) * time.Millisecond
		if jitterMs > 0 { sleep += time.Duration(rand.Intn(jitterMs)) * time.Millisecond }
		time.Sleep(sleep)
	}
	return nil
}

func matchPath(pattern, path string) bool {
	// very small matcher: suffix "**" means prefix match
	if strings.HasSuffix(pattern, "**") {
		prefix := strings.TrimSuffix(pattern, "**")
		return strings.HasPrefix(path, strings.TrimRight(prefix, "/"))
	}
	return path == strings.TrimRight(pattern, "/")
}


