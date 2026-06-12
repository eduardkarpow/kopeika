package domain

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type ProjectStatus string

const (
	StatusIdle     ProjectStatus = "idle"
	StatusBuilding ProjectStatus = "building"
	StatusDeployed ProjectStatus = "deployed"
	StatusFailed   ProjectStatus = "failed"
)

type App struct {
	ID        string    `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Name      string    `json:"name" db:"name"`
	RepoURL   string    `json:"repo_url" db:"repo_url"`
	Branch    string    `json:"branch" db:"branch"`
	Status    string    `json:"status" db:"status"`
	EnvVars   EnvVars   `json:"env_vars" db:"env_vars"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type AppRepository interface {
	Create(ctx context.Context, app *App) error
	GetByID(ctx context.Context, id string) (*App, error)
	GetByName(ctx context.Context, name string) (*App, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	UpdateEnvVars(ctx context.Context, id string, envs EnvVars) error
}

type K8sService interface {
	CreateDeployment(ctx context.Context, app *App, imageTag string) error
	CreateService(ctx context.Context, app *App) error
	CreateIngress(ctx context.Context, app *App) error
	StreamLogs(ctx context.Context, appName string) (chan string, error)
}

type EnvVars map[string]string

func (e *EnvVars) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to byte[] failed")
	}
	return json.Unmarshal(bytes, e)
}

func (e EnvVars) Value() (driver.Value, error) {
	return json.Marshal(e)
}

func (s ProjectStatus) Validate() error {
	switch s {
	case StatusBuilding, StatusDeployed, StatusFailed, StatusIdle:
		return nil
	default:
		return fmt.Errorf("%w: %s", errors.New("invalid status value"), s)
	}
}

var (
	NameRegex         = "^[a-z0-9]([-a-z0-9]*[a-z0-9])?$"
	ErrInvalidName    = errors.New("Name is not valid")
	ErrBranchInvalid  = errors.New("Branch is Empty")
	ErrRepoUrlInvalid = errors.New("repo url is Empty")
)

func (a *App) Validate() error {
	ok, _ := regexp.MatchString(NameRegex, a.Name)
	if !ok {
		return ErrInvalidName
	}
	if strings.TrimSpace(a.Branch) == "" {
		return ErrBranchInvalid
	}
	if strings.TrimSpace(a.RepoURL) == "" {
		return ErrRepoUrlInvalid
	}
	return nil
}
