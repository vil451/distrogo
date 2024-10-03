package container

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/strslice"
	"github.com/pkg/errors"
)

const (
	ErrCreatingExecInstance    = "error creating exec instance"
	ErrCAttachingToExecSession = "error attaching to exec session"
)

func (s *Service) Attach(containerName string) (*types.HijackedResponse, error) {
	cli, err := s.cliService.GetCLI()
	if err != nil {
		return nil, err
	}

	execConfig := types.ExecConfig{
		Cmd:          strslice.StrSlice([]string{"/bin/sh"}),
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
	}

	execIDResp, err := cli.ContainerExecCreate(s.ctx, containerName, execConfig)
	if err != nil {
		return nil, errors.Wrap(err, ErrCreatingExecInstance)
	}

	attachResp, err := cli.ContainerExecAttach(s.ctx, execIDResp.ID, types.ExecStartCheck{Tty: true})
	if err != nil {
		return nil, errors.Wrap(err, ErrCAttachingToExecSession)
	}

	return &attachResp, nil
}
