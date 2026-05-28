# ClickHouse Analytics Store

ClickHouse is the analytics and telemetry store for Headlamp.

SurrealDB remains the operational graph, artifact, relationship, memory, incident, and live-event store.

ClickHouse stores high-volume append-only operational data for fast analytics across hosts, clusters, agents, pods, logs, traces, metrics, incidents, and auto-healing loops.

## Role split

```txt
SurrealDB
  = operational graph + live events + artifacts + relationships + current state

ClickHouse
  = high-volume analytics + telemetry + historical queries + time-series aggregation
```

Do not replace SurrealDB with ClickHouse.

Use both:

- SurrealDB for graph-native operational state
- ClickHouse for observability-scale analytics

## Why ClickHouse

Headlamp will produce high-cardinality, high-volume events:

- agent heartbeats
- host metrics
- pod metrics
- container restarts
- Kubernetes events
- logs
- OPA decisions
- auto-healing attempts
- remediation outcomes
- SSH fallback events
- audit trails
- map/reduce findings
- cluster scans

SurrealDB is excellent for connected operational objects. ClickHouse is better for analytical scans, rollups, dashboards, and historical trend queries.

## Data flow

```txt
Node Agent
  → emits structured events
  → writes current state / graph to SurrealDB
  → writes telemetry stream to ClickHouse

Cortex
  → reads SurrealDB live streams
  → writes health decisions and incident state to SurrealDB
  → writes analytical event history to ClickHouse

Headlamp Desktop
  → reads live state from SurrealDB
  → reads dashboards and history from ClickHouse
```

## Storage responsibilities

| Data | SurrealDB | ClickHouse |
|---|---:|---:|
| Current agent state | yes | optional snapshot |
| Artifact graph | yes | no |
| Relationships | yes | no |
| Live notifications | yes | no |
| Incident current state | yes | historical copy |
| Heartbeat stream | recent/current | full history |
| Metrics | no/minimal | yes |
| Logs | no/minimal | yes |
| Audit history | important records | full append-only stream |
| OPA decisions | current/linked | full history |
| MapReduce findings | current linked findings | analytical stream |
| Auto-healing loops | current workflow | historical analytics |

## Local Podman service

ClickHouse should run as an optional local service in Podman:

```txt
headlamp-clickhouse
  http: 127.0.0.1:8123
  native: 127.0.0.1:9000
```

It should be optional for MVP, but strongly recommended for fleet analytics.

## Suggested database

```sql
CREATE DATABASE IF NOT EXISTS headlamp;
```

## Core tables

### agent_heartbeats

```sql
CREATE TABLE IF NOT EXISTS headlamp.agent_heartbeats
(
  timestamp DateTime64(3),
  agent_id String,
  host_id String,
  status LowCardinality(String),
  version String,
  capabilities Array(String),
  health_json String
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(timestamp)
ORDER BY (agent_id, timestamp);
```

### host_metrics

```sql
CREATE TABLE IF NOT EXISTS headlamp.host_metrics
(
  timestamp DateTime64(3),
  host_id String,
  cpu_usage Float64,
  memory_used_bytes UInt64,
  memory_total_bytes UInt64,
  disk_used_bytes UInt64,
  disk_total_bytes UInt64,
  load1 Float64,
  load5 Float64,
  load15 Float64
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(timestamp)
ORDER BY (host_id, timestamp);
```

### kubernetes_events

```sql
CREATE TABLE IF NOT EXISTS headlamp.kubernetes_events
(
  timestamp DateTime64(3),
  cluster_id String,
  namespace String,
  kind LowCardinality(String),
  name String,
  reason String,
  event_type LowCardinality(String),
  message String,
  source String
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(timestamp)
ORDER BY (cluster_id, namespace, kind, name, timestamp);
```

### opa_decisions

```sql
CREATE TABLE IF NOT EXISTS headlamp.opa_decisions
(
  timestamp DateTime64(3),
  decision_id String,
  policy_bundle String,
  subject_id String,
  action String,
  result LowCardinality(String),
  reasons Array(String),
  input_json String
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(timestamp)
ORDER BY (subject_id, action, timestamp);
```

### remediation_runs

```sql
CREATE TABLE IF NOT EXISTS headlamp.remediation_runs
(
  timestamp DateTime64(3),
  run_id String,
  cluster_id String,
  host_id String,
  playbook_id String,
  status LowCardinality(String),
  action_count UInt32,
  approved Bool,
  auto_approved Bool,
  verification_status LowCardinality(String),
  duration_ms UInt64
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(timestamp)
ORDER BY (cluster_id, host_id, timestamp);
```

### audit_events

```sql
CREATE TABLE IF NOT EXISTS headlamp.audit_events
(
  timestamp DateTime64(3),
  audit_id String,
  actor String,
  channel LowCardinality(String),
  subject_id String,
  action String,
  decision LowCardinality(String),
  result LowCardinality(String),
  metadata_json String
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(timestamp)
ORDER BY (subject_id, actor, timestamp);
```

## Query examples

Agent downtime over time:

```sql
SELECT
  agent_id,
  countIf(status != 'online') AS unhealthy_count,
  max(timestamp) AS last_seen
FROM headlamp.agent_heartbeats
GROUP BY agent_id
ORDER BY unhealthy_count DESC;
```

Most common Kubernetes warning reasons:

```sql
SELECT
  reason,
  count() AS count
FROM headlamp.kubernetes_events
WHERE event_type = 'Warning'
GROUP BY reason
ORDER BY count DESC
LIMIT 20;
```

Auto-healing success rate:

```sql
SELECT
  playbook_id,
  count() AS runs,
  countIf(verification_status = 'passed') AS successful,
  successful / runs AS success_rate
FROM headlamp.remediation_runs
GROUP BY playbook_id
ORDER BY runs DESC;
```

## Go package boundary

```txt
pkg/analytics/clickhouse/
  client.go
  schema.go
  writer.go
  queries.go
```

The runtime should treat ClickHouse as an append-only analytics sink.

If ClickHouse is unavailable, Headlamp should continue operating with SurrealDB and buffer/retry telemetry where possible.

## Failure behavior

ClickHouse failure must not stop operations.

```txt
SurrealDB unavailable
  → operational state degraded

ClickHouse unavailable
  → analytics degraded, operations continue
```

## Product rule

Use ClickHouse for history, metrics, and analytical questions.

Use SurrealDB for current state, graph relationships, live events, and workflow coordination.
