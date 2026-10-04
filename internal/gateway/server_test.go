package gateway

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"supplx-gateway-marketplace/internal/config"
)

func TestRoutesToConfiguredUpstream(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/products/42" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte("product-42"))
	}))
	defer upstream.Close()
	server := NewServer(&config.GatewayConfig{Routes: []config.Route{{Match: "/v1/products/**", Upstreams: []string{upstream.URL}}}})
	handler := server.Handler()
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/v1/products/42", nil))
	if recorder.Code != http.StatusAccepted {
		t.Fatalf("status = %d", recorder.Code)
	}
	if body, _ := io.ReadAll(recorder.Body); string(body) != "product-42" {
		t.Fatalf("unexpected body: %s", body)
	}
	missing := httptest.NewRecorder()
	handler.ServeHTTP(missing, httptest.NewRequest(http.MethodGet, "/not-configured", nil))
	if missing.Code != http.StatusNotFound {
		t.Fatalf("unmatched route status = %d", missing.Code)
	}
}

func BenchmarkGatewayProxy(b *testing.B) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))
	defer upstream.Close()
	server := NewServer(&config.GatewayConfig{Routes: []config.Route{{Match: "/v1/products/**", Upstreams: []string{upstream.URL}}}})
	handler := server.Handler()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/v1/products/42", nil))
		if recorder.Code != http.StatusOK {
			b.Fatalf("status = %d", recorder.Code)
		}
	}
}
