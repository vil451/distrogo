package container

import (
	"fmt"

	distrogoLabels "distrogo/internal/services/container/labels"
	distrogoState "distrogo/internal/services/container/state"
	"github.com/docker/docker/api/types/container"
	"github.com/pkg/errors"
)

const (
	ErrRunContainer = "error run container"
)

func (s *Service) Run(containerName string) error {
	cli, err := s.cliService.GetCLI()
	if err != nil {
		return errors.Wrap(err, ErrRunContainer)
	}

	containers, err := s.List(true)
	if err != nil {
		return errors.Wrap(err, ErrRunContainer)
	}
	distrogoContainers := s.FilterByLabelValue(containers, distrogoLabels.LabelManager, distrogoLabels.LabelValueDistrogo)

	containersWithName := s.FilterByName(distrogoContainers, containerName)
	if len(containersWithName) == 0 {
		return errors.Wrap(fmt.Errorf("container not found: %s", containerName), ErrRunContainer)
	}

	cont := containersWithName[0]
	if cont.State == distrogoState.Running {
		return nil
	}

	startOptions := container.StartOptions{}
	if errStart := cli.ContainerStart(s.ctx, cont.ID, startOptions); errStart != nil {
		return errors.Wrap(fmt.Errorf("error starting container: %v", errStart), ErrRunContainer)
	}

	fmt.Printf("Container %s is started with ID: %s\n", containerName, cont.ID)
	return nil
}
