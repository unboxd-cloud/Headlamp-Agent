package setup

// InstallPlan is the reviewable plan shown before Headlamp downloads,
// starts, or configures any local dependency.
type InstallPlan struct {
	ProfileID        ProfileID   `json:"profileId"`
	ProfileName      string      `json:"profileName"`
	Summary          string      `json:"summary"`
	Required         []Component `json:"required"`
	Optional         []Component `json:"optional"`
	Ports            []string    `json:"ports"`
	CloudDisabled    bool        `json:"cloudDisabled"`
	RequiresApproval bool        `json:"requiresApproval"`
	Notes            []string    `json:"notes"`
}

// BuildInstallPlan turns a setup profile into a user-reviewable install plan.
func BuildInstallPlan(profile Profile) InstallPlan {
	plan := InstallPlan{
		ProfileID:        profile.ID,
		ProfileName:      profile.Name,
		Summary:          profile.Summary,
		CloudDisabled:    true,
		RequiresApproval: true,
		Notes: []string{
			"No downloads or services should start until the operator approves this plan.",
			"All services are local-first and should bind to localhost by default.",
			"Kubernetes mutations require OPA policy checks and explicit approval unless a policy allows automation.",
			"SSH fallback is for recovery when the node agent is unavailable.",
		},
	}

	seenPorts := map[string]bool{}
	for _, component := range profile.Components {
		if component.Required {
			plan.Required = append(plan.Required, component)
		} else {
			plan.Optional = append(plan.Optional, component)
		}

		for _, port := range component.Ports {
			if !seenPorts[port] {
				seenPorts[port] = true
				plan.Ports = append(plan.Ports, port)
			}
		}
	}

	return plan
}

// DefaultInstallPlan returns the recommended minimal-but-scalable Kubernetes operator setup.
func DefaultInstallPlan() InstallPlan {
	return BuildInstallPlan(LocalKubernetesOperatorProfile())
}
