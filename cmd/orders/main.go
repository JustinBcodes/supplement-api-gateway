package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type checkoutReq struct {
	UserID string `json:"user_id"`
	TotalCents int `json:"total_cents"`
}

type checkoutResp struct {
	Status string `json:"status"`
}

func main() {
	addr := ":8080"
	if v := os.Getenv("ORDERS_ADDR"); v != "" {
		addr = v
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/v1/orders/checkout", func(w http.ResponseWriter, r *http.Request) {
		var req checkoutReq
		_ = json.NewDecoder(r.Body).Decode(&req)
		body, _ := json.Marshal(map[string]any{"amount_cents": req.TotalCents})
		resp, err := http.Post("http://svc-payments-1:8080/v1/charge", "application/json", bytes.NewReader(body))
		if err != nil || resp.StatusCode >= 500 {
			w.WriteHeader(http.StatusBadGateway)
			_ = json.NewEncoder(w).Encode(checkoutResp{Status: "failed"})
			return
		}
		_ = json.NewEncoder(w).Encode(checkoutResp{Status: "ok"})
	})

	log.Printf("svc-orders listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}


