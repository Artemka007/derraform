package core

import (
    "context"
    "fmt"
    
    "github.com/yourname/myterraform/internal/errors"
    "github.com/yourname/myterraform/internal/logging"
    "github.com/yourname/myterraform/internal/ui"
)

type Engine struct {
    config      *config.Config
    state       *state.StateManager
    docker      *docker.DockerClient
    logger      *logging.Logger
    progress    *ui.ProgressTracker
    errorReporter *ui.ErrorReporter
}

func NewEngine() (*Engine, error) {
    dockerClient, err := docker.NewDockerClient()
    if err != nil {
        return nil, err
    }

    logger := logging.NewLogger(logging.INFO)
    progress := ui.NewProgressTracker()
    errorReporter := ui.NewErrorReporter()

    return &Engine{
        state:        state.NewStateManager("terraform.tfstate"),
        docker:       dockerClient,
        logger:       logger,
        progress:     progress,
        errorReporter: errorReporter,
    }, nil
}

func (e *Engine) Apply(configFile string) error {
    e.logger.Info("Starting deployment...")
    e.progress.StartStep("Parsing configuration")
    
    cfg, err := config.ParseFile(configFile)
    if err != nil {
        e.progress.EndStep(false, "Configuration parsing failed")
        return errors.WrapError(err, "CONFIG_ERROR", "Failed to parse configuration file")
    }
    e.config = cfg
    e.progress.EndStep(true, fmt.Sprintf("Found %d resources", len(cfg.Resources)))

    // Применяем изменения для каждого ресурса
    for _, resource := range cfg.Resources {
        resourceID := fmt.Sprintf("%s.%s", resource.Type, resource.Name)
        e.progress.StartStep(fmt.Sprintf("Processing %s", resourceID))
        
        if err := e.applyResource(resource); err != nil {
            e.progress.EndStep(false, "Failed")
            terraErr := errors.ResourceError(resourceID, "Failed to apply resource", err)
            e.errorReporter.AddError(terraErr)
            e.logger.Error("Resource application failed: %v", err)
        } else {
            e.progress.EndStep(true, "Success")
            e.logger.Info("Resource %s applied successfully", resourceID)
        }
    }

    e.errorReporter.PrintSummary()
    if e.errorReporter.HasErrors() {
        return errors.NewError("DEPLOYMENT_FAILED", "Deployment completed with errors")
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
    default:
        return errors.NewError("UNKNOWN_RESOURCE", 
            fmt.Sprintf("Unknown resource type: %s", resource.Type))
    }
}