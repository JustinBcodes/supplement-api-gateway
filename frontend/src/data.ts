export type Point = { time: string; traffic: number; errors: number };
export type Route = {
  path: string;
  service: string;
  instances: number;
  policy: string;
  limit: string;
};
export type Snapshot = {
  rps: number;
  p95: number;
  inFlight: number;
  blocked: number;
  allowed: number;
  totalRequests: number;
  circuitOpen: number;
  healthyUpstreams: number;
  series: Point[];
};

export const routes: Route[] = [
  {
    path: "/v1/products/**",
    service: "Products",
    instances: 2,
    policy: "IP token bucket",
    limit: "300 capacity · 100/s",
  },
  {
    path: "/v1/orders/**",
    service: "Orders",
    instances: 2,
    policy: "JWT + user token bucket",
    limit: "60 capacity · 1/s",
  },
];

export const demoSeries: Point[] = Array.from({ length: 24 }, (_, i) => ({
  time: `${String(9 + Math.floor(i / 6)).padStart(2, "0")}:${String((i % 6) * 10).padStart(2, "0")}`,
  traffic: Math.round(
    730 + Math.sin(i * 0.63) * 125 + Math.cos(i * 0.28) * 74 + i * 11,
  ),
  errors: Math.round(10 + Math.sin(i * 0.7 + 1) * 5 + (i === 17 ? 12 : 0)),
}));

export const demoSnapshot: Snapshot = {
  rps: 1048,
  p95: 18.4,
  inFlight: 37,
  blocked: 284,
  allowed: 96342,
  totalRequests: 96626,
  circuitOpen: 0,
  healthyUpstreams: 6,
  series: demoSeries,
};
