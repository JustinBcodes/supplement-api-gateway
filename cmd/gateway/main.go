package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"supplx-gateway-marketplace/internal/config"
	"supplx-gateway-marketplace/internal/gateway"
	"supplx-gateway-marketplace/internal/middleware"
    "supplx-gateway-marketplace/internal/common"

	"github.com/redis/go-redis/v9"
)

func main() {
	addr := ":8080"
	if v := os.Getenv("GATEWAY_ADDR"); v != "" {
		addr = v
	}
	cfgPath := os.Getenv("GATEWAY_CONFIG")
	if cfgPath == "" {
		cfgPath = "configs/routes.example.yaml"
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	gw := gateway.NewServer(cfg)

	// Redis client for rate limiting
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "redis:6379"
	}
	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})
    gw.SetRateLimiter(middleware.NewRateLimiter(rdb))
    gw.SetJWKSCache(common.NewJWKSCache())

	srv := &http.Server{
		Addr:              addr,
		Handler:           gw.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("gateway listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Hot reload config
	go func() {
		_ = config.Watch(cfgPath, func(newCfg *config.GatewayConfig) {
			log.Printf("config reloaded: %d routes", len(newCfg.Routes))
			gw.UpdateConfig(newCfg)
		})
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}


