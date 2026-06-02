package domain

import (
	"context"
	"time"
)

const (
	StatusIdle = "idle"
	StatusBuilding = "building"
	StatusDeployed = "deployed"
	StatusFailed = "failed"
)

type App struct {
	ID string `json:"id"`
	UserID int `json:"user_id"`
	Name string `json:"name"`
	RepoURL string `json:"repo_url"`
	Branch string `json:"branch"`
	Status string `json:"status"`
	EnvVars map[string]string `json:"env_vars"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AppRepository interface {
	Create(ctx context.Context, app *App) error
	GetByID(ctx context.Context, id string) (*App, error)
	GetByName(ctx context.Context, name string) (*App, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	UpdateEnvVars(ctx context.Context, id string, envs map[string]string) error
}

type K8sService interface {
	CreateDeployment(ctx context.Context, app *App, imageTag string) error
	CreateService(ctx context.Context, app *App) error
	CreateIngress(ctx context.Context, app *App) error
	StreamLogs(ctx context.Context, appName string) (chan string, error)
}