CREATE TABLE IF NOT EXISTS products (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  category TEXT NOT NULL,
  price_cents INT NOT NULL,
  inventory INT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);


