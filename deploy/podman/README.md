# Podman Local Runtime

Headlamp runs as a native desktop app, but supporting local services should run in Podman.

This keeps the desktop experience lightweight while isolating local infrastructure such as local LLM gateways, policy services, test clusters, registries, and artifact stores.

## Principle

```txt
Headlamp Desktop App
  = native Go/Wails desktop app

Local support services
  = Podman containers/pods
```

Headlamp should not require Docker Desktop.

## Local services

The first local Podman profile should support:

| Service | Purpose | Required |
|---|---|---|
| `headlamp-llm-gateway` | OpenAI-compatible local LLM endpoint adapter | optional |
| `headlamp-opa` | OPA policy evaluation server | optional, embedded Go eval also allowed |
| `headlamp-artifacts` | Local artifact workspace volume | yes |
| `headlamp-registry` | Local OCI registry for test artifacts/images | optional |
| `headlamp-k3s-test` | Local K3s test cluster | optional |
| `headlamp-db` | Local metadata/audit store if not embedded SQLite | optional |

Default MVP should use embedded SQLite and embedded OPA evaluation where possible. Podman is for heavier local services.

## Podman pod layout

```txt
pod headlamp-local
  ├─ llm-gateway
  ├─ opa
  ├─ registry
  └─ optional test cluster services
```

## Volumes

```txt
headlamp-data       persistent app data
headlamp-artifacts  local artifact registry mirror/workspace
headlamp-models     optional local model storage
headlamp-opa        OPA bundles and policy data
headlamp-kube       kubeconfigs and test cluster state
```

## Network

Use a named Podman network:

```txt
headlamp-local
```

Suggested local endpoints:

| Endpoint | Service |
|---|---|
| `http://localhost:11434` | Ollama-compatible local LLM |
| `http://localhost:1234/v1` | LM Studio/OpenAI-compatible local endpoint |
| `http://localhost:8181` | OPA server |
| `http://localhost:5000` | local OCI registry |

## Security notes

- Do not mount the full host filesystem by default.
- Do not mount SSH keys by default.
- Do not mount kubeconfig into containers unless explicitly requested.
- Prefer read-only mounts where possible.
- Keep credentials in the desktop app's secure local store, not container env vars.
- Require user approval before any container gets access to kubeconfig or host paths.

## MVP command targets

The repo should provide:

```txt
make podman-up
make podman-down
make podman-logs
make podman-reset
```

These should manage only local support services, not the native desktop app.
