# SSH Fallback Control Model

Headlamp must not depend only on the node agent.

If the Headlamp Node Agent is dead, unreachable, misconfigured, or degraded, Headlamp should support SSH as a break-glass fallback control channel.

## Control channels

```txt
Primary channel:
  Headlamp Desktop → Headlamp Node Agent API

Fallback channel:
  Headlamp Desktop → SSH → VPS/host
```

The node agent is preferred because it can enforce local policy, expose structured APIs, and write local audit records.

SSH is required because the node agent itself can fail.

## When SSH fallback is allowed

SSH fallback should be used when:

- node agent `/healthz` is unreachable
- node agent service is stopped
- node agent config is broken
- VPS is reachable but agent is not
- emergency repair is needed
- user explicitly selects break-glass mode

## SSH fallback is not unrestricted automation

SSH fallback must still be governed.

Headlamp should not silently run arbitrary commands over SSH.

The fallback flow is:

```txt
Detect agent unavailable
  → offer SSH fallback
  → show exact command or command bundle
  → require user approval
  → execute over SSH
  → capture stdout/stderr/exit code
  → write local audit event
  → attempt to restore node agent
```

## First SSH use cases

Safe first SSH capabilities:

```txt
ssh.health.check
ssh.agent.status
ssh.agent.restart
ssh.agent.logs
ssh.host.inventory
ssh.systemd.status
ssh.podman.ps
ssh.k3s.status
ssh.kubernetes.readonly
```

Examples:

```bash
systemctl status headlamp-node-agent --no-pager
journalctl -u headlamp-node-agent -n 200 --no-pager
systemctl restart headlamp-node-agent
uname -a
podman ps --format json
kubectl get nodes -o json
```

## Dangerous SSH capabilities

These require stronger confirmation or policy approval:

```txt
ssh.systemd.restart-any
ssh.file.write
ssh.file.delete
ssh.package.install
ssh.kubectl.apply
ssh.kubectl.delete
ssh.k3s.restart
ssh.node.reboot
ssh.run.arbitrary
```

## SSH identity model

Each host artifact should include SSH access metadata without storing raw private keys in the artifact.

Example:

```json
{
  "kind": "HostServerArtifact",
  "id": "host:vps-prod-01",
  "name": "vps-prod-01",
  "type": "host-server",
  "network": {
    "sshHost": "203.0.113.10",
    "sshPort": 22
  },
  "ssh": {
    "user": "ubuntu",
    "authRef": "credential:ssh:vps-prod-01",
    "connectionMode": "direct",
    "allowedCapabilities": [
      "ssh.agent.status",
      "ssh.agent.restart",
      "ssh.agent.logs",
      "ssh.host.inventory"
    ]
  }
}
```

Credential material should live in the desktop secure store or user-provided SSH agent, not in the artifact JSON.

## Break-glass audit event

Every SSH fallback action should write an audit event locally in Headlamp.

Required fields:

```json
{
  "kind": "AuditEventArtifact",
  "type": "audit-event",
  "channel": "ssh-fallback",
  "hostId": "host:vps-prod-01",
  "reason": "node-agent-unreachable",
  "command": "systemctl restart headlamp-node-agent",
  "approvedBy": "user:local",
  "startedAt": "2026-05-28T00:00:00Z",
  "completedAt": "2026-05-28T00:00:02Z",
  "exitCode": 0
}
```

## Recovery-first behavior

When SSH fallback is used because the node agent is down, the first priority should be recovery:

1. check service status
2. read recent logs
3. validate config file
4. restart service
5. verify `/healthz`
6. return to normal node-agent API channel

## Architecture boundary

```txt
Headlamp Desktop
  ├─ Agent API client
  ├─ SSH fallback client
  ├─ OPA policy evaluator
  ├─ approval UI
  └─ audit writer

VPS/Host
  ├─ headlamp-node-agent
  ├─ systemd
  ├─ podman
  ├─ k3s/kubernetes
  └─ sshd
```

## Implementation packages

```txt
pkg/sshcontrol/
  client.go
  command.go
  capabilities.go
  audit.go

pkg/nodeagent/recovery/
  status.go
  logs.go
  restart.go
```

## MVP target

The first SSH fallback implementation should support only:

1. check node agent status
2. fetch node agent logs
3. restart node agent
4. verify `/healthz`
5. collect basic host inventory

All other SSH actions should be blocked until explicitly added as governed capabilities.
