package dockerCLI

import (
	"context"
	"kopeika/internal/domain"
	"strings"
	"testing"
	"time"

	"github.com/moby/moby/client"
)

func TestOrchestrator_Deploy_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	cli, err := client.New(client.FromEnv)
	if err != nil {
		t.Fatalf("Failed to create docker cli client")
	}
	defer cli.Close()

	manager, err := NewInfrastructureManager()
	if err != nil {
		t.Fatalf("Failed to create infrastructure manager")
	}
	app := domain.App{
		ID:        "a1fde0b7-f853-4172-9bd4-16da988e580a",
		UserID:    0,
		Name:      "some-app",
		RepoURL:   "https://github.com/eduardkarpow/123",
		Branch:    "main",
		Status:    string(domain.StatusIdle),
		EnvVars:   domain.EnvVars{"var1": "123"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := manager.Deploy(ctx, &app)
	if err != nil {
		t.Fatalf("container deploy error: %v", err)
	}

	defer cli.ContainerStop(ctx, id, client.ContainerStopOptions{})
	defer cli.ContainerRemove(ctx, id, client.ContainerRemoveOptions{Force: true})

	inspect, err := cli.ContainerInspect(ctx, "some-app", client.ContainerInspectOptions{})
	if err != nil {
		t.Fatalf("Container not found after deploy: %v", err)
	}

	c := inspect.Container

	if got := strings.TrimPrefix(c.Name, "/"); got != "some-app" {
		t.Errorf("container name mismatch %q", got)
	}

	if c.State == nil || !c.State.Running {
		t.Errorf("Container is not running: state=%+v", c.State)
	}

	if c.Config == nil {
		t.Fatal("Container has no config: labels unavailable")
	}

	labels := c.Config.Labels

	wantLabels := map[string]string{
		"traefik.enable":                           "true",
		"traefik.http.routes.some-app.rule":        "Host(`some-app.localhost`)",
		"traefik.http.services.some-app.localhost": "80",
	}

	for key, want := range wantLabels {
		got, ok := labels[key]
		if !ok {
			t.Errorf("missing traefik label %q", key)
			continue
		}
		if got != want {
			t.Errorf("label %q = %q want %q", key, got, want)
		}
	}
}
