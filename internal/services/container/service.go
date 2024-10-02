package container

import (
	"distrogo/internal/logger"
	"fmt"
	"github.com/docker/docker/client"
)

type Service struct {
	cli *client.Client
}

func New() (*Service, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		logger.Error(fmt.Errorf("error creating docker client: %v", err))
		return nil, err
	}
	return &Service{
		cli: cli,
	}, nil
}

func (s *Service) Terminate() {
	if s.cli == nil {
		err := s.cli.Close()
		if err != nil {
			logger.Error(fmt.Errorf("error closing docker client: %v", err))
		}
	}
}
