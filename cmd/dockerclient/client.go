package dockerclient

import (
	"fmt"
	"github.com/docker/docker/client"
)

func InitDockerClient() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return cli, nil
}

func CloseDockerClient(cli *client.Client) error {
	if err := cli.Close(); err != nil {
		return fmt.Errorf("error closing Docker client: %w", err)
	}
	return nil
}
