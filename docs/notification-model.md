# Notification Model

Headlamp must actively tell the operator when an agent, watcher, host, cluster, or recovery flow needs attention.

Notification cannot depend only on the failing agent. If the agent is dead, the notification must come from another layer such as Headlamp Desktop, a watcher, external monitor, SSH fallback, or provider-level monitoring.

## Notification sources

```txt
Node Agent
  → reports health when alive

Watcher
  → detects stale/dead agents

Headlamp Desktop
  → checks health while open

External Monitor
  → checks endpoints from outside

SSH Fallback
  → confirms whether host is reachable when API is dead
```

## Notification channels

Headlamp should support these channels:

| Channel | Purpose |
|---|---|
| Desktop notification | Immediate local alert when Headlamp is open/running |
| Dashboard badge | Persistent visual state inside Headlamp |
| Incident artifact | Durable record in SurrealDB |
| Audit event | Immutable operational history |
| Webhook | Slack, Discord, Teams, custom automation |
| Email | Operator notification for important events |
| External monitor | Uptime Kuma/Prometheus/Grafana alerting |
| SSH fallback prompt | Break-glass recovery when agent API is unreachable |

## Severity levels

```txt
info
warning
critical
emergency
```

Examples:

| Event | Severity |
|---|---|
| Agent heartbeat delayed | warning |
| Agent unreachable but SSH reachable | critical |
| Agent and SSH both unreachable | emergency |
| Auto-healing succeeded | info |
| Auto-healing failed verification | critical |
| Production mutation requires approval | warning |
| OPA denied dangerous action | warning |

## Notification event artifact

Notifications should be stored as artifacts.

```json
{
  "kind": "NotificationEventArtifact",
  "id": "notification:agent:vps-prod-01:unreachable:2026-05-28T23:55:00Z",
  "type": "notification-event",
  "severity": "critical",
  "source": "headlamp-desktop",
  "subject": {
    "kind": "AgentArtifact",
    "id": "agent:vps-prod-01"
  },
  "title": "Node agent unreachable",
  "message": "Headlamp could not reach /healthz for agent:vps-prod-01. SSH is reachable, recovery is available.",
  "status": "open",
  "createdAt": "2026-05-28T23:55:00Z",
  "recommendedActions": [
    "Run SSH recovery: check service status",
    "Fetch recent node-agent logs",
    "Restart headlamp-node-agent"
  ],
  "channels": [
    {
      "type": "desktop",
      "status": "sent"
    },
    {
      "type": "dashboard",
      "status": "visible"
    },
    {
      "type": "webhook",
      "status": "not_configured"
    }
  ]
}
```

## Incident artifact

For important failures, create an incident artifact.

```json
{
  "kind": "IncidentArtifact",
  "id": "incident:agent:vps-prod-01:2026-05-28T23:55:00Z",
  "type": "incident",
  "severity": "critical",
  "status": "open",
  "title": "Headlamp node agent unreachable on vps-prod-01",
  "affectedArtifacts": [
    "host:vps-prod-01",
    "agent:vps-prod-01"
  ],
  "detectedBy": "headlamp-desktop",
  "createdAt": "2026-05-28T23:55:00Z",
  "recoveryState": "ssh-reachable",
  "nextAction": "restart-node-agent-via-ssh"
}
```

## How the operator is notified

### When Headlamp Desktop is open

```txt
Desktop health loop detects stale agent
  → show red status badge
  → create NotificationEventArtifact
  → create IncidentArtifact if critical
  → show desktop notification
  → offer SSH recovery action
```

### When Headlamp Desktop is closed

Use one or more background channels:

```txt
headlamp-watcher
external monitor
webhook/email integration
provider alerting
```

For the MVP, Headlamp Desktop can notify only while open. For always-on alerts, run `headlamp-watcher` on a reliable host or configure Uptime Kuma/Prometheus.

## Escalation policy

```txt
warning after 30s stale
critical after 120s unreachable
emergency when API and SSH are both unreachable
```

## Dedupe and rate limiting

Headlamp must avoid notification storms.

Rules:

- group repeated alerts by subject and failure code
- update existing incident instead of creating duplicates
- rate-limit repeated notifications
- escalate severity if the condition remains unresolved
- close notification when health is restored

## Recovery notification flow

```txt
Agent unreachable
  → notify operator
  → operator approves SSH recovery
  → Headlamp runs recovery command
  → recovery succeeds or fails
  → Headlamp updates incident
  → Headlamp sends resolution notification
```

## MVP implementation

First implementation should support:

1. in-app dashboard badge
2. desktop notification while Headlamp is open
3. local NotificationEventArtifact in SurrealDB
4. local IncidentArtifact in SurrealDB
5. SSH recovery prompt

Later versions should add:

- webhook notifier
- email notifier
- Uptime Kuma integration
- Prometheus Alertmanager integration
- mobile push through an external service
