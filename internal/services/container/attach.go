package container

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/strslice"
	"github.com/pkg/errors"
)

const (
	ErrCreatingExecInstance    = "error creating exec instance"
	ErrCAttachingToExecSession = "error attaching to exec session"
)

func (s *Service) Attach(containerName string) (*types.HijackedResponse, context.Context, context.CancelFunc, error) {
	ctx, cancel := context.WithCancel(context.Background())

	execConfig := types.ExecConfig{
		Cmd:          strslice.StrSlice([]string{"/bin/sh"}),
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
	}

	execIDResp, err := s.cli.ContainerExecCreate(ctx, containerName, execConfig)
	if err != nil {
		cancel()
		return nil, nil, nil, errors.Wrap(err, ErrCreatingExecInstance)
	}

	attachResp, err := s.cli.ContainerExecAttach(ctx, execIDResp.ID, types.ExecStartCheck{Tty: true})
	if err != nil {
		cancel()
		return nil, nil, nil, errors.Wrap(err, ErrCAttachingToExecSession)
	}

	return &attachResp, ctx, cancel, nil
}
