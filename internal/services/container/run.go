package container

import (
	"context"
	"distrogo/cmd/listcontainers"
	"fmt"
	"github.com/docker/docker/api/types/container"
)

func (s *Service) Run(containerName string) error {
	ctx := context.Background()
	containers, err := listcontainers.GetContainers(ctx, s.cli, true)
	if err != nil {
		return fmt.Errorf("error listing containers: %v", err)
	}

	containers = listcontainers.FilterContainersByLabel(containers, "manager", "distrogo")
	var resultContainerID, state string
	for _, container := range containers {
		if container.Names[0][1:] == containerName {
			resultContainerID = container.ID
			state = container.State
		}
	}
	if state == "running" {
		return nil
	}
	if resultContainerID == "" {
		return fmt.Errorf("container not found: %s", containerName)
	}

	startOptions := container.StartOptions{}
	if err := s.cli.ContainerStart(ctx, resultContainerID, startOptions); err != nil {
		return fmt.Errorf("error starting container: %v", err)
	}

	fmt.Printf("Container %s is started with ID: %s\n", containerName, resultContainerID)
	return nil
}
