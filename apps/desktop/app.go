package main

import (
	"context"

	"github.com/unboxd-cloud/headlamp-agent/pkg/setup"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetSetupPlan() setup.InstallPlan {
	return setup.DefaultInstallPlan()
}

func (a *App) VerifySetup() setup.VerifyReport {
	return setup.VerifyLocalKubernetesOperator()
}
