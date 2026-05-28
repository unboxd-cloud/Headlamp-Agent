package setup

func LocalKubernetesOperatorProfile() Profile {
	return Profile{
		ID:          ProfileLocalKubernetesOperator,
		Name:        "Local Kubernetes Operator",
		Summary:     "Minimal local setup for Kubernetes/K3s operations with a scalable path to fleet monitoring and automation.",
		Recommended: true,
		Components: []Component{
			{ID: "headlamp-desktop", Name: "Headlamp Desktop", Required: true, Purpose: "Desktop dashboard, chat UI, approvals, cluster explorer, and operator control surface."},
			{ID: "podman", Name: "Podman", Required: true, Purpose: "Runs local support services without Docker Desktop."},
			{ID: "surrealdb", Name: "SurrealDB", Required: true, Purpose: "Operational graph, artifacts, live events, heartbeats, incidents, and current state.", Ports: []string{"127.0.0.1:8000"}},
			{ID: "opa", Name: "OPA", Required: true, Purpose: "Policy gate for Kubernetes actions, SSH fallback, auto-healing, and approvals.", Ports: []string{"127.0.0.1:8181"}},
			{ID: "node-agent", Name: "Headlamp Node Agent", Required: true, Purpose: "Bounded host/VPS/Kubernetes/K3s inventory and approved actions.", Ports: []string{"127.0.0.1:9080"}},
			{ID: "ssh-fallback", Name: "SSH Fallback", Required: true, Purpose: "Break-glass recovery when the node agent is unavailable."},
			{ID: "local-llm", Name: "Local LLM", Required: false, Purpose: "Optional local reasoning, explanation, and YAML/Rego/SurrealQL generation.", Ports: []string{"127.0.0.1:11434", "127.0.0.1:1234"}},
			{ID: "cortex", Name: "Cortex Watchdog", Required: false, Purpose: "Optional always-on watchdog for fleet health, incidents, and recovery coordination."},
			{ID: "clickhouse", Name: "ClickHouse", Required: false, Purpose: "Optional high-volume telemetry and historical analytics.", Ports: []string{"127.0.0.1:8123", "127.0.0.1:9000"}},
			{ID: "temporal", Name: "Temporal", Required: false, Purpose: "Optional durable workflows for recovery, remediation, and auto-healing."},
		},
	}
}
