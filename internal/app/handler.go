package app

import (
	"kopeika/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
	var app domain.App

	if err := c.ShouldBindJSON(&app); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	ctx := c.Request.Context()
	err := h.service.Create(ctx, app)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true})
}

func (h *Handler) GetByID(c *gin.Context) {
	id := string(c.Query("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id cannot be empty"})
		return
	}
	ctx := c.Request.Context()
	app, err := h.service.GetByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "App not found"})
		return
	}
	c.JSON(http.StatusOK, app)
}

func (h *Handler) GetByName(c *gin.Context) {
	name := string(c.Query("name"))
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name cannot be empty"})
		return
	}
	ctx := c.Request.Context()
	app, err := h.service.GetByName(ctx, name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "App not found"})
		return
	}
	c.JSON(http.StatusOK, app)
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	status := string(c.Query("status"))
	id := string(c.Query("id"))
	if status == "" || id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status or id cannot be empty"})
		return
	}
	ctx := c.Request.Context()
	err := h.service.UpdateStatus(ctx, id, status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "App not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) UpdateEnvVars(c *gin.Context) {
	id := string(c.Query("id"))

	var envs domain.EnvVars

	if err := c.ShouldBindJSON(&envs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body structure"})
		return
	}
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id cannot be empty"})
		return
	}
	ctx := c.Request.Context()
	err := h.service.UpdateEnvVars(ctx, id, envs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
