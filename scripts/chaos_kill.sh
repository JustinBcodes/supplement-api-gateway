#!/usr/bin/env bash
set -euo pipefail

svc=${1:-svc-payments-1}
echo "Killing $svc..."
docker compose kill $svc || true
sleep 2
echo "Restarting $svc..."
docker compose up -d $svc


