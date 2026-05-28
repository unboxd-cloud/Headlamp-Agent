package artifacts

import "time"

// HostServerArtifact represents a physical or virtual machine
// that may participate in Kubernetes, K3s, storage, registry,
// gateway, or other operational roles.
type HostServerArtifact struct {
	Kind        string            `json:"kind"`
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Environment string            `json:"environment,omitempty"`
	Provider    string            `json:"provider,omitempty"`
	Region      string            `json:"region,omitempty"`
	Location    *LocationInfo     `json:"location,omitempty"`
	OS          *OSInfo           `json:"os,omitempty"`
	Hardware    *HardwareInfo     `json:"hardware,omitempty"`
	Network     *NetworkInfo      `json:"network,omitempty"`
	Roles       []string          `json:"roles,omitempty"`
	Runtimes    *RuntimeInfo      `json:"runtimes,omitempty"`
	ClusterRefs []string          `json:"clusterRefs,omitempty"`
	Status      *ArtifactStatus   `json:"status,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}

type LocationInfo struct {
	Country   string `json:"country,omitempty"`
	Region    string `json:"region,omitempty"`
	City      string `json:"city,omitempty"`
	Datacenter string `json:"datacenter,omitempty"`
	Rack      string `json:"rack,omitempty"`
}

type OSInfo struct {
	Family      string `json:"family,omitempty"`
	Distribution string `json:"distribution,omitempty"`
	Version     string `json:"version,omitempty"`
	Kernel      string `json:"kernel,omitempty"`
	Architecture string `json:"architecture,omitempty"`
}

type HardwareInfo struct {
	CPUCores   int   `json:"cpuCores,omitempty"`
	MemoryBytes int64 `json:"memoryBytes,omitempty"`
	DiskBytes  int64 `json:"diskBytes,omitempty"`
	GPUCount   int   `json:"gpuCount,omitempty"`
}

type NetworkInfo struct {
	Hostname  string   `json:"hostname,omitempty"`
	PrivateIPs []string `json:"privateIps,omitempty"`
	PublicIPs []string `json:"publicIps,omitempty"`
	SSHHost   string   `json:"sshHost,omitempty"`
	SSHPort   int      `json:"sshPort,omitempty"`
}

type RuntimeInfo struct {
	ContainerRuntime string `json:"containerRuntime,omitempty"`
	DockerVersion    string `json:"dockerVersion,omitempty"`
	ContainerdVersion string `json:"containerdVersion,omitempty"`
	KubernetesVersion string `json:"kubernetesVersion,omitempty"`
	K3sVersion       string `json:"k3sVersion,omitempty"`
}

type ArtifactStatus struct {
	Phase          string    `json:"phase,omitempty"`
	Summary        string    `json:"summary,omitempty"`
	LastObservedAt time.Time `json:"lastObservedAt,omitempty"`
}
