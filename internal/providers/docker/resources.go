package docker

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
)

type ContainerConfig struct {
	Name        string
	Image       string
	Ports       map[string]string
	Env         map[string]string
	Networks    []string
	Volumes     []VolumeMount
	HealthCheck *HealthCheck
	Command     []string // Добавляем поле Command
}

type VolumeMount struct {
	Source   string
	Target   string
	ReadOnly bool
}

type HealthCheck struct {
	Test     []string
	Interval time.Duration
	Timeout  time.Duration
	Retries  int
}

func (d *DockerClient) CreateContainer(ctx context.Context, config *ContainerConfig) (string, error) {
	// Pull image
	d.logger.Info("Pulling image: %s", config.Image)
	reader, err := d.cli.ImagePull(ctx, config.Image, image.PullOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to pull image: %w", err)
	}
	defer reader.Close()

	// Parse port bindings
	portBindings := make(nat.PortMap)
	exposedPorts := make(nat.PortSet)

	for internal, external := range config.Ports {
		port := nat.Port(internal + "/tcp")
		exposedPorts[port] = struct{}{}
		portBindings[port] = []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: external,
			},
		}
	}

	// Prepare environment variables
	var envVars []string
	for key, value := range config.Env {
		envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
	}

	// Prepare volume mounts
	var mounts []mount.Mount
	for _, vol := range config.Volumes {
		mounts = append(mounts, mount.Mount{
			Type:     mount.TypeBind,
			Source:   vol.Source,
			Target:   vol.Target,
			ReadOnly: vol.ReadOnly,
		})
	}

	// Prepare health check
	var healthConfig *container.HealthConfig
	if config.HealthCheck != nil {
		healthConfig = &container.HealthConfig{
			Test:     config.HealthCheck.Test,
			Interval: config.HealthCheck.Interval,
			Timeout:  config.HealthCheck.Timeout,
			Retries:  config.HealthCheck.Retries,
		}
	}

	// Create container
	resp, err := d.cli.ContainerCreate(ctx,
		&container.Config{
			Image:        config.Image,
			Env:          envVars,
			ExposedPorts: exposedPorts,
			Healthcheck:  healthConfig,
		},
		&container.HostConfig{
			PortBindings: portBindings,
			Mounts:       mounts,
			RestartPolicy: container.RestartPolicy{
				Name: "unless-stopped",
			},
		},
		&network.NetworkingConfig{},
		nil,
		config.Name,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	// Connect to networks
	for _, networkName := range config.Networks {
		if err := d.cli.NetworkConnect(ctx, networkName, resp.ID, nil); err != nil {
			d.logger.Warn("Failed to connect container to network %s: %v", networkName, err)
		}
	}

	// Start container
	if err := d.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("failed to start container: %w", err)
	}

	d.logger.Info("Container %s created successfully with ID: %s", config.Name, resp.ID[:12])
	return resp.ID, nil
}
