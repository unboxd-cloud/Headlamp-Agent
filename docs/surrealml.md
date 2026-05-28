# SurrealML Intelligence Layer

SurrealML should be part of Headlamp's intelligence stack.

It does not replace local LLMs. It complements them.

## Role split

```txt
Local LLM
  = chat, reasoning, explanation, remediation drafting, YAML/code assistance

SurrealML
  = data-native predictions, classification, scoring, anomaly detection, learned operational signals inside SurrealDB

OPA
  = deterministic governance and allow/deny/approval decisions

ClickHouse
  = high-volume historical analytics and feature extraction
```

## Why SurrealML

Headlamp stores operational graph and event data in SurrealDB:

- agents
- hosts
- clusters
- Kubernetes/K3s artifacts
- heartbeats
- incidents
- notifications
- OPA decisions
- remediation plans
- audit events
- recovery outcomes

SurrealML is useful when intelligence should run close to this data instead of sending everything to an external model service.

## Good SurrealML use cases

SurrealML should be used for:

| Use case | Example |
|---|---|
| Health scoring | score agent/host/cluster health from recent signals |
| Incident risk scoring | predict whether warning state will become critical |
| Auto-healing confidence | estimate probability a playbook will succeed |
| Anomaly detection | detect unusual restart/heartbeat/event patterns |
| Classification | classify incident type from structured fields |
| Prioritization | rank remediation backlog |
| Drift scoring | score runtime vs desired artifact drift |
| Fleet risk | identify hosts/clusters likely to fail soon |

## Things SurrealML should not do

SurrealML should not be the authority for:

- Kubernetes source-of-truth state
- policy allow/deny decisions
- unrestricted cluster mutation
- secret interpretation
- direct shell command generation

Those remain with Kubernetes APIs, OPA, and approval workflows.

## Intelligence flow

```txt
Node Agent / Kubernetes collectors
  → structured events and artifacts
  → SurrealDB
  → SurrealML scoring/classification
  → Cortex evaluates health/risk
  → OPA gates actions
  → Headlamp displays recommendation
```

## Example model outputs

```json
{
  "kind": "MLSignalArtifact",
  "type": "ml-signal",
  "source": "surrealml",
  "subjectId": "cluster:prod-k3s-01",
  "signal": "cluster_failure_risk",
  "score": 0.82,
  "label": "high-risk",
  "modelRef": "surrealml:model:cluster-risk:v1",
  "createdAt": "2026-05-28T00:00:00Z"
}
```

## SurrealQL tables

Suggested tables:

```surql
DEFINE TABLE ml_models SCHEMAFULL;
DEFINE TABLE ml_signals SCHEMAFULL;
DEFINE TABLE ml_training_runs SCHEMAFULL;
DEFINE TABLE ml_feature_snapshots SCHEMAFULL;
```

## Relationship with local LLM

The local LLM can explain SurrealML signals to the operator.

Example:

```txt
SurrealML says cluster risk = 0.82
  → local LLM explains likely causes using evidence
  → Cortex opens incident if threshold passes
  → OPA decides whether healing can run
```

## Product rule

Use SurrealML for operational scoring and prediction close to SurrealDB data.

Use local LLMs for natural-language interaction and reasoning over gathered evidence.

Use OPA for decisions that must be deterministic and auditable.
