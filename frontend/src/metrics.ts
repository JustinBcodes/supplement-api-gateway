import type { Point, Snapshot } from "./data";

type Vector = { value: [number, string] };
type Matrix = { values: [number, string][] };

async function prometheus<T>(
  path: string,
  params: Record<string, string>,
): Promise<T[]> {
  const url = `/api/prometheus/api/v1/${path}?${new URLSearchParams(params)}`;
  const response = await fetch(url, { signal: AbortSignal.timeout(5000) });
  if (!response.ok) throw new Error(`Prometheus returned ${response.status}`);
  const data = await response.json();
  if (data.status !== "success") throw new Error("Prometheus query failed");
  return data.data.result as T[];
}

async function scalar(query: string): Promise<number> {
  const result = await prometheus<Vector>("query", { query });
  return Number(result[0]?.value?.[1] || 0);
}

export async function loadSnapshot(minutes: number): Promise<Snapshot> {
  const end = Math.floor(Date.now() / 1000);
  const start = end - minutes * 60;
  const [
    rps,
    p95,
    inFlight,
    blocked,
    allowed,
    totalRequests,
    circuitOpen,
    healthyUpstreams,
    trafficMatrix,
    errorMatrix,
  ] = await Promise.all([
    scalar("sum(rate(gateway_requests_total[5m]))"),
    scalar(
      "histogram_quantile(0.95, sum(rate(gateway_request_duration_seconds_bucket[5m])) by (le))",
    ),
    scalar("sum(gateway_requests_in_flight)"),
    scalar("gateway_ratelimit_blocked_total"),
    scalar("gateway_ratelimit_allowed_total"),
    scalar("sum(gateway_requests_total)"),
    scalar("count(gateway_cb_state == 1)"),
    scalar("count(gateway_upstream_healthy == 1)"),
    prometheus<Matrix>("query_range", {
      query: "sum(rate(gateway_requests_total[5m]))",
      start: String(start),
      end: String(end),
      step: String(Math.max(15, Math.floor((minutes * 60) / 23))),
    }),
    prometheus<Matrix>("query_range", {
      query: 'sum(rate(gateway_requests_total{status_class="5xx"}[5m]))',
      start: String(start),
      end: String(end),
      step: String(Math.max(15, Math.floor((minutes * 60) / 23))),
    }),
  ]);
  const errorByTime = new Map(
    (errorMatrix[0]?.values || []).map(([time, value]) => [
      time,
      Number(value),
    ]),
  );
  const series: Point[] = (trafficMatrix[0]?.values || []).map(
    ([time, value]) => ({
      time: new Date(time * 1000).toLocaleTimeString([], {
        hour: "2-digit",
        minute: "2-digit",
      }),
      traffic: Number(value),
      errors: errorByTime.get(time) || 0,
    }),
  );
  return {
    rps,
    p95: p95 * 1000,
    inFlight,
    blocked,
    allowed,
    totalRequests,
    circuitOpen,
    healthyUpstreams,
    series,
  };
}
