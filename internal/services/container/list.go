package container

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/pkg/errors"
)

const (
	ErrGetContainersList = "error get containers list"
)

func (s *Service) List(all bool) ([]types.Container, error) {
	cli, err := s.cliService.GetCLI()
	if err != nil {
		return nil, err
	}

	options := container.ListOptions{}
	if all {
		options.All = true
	}
	containers, err := cli.ContainerList(s.ctx, options)
	if err != nil {
		return nil, errors.Wrap(err, ErrGetContainersList)
	}
	return containers, nil
}
