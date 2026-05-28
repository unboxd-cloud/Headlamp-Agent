-- Seed Headlamp install plans into SurrealDB.
-- Run after surreal/schema.surql.

USE NS headlamp DB headlamp;

UPSERT install_plans:local_kubernetes_operator CONTENT {
  kind: "InstallPlanArtifact",
  type: "install-plan",
  version: "0.1.0",
  profile_id: "local-kubernetes-operator",
  profile_name: "Local Kubernetes Operator",
  summary: "Minimal local setup for Kubernetes/K3s operations with a scalable path to fleet monitoring and automation.",
  cloud_disabled: true,
  requires_approval: true,
  required: [
    {
      id: "headlamp-desktop",
      name: "Headlamp Desktop",
      purpose: "Desktop dashboard, chat UI, approvals, cluster explorer, and operator control surface.",
      ports: []
    },
    {
      id: "podman",
      name: "Podman",
      purpose: "Runs local support services without Docker Desktop.",
      ports: []
    },
    {
      id: "surrealdb",
      name: "SurrealDB",
      purpose: "Operational graph, artifacts, live events, heartbeats, incidents, and current state.",
      ports: ["127.0.0.1:8000"]
    },
    {
      id: "opa",
      name: "OPA",
      purpose: "Policy gate for Kubernetes actions, SSH fallback, auto-healing, and approvals.",
      ports: ["127.0.0.1:8181"]
    },
    {
      id: "node-agent",
      name: "Headlamp Node Agent",
      purpose: "Bounded host/VPS/Kubernetes/K3s inventory and approved actions.",
      ports: ["127.0.0.1:9080"]
    },
    {
      id: "ssh-fallback",
      name: "SSH Fallback",
      purpose: "Break-glass recovery when the node agent is unavailable.",
      ports: []
    }
  ],
  optional: [
    {
      id: "local-llm",
      name: "Local LLM",
      purpose: "Optional local reasoning, explanation, and YAML/Rego/SurrealQL generation.",
      ports: ["127.0.0.1:11434", "127.0.0.1:1234"]
    },
    {
      id: "cortex",
      name: "Cortex Watchdog",
      purpose: "Optional always-on watchdog for fleet health, incidents, and recovery coordination.",
      ports: []
    },
    {
      id: "clickhouse",
      name: "ClickHouse",
      purpose: "Optional high-volume telemetry and historical analytics.",
      ports: ["127.0.0.1:8123", "127.0.0.1:9000"]
    },
    {
      id: "temporal",
      name: "Temporal",
      purpose: "Optional durable workflows for recovery, remediation, and auto-healing.",
      ports: []
    }
  ],
  notes: [
    "No downloads or services should start until the operator approves this plan.",
    "All services are local-first and should bind to localhost by default.",
    "Kubernetes mutations require OPA policy checks and explicit approval unless a policy allows automation.",
    "SSH fallback is for recovery when the node agent is unavailable."
  ],
  created_at: time::now(),
  updated_at: time::now()
};
