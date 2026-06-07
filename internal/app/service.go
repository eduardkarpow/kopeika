package app

import (
	"context"
	"kopeika/internal/domain"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo domain.AppRepository
}

func NewService(repo domain.AppRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, app domain.App) error {
	err := app.Validate()
	if err != nil {
		return err
	}
	app.ID = uuid.New().String()
	app.Status = string(domain.StatusIdle)
	app.CreatedAt = time.Now().UTC()
	app.UpdatedAt = time.Now().UTC()

	err = s.repo.Create(ctx, &app)
	return err
}

func (s *Service) GetByID(ctx context.Context, id string) (*domain.App, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetByName(ctx context.Context, name string) (*domain.App, error) {
	return s.repo.GetByName(ctx, name)
}

func (s *Service) UpdateStatus(ctx context.Context, id string, status string) error {
	statusValidate := domain.ProjectStatus(status)
	if err := statusValidate.Validate(); err != nil {
		return err
	}
	return s.repo.UpdateStatus(ctx, id, status)
}

func (s *Service) UpdateEnvVars(ctx context.Context, id string, envs domain.EnvVars) error {
	return s.repo.UpdateEnvVars(ctx, id, envs)
}
