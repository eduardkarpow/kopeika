package dockerCLI

import (
	"context"
	"fmt"
	"kopeika/internal/domain"
	"strings"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

type InfrastructureManager struct {
	dockerCli *client.Client
}

func NewInfrastructureManager() (*InfrastructureManager, error) {
	cli, err := client.New(client.FromEnv)

	if err != nil {
		return nil, err
	}

	return &InfrastructureManager{dockerCli: cli}, nil
}

func (im *InfrastructureManager) Deploy(ctx context.Context, appEntity *domain.App) (string, error) {
	var env []string
	for k, v := range appEntity.EnvVars {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	labels := map[string]string{
		"traefik.enable": "true",
		fmt.Sprintf("traefik.http.routes.%s.rule", appEntity.Name):        fmt.Sprintf("Host(`%s.localhost`)", strings.ToLower(appEntity.Name)),
		fmt.Sprintf("traefik.http.services.%s.localhost", appEntity.Name): "80",
	}

	config := &container.Config{
		Image:  "nginx:alpine",
		Env:    env,
		Labels: labels,
	}

	hostConfig := &container.HostConfig{
		NetworkMode: "default",
	}

	resp, err := im.dockerCli.ContainerCreate(ctx, client.ContainerCreateOptions{
		Config:     config,
		HostConfig: hostConfig,
		Name:       appEntity.Name,
	},
	)

	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}
	_, err = im.dockerCli.ContainerStart(ctx, resp.ID, client.ContainerStartOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to start container: %w", err)
	}
	return resp.ID, nil
}
