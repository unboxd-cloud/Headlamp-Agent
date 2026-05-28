package setup

// ProfileID identifies a local Headlamp setup option.
type ProfileID string

const (
	ProfileMinimalLocal            ProfileID = "minimal-local"
	ProfileLocalKubernetesOperator ProfileID = "local-kubernetes-operator"
	ProfileFullLocalOpsLab         ProfileID = "full-local-ops-lab"
	ProfileCustom                  ProfileID = "custom"
)

// Component describes a local service or capability shown before install.
type Component struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Required bool     `json:"required"`
	Purpose  string   `json:"purpose"`
	Ports    []string `json:"ports,omitempty"`
}

// Profile is the first-run install choice shown to the user.
type Profile struct {
	ID          ProfileID   `json:"id"`
	Name        string      `json:"name"`
	Summary     string      `json:"summary"`
	Recommended bool        `json:"recommended"`
	Components  []Component `json:"components"`
}

func Profiles() []Profile {
	return []Profile{
		LocalKubernetesOperatorProfile(),
	}
}
