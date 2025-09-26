package docker

import (
	"context"
	"fmt"

	"github.com/Artemka007/derraform/internal/logging"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

type DockerClient struct {
	cli    *client.Client
	logger *logging.Logger
}

func NewDockerClient() (*DockerClient, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	logger := logging.NewLogger(logging.DEBUG)
	if err != nil {
		return nil, err
	}
	return &DockerClient{cli: cli, logger: logger}, nil
}

// internal/providers/docker/client.go

// NetworkConfig конфигурация для создания сети
type NetworkConfig struct {
	Name   string
	Driver string
}

// CreateNetwork создает Docker сеть
func (d *DockerClient) CreateNetwork(ctx context.Context, config *NetworkConfig) (string, error) {
	if d.logger == nil {
		d.logger = logging.NewLogger(logging.INFO)
	}

	d.logger.Info("Creating network: %s", config.Name)

	resp, err := d.cli.NetworkCreate(ctx, config.Name, network.CreateOptions{
		Driver: config.Driver,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create network: %w", err)
	}

	d.logger.Info("Network %s created successfully with ID: %s", config.Name, resp.ID[:12])
	return resp.ID, nil
}

// DestroyNetwork удаляет Docker сеть
func (d *DockerClient) DestroyNetwork(ctx context.Context, networkID string) error {
	if d.logger == nil {
		d.logger = logging.NewLogger(logging.INFO)
	}

	d.logger.Info("Destroying network: %s", networkID[:12])

	if err := d.cli.NetworkRemove(ctx, networkID); err != nil {
		return fmt.Errorf("failed to remove network: %w", err)
	}

	d.logger.Info("Network %s destroyed successfully", networkID[:12])
	return nil
}

func (d *DockerClient) DestroyContainer(ctx context.Context, containerID string) error {
	if d.logger == nil {
		d.logger = logging.NewLogger(logging.INFO)
	}

	d.logger.Info("Destroying container: %s", containerID[:12])

	// Останавливаем контейнер
	d.logger.Debug("Stopping container: %s", containerID[:12])
	timeout := 30 // seconds
	if err := d.cli.ContainerStop(ctx, containerID, container.StopOptions{Timeout: &timeout}); err != nil {
		d.logger.Warn("Failed to stop container %s: %v", containerID[:12], err)
		// Продолжаем удаление даже если не удалось остановить
	} else {
		d.logger.Debug("Container stopped successfully: %s", containerID[:12])
	}

	// Удаляем контейнер
	d.logger.Debug("Removing container: %s", containerID[:12])
	if err := d.cli.ContainerRemove(ctx, containerID, container.RemoveOptions{
		Force:         true,  // Принудительное удаление
		RemoveVolumes: true,  // Удаляем связанные тома
		RemoveLinks:   false, // Не удаляем линки
	}); err != nil {
		d.logger.Error("Failed to remove container %s: %v", containerID[:12], err)
		return fmt.Errorf("failed to remove container: %w", err)
	}

	d.logger.Info("Container destroyed successfully: %s", containerID[:12])
	return nil
}
