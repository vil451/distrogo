package application

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"distrogo/internal/services/container"
	"distrogo/internal/services/docker_cli"
)

type Application struct {
	ctx       context.Context
	ctxCancel context.CancelFunc

	dockerCliService dockerCliService
	containerService ContainerService
}

func New() (*Application, error) {
	ctx, cancel := context.WithCancel(context.Background())
	app := &Application{
		ctx:       ctx,
		ctxCancel: cancel,
	}
	go app.signalsHandler()
	return app, nil
}

func (a *Application) Terminate() {
	a.ctxCancel()
}

func (a *Application) GetContainerService() (ContainerService, error) {
	if a.containerService != nil {
		return a.containerService, nil
	}

	cliSvc, err := a.getCliService()
	if err != nil {
		return nil, err
	}

	containerService, err := container.New(a.ctx, cliSvc)
	if err != nil {
		return nil, err
	}
	a.containerService = containerService
	return a.containerService, nil
}

func (a *Application) getCliService() (dockerCliService, error) {
	if a.dockerCliService != nil {
		return a.dockerCliService, nil
	}

	cliSvc, err := docker_cli.New()
	if err != nil {
		return nil, err
	}
	a.dockerCliService = cliSvc
	return a.dockerCliService, nil
}

func (a *Application) signalsHandler() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-a.ctx.Done():
			a.cleanUp()
			return
		case sig := <-sigs:
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				if a.containerService != nil && a.containerService.IsAttached() {
					a.containerService.SendSignalToAttach(sig)
					continue
				}

				a.Terminate()
			}
		}
	}
}

func (a *Application) cleanUp() {
	if a.dockerCliService != nil {
		a.dockerCliService.Terminate()
	}
}
