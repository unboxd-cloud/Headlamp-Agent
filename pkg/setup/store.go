package setup

import "context"

// PlanStore loads and saves install plans.
// Implementations may use SurrealDB, files, or embedded defaults.
type PlanStore interface {
	GetInstallPlan(ctx context.Context, profileID ProfileID) (*InstallPlan, error)
	SaveInstallPlan(ctx context.Context, plan InstallPlan) error
}

// LoadPlanWithFallback tries the store first, then falls back to the embedded default.
func LoadPlanWithFallback(ctx context.Context, store PlanStore, profileID ProfileID) InstallPlan {
	if store != nil {
		if plan, err := store.GetInstallPlan(ctx, profileID); err == nil && plan != nil {
			return *plan
		}
	}

	if profileID == ProfileLocalKubernetesOperator || profileID == "" {
		return DefaultInstallPlan()
	}

	return DefaultInstallPlan()
}
