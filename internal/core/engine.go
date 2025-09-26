// internal/core/engine.go
package core

import (
	"context"
	"fmt"

	"github.com/Artemka007/derraform/internal/config"
	"github.com/Artemka007/derraform/internal/errors"
	"github.com/Artemka007/derraform/internal/logging"
	"github.com/Artemka007/derraform/internal/providers/docker"
	"github.com/Artemka007/derraform/internal/state"
	"github.com/zclconf/go-cty/cty"
)

type Engine struct {
	config       *config.Config
	stateManager *state.StateManager
	dockerClient *docker.DockerClient
	logger       *logging.Logger
}

func NewEngine() (*Engine, error) {
	// Создаем Docker клиент
	dockerClient, err := docker.NewDockerClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}

	// Создаем state manager
	stateManager := state.NewStateManager("terraform.tfstate")

	// Создаем логгер
	logger := logging.NewLogger(logging.INFO)

	return &Engine{
		stateManager: stateManager,
		dockerClient: dockerClient,
		logger:       logger,
	}, nil
}

func (e *Engine) Apply(configFile string) error {
	e.logger.Info("Starting deployment...")

	// Парсим конфигурацию
	cfg, err := config.ParseFile(configFile)
	if err != nil {
		return errors.WrapError(err, "CONFIG_ERROR", "Failed to parse configuration file")
	}
	e.config = cfg

	e.logger.Info("Found %d resources to process", len(cfg.Resources))

	// Применяем каждый ресурс
	for _, resource := range cfg.Resources {
		resourceID := fmt.Sprintf("%s.%s", resource.Type, resource.Name)
		e.logger.Info("Processing resource: %s", resourceID)

		if err := e.applyResource(resource); err != nil {
			e.logger.Error("Failed to apply resource %s: %v", resourceID, err)
			return errors.ResourceError(resourceID, "Failed to apply resource", err)
		}

		e.logger.Info("Resource %s applied successfully", resourceID)
	}

	e.logger.Info("Deployment completed successfully!")
	return nil
}

func (e *Engine) applyResource(resource config.Resource) error {
	switch resource.Type {
	case "docker_container":
		return e.applyDockerContainer(resource)
	case "docker_network":
		return e.applyDockerNetwork(resource)
	case "docker_volume":
		return e.applyDockerVolume(resource)
	case "docker_image":
		return e.applyDockerImage(resource)
	default:
		return errors.NewError("UNKNOWN_RESOURCE",
			fmt.Sprintf("Unknown resource type: %s", resource.Type))
	}
}

// applyDockerContainer применяет конфигурацию Docker контейнера
func (e *Engine) applyDockerContainer(resource config.Resource) error {
	e.logger.Debug("Applying Docker container: %s", resource.Name)

	// Преобразуем атрибуты в Docker конфиг
	containerConfig, err := e.resourceToContainerConfig(resource)
	if err != nil {
		return fmt.Errorf("failed to parse container config: %w", err)
	}

	// Создаем контейнер
	containerID, err := e.dockerClient.CreateContainer(context.Background(), containerConfig)
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	// Сохраняем состояние
	if err := e.stateManager.SaveResourceState(resource.Type, resource.Name, map[string]interface{}{
		"id":    containerID,
		"name":  containerConfig.Name,
		"image": containerConfig.Image,
	}); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	e.logger.Info("Docker container %s applied successfully", resource.Name)
	return nil
}

// applyDockerNetwork применяет конфигурацию Docker сети
func (e *Engine) applyDockerNetwork(resource config.Resource) error {
	e.logger.Debug("Applying Docker network: %s", resource.Name)

	// Преобразуем атрибуты в сетевой конфиг
	networkConfig, err := e.resourceToNetworkConfig(resource)
	if err != nil {
		return fmt.Errorf("failed to parse network config: %w", err)
	}

	// Создаем сеть
	networkID, err := e.dockerClient.CreateNetwork(context.Background(), networkConfig)
	if err != nil {
		return fmt.Errorf("failed to create network: %w", err)
	}

	// Сохраняем состояние
	if err := e.stateManager.SaveResourceState(resource.Type, resource.Name, map[string]interface{}{
		"id":   networkID,
		"name": networkConfig.Name,
	}); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	e.logger.Info("Docker network %s applied successfully", resource.Name)
	return nil
}

// applyDockerVolume применяет конфигурацию Docker тома
func (e *Engine) applyDockerVolume(resource config.Resource) error {
	e.logger.Debug("Applying Docker volume: %s", resource.Name)

	// Для томов пока просто логируем
	e.logger.Info("Volume support will be implemented later: %s", resource.Name)
	return nil
}

// applyDockerImage применяет конфигурацию Docker образа
func (e *Engine) applyDockerImage(resource config.Resource) error {
	e.logger.Debug("Applying Docker image: %s", resource.Name)

	// Для образов пока просто логируем
	e.logger.Info("Image support will be implemented later: %s", resource.Name)
	return nil
}

func (e *Engine) resourceToContainerConfig(resource config.Resource) (*docker.ContainerConfig, error) {
	config := &docker.ContainerConfig{
		Name: resource.Name,
	}

	// Извлекаем image (обязательный атрибут)
	imageVal, exists := resource.Attributes["image"]
	if !exists {
		return nil, fmt.Errorf("missing required attribute 'image'")
	}

	if imageVal.Type() == cty.String {
		config.Image = imageVal.AsString()
	} else {
		return nil, fmt.Errorf("attribute 'image' must be a string")
	}

	// Обрабатываем порты
	if portsVal, exists := resource.Attributes["ports"]; exists {
		if portsVal.Type().IsObjectType() || portsVal.Type().IsMapType() {
			config.Ports = make(map[string]string)
			portsMap := portsVal.AsValueMap()

			for key, value := range portsMap {
				if value.Type() == cty.String {
					config.Ports[key] = value.AsString()
				}
			}
		}
	}

	// Обрабатываем environment variables
	if envVal, exists := resource.Attributes["env"]; exists {
		if envVal.Type().IsObjectType() || envVal.Type().IsMapType() {
			config.Env = make(map[string]string)
			envMap := envVal.AsValueMap()

			for key, value := range envMap {
				if value.Type() == cty.String {
					config.Env[key] = value.AsString()
				}
			}
		}
	}

	// Обрабатываем сети
	if networksVal, exists := resource.Attributes["networks"]; exists {
		if networksVal.Type().IsListType() || networksVal.Type().IsTupleType() {
			networksList := networksVal.AsValueSlice()
			config.Networks = make([]string, len(networksList))

			for i, netVal := range networksList {
				if netVal.Type() == cty.String {
					config.Networks[i] = netVal.AsString()
				}
			}
		}
	}

	// Обрабатываем команду
	if commandVal, exists := resource.Attributes["command"]; exists {
		if commandVal.Type().IsListType() {
			commandList := commandVal.AsValueSlice()
			config.Command = make([]string, len(commandList))

			for i, cmdVal := range commandList {
				if cmdVal.Type() == cty.String {
					config.Command[i] = cmdVal.AsString()
				}
			}
		}
	}

	return config, nil
}

// resourceToNetworkConfig преобразует Resource в Docker NetworkConfig
func (e *Engine) resourceToNetworkConfig(resource config.Resource) (*docker.NetworkConfig, error) {
	config := &docker.NetworkConfig{
		Name:   resource.Name,
		Driver: "bridge", // Значение по умолчанию
	}

	// Извлекаем driver если есть
	if driverVal, exists := resource.Attributes["driver"]; exists {
		if driverVal.Type() == cty.String {
			config.Driver = driverVal.AsString()
		}
	}

	return config, nil
}

// Plan показывает план изменений
func (e *Engine) Plan(configFile string) error {
	e.logger.Info("Generating execution plan...")

	cfg, err := config.ParseFile(configFile)
	if err != nil {
		return err
	}

	e.logger.Info("Plan:")
	for _, resource := range cfg.Resources {
		e.logger.Info("  + create %s.%s", resource.Type, resource.Name)
	}

	e.logger.Info("This plan would create %d resources.", len(cfg.Resources))
	return nil
}

// Destroy удаляет все ресурсы
func (e *Engine) Destroy() error {
	e.logger.Info("Destroying all resources...")

	// Загружаем состояние
	state, err := e.stateManager.Load()
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	// Удаляем ресурсы в обратном порядке
	for resourceID, resourceState := range state.Resources {
		e.logger.Info("Destroying resource: %s", resourceID)

		switch resourceState.Type {
		case "docker_container":
			if err := e.dockerClient.DestroyContainer(context.Background(), resourceState.ID); err != nil {
				e.logger.Error("Failed to destroy container %s: %v", resourceID, err)
			}
		case "docker_network":
			if err := e.dockerClient.DestroyNetwork(context.Background(), resourceState.ID); err != nil {
				e.logger.Error("Failed to destroy network %s: %v", resourceID, err)
			}
		}
	}

	// Очищаем состояние
	if err := e.stateManager.Clear(); err != nil {
		return fmt.Errorf("failed to clear state: %w", err)
	}

	e.logger.Info("Destruction completed!")
	return nil
}
