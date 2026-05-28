# Auto-Healing Model

Headlamp can support auto-healing, but it must be policy-bounded, evidence-driven, and reversible where possible.

Auto-healing does not mean the agent can freely mutate a Kubernetes cluster. It means the system can detect known failure modes, match them to approved remediation playbooks, and execute only actions that are allowed by policy.

## Healing maturity levels

### Level 0: Observe only

The system detects issues and reports them.

No remediation is proposed or executed.

### Level 1: Recommend only

The system diagnoses the issue and recommends a remediation plan.

No changes are made.

### Level 2: Supervised healing

The system prepares an exact action plan and manifest/command diff.

A human must approve before execution.

This is the default MVP target.

### Level 3: Policy-bounded auto-healing

The system may automatically execute pre-approved low-risk remediations when all guardrails pass.

Examples:

- restart a failed pod owned by a Deployment
- trigger rollout restart for a non-production namespace
- scale within configured min/max bounds
- recreate a stuck Job when retry policy allows it
- cordon a node only when policy explicitly permits it

### Level 4: Closed-loop autonomous remediation

The system diagnoses, remediates, verifies, rolls back if needed, and records the full loop.

This level should only be enabled for narrow, tested, reversible playbooks.

## Default posture

The default mode is **supervised healing**, not full auto-healing.

```txt
Detect → Diagnose → Plan → Ask Approval → Apply → Verify → Log
```

Policy-bounded auto-healing can be enabled per cluster, namespace, resource type, and remediation playbook.

## Auto-healing decision gate

Every healing action must pass this gate:

```txt
Known failure pattern?
  → enough evidence?
  → supported remediation playbook?
  → policy allows automation?
  → RBAC allows action?
  → blast radius acceptable?
  → rollback or recovery path exists?
  → rate limit not exceeded?
  → execute
  → verify
  → audit
```

If any gate fails, the system must fall back to supervised approval.

## Safe auto-healing candidates

Good first candidates:

| Failure | Possible healing | Notes |
|---|---|---|
| Pod stuck terminating | delete pod after grace checks | Only if controller-owned |
| CrashLoopBackOff caused by transient dependency | restart rollout | Only if recent config did not change |
| ImagePullBackOff from temporary registry error | retry/recreate pod | Do not change image automatically |
| Job failed due transient node issue | recreate Job | Respect backoff and idempotency |
| Deployment rollout stuck | pause/rollback suggestion | Auto-rollback only if policy allows |
| HPA at max replicas | recommend scaling limit review | Usually not auto-fix |
| Node NotReady | cordon or drain suggestion | Auto-cordon only with explicit policy |
| PVC pending | diagnose storage class/binding | Usually recommend-only |
| Service has no endpoints | diagnose selector mismatch | Usually supervised patch only |

## Actions that must not be auto-healed by default

These require human approval unless explicitly allowed by a strong policy:

- deleting namespaces
- deleting PVCs or PVs
- modifying cluster-wide RBAC
- modifying CRDs
- modifying admission policies
- modifying secrets
- changing production images
- changing database/stateful workloads
- draining nodes
- changing network policies in production
- scaling to zero
- applying unknown generated manifests

## Healing policy artifact

Auto-healing rules should be modeled as artifacts.

Example:

```yaml
kind: HealingPolicyArtifact
id: healing-policy:dev-safe-defaults
name: Dev Safe Defaults
type: healing-policy
scope:
  clusters:
    - dev-*
  namespaces:
    - sandbox
    - dev
allowedActions:
  - restart-rollout
  - delete-controller-owned-pod
  - recreate-failed-job
limits:
  maxActionsPerHour: 10
  maxAffectedPods: 5
  requireHumanApprovalFor:
    - production
    - statefulset
    - secret
    - persistentvolumeclaim
verification:
  required: true
  timeoutSeconds: 300
rollback:
  requiredWhenAvailable: true
```

## Healing playbook artifact

Each remediation should be represented as a playbook artifact.

Example:

```yaml
kind: HealingPlaybookArtifact
id: healing-playbook:restart-crashlooping-deployment
name: Restart CrashLooping Deployment
type: healing-playbook
failurePatterns:
  - CrashLoopBackOff
requiredEvidence:
  - pod.status.containerStatuses.state.waiting.reason
  - pod events
  - owning deployment
allowedResourceKinds:
  - Deployment
actions:
  - type: rollout-restart
verification:
  - rollout-status
  - pod-readiness
  - recent-events
riskLevel: low
requiresApprovalByDefault: true
```

## MapReduce auto-healing flow

Auto-healing should reuse the MapReduce model.

```txt
Map:
  inspect clusters/namespaces/workloads/pods

Reduce:
  group findings by failure pattern and blast radius

Plan:
  match findings to healing playbooks

Gate:
  evaluate healing policy

Execute:
  run only approved or policy-allowed actions

Verify:
  map affected resources again

Audit:
  record full loop
```

## Verification requirements

Every healing action must verify outcome.

Examples:

- Deployment rollout completed
- pods became Ready
- failed pods reduced
- service endpoints exist
- no new warning events appeared
- logs no longer show the triggering error

If verification fails:

1. stop further automated actions in that scope
2. record failure
3. recommend rollback or escalation
4. require human review

## Anti-loop controls

Auto-healing must avoid retry storms.

Controls:

- max actions per cluster per hour
- max actions per namespace per hour
- max repeated action per resource
- cooldown windows
- stop after failed verification
- require approval after repeated recurrence
- incident mode escalation

## Audit fields

Each auto-healing loop should log:

- trigger
- evidence
- matched playbook
- policy decision
- action diff/command
- approval status or auto-approval reason
- execution result
- verification result
- rollback state
- recurrence counter

## MVP implementation target

The first implementation should support supervised healing only:

```txt
Detect unhealthy pod
  → diagnose cause
  → propose remediation
  → show exact diff/command
  → ask approval
  → execute after approval
  → verify
  → log
```

After that, introduce Level 3 policy-bounded auto-healing for dev/test clusters only.
