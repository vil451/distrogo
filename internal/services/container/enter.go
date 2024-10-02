package container

import (
	"distrogo/internal/logger"
	"distrogo/internal/tty"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"os/signal"
	"syscall"
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

	attachResp, ctx, ctxCancel, err := s.Attach(containerName)
	if err != nil {
		return errors.Wrap(err, errAttachingToContainer)
	}

	detach := func(err error) {
		if err != nil {
			logger.Debug(err)
		}
		ctxCancel()
		attachResp.Close()
	}

	tty.NewTTY(ctx, attachResp.Conn, attachResp.Reader, detach)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		return nil
	case sig := <-sigs:
		fmt.Printf("\nReceived signal: %v. Exiting...\n", sig)
		detach(nil)
	}

	fmt.Println("Session terminated.")
	return nil
}
