# Cortex Watchdog

Cortex is the watchdog layer for Headlamp.

Headlamp is the desktop dashboard/control surface.

The Kubernetes Orchestrator Agent performs cluster and host operations.

The Headlamp Node Agent runs on each VPS/host.

Cortex watches agents, watchers, hosts, clusters, heartbeats, incidents, and recovery workflows.

## Role split

```txt
Headlamp Desktop
  = operator dashboard and control surface

Kubernetes Orchestrator Agent
  = reasoning/planning/remediation capability

Headlamp Node Agent
  = local host/VPS execution and inventory agent

Cortex
  = watchdog, health observer, incident detector, recovery coordinator

SurrealDB
  = event stream, artifact graph, heartbeat store, incident store

SSH
  = break-glass recovery channel
```

## Why Cortex

A node agent cannot be the only monitor for itself.

Cortex provides an external view:

```txt
Node Agent emits heartbeat
  → SurrealDB event stream
  → Cortex evaluates freshness and health
  → Cortex creates notification/incident events
  → Headlamp Desktop shows alert and recovery action
```

If the node agent is dead, Cortex can detect missing heartbeats.

If Cortex is dead, systemd/external monitors should detect Cortex failure.

## Cortex responsibilities

Cortex should:

- monitor node-agent heartbeats
- monitor watcher heartbeats, if watchers exist
- monitor cluster health summaries
- monitor SurrealDB connectivity
- monitor OPA decision failures
- monitor auto-healing verification failures
- detect stale/unreachable agents
- create notification events
- create incident events
- trigger recovery workflows when policy allows
- escalate to SSH fallback when agent API fails
- update incident state after recovery

Cortex should not become a broad mutation agent.

It coordinates recovery, but actual actions still pass through policy and approval gates.

## Cortex deployment modes

### Local desktop mode

Cortex runs inside or next to Headlamp Desktop.

Best for:

- single machine
- local lab
- small VPS set
- development

Limitation:

- alerts work only while Headlamp Desktop is running.

### VPS controller mode

Cortex runs as a small service on a trusted VPS.

Best for:

- always-on monitoring
- multiple hosts
- K3s fleet
- production environments

### Kubernetes mode

Cortex runs as a Deployment inside a management cluster.

Best for:

- Kubernetes-native fleet management
- multiple managed clusters
- cluster operations teams

## Recommended MVP

Start with Cortex as a Go service:

```txt
cmd/headlamp-cortex
```

It should:

1. connect to SurrealDB
2. read/subscribe to heartbeat events
3. mark stale agents
4. create NotificationEventArtifact records
5. create IncidentArtifact records for critical failures
6. expose `/healthz`
7. expose `/status`

## Cortex event flow

```txt
AgentHeartbeatEvent
  → Cortex health evaluator
  → HealthChangeEvent
  → NotificationEvent
  → IncidentEvent when critical
  → Headlamp Desktop live subscription
```

## Heartbeat freshness rules

Default rules:

```txt
< 30s       online
30s-120s    stale/degraded
> 120s      unreachable
API failed  unreachable
SSH works   recoverable
SSH failed  host unreachable
```

## Cortex and OPA

Cortex must use OPA for recovery decisions.

```txt
stale agent detected
  → Cortex prepares recovery plan
  → OPA evaluates recovery plan
  → allow / require approval / deny / escalate
  → Headlamp notifies operator
```

## Cortex and SSH

Cortex may request SSH fallback, but should not silently execute broad SSH commands.

Allowed initial SSH recovery actions:

- check `headlamp-node-agent` systemd status
- fetch recent logs
- restart `headlamp-node-agent`
- verify `/healthz`

## Cortex failure handling

Cortex is also monitored.

```txt
systemd watches Cortex service
external monitor watches Cortex /healthz
Headlamp Desktop watches Cortex when open
```

If Cortex is dead but node agents are alive, Headlamp can still talk directly to node agents.

If Cortex and node agents are dead, Headlamp uses SSH fallback.

## Systemd service

Cortex should run with:

```ini
Restart=always
RestartSec=5
```

Later, add systemd watchdog support:

```ini
Type=notify
WatchdogSec=30
```

## SurrealDB tables

Suggested tables:

```sql
agent_heartbeats
watcher_heartbeats
health_events
notification_events
incident_events
audit_events
recovery_workflows
```

Headlamp Desktop subscribes to live queries over these tables.

## Product sentence

Cortex is the always-on watchdog and recovery coordinator for Headlamp-managed hosts, agents, Kubernetes clusters, and artifact operations.
