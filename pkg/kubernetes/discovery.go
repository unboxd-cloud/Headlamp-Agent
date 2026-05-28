package kubernetes

import (
	"context"
	"fmt"
	"time"
)

// DiscoveryService discovers Kubernetes and K3s clusters
// and converts them into artifact models.
type DiscoveryService struct{}

// DiscoverCluster creates a minimal cluster artifact from
// kubeconfig context information.
//
// Full client-go integration will be added next.
func (d *DiscoveryService) DiscoverCluster(ctx context.Context, source KubeConfigSource) (*ClusterInstanceArtifact, error) {
	path, err := ResolveKubeConfigPath(source)
	if err != nil {
		return nil, err
	}

	artifact := &ClusterInstanceArtifact{
		Kind:         "ClusterInstanceArtifact",
		ID:           fmt.Sprintf("cluster:%s", source.Context),
		Name:         source.Context,
		Type:         "cluster-instance",
		Distribution: "distribution:kubernetes",
		Endpoint: ClusterEndpoint{
			AccessMode:  "kubeconfig",
			ContextName: source.Context,
		},
		Status: ClusterStatus{
			Phase:         "unknown",
			Summary:       fmt.Sprintf("cluster discovered from kubeconfig %s", path),
			LastObservedAt: time.Now().UTC(),
		},
	}

	return artifact, nil
}
