# Artifact Model

Headlamp is the desktop dashboard/control surface. The Kubernetes Orchestrator Agent operates Kubernetes-related artifacts through safe, approval-first workflows.

In this project, Kubernetes and K3s are treated as artifacts too. They are not only runtime targets.

## Artifact principle

An artifact is any versioned, inspectable, governable object used to define, operate, validate, or explain an environment.

Artifacts can be:

- source artifacts stored in Git
- generated artifacts produced by tools
- runtime artifacts discovered from clusters
- policy artifacts used to govern actions
- operational artifacts such as runbooks, audits, and remediation plans

## Core artifact categories

| Category | Examples | Purpose |
|---|---|---|
| Cluster Distribution | Kubernetes, K3s, MicroK8s, kind, Minikube | Defines the cluster runtime family |
| Cluster Instance | dev cluster, prod cluster, edge cluster | A concrete environment operated by Headlamp |
| Kubernetes Resource | Pod, Deployment, Service, Ingress, Node, Namespace | Runtime/resource objects discovered or applied |
| Custom Resource | CRDs and CR instances | Extension APIs installed into a cluster |
| Manifest | YAML/JSON Kubernetes manifests | Declarative desired state |
| Package | Helm chart, Kustomize overlay | Deployable app/config package |
| Policy | OPA, Gatekeeper, Kyverno, admission policy | Governance and guardrails |
| Runbook | Diagnose/recover procedures | Human and agent operational guidance |
| Remediation Plan | Proposed fix with evidence and diff | Approval-ready action plan |
| Audit Event | Recorded action/recommendation/result | Governance history |

## Kubernetes as artifact

Kubernetes should be modeled as a cluster distribution artifact.

Example fields:

```yaml
kind: ClusterDistributionArtifact
id: distribution:kubernetes
name: Kubernetes
type: cluster-distribution
family: kubernetes
supportedApis:
  - core/v1
  - apps/v1
  - batch/v1
  - networking.k8s.io/v1
  - rbac.authorization.k8s.io/v1
capabilities:
  - workloads
  - networking
  - storage
  - rbac
  - crds
  - admission
  - autoscaling
```

## K3s as artifact

K3s should be modeled as a Kubernetes-compatible distribution artifact with its own defaults and constraints.

Example fields:

```yaml
kind: ClusterDistributionArtifact
id: distribution:k3s
name: K3s
type: cluster-distribution
family: kubernetes
compatibleWith: distribution:kubernetes
traits:
  - lightweight
  - edge-friendly
  - single-binary
  - embedded-defaults
commonDefaults:
  ingressController: traefik
  serviceLoadBalancer: servicelb
  datastore: sqlite
capabilities:
  - workloads
  - networking
  - storage
  - rbac
  - crds
  - edge-clusters
```

## Artifact graph

Artifacts should connect into an operational graph:

```txt
ClusterDistributionArtifact
  └─ ClusterInstanceArtifact
       ├─ KubernetesResourceArtifact
       ├─ CustomResourceArtifact
       ├─ PolicyArtifact
       ├─ PackageArtifact
       ├─ RunbookArtifact
       ├─ RemediationPlanArtifact
       └─ AuditEventArtifact
```

Example:

```txt
K3s
  → powers edge-prod-cluster
  → contains namespace payments
  → runs deployment payment-api
  → governed by kyverno-policy require-resource-limits
  → diagnosed by runbook crashloopbackoff
  → changed by remediation-plan fix-payment-api-probe
  → recorded by audit-event 2026-05-28T22:40Z
```

## Headlamp behavior

Headlamp should be able to:

1. discover artifacts from connected clusters, local folders, and repositories
2. classify artifacts by type
3. explain artifact purpose and relationships
4. validate artifacts against schema and policy
5. compare desired artifacts with runtime artifacts
6. propose remediation plans
7. ask approval before changing artifacts or clusters
8. log artifact changes and runtime actions

## First artifact types to implement

```txt
ClusterDistributionArtifact
ClusterInstanceArtifact
KubernetesResourceArtifact
CustomResourceArtifact
ManifestArtifact
PackageArtifact
PolicyArtifact
RunbookArtifact
RemediationPlanArtifact
AuditEventArtifact
```

## Relationship to artifact-registry

The artifact registry should be the source of truth for artifact definitions, metadata, versions, relationships, and policies.

Headlamp should operate those artifacts locally through a dashboard and agent runtime.

```txt
artifact-registry
  = artifact source of truth

Headlamp
  = desktop operator dashboard

Kubernetes Orchestrator Agent
  = agent capability inside Headlamp
```
