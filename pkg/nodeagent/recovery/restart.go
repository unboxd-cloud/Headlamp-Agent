package recovery

// RecoveryPlan describes how Headlamp should restore a failed node agent.
type RecoveryPlan struct {
	HostID string   `json:"hostId"`
	Reason string   `json:"reason"`
	Steps  []string `json:"steps"`
}

func DefaultAgentRecoveryPlan(hostID string) RecoveryPlan {
	return RecoveryPlan{
		HostID: hostID,
		Reason: "node-agent-unreachable",
		Steps: []string{
			"check service status",
			"inspect recent logs",
			"validate config",
			"restart headlamp-node-agent",
			"verify /healthz",
		},
	}
}
