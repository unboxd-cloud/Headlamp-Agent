# Agent Health Monitoring

The Headlamp Node Agent must not be responsible for monitoring only itself.

A healthy design uses multiple layers of monitoring so that a failed agent can still be detected and recovered.

## Monitoring layers

```txt
Layer 1: systemd
  monitors the local node-agent process and restarts it if it exits

Layer 2: Headlamp Desktop
  checks /healthz and heartbeat age when the user is connected

Layer 3: Peer / Controller Watcher
  checks registered agents and marks stale agents unhealthy

Layer 4: SSH Fallback
  repairs the node-agent when API health checks fail

Layer 5: External Monitoring
  optional Prometheus/Uptime Kuma/healthcheck service
```

## Why not self-monitor only

An agent cannot reliably detect or report its own complete failure.

If the process is dead, network is broken, config is invalid, or the agent is wedged, it cannot send a useful status.

Therefore health must be checked from outside the agent.

## Responsibilities

| Monitor | Responsibility |
|---|---|
| systemd | process restart on crash |
| Headlamp Desktop | user-visible health and recovery controls |
| Controller Watcher | fleet heartbeat tracking and stale detection |
| SSH Fallback | break-glass recovery |
| External monitoring | independent alerting |

## Heartbeat model

Each node agent should emit a heartbeat artifact.

```json
{
  "kind": "AgentHeartbeatArtifact",
  "id": "heartbeat:agent:vps-prod-01:2026-05-28T23:30:00Z",
  "type": "agent-heartbeat",
  "agentId": "agent:vps-prod-01",
  "hostId": "host:vps-prod-01",
  "status": "online",
  "observedAt": "2026-05-28T23:30:00Z",
  "version": "0.1.0",
  "capabilities": [
    "observe.host",
    "observe.podman",
    "observe.kubernetes",
    "action.restart_service"
  ],
  "health": {
    "process": "ok",
    "podman": "ok",
    "kubernetes": "unknown",
    "surrealdb": "ok",
    "opa": "ok"
  }
}
```

## Health states

```txt
online
stale
degraded
unreachable
recovering
disabled
unknown
```

## Stale detection

A controller or Headlamp Desktop should calculate health from heartbeat age.

Example:

```txt
last heartbeat < 30s       → online
last heartbeat 30s-120s    → degraded/stale
last heartbeat > 120s      → unreachable
agent API unavailable      → unreachable
SSH reachable              → recoverable
SSH unreachable            → host unreachable or network failure
```

## Recovery flow

```txt
Heartbeat stale
  → call /healthz
  → if /healthz fails, try SSH fallback
  → check systemd service
  → read logs
  → restart node agent
  → verify /healthz
  → record recovery audit event
```

## systemd role

systemd is the first local monitor.

The node agent service should use:

```ini
Restart=always
RestartSec=5
WatchdogSec=30
```

The agent can notify systemd when alive in a later implementation.

## Controller watcher

The watcher can run in one of three places:

1. Headlamp Desktop while open
2. a small controller service on a trusted VPS
3. a Kubernetes CronJob/Deployment if managing clusters from inside Kubernetes

The watcher should not need broad mutation permissions. It only needs to read heartbeats and trigger recovery workflows when policy allows.

## External monitoring

Optional external systems can monitor:

- `GET /healthz`
- TCP port reachability
- systemd unit status through SSH
- heartbeat freshness in SurrealDB

Suggested tools:

- Uptime Kuma
- Prometheus Blackbox Exporter
- Grafana Agent/Alloy
- simple cron-based healthcheck

## MVP answer

For the first implementation:

```txt
systemd monitors and restarts the local node-agent process.
Headlamp Desktop checks /healthz and heartbeat freshness.
SSH fallback repairs the agent if /healthz fails.
```

Later, add a dedicated `headlamp-watcher` service for fleet monitoring.
