# MapReduce Orchestration Model

Headlamp uses MapReduce as a core execution pattern for operating large Kubernetes and artifact environments.

The idea is simple:

```txt
Map many things → analyze each thing independently → reduce results into decisions, plans, risks, and actions
```

This is useful because Kubernetes environments are naturally distributed across:

- clusters
- contexts
- namespaces
- workloads
- pods
- nodes
- CRDs
- policies
- Helm releases
- GitOps applications
- artifact registries

## Why MapReduce

A Kubernetes orchestrator agent must avoid treating a whole environment as one giant prompt.

Instead, it should:

1. split the environment into safe work units
2. inspect each unit with bounded context
3. produce structured findings
4. aggregate findings into summaries and plans
5. ask for approval before any mutation

## General pattern

```txt
Input Scope
  → Planner creates map tasks
  → Map workers inspect each unit
  → Each worker emits structured findings
  → Reducer groups and ranks findings
  → Agent generates plan
  → User approves action
  → Executor applies approved changes
  → Verifier maps again to confirm outcome
  → Audit log records everything
```

## Kubernetes map units

Map tasks can run across these units:

| Unit | Examples |
|---|---|
| Cluster | dev, staging, prod, edge |
| Namespace | default, kube-system, payments |
| Workload | Deployment, StatefulSet, DaemonSet, Job, CronJob |
| Pod | health, logs, restarts, events |
| Node | pressure, capacity, taints, conditions |
| Service path | Service, endpoints, ingress, gateway |
| Storage path | PVC, PV, StorageClass, CSI events |
| Policy surface | RBAC, NetworkPolicy, Kyverno, Gatekeeper, OPA |
| CRD family | CRD schema and custom resource instances |

## Artifact map units

Artifact operations can map across:

| Unit | Examples |
|---|---|
| Distribution artifact | Kubernetes, K3s |
| Cluster artifact | prod cluster, edge cluster |
| Manifest artifact | deployment.yaml, service.yaml |
| Package artifact | Helm chart, Kustomize overlay |
| Policy artifact | Gatekeeper/Kyverno/OPA policy |
| Runbook artifact | incident response, rollout recovery |
| Audit artifact | historical operation record |

## Structured map output

Each map task should emit structured output, not free text.

Example:

```json
{
  "scope": {
    "cluster": "prod-us-east",
    "namespace": "payments",
    "kind": "Deployment",
    "name": "payment-api"
  },
  "status": "warning",
  "findings": [
    {
      "code": "MISSING_RESOURCE_LIMITS",
      "severity": "medium",
      "evidence": [
        "container api has cpu request but no memory limit"
      ],
      "recommendation": "Add memory limit and request based on observed usage."
    }
  ],
  "proposedActions": []
}
```

## Reducer output

Reducers combine map results into operational summaries.

Example reducer outputs:

- cluster health summary
- namespace risk ranking
- failed workload list
- policy violation summary
- remediation backlog
- blast radius estimate
- rollout readiness report
- drift report between Git and runtime
- K3s edge fleet summary

Example:

```json
{
  "summary": "12 workloads need resource limits across 4 namespaces.",
  "riskLevel": "medium",
  "groups": [
    {
      "findingCode": "MISSING_RESOURCE_LIMITS",
      "count": 12,
      "affectedNamespaces": ["payments", "orders", "billing", "search"],
      "recommendedNextStep": "Generate patch plan for affected deployments."
    }
  ]
}
```

## Agent roles

MapReduce can be implemented with logical agent roles:

| Role | Responsibility |
|---|---|
| Planner | Converts user intent into bounded map tasks |
| Mapper | Inspects one bounded unit and emits structured findings |
| Reducer | Aggregates findings and identifies patterns |
| Planner/Reviewer | Builds remediation options and risk notes |
| Executor | Runs only approved actions |
| Verifier | Re-runs focused map tasks after change |
| Auditor | Records input, evidence, approval, action, and result |

These roles do not need to be separate processes initially. They can be functions in the local runtime and later become parallel workers.

## Example: diagnose all unhealthy pods

```txt
User: Diagnose unhealthy pods across all namespaces.

Planner:
  - list namespaces
  - map each namespace
  - inspect pods not Ready or with restart count > threshold

Mapper:
  - read pod status
  - read recent events
  - fetch logs when safe
  - classify cause

Reducer:
  - group by cause
  - rank by severity and blast radius
  - produce summary

Agent:
  - explains top issues
  - proposes fixes
  - asks approval before changes
```

## Example: K3s edge fleet scan

```txt
User: Check all K3s edge clusters for risky drift.

Planner:
  - discover K3s cluster artifacts
  - map each cluster

Mapper:
  - inspect version
  - inspect node conditions
  - inspect bundled components
  - inspect workloads and policies

Reducer:
  - group clusters by version
  - flag outdated clusters
  - flag missing policies
  - produce edge fleet report
```

## Parallelism and limits

MapReduce must be safe by design.

Execution controls:

- max concurrent map tasks
- namespace allowlist/denylist
- cluster allowlist/denylist
- timeout per task
- token budget per task
- log redaction
- secret redaction
- network restrictions
- dry-run option for mutations

## Mutation policy

MapReduce may produce many proposed actions, but it must not execute them automatically.

Allowed pattern:

```txt
Map → Reduce → Plan → Diff → Approval → Execute selected actions → Verify
```

Disallowed pattern:

```txt
Map → Reduce → Bulk mutate silently
```

## Implementation sketch

```txt
packages/orchestration/
  map-task.ts
  map-result.ts
  reducer.ts
  planner.ts
  executor.ts
  verifier.ts

packages/kubernetes/
  map-units.ts
  mappers/
    namespace-health.mapper.ts
    workload-health.mapper.ts
    pod-diagnostics.mapper.ts
    node-health.mapper.ts
    service-path.mapper.ts
    storage-path.mapper.ts
    policy-scan.mapper.ts
    crd-scan.mapper.ts

packages/artifacts/
  artifact-map-units.ts
  artifact-mappers.ts
  artifact-reducers.ts
```

## First implementation target

The first MapReduce workflow should be:

```txt
Scan cluster health
  → map namespaces
  → map workloads in each namespace
  → map unhealthy pods
  → reduce by severity and cause
  → produce remediation plan
```

No mutation in the first implementation. Read-only diagnostics first.
