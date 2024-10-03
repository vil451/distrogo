package container

import (
	"distrogo/internal/logger"
	"distrogo/internal/tty"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	errRunningContainer     = "error running container"
	errAttachingToContainer = "error attaching to container"
	sigintTimeout           = 5 * time.Second
)

func (s *Service) Enter(containerName string) error {
	count := 0
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

	var lastSigintTime time.Time

	for {
		select {
		case <-ctx.Done():
			return nil
		case sig := <-sigs:
			if sig == syscall.SIGINT {
				// Если это первый `Ctrl+C` или прошло более 5 секунд с предыдущего
				if time.Since(lastSigintTime) > sigintTimeout {
					fmt.Printf("\nReceived signal: %v. Press again within 5 seconds to exit.\n", sig)
					lastSigintTime = time.Now()
					count = 1 // сбрасываем счётчик
				} else {
					count++
				}

				// Если нажато дважды за 5 секунд
				if count >= 2 {
					fmt.Println("\nExiting...")
					detach(nil)
					return nil
				}
			} else {
				// Получен другой сигнал
				fmt.Printf("\nReceived signal: %v. Exiting...\n", sig)
				detach(nil)
				return nil
			}
		}
	}
}
