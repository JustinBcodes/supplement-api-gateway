package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type chargeReq struct {
	AmountCents int `json:"amount_cents"`
}

type chargeResp struct {
	Status string `json:"status"`
}

func main() {
	addr := ":8080"
	if v := os.Getenv("PAYMENTS_ADDR"); v != "" {
		addr = v
	}
	latMs := 8
	if v := os.Getenv("LATENCY_MS"); v != "" { if n, err := strconv.Atoi(v); err == nil { latMs = n } }
	failRate := 0.2
	if v := os.Getenv("FAIL_RATE"); v != "" { if f, err := strconv.ParseFloat(v, 64); err == nil { failRate = f } }

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/v1/charge", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(latMs) * time.Millisecond)
		if rng.Float64() < failRate {
			w.WriteHeader(http.StatusBadGateway)
			_ = json.NewEncoder(w).Encode(chargeResp{Status: "failed"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(chargeResp{Status: "ok"})
	})

	log.Printf("svc-payments listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}


