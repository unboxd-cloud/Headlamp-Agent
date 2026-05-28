# OPA Policy Artifacts

Headlamp treats OPA policies as first-class artifacts.

OPA artifacts define what the Kubernetes Orchestrator Agent is allowed to recommend, auto-heal, mutate, publish, approve, or block.

OPA is not only an implementation detail. In this architecture, OPA policy bundles, Rego modules, test cases, decisions, and exceptions are all governed artifacts.

## Why OPA artifacts

Auto-healing and orchestration need clear boundaries.

OPA artifacts provide those boundaries as versioned, reviewable, testable policy objects.

```txt
User intent
  → Agent plan
  → Kubernetes/action context
  → OPA policy decision
  → allow / deny / require approval / require escalation
```

## Core OPA artifact types

| Artifact | Purpose |
|---|---|
| `OPAPolicyBundleArtifact` | Versioned bundle of Rego modules and data files |
| `RegoModuleArtifact` | One Rego policy module |
| `OPADataArtifact` | Static policy data such as environment labels, risk rules, namespaces |
| `OPATestArtifact` | Rego tests for policies |
| `OPADecisionArtifact` | Recorded policy decision for an attempted action |
| `OPAExceptionArtifact` | Time-bound, scoped exception to a policy |
| `HealingPolicyArtifact` | Auto-healing rules evaluated by OPA |
| `ActionPolicyArtifact` | Rules for shell, Kubernetes, Git, file, and network actions |
| `KubernetesGuardrailArtifact` | Cluster and resource safety policy |

## Policy decision types

OPA should return one of these decision outcomes:

```txt
allow
require_approval
deny
require_escalation
```

Meaning:

| Decision | Meaning |
|---|---|
| `allow` | Action can proceed automatically |
| `require_approval` | Human approval is required before execution |
| `deny` | Action is blocked |
| `require_escalation` | Higher-trust operator or incident process required |

## OPA artifact graph

```txt
OPAPolicyBundleArtifact
  ├─ RegoModuleArtifact
  ├─ OPADataArtifact
  ├─ OPATestArtifact
  └─ produces OPADecisionArtifact

HealingPolicyArtifact
  └─ evaluated_by OPAPolicyBundleArtifact

KubernetesGuardrailArtifact
  └─ evaluated_by OPAPolicyBundleArtifact

ActionPolicyArtifact
  └─ evaluated_by OPAPolicyBundleArtifact
```

## Example: OPA policy bundle artifact

```yaml
kind: OPAPolicyBundleArtifact
id: opa-bundle:kubernetes-orchestrator-defaults
name: Kubernetes Orchestrator Defaults
type: opa-policy-bundle
version: 0.1.0
entrypoints:
  - data.headlamp.kubernetes.action.decision
  - data.headlamp.healing.decision
modules:
  - policies/kubernetes/action.rego
  - policies/healing/auto_heal.rego
  - policies/common/risk.rego
data:
  - policies/data/environments.json
tests:
  - policies/kubernetes/action_test.rego
  - policies/healing/auto_heal_test.rego
appliesTo:
  - KubernetesResourceArtifact
  - ClusterInstanceArtifact
  - HealingPlaybookArtifact
  - RemediationPlanArtifact
```

## Example: Rego module artifact

```yaml
kind: RegoModuleArtifact
id: rego-module:headlamp-kubernetes-action
name: Headlamp Kubernetes Action Policy
type: rego-module
package: headlamp.kubernetes.action
path: policies/kubernetes/action.rego
entrypoint: data.headlamp.kubernetes.action.decision
inputs:
  - user
  - cluster
  - resource
  - action
  - risk
outputs:
  - decision
  - reasons
  - requiredApprovals
```

## Example: OPA input for Kubernetes action

```json
{
  "user": {
    "id": "user:local",
    "roles": ["operator"]
  },
  "cluster": {
    "id": "cluster:prod-us-east",
    "environment": "production",
    "distribution": "kubernetes"
  },
  "resource": {
    "apiVersion": "apps/v1",
    "kind": "Deployment",
    "namespace": "payments",
    "name": "payment-api"
  },
  "action": {
    "type": "rollout_restart",
    "mutation": true,
    "source": "auto_healing",
    "hasDiff": true
  },
  "risk": {
    "level": "medium",
    "stateful": false,
    "destructive": false,
    "blastRadius": "namespace"
  }
}
```

## Example: OPA decision artifact

```yaml
kind: OPADecisionArtifact
id: opa-decision:2026-05-28T23-10-00Z:payment-api-restart
name: Policy decision for payment-api rollout restart
type: opa-decision
timestamp: 2026-05-28T23:10:00Z
policyBundle: opa-bundle:kubernetes-orchestrator-defaults
inputRef: remediation-plan:restart-payment-api
decision: require_approval
reasons:
  - production cluster requires human approval
  - auto-healing source cannot mutate production workloads automatically
requiredApprovals:
  - role:sre
```

## Example Rego policy

```rego
package headlamp.kubernetes.action

default decision := {
  "result": "deny",
  "reasons": ["no matching allow rule"]
}

production := input.cluster.environment == "production"
mutation := input.action.mutation == true
destructive := input.risk.destructive == true

# Never allow destructive production actions automatically.
decision := {
  "result": "require_escalation",
  "reasons": ["destructive production action requires escalation"]
} if {
  production
  destructive
}

# Production mutations require approval.
decision := {
  "result": "require_approval",
  "reasons": ["production mutations require approval"],
  "requiredApprovals": ["role:sre"]
} if {
  production
  mutation
  not destructive
}

# Low-risk non-production healing can auto-run.
decision := {
  "result": "allow",
  "reasons": ["low-risk non-production healing is allowed"]
} if {
  not production
  input.action.source == "auto_healing"
  input.risk.level == "low"
  not destructive
}
```

## Auto-healing with OPA

Auto-healing must be evaluated through OPA.

```txt
Finding
  → HealingPlaybookArtifact
  → RemediationPlanArtifact
  → OPA input
  → OPA decision
  → allow / approval / deny / escalation
  → execute or stop
  → OPADecisionArtifact logged
```

## Policy artifact storage

Suggested repo paths:

```txt
artifacts/opa/bundles/
artifacts/opa/modules/
artifacts/opa/data/
artifacts/opa/tests/
artifacts/opa/decisions/
artifacts/opa/exceptions/
```

Suggested local Headlamp paths:

```txt
.headlamp/policies/opa/
.headlamp/audit/opa-decisions/
```

## Required checks before using an OPA artifact

Before a policy bundle can govern actions, Headlamp should verify:

1. bundle has a declared entrypoint
2. Rego modules parse successfully
3. tests pass
4. data files are valid JSON/YAML
5. artifact version is pinned
6. artifact provenance is recorded
7. policy bundle is active for the selected cluster/scope

## Relationship to artifact-registry

The artifact registry stores OPA policy artifacts as source-of-truth governance objects.

Headlamp loads and evaluates them locally before executing actions.

```txt
artifact-registry
  → stores policy artifacts

Headlamp
  → loads policy artifacts
  → evaluates agent plans through OPA
  → records OPA decision artifacts

Kubernetes Orchestrator Agent
  → cannot bypass OPA decisions
```
