package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	keyMu sync.RWMutex
	priv  *rsa.PrivateKey
)

func ensureKey() {
	keyMu.Lock()
	defer keyMu.Unlock()
	if priv != nil { return }
	pp, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil { panic(err) }
	priv = pp
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func main() {
	addr := ":8080"
	if v := os.Getenv("USERS_ADDR"); v != "" {
		addr = v
	}
	ensureKey()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("/v1/auth/login", func(w http.ResponseWriter, r *http.Request) {
		var req loginReq
		_ = json.NewDecoder(r.Body).Decode(&req)
		// accept any credentials for demo
		claims := jwt.MapClaims{
			"iss": "svc-users",
			"aud": "supplx",
			"sub": "user-1",
			"scope": "user",
			"exp": time.Now().Add(1 * time.Hour).Unix(),
		}
		tok := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		tok.Header["kid"] = "demo-key"
		s, _ := tok.SignedString(priv)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"token": s})
	})

	mux.HandleFunc("/.well-known/jwks.json", func(w http.ResponseWriter, r *http.Request) {
		keyMu.RLock()
		pub := &priv.PublicKey
		n := base64.RawURLEncoding.EncodeToString(pub.N.Bytes())
		// exponent typically 65537 -> 0x10001 -> AQAB
		e := base64.RawURLEncoding.EncodeToString([]byte{0x01, 0x00, 0x01})
		keyMu.RUnlock()
		_ = json.NewEncoder(w).Encode(map[string]any{"keys": []map[string]any{{"kty": "RSA", "kid": "demo-key", "n": n, "e": e}}})
	})

	log.Printf("svc-users listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}


