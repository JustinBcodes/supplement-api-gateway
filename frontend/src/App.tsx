import { useCallback, useEffect, useState } from "react";
import {
  Activity,
  ArrowRight,
  ChevronDown,
  ChevronRight,
  CircleAlert,
  CircleCheck,
  Clock3,
  Database,
  ExternalLink,
  Gauge,
  GitBranch,
  Globe2,
  LayoutDashboard,
  LockKeyhole,
  Network,
  Pause,
  Play,
  Radio,
  RefreshCw,
  Route as RouteIcon,
  Shield,
  ShieldCheck,
  SlidersHorizontal,
  TrendingUp,
  Zap,
} from "lucide-react";
import {
  Area,
  AreaChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import { demoSnapshot, routes, type Snapshot } from "./data";
import { loadSnapshot } from "./metrics";
import "./App.css";

type Tab = "overview" | "traffic" | "reliability" | "routes";
const tabs: { id: Tab; label: string; icon: typeof Activity }[] = [
  { id: "overview", label: "Overview", icon: LayoutDashboard },
  { id: "traffic", label: "Traffic analytics", icon: Activity },
  { id: "reliability", label: "Reliability", icon: ShieldCheck },
  { id: "routes", label: "Routes & policies", icon: RouteIcon },
];
const nf = (value: number, digits = 0) =>
  new Intl.NumberFormat("en-US", { maximumFractionDigits: digits }).format(
    value,
  );

function App() {
  const [tab, setTab] = useState<Tab>("overview");
  const [snapshot, setSnapshot] = useState<Snapshot>(demoSnapshot);
  const [source, setSource] = useState<"demo" | "live">("demo");
  const [range, setRange] = useState(60);
  const [loading, setLoading] = useState(false);
  const [lastUpdated, setLastUpdated] = useState<Date | null>(null);
  const [notice, setNotice] = useState("");
  const [paused, setPaused] = useState(false);
  const refresh = useCallback(async () => {
    setLoading(true);
    try {
      const data = await loadSnapshot(range);
      setSnapshot(data);
      setSource("live");
      setLastUpdated(new Date());
      setNotice("");
    } catch {
      setSnapshot(demoSnapshot);
      setSource("demo");
      setNotice("Prometheus is unavailable. Showing labeled preview data.");
    } finally {
      setLoading(false);
    }
  }, [range]);
  useEffect(() => {
    void refresh();
  }, [refresh]);
  useEffect(() => {
    if (paused) return;
    const id = setInterval(() => void refresh(), 15000);
    return () => clearInterval(id);
  }, [paused, refresh]);

  const pageTitle = tabs.find((t) => t.id === tab)?.label;
  return (
    <div className="shell">
      <aside className="sidebar">
        <div className="brand">
          <span className="brand-icon">
            <Zap size={19} fill="currentColor" />
          </span>
          <span>
            supplex<span className="brand-pulse">.</span>
          </span>
          <small>CONTROL</small>
        </div>
        <div className="workspace">
          <span className="workspace-icon">
            <Globe2 size={18} />
          </span>
          <span>
            <strong>Gateway Marketplace</strong>
            <small>Operations workspace</small>
          </span>
          <ChevronDown size={14} />
        </div>
        <div className="nav-title">WORKSPACE</div>
        <nav>
          {tabs.map((item) => (
            <button
              key={item.id}
              className={tab === item.id ? "active" : ""}
              onClick={() => setTab(item.id)}
            >
              <item.icon size={18} />
              {item.label}
              {item.id === "overview" && <span className="nav-dot" />}
            </button>
          ))}
        </nav>
        <div className="nav-title system-title">SYSTEM</div>
        <div className="system-list">
          <div>
            <span className="system-mark mint">
              <Database size={14} />
            </span>
            <span>Redis</span>
            <small>Rate limits</small>
          </div>
          <div>
            <span className="system-mark blue">
              <Activity size={14} />
            </span>
            <span>Prometheus</span>
            <small>Metrics</small>
          </div>
          <div>
            <span className="system-mark orange">
              <Network size={14} />
            </span>
            <span>Upstreams</span>
            <small>Gateway</small>
          </div>
        </div>
        <div className="side-bottom">
          <a href="http://localhost:3000" target="_blank" rel="noreferrer">
            <Gauge size={17} /> Grafana dashboard <ExternalLink size={13} />
          </a>
          <a
            href="https://github.com/JustinBcodes/supplement-api-gateway"
            target="_blank"
            rel="noreferrer"
          >
            <GitBranch size={17} /> View source <ExternalLink size={13} />
          </a>
          <div className="operator">
            <span>JB</span>
            <div>
              <strong>Operator console</strong>
              <small>Local environment</small>
            </div>
            <ChevronDown size={14} />
          </div>
        </div>
      </aside>
      <main className="main">
        <header className="topbar">
          <div className="breadcrumbs">
            <span>Supplex</span>
            <ChevronRight size={14} />
            <strong>{pageTitle}</strong>
          </div>
          <div className="top-actions">
            <span className={`connection ${source}`}>
              <span />
              {source === "live" ? "Live telemetry" : "Preview data"}
            </span>
            <button title="Refresh telemetry" onClick={() => void refresh()}>
              <RefreshCw className={loading ? "spin" : ""} size={17} />
            </button>
            <span className="bar-divider" />
            <div className="top-avatar">JB</div>
          </div>
        </header>
        <div className="page">
          <div className="page-head">
            <div>
              <div className="eyebrow">
                <span /> GATEWAY OPERATIONS{" "}
                <span className="eyebrow-sep">/</span>{" "}
                {source === "live" ? "LIVE" : "PREVIEW"}
              </div>
              <h1>
                {pageTitle}
                <span className="title-star">✳</span>
              </h1>
              <p>
                {tab === "overview"
                  ? "A clear view of every request, route, and safeguard."
                  : tab === "traffic"
                    ? "Understand request volume, latency, and throughput."
                    : tab === "reliability"
                      ? "Monitor the systems protecting your services."
                      : "Inspect route definitions and gateway policies."}
              </p>
            </div>
            <div className="head-actions">
              <select
                aria-label="Time range"
                value={range}
                onChange={(e) => setRange(Number(e.target.value))}
              >
                <option value={15}>Last 15 min</option>
                <option value={60}>Last 1 hour</option>
                <option value={360}>Last 6 hours</option>
              </select>
              <button onClick={() => setPaused(!paused)}>
                {paused ? <Play size={15} /> : <Pause size={15} />}{" "}
                {paused ? "Resume" : "Auto refresh"}
              </button>
            </div>
          </div>
          {notice && (
            <div className="notice">
              <CircleAlert size={15} />
              {notice}
              <button onClick={() => setNotice("")}>Dismiss</button>
            </div>
          )}
          {tab === "overview" && (
            <>
              <div className="hero">
                <div className="hero-copy">
                  <span>
                    <Radio size={14} /> TRAFFIC AT A GLANCE
                  </span>
                  <h2>
                    {nf(snapshot.rps, 1)}
                    <small>req / sec</small>
                  </h2>
                  <p>
                    {source === "live"
                      ? "Requests per second from Prometheus"
                      : "Illustrative traffic for interface preview"}
                  </p>
                  <div className="hero-trend">
                    <TrendingUp size={15} /> Gateway ready to route traffic{" "}
                    <ArrowRight size={15} />
                  </div>
                </div>
                <div className="hero-chart">
                  <TrafficChart data={snapshot.series} dark />
                </div>
              </div>
              <div className="stats">
                <Metric
                  icon={Clock3}
                  color="violet"
                  label="P95 latency"
                  value={`${nf(snapshot.p95, 1)} ms`}
                  detail="Request duration"
                />
                <Metric
                  icon={Activity}
                  color="blue"
                  label="In flight"
                  value={nf(snapshot.inFlight)}
                  detail="Active requests"
                />
                <Metric
                  icon={Shield}
                  color="peach"
                  label="Rate limited"
                  value={nf(snapshot.blocked)}
                  detail="Blocked, cumulative"
                />
                <Metric
                  icon={CircleCheck}
                  color="mint"
                  label="Healthy upstreams"
                  value={nf(snapshot.healthyUpstreams)}
                  detail="Reported targets"
                />
              </div>
              <div className="content-grid">
                <section className="card traffic-card">
                  <CardTitle
                    icon={Activity}
                    title="Request volume"
                    subtitle="Traffic and server errors over time"
                  />
                  <div className="legend">
                    <span>
                      <i className="legend-green" /> Requests
                    </span>
                    <span>
                      <i className="legend-red" /> 5xx errors
                    </span>
                  </div>
                  <div className="chart-wrap">
                    <TrafficChart data={snapshot.series} />
                  </div>
                </section>
                <section className="card health-card">
                  <CardTitle
                    icon={ShieldCheck}
                    title="System health"
                    subtitle="Protection layers and circuit state"
                  />
                  <HealthRow
                    icon={CircleCheck}
                    name="Circuit breakers"
                    value={
                      snapshot.circuitOpen
                        ? `${snapshot.circuitOpen} open`
                        : "All closed"
                    }
                    tone={snapshot.circuitOpen ? "warn" : "good"}
                  />
                  <HealthRow
                    icon={Shield}
                    name="Rate limiter"
                    value="Configured"
                    tone="good"
                  />
                  <HealthRow
                    icon={Network}
                    name="Load balancer"
                    value="Least connections"
                    tone="neutral"
                  />
                  <div className="health-footer">
                    <span className="green-dot" />{" "}
                    {source === "live"
                      ? `Updated ${lastUpdated?.toLocaleTimeString()}`
                      : "Sample status · connect stack for live data"}
                  </div>
                </section>
              </div>
              <RouteSection />
            </>
          )}
          {tab === "traffic" && (
            <>
              <div className="stats">
                <Metric
                  icon={Activity}
                  color="violet"
                  label="Throughput"
                  value={`${nf(snapshot.rps, 1)} /s`}
                  detail="5 minute rate"
                />
                <Metric
                  icon={Clock3}
                  color="blue"
                  label="P95 latency"
                  value={`${nf(snapshot.p95, 1)} ms`}
                  detail="5 minute window"
                />
                <Metric
                  icon={Globe2}
                  color="peach"
                  label="Total requests"
                  value={nf(snapshot.totalRequests)}
                  detail="Since start"
                />
                <Metric
                  icon={Activity}
                  color="mint"
                  label="In flight"
                  value={nf(snapshot.inFlight)}
                  detail="Active requests"
                />
              </div>
              <section className="card wide-chart">
                <CardTitle
                  icon={TrendingUp}
                  title="Traffic over time"
                  subtitle="Incoming requests and 5xx responses"
                />
                <div className="chart-wrap">
                  <TrafficChart data={snapshot.series} />
                </div>
              </section>
              <RouteSection />
            </>
          )}
          {tab === "reliability" && (
            <>
              <div className="stats">
                <Metric
                  icon={Shield}
                  color="violet"
                  label="Allowed requests"
                  value={nf(snapshot.allowed)}
                  detail="Token bucket decisions"
                />
                <Metric
                  icon={CircleAlert}
                  color="peach"
                  label="Blocked requests"
                  value={nf(snapshot.blocked)}
                  detail="Cumulative count"
                />
                <Metric
                  icon={ShieldCheck}
                  color="mint"
                  label="Open circuits"
                  value={nf(snapshot.circuitOpen)}
                  detail="Reported states"
                />
                <Metric
                  icon={Network}
                  color="blue"
                  label="Healthy upstreams"
                  value={nf(snapshot.healthyUpstreams)}
                  detail="Reported targets"
                />
              </div>
              <div className="content-grid">
                <section className="card protection-card">
                  <CardTitle
                    icon={Shield}
                    title="Protection policies"
                    subtitle="Per-route safeguards configured in the gateway"
                  />
                  {routes.map((route) => (
                    <div className="policy-row" key={route.path}>
                      <span className="policy-symbol">
                        <LockKeyhole size={17} />
                      </span>
                      <div>
                        <strong>{route.service}</strong>
                        <small>{route.policy}</small>
                      </div>
                      <span>{route.limit}</span>
                    </div>
                  ))}
                </section>
                <section className="card health-card">
                  <CardTitle
                    icon={Network}
                    title="Circuit states"
                    subtitle="0 closed · 1 open · 2 half-open"
                  />
                  <div className="circuit-visual">
                    <div className="circuit-ring">
                      <ShieldCheck size={30} />
                    </div>
                    <strong>
                      {snapshot.circuitOpen === 0
                        ? "All circuits closed"
                        : `${snapshot.circuitOpen} open circuits`}
                    </strong>
                    <p>
                      {source === "demo"
                        ? "Preview state. Start the stack for live readings."
                        : "Current gateway upstream state."}
                    </p>
                  </div>
                </section>
              </div>
            </>
          )}
          {tab === "routes" && (
            <>
              <div className="route-intro">
                <div>
                  <span>
                    <RouteIcon size={21} />
                  </span>
                  <h2>Routing map</h2>
                  <p>
                    Traffic policies configured in{" "}
                    <code>configs/routes.example.yaml</code>.
                  </p>
                </div>
                <strong>{routes.length} configured routes</strong>
              </div>
              <RouteSection detailed />
            </>
          )}
          <footer className="page-footer">
            <span>
              SUPPLEX GATEWAY <i /> OPERATOR CONSOLE
            </span>
            <span>
              {source === "live"
                ? `Telemetry refreshed ${lastUpdated?.toLocaleTimeString()}`
                : "Preview data · no live benchmark claims"}
            </span>
          </footer>
        </div>
      </main>
    </div>
  );
}

function Metric({
  icon: Icon,
  color,
  label,
  value,
  detail,
}: {
  icon: typeof Activity;
  color: string;
  label: string;
  value: string;
  detail: string;
}) {
  return (
    <div className="metric">
      <span className={`metric-icon ${color}`}>
        <Icon size={19} />
      </span>
      <small>{label}</small>
      <strong>{value}</strong>
      <p>{detail}</p>
    </div>
  );
}
function CardTitle({
  icon: Icon,
  title,
  subtitle,
}: {
  icon: typeof Activity;
  title: string;
  subtitle: string;
}) {
  return (
    <div className="card-title">
      <span>
        <Icon size={18} />
      </span>
      <div>
        <h3>{title}</h3>
        <p>{subtitle}</p>
      </div>
      <SlidersHorizontal size={16} />
    </div>
  );
}
function HealthRow({
  icon: Icon,
  name,
  value,
  tone,
}: {
  icon: typeof Activity;
  name: string;
  value: string;
  tone: string;
}) {
  return (
    <div className="health-row">
      <span>
        <Icon size={17} />
      </span>
      <strong>{name}</strong>
      <em className={tone}>{value}</em>
    </div>
  );
}
function TrafficChart({
  data,
  dark = false,
}: {
  data: Snapshot["series"];
  dark?: boolean;
}) {
  return (
    <ResponsiveContainer width="100%" height="100%">
      <AreaChart
        data={data}
        margin={{ top: 15, right: 4, bottom: 0, left: -24 }}
      >
        <defs>
          <linearGradient
            id={dark ? "heroFill" : "areaFill"}
            x1="0"
            y1="0"
            x2="0"
            y2="1"
          >
            <stop
              offset="0%"
              stopColor={dark ? "#8ce8bc" : "#5ccf9a"}
              stopOpacity={dark ? 0.42 : 0.25}
            />
            <stop
              offset="95%"
              stopColor={dark ? "#8ce8bc" : "#5ccf9a"}
              stopOpacity="0"
            />
          </linearGradient>
        </defs>
        <CartesianGrid
          vertical={false}
          stroke={dark ? "#ffffff18" : "#e9edf0"}
          strokeDasharray="4 4"
        />
        <XAxis
          dataKey="time"
          tick={{ fill: dark ? "#82978e" : "#9aa4ad", fontSize: 10 }}
          axisLine={false}
          tickLine={false}
          minTickGap={25}
        />
        {!dark && (
          <YAxis
            tick={{ fill: "#9aa4ad", fontSize: 10 }}
            axisLine={false}
            tickLine={false}
          />
        )}
        <Tooltip
          contentStyle={{
            borderRadius: 9,
            fontSize: 11,
            border: "1px solid #e5ebe8",
          }}
        />
        <Area
          type="monotone"
          dataKey="traffic"
          stroke={dark ? "#99edc1" : "#55c792"}
          strokeWidth={2.7}
          fill={`url(#${dark ? "heroFill" : "areaFill"})`}
        />
        {!dark && (
          <Area
            type="monotone"
            dataKey="errors"
            stroke="#f1a38d"
            strokeWidth={2}
            fill="transparent"
          />
        )}
      </AreaChart>
    </ResponsiveContainer>
  );
}
function RouteSection({ detailed = false }: { detailed?: boolean }) {
  return (
    <section className="card route-card">
      <CardTitle
        icon={RouteIcon}
        title="Active routes"
        subtitle="Gateway routing and protection policies"
      />
      <div className="route-table">
        <div className="route-header">
          <span>ROUTE</span>
          <span>SERVICE</span>
          <span>INSTANCES</span>
          <span>POLICY</span>
          <span>STATUS</span>
        </div>
        {routes.map((route) => (
          <div className="route-row" key={route.path}>
            <span className="route-path">
              <span className="path-icon">
                <RouteIcon size={16} />
              </span>
              <code>{route.path}</code>
            </span>
            <span>{route.service}</span>
            <span>
              <span className="instances">{route.instances} targets</span>
            </span>
            <span>
              {route.policy}
              {detailed && <small>{route.limit}</small>}
            </span>
            <span>
              <span className="healthy-dot" /> Configured
            </span>
          </div>
        ))}
      </div>
      {detailed && (
        <p className="route-note">
          Route definitions are read from the checked-in configuration. Live
          health metrics appear in Reliability when Prometheus is available.
        </p>
      )}
    </section>
  );
}

export default App;
