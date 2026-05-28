package sshcontrol

const (
	CapabilityAgentStatus    = "ssh.agent.status"
	CapabilityAgentLogs      = "ssh.agent.logs"
	CapabilityAgentRestart   = "ssh.agent.restart"
	CapabilityHostInventory  = "ssh.host.inventory"
	CapabilitySystemdStatus  = "ssh.systemd.status"
	CapabilityPodmanPS       = "ssh.podman.ps"
	CapabilityK3sStatus      = "ssh.k3s.status"
	CapabilityKubernetesRead = "ssh.kubernetes.readonly"
)

var DefaultAllowedCapabilities = map[string]bool{
	CapabilityAgentStatus:    true,
	CapabilityAgentLogs:      true,
	CapabilityAgentRestart:   true,
	CapabilityHostInventory:  true,
	CapabilitySystemdStatus:  true,
	CapabilityPodmanPS:       true,
	CapabilityK3sStatus:      true,
	CapabilityKubernetesRead: true,
}
