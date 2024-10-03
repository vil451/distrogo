package application

import (
	"os"

	"github.com/docker/docker/client"
)

type ContainerService interface {
	IsAttached() bool
	SendSignalToAttach(signal os.Signal)

	Enter(containerName string) error
}

type dockerCliService interface {
	GetCLI() (*client.Client, error)
	Terminate()
}
