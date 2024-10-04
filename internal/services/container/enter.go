package container

import (
	"distrogo/internal/logger"
	"distrogo/internal/tty"
	"github.com/pkg/errors"
)

const (
	errRunningContainer     = "error running container"
	errAttachingToContainer = "error attaching to container"
)

func (s *Service) Enter(containerName string) error {
	err := s.Run(containerName)
	if err != nil {
		return errors.Wrap(err, errRunningContainer)
	}

	attachResp, err := s.Attach(containerName)
	if err != nil {
		return errors.Wrap(err, errAttachingToContainer)
	}

	done := make(chan struct{})
	detach := func(err error) {
		if err != nil {
			logger.Debug(err)
		}
		attachResp.Close()
		s.attachTty = nil
		done <- struct{}{}
	}

	s.attachTty = tty.NewTTY(s.ctx, attachResp.Conn, attachResp.Reader, detach)

	select {
	case <-s.ctx.Done():
		detach(s.ctx.Err())
	case <-done:
		break
	}
	return nil
}
