#!/usr/bin/env bash
set -euo pipefail

if ! command -v wrk >/dev/null 2>&1; then
  echo "Install wrk to run the full-stack load test" >&2
  exit 1
fi

BASE=${BASE:-http://localhost:8080}
TARGET=${TARGET:-/v1/products}
THREADS=${THREADS:-4}
CONNECTIONS=${CONNECTIONS:-64}
DURATION=${DURATION:-30s}
mkdir -p bench/results
output="bench/results/wrk-$(date -u +%Y%m%dT%H%M%SZ).txt"

echo "Target: ${BASE}${TARGET}"
echo "Threads: ${THREADS}; connections: ${CONNECTIONS}; duration: ${DURATION}"
echo "The product route has an IP-based rate limit. HTTP 429 responses are expected during load."
wrk -t"${THREADS}" -c"${CONNECTIONS}" -d"${DURATION}" --latency "${BASE}${TARGET}" | tee "$output"
echo "Saved raw output to $output"
