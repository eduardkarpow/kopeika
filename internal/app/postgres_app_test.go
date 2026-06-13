package app_test

import (
	"context"
	"kopeika/internal/app"
	"kopeika/internal/domain"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestPostgresRepository_GetByID(t *testing.T) {
	ctx := context.Background()

	dbName := "test_db"
	dbUser := "user"
	dbPassword := "pass"

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(30*time.Second),
		),
	)
	require.NoError(t, err)

	defer func() {
		err := postgresContainer.Terminate(ctx)
		assert.NoError(t, err)
	}()

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	db, err := sqlx.Connect("postgres", connStr)
	require.NoError(t, err)
	defer db.Close()

	schema := `
	CREATE TABLE IF NOT EXISTS app (
    id VARCHAR(255) PRIMARY KEY,
    user_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    repo_url TEXT NOT NULL,
    branch VARCHAR(255),
    status VARCHAR(50) NOT NULL,
    env_vars JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);`
	_, err = db.Exec(schema)
	require.NoError(t, err)

	insertQuery := `
		INSERT INTO app (id, user_id, name, repo_url, branch, status, env_vars, created_at, updated_at)
		VALUES (:id, :user_id, :name, :repo_url, :branch, :status, :env_vars, :created_at, :updated_at)
	`

	idStub := "a1fde0b7-f853-4172-9bd4-16da988e580a"
	fakeApp := domain.App{
		ID:        idStub,
		UserID:    0,
		Name:      "some-app",
		RepoURL:   "https://github.com/eduardkarpow/123",
		Branch:    "main",
		Status:    string(domain.StatusIdle),
		EnvVars:   domain.EnvVars{"var1": "123"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err = db.NamedExecContext(ctx, insertQuery, fakeApp)
	require.NoError(t, err)

	repo := app.NewAppRepository(db)
	t.Run("should return app", func(t *testing.T) {
		app, err := repo.GetByID(ctx, idStub)

		require.NoError(t, err)
		assert.NotNil(t, app)
		assert.Equal(t, idStub, app.ID)
		assert.Equal(t, "some-app", app.Name)
		assert.Equal(t, "main", app.Branch)
		assert.Equal(t, "123", app.EnvVars["var1"])
	})

	t.Run("should return nothing", func(t *testing.T) {
		app, err := repo.GetByID(ctx, "some")
		assert.Error(t, err)
		assert.Nil(t, app)
	})
}
