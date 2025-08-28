#!/usr/bin/env bash
set -euo pipefail

if ! command -v wrk >/dev/null 2>&1; then
  echo "wrk is required" >&2
  exit 1
fi

BASE=${BASE:-http://localhost:8080}
DUR=60s
WARM=15s

echo "Baseline: GET /v1/products"
wrk -t8 -c256 -d${DUR} -L -s <(cat <<'LUA'
done = function(summary, latency, requests)
  -- placeholder output
end
LUA
) ${BASE}/v1/products > /dev/null


