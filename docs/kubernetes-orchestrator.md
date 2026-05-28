# Kubernetes Orchestrator Agent

Headlamp Agent is a local-first desktop operator tool for Kubernetes orchestration.

The goal is to support every Kubernetes resource category through a safe, explainable, approval-first agent workflow.

## Scope

The Kubernetes Orchestrator Agent should help users operate:

- core workloads: Pods, Deployments, StatefulSets, DaemonSets, ReplicaSets, Jobs, CronJobs
- networking: Services, Ingress, Gateway API, NetworkPolicies, Endpoints, EndpointSlices
- storage: PersistentVolumes, PersistentVolumeClaims, StorageClasses, CSI resources
- configuration: ConfigMaps, Secrets, ServiceAccounts
- security and access: Roles, RoleBindings, ClusterRoles, ClusterRoleBindings, PodSecurity, RBAC reviews
- cluster operations: Nodes, Namespaces, Events, ResourceQuotas, LimitRanges
- autoscaling: HPA, VPA, Cluster Autoscaler signals where available
- policy resources: admission policies, OPA/Gatekeeper/Kyverno resources where installed
- custom resources: CRDs and any installed custom resource instances
- Helm/Kustomize/GitOps objects where detectable

## Operating model

The agent must follow this loop:

```txt
Discover → Inspect → Explain → Diagnose → Plan → Ask Approval → Apply → Verify → Log
```

The default mode is read-only.

Mutating operations require explicit user approval and should present:

- target cluster and namespace
- affected resource names
- planned manifest diff
- risk level
- rollback path where possible
- exact command/API operation to be executed

## Capability modes

### Observe mode

Read cluster state without mutation.

Examples:

- list resources
- fetch manifests
- read events
- inspect pod logs
- inspect metrics when available
- summarize namespaces
- map workload dependencies

### Diagnose mode

Identify likely causes of operational problems.

Examples:

- CrashLoopBackOff
- ImagePullBackOff
- Pending pods
- failed probes
- unschedulable workloads
- service routing failures
- ingress/backend mismatch
- PVC binding issues
- node pressure
- RBAC denial
- quota exhaustion

### Plan mode

Create a safe remediation plan with evidence.

Examples:

- scale deployment
- restart rollout
- patch image
- update resource requests/limits
- fix service selector
- adjust probe config
- create missing secret/configmap reference
- propose network policy changes

### Apply mode

Execute approved changes only.

Supported mutation categories:

- apply manifest
- patch resource
- scale workload
- restart rollout
- delete pod for controlled recreation
- create namespace-scoped resource
- update config resources

Dangerous operations require stronger confirmation:

- delete namespace
- delete PVC/PV
- modify cluster-wide RBAC
- modify CRDs
- modify admission policies
- drain node
- delete production workloads

### Verify mode

After an approved change, verify outcome.

Examples:

- rollout status
- pod readiness
- event changes
- service endpoints
- ingress health
- metrics checks
- logs after change

## CRD support

The orchestrator must support arbitrary Kubernetes kinds, not only built-in resources.

For CRDs, the agent should:

1. discover API groups and resources
2. read CRD OpenAPI schema when available
3. infer namespaced vs cluster-scoped behavior
4. inspect status conditions
5. avoid mutation unless schema and policy are understood
6. require approval before applying changes

## Safety model

Never mutate the cluster silently.

All actions must pass:

```txt
User intent
  → capability check
  → RBAC check where possible
  → policy check
  → diff generation
  → approval
  → execution
  → audit log
```

## Audit event fields

Each Kubernetes operation should log:

- timestamp
- cluster context
- namespace
- resource kind
- resource name
- user request
- evidence collected
- recommendation
- proposed action
- approval status
- execution result
- verification result
- rollback note

## Local LLM support

The agent should work with local OpenAI-compatible model endpoints first.

The Kubernetes client and tools provide grounded facts. The LLM should explain, rank, summarize, and draft plans, but cluster state must come from Kubernetes APIs, not from model memory.

## Implementation modules

```txt
packages/kubernetes/
  discovery.ts
  client.ts
  resources.ts
  logs.ts
  events.ts
  metrics.ts
  manifests.ts
  diff.ts
  diagnostics.ts
  remediation.ts
  approvals.ts
  audit.ts
```

## Initial MVP

1. Load kubeconfig contexts.
2. Connect to selected cluster.
3. List namespaces and workloads.
4. Chat with cluster context.
5. Diagnose pod/workload failures.
6. Generate safe remediation plan.
7. Show manifest diff.
8. Require approval before apply/patch/scale/restart.
9. Verify change.
10. Save audit history locally.
