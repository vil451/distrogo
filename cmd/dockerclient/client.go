package dockerclient

import (
	"fmt"
	"github.com/docker/docker/client"
	"os"
)

func InitDockerClient() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return cli, nil
}

func CloseDockerClient(cli *client.Client) {
	if err := cli.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Error closing client: %v\n", err)
	}
}
