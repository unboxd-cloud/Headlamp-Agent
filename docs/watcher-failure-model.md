# Watcher Failure Model

A health checker can also fail.

Headlamp must not rely on a single watcher, single node agent, or single control channel.

## Core principle

```txt
No component should be the only proof of its own health.
```

If the node agent dies, the watcher detects it.

If the watcher dies, systemd and external monitoring detect it.

If both agent and watcher are dead, SSH remains the break-glass channel.

If SSH is unavailable, the host is considered unreachable and must be handled through provider console, hypervisor, or physical access.

## Failure ladder

```txt
1. Node Agent healthy
   → use node-agent API

2. Node Agent unhealthy, Watcher healthy
   → watcher detects stale heartbeat
   → Headlamp offers recovery
   → SSH fallback can restart node agent

3. Watcher unhealthy, Node Agent healthy
   → Headlamp Desktop detects watcher stale
   → node-agent API still works
   → restart watcher

4. Node Agent and Watcher both unhealthy
   → Headlamp uses SSH fallback
   → check systemd status for both services
   → restart services
   → verify health endpoints

5. SSH unavailable too
   → host/network/provider-level failure
   → mark host unreachable
   → escalate to cloud/VPS console or physical access
```

## Watcher deployment options

### Local single-host mode

For one VPS:

```txt
systemd
  ├─ headlamp-node-agent
  └─ headlamp-watcher
```

systemd restarts both processes if they crash.

### Multi-host mode

For multiple VPS servers:

```txt
headlamp-watcher-a watches agents b/c/d
headlamp-watcher-b watches agents a/c/d
headlamp-watcher-c watches agents a/b/d
```

This avoids one central watcher becoming a single point of failure.

### Desktop mode

When the Headlamp desktop app is open, it can also check:

- watcher heartbeat
- node agent heartbeat
- direct `/healthz`
- SSH reachability

### External monitor mode

Optional tools can monitor both node-agent and watcher:

- Uptime Kuma
- Prometheus Blackbox Exporter
- Grafana Alloy
- cron healthcheck
- provider monitoring

## Watcher heartbeat artifact

Watchers must emit their own heartbeat.

```json
{
  "kind": "WatcherHeartbeatArtifact",
  "id": "watcher-heartbeat:watcher:vps-prod-01:2026-05-28T23:45:00Z",
  "type": "watcher-heartbeat",
  "watcherId": "watcher:vps-prod-01",
  "hostId": "host:vps-prod-01",
  "status": "online",
  "observedAt": "2026-05-28T23:45:00Z",
  "watching": [
    "agent:vps-prod-01",
    "agent:vps-prod-02"
  ]
}
```

## Watcher responsibilities

The watcher should be smaller and simpler than the node agent.

It should only:

- check node-agent `/healthz`
- check heartbeat freshness
- optionally check SSH reachability
- record watcher heartbeat
- trigger recovery workflow when allowed
- alert Headlamp Desktop or external monitor

It should not perform broad cluster mutations.

## systemd role

systemd is the local process supervisor.

Both services should use:

```ini
Restart=always
RestartSec=5
```

For watchdog integration later:

```ini
WatchdogSec=30
Type=notify
```

## If systemd is unhealthy

If systemd cannot keep services alive, the issue is at the host level.

Headlamp should escalate to SSH diagnostics:

```bash
systemctl status headlamp-node-agent --no-pager
systemctl status headlamp-watcher --no-pager
journalctl -u headlamp-node-agent -n 200 --no-pager
journalctl -u headlamp-watcher -n 200 --no-pager
systemctl reset-failed headlamp-node-agent
systemctl reset-failed headlamp-watcher
systemctl restart headlamp-node-agent headlamp-watcher
```

## Minimal robust setup

For the first reliable version:

```txt
On every VPS:
  - headlamp-node-agent under systemd
  - headlamp-watcher under systemd
  - SSH fallback configured

On desktop:
  - checks both agent and watcher
  - can recover both through SSH
```

## Final escalation model

```txt
Agent API
  ↓ if failed
Watcher
  ↓ if failed
systemd
  ↓ if failed
SSH
  ↓ if failed
provider console / hypervisor / physical access
```

This is the correct operational model: each layer is useful, but no layer is trusted as the only recovery path.
