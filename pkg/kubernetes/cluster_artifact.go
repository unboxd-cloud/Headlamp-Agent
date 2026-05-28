package kubernetes

import "time"

// ClusterInstanceArtifact represents a Kubernetes or K3s cluster
// discovered and operated by Headlamp.
type ClusterInstanceArtifact struct {
	Kind         string                 `json:"kind"`
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Type         string                 `json:"type"`
	Distribution string                 `json:"distribution"`
	Environment  string                 `json:"environment,omitempty"`
	Provider     string                 `json:"provider,omitempty"`
	Region       string                 `json:"region,omitempty"`
	Endpoint     ClusterEndpoint        `json:"endpoint"`
	Version      *ClusterVersionInfo    `json:"version,omitempty"`
	Topology     *ClusterTopology       `json:"topology,omitempty"`
	Network      *ClusterNetworkInfo    `json:"network,omitempty"`
	Security     *ClusterSecurityInfo   `json:"security,omitempty"`
	Observability *ClusterObservability `json:"observability,omitempty"`
	Status       ClusterStatus          `json:"status"`
	Relationships map[string][]string   `json:"relationships,omitempty"`
}

type ClusterEndpoint struct {
	APIServer    string `json:"apiServer,omitempty"`
	AccessMode   string `json:"accessMode,omitempty"`
	ContextName  string `json:"contextName,omitempty"`
	CredentialRef string `json:"credentialRef,omitempty"`
}

type ClusterVersionInfo struct {
	Kubernetes string `json:"kubernetes,omitempty"`
	K3s        string `json:"k3s,omitempty"`
	Platform   string `json:"platform,omitempty"`
}

type ClusterTopology struct {
	Mode              string `json:"mode,omitempty"`
	ControlPlaneNodes int    `json:"controlPlaneNodes,omitempty"`
	WorkerNodes       int    `json:"workerNodes,omitempty"`
}

type ClusterNetworkInfo struct {
	CNI         string `json:"cni,omitempty"`
	ServiceCIDR string `json:"serviceCIDR,omitempty"`
	PodCIDR     string `json:"podCIDR,omitempty"`
	DNSDomain   string `json:"dnsDomain,omitempty"`
}

type ClusterSecurityInfo struct {
	RBACEnabled  bool     `json:"rbacEnabled,omitempty"`
	PodSecurity  string   `json:"podSecurity,omitempty"`
	PolicyEngines []string `json:"policyEngines,omitempty"`
}

type ClusterObservability struct {
	MetricsServer bool   `json:"metricsServer,omitempty"`
	Logging       string `json:"logging,omitempty"`
	Tracing       string `json:"tracing,omitempty"`
}

type ClusterStatus struct {
	Phase         string    `json:"phase,omitempty"`
	Summary       string    `json:"summary,omitempty"`
	LastObservedAt time.Time `json:"lastObservedAt,omitempty"`
}
