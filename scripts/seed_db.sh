#!/usr/bin/env bash
set -euo pipefail

echo "Seeding products..."
docker compose exec -T postgres-products psql -U products -d products -c "CREATE EXTENSION IF NOT EXISTS pgcrypto;"
ids=$(jq -c '.[]' seed/products_seed.json)
for row in $ids; do
  name=$(echo "$row" | jq -r .name)
  category=$(echo "$row" | jq -r .category)
  price=$(echo "$row" | jq -r .price_cents)
  inv=$(echo "$row" | jq -r .inventory)
  docker compose exec -T postgres-products psql -U products -d products -c \
    "INSERT INTO products(name, category, price_cents, inventory) VALUES ('$name','$category',$price,$inv) ON CONFLICT DO NOTHING;"
done
echo "Done."


