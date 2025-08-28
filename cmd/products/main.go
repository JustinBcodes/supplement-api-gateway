package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

type Product struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Category   string `json:"category"`
	PriceCents int    `json:"price_cents"`
	Inventory  int    `json:"inventory"`
}

var products = []Product{
	{ID: "1", Name: "Whey Protein Isolate", Category: "protein", PriceCents: 2999, Inventory: 200},
	{ID: "2", Name: "Creatine Monohydrate", Category: "performance", PriceCents: 1499, Inventory: 300},
	{ID: "3", Name: "Fish Oil (Omega-3)", Category: "health", PriceCents: 1299, Inventory: 180},
	{ID: "4", Name: "Magnesium Glycinate", Category: "health", PriceCents: 1599, Inventory: 160},
}

func main() {
	addr := ":8080"
	if v := os.Getenv("PRODUCTS_ADDR"); v != "" {
		addr = v
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/v1/products", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(products)
	})
	mux.HandleFunc("/v1/products/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/v1/products/")
		for _, p := range products {
			if p.ID == id {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(p)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
	})

	log.Printf("svc-products listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}


