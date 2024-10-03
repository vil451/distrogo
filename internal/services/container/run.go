package container

import (
	"fmt"

	"distrogo/cmd/listcontainers"
	"github.com/docker/docker/api/types/container"
)

func (s *Service) Run(containerName string) error {
	cli, err := s.cliService.GetCLI()
	if err != nil {
		return err
	}

	containers, err := listcontainers.GetContainers(s.ctx, cli, true)
	if err != nil {
		return fmt.Errorf("error listing containers: %v", err)
	}

	containers = listcontainers.FilterContainersByLabel(containers, "manager", "distrogo")
	var resultContainerID, state string
	for _, cont := range containers {
		if cont.Names[0][1:] == containerName {
			resultContainerID = cont.ID
			state = cont.State
		}
	}
	if state == "running" {
		return nil
	}
	if resultContainerID == "" {
		return fmt.Errorf("container not found: %s", containerName)
	}

	startOptions := container.StartOptions{}
	if errStart := cli.ContainerStart(s.ctx, resultContainerID, startOptions); errStart != nil {
		return fmt.Errorf("error starting container: %v", errStart)
	}

	fmt.Printf("Container %s is started with ID: %s\n", containerName, resultContainerID)
	return nil
}
