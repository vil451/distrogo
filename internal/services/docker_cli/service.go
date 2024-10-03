package docker_cli

import (
	"fmt"

	"distrogo/internal/logger"
	"github.com/docker/docker/client"
)

type Service struct {
	cli *client.Client
}

func New() (*Service, error) {
	return &Service{}, nil
}

func (s *Service) GetCLI() (*client.Client, error) {
	if s.cli != nil {
		return s.cli, nil
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		logger.Error(fmt.Errorf("error creating docker client: %v", err))
		return nil, err
	}
	s.cli = cli
	return s.cli, nil
}

func (s *Service) Terminate() {
	if s.cli == nil {
		err := s.cli.Close()
		if err != nil {
			logger.Error(fmt.Errorf("error closing docker client: %v", err))
		}
	}
}
