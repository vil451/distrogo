package application

import (
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type ContainerService interface {
	IsAttached() bool
	SendSignalToAttach(signal os.Signal)

	List(all bool) ([]types.Container, error)
	Enter(containerName string) error

	FilterByLabelValue(containers []types.Container, labelKey, labelValue string) []types.Container
	FilterByName(containers []types.Container, name string) []types.Container
	FilterByState(containers []types.Container, state string) []types.Container

	RenderTable(containers []types.Container)
}

type dockerCliService interface {
	GetCLI() (*client.Client, error)
	Terminate()
}
