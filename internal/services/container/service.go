package container

import (
	"context"
	"os"

	"github.com/docker/docker/client"
)

type CliService interface {
	GetCLI() (*client.Client, error)
}

type Tty interface {
	SendSignal(signal os.Signal)
}

type Service struct {
	ctx context.Context

	cliService CliService
	attachTty  Tty
}

func New(ctx context.Context, cliService CliService) (*Service, error) {
	return &Service{
		ctx:        ctx,
		cliService: cliService,
	}, nil
}

func (s *Service) IsAttached() bool {
	return s.attachTty != nil
}

func (s *Service) SendSignalToAttach(signal os.Signal) {
	if s.IsAttached() {
		s.attachTty.SendSignal(signal)
	}
}
