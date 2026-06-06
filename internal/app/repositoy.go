package app

import (
	"context"
	"kopeika/internal/domain"

	"github.com/jmoiron/sqlx"
)

type AppRepository struct {
	db *sqlx.DB
}

func NewAppRepository(db *sqlx.DB) *AppRepository {
	return &AppRepository{
		db: db,
	}
}

func (r *AppRepository) GetByID(ctx context.Context, id string) (*domain.App, error) {
	query := "SELECT * FROM app WHERE id = $1"

	var a domain.App

	err := r.db.GetContext(ctx, &a, query, id)

	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AppRepository) GetByName(ctx context.Context, name string) (*domain.App, error) {
	query := "SELECT * FROM app WHERE name = $1"

	var a domain.App

	err := r.db.GetContext(ctx, &a, query, name)

	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AppRepository) UpdateEnvVars(ctx context.Context, id string, envs domain.EnvVars) error {
	query := "UPDATE app SET env_vars = env_vars || $2 WHERE id = $1"

	_, err := r.db.ExecContext(ctx, query, id, envs)
	return err
}

func (r *AppRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	query := "UPDATE app SET status = $2 WHERE id = $1"

	_, err := r.db.ExecContext(ctx, query, id, status)
	return err
}

func (r *AppRepository) Create(ctx context.Context, app *domain.App) error {
	query := `
		INSERT INTO app (id, user_id, name, repo_url, branch, status, env_vars, created_at, updated_at)
		VALUES (:id, :user_id, :name, :repo_url, :branch, :status, :env_vars, :created_at, :updated_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, app)
	return err
}
