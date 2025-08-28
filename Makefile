.PHONY: up down seed bench chaos test pprof tidy

COMPOSE = docker compose -f docker-compose.yaml

up:
	$(COMPOSE) up -d

down:
	$(COMPOSE) down -v

seed:
	bash scripts/seed_db.sh

bench:
	bash scripts/bench_wrk.sh

chaos:
	bash scripts/chaos_kill.sh

test:
	go test ./...

pprof:
	@echo "Enable pprof via env var or flags as documented in README"

tidy:
	go mod tidy


