package app

import (
	"bytes"
	"context"
	"kopeika/internal/domain"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAppRepo struct {
	mock.Mock
}

func (m *MockAppRepo) Create(ctx context.Context, app *domain.App) error {
	args := m.Called(ctx, app)
	return args.Error(0)
}

func (m *MockAppRepo) GetByID(ctx context.Context, id string) (*domain.App, error) {
	return nil, nil
}

func (m *MockAppRepo) GetByName(ctx context.Context, name string) (*domain.App, error) {
	return nil, nil
}

func (m *MockAppRepo) UpdateEnvVars(ctx context.Context, name string, envs domain.EnvVars) error {
	return nil
}

func (m *MockAppRepo) UpdateStatus(ctx context.Context, name string, status string) error {
	return nil
}

func TestCreateAppHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Should return 400 when JSON is invalid", func(t *testing.T) {
		mockRepo := new(MockAppRepo)
		appService := NewService(mockRepo)
		appHandler := NewHandler(appService)

		r := gin.New()
		r.POST("/api/v1", appHandler.Create)

		badJSON := []byte(`{
			"name": "some_app", 
			"repo_url": "https://github.com/some",
    		"branch": "main",
    		"env_vars": {
        		"env": "1",
        		"env2": "2"
    		}}`)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1", bytes.NewBuffer(badJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("should return 201 Created when JSON is valid", func(t *testing.T) {
		mockRepo := new(MockAppRepo)
		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
		appService := NewService(mockRepo)
		appHandler := NewHandler(appService)

		r := gin.New()
		r.POST("/api/v1", appHandler.Create)

		validJSON := []byte(`{
			"name": "some-app", 
			"repo_url": "https://github.com/some",
    		"branch": "main",
    		"env_vars": {
        		"env": "1",
        		"env2": "2"
    		}}`)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1", bytes.NewBuffer(validJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockRepo.AssertNumberOfCalls(t, "Create", 1)
	})
}
