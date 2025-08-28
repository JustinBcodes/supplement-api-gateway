package middleware

import (
	"context"
	"net/http"
	"time"
)

type ctxKey int

const requestIDKey ctxKey = 1

func WithRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = generateID()
		}
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), requestIDKey, id)))
	})
}

func RequestIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(requestIDKey).(string); ok {
		return v
	}
	return ""
}

func generateID() string {
	// simple timestamp-based ID to start; will replace with faster rand
	return time.Now().UTC().Format("20060102T150405.000000000Z07:00")
}


