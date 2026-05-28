package kubernetes

import (
	"fmt"
	"os"
	"path/filepath"
)

// KubeConfigSource describes where Headlamp should load Kubernetes access from.
type KubeConfigSource struct {
	Path    string `json:"path,omitempty"`
	Context string `json:"context,omitempty"`
}

// ResolveKubeConfigPath returns the kubeconfig path to use.
//
// It prefers an explicit path, then KUBECONFIG, then ~/.kube/config.
func ResolveKubeConfigPath(source KubeConfigSource) (string, error) {
	if source.Path != "" {
		return source.Path, nil
	}

	if env := os.Getenv("KUBECONFIG"); env != "" {
		return env, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve user home: %w", err)
	}

	return filepath.Join(home, ".kube", "config"), nil
}
