package createcontainers

import (
	"context"
	"distrogo/cmd/dockerclient"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"io"
	"os"
)

func CreateContainer() *cobra.Command {
	var containerName string
	var imageName string
	var pullImage bool
	command := &cobra.Command{
		Use:     "create",
		Short:   "Create a container",
		Aliases: []string{"c"},
		Run: func(cmd *cobra.Command, args []string) {
			create(containerName, pullImage)
		},
	}

	command.Flags().StringVarP(
		&imageName,
		"image",
		"i",
		"",
		"image name of a container",
	)

	command.Flags().StringVarP(
		&containerName,
		"name",
		"n",
		"",
		"container name",
	)

	command.Flags().BoolVarP(
		&pullImage,
		"pull",
		"p",
		true,
		"pull image",
	)
	return command
}

func create(containerName string, pull bool) {
	ctx := context.Background()
	cli, err := dockerclient.InitDockerClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating Docker client: %v\n", err)
		os.Exit(1)
	}
	defer dockerclient.CloseDockerClient(cli)
	if pull {
		_, err := pullImage(ctx, cli, containerName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error pulling image: %v\n", err)
		}
	}
	_, err = createContainer(ctx, cli, containerName)
}

func pullImage(ctx context.Context, cli *client.Client, name string) (io.ReadCloser, error) {
	config := &container.Config{
		Image: name,
	}
	resp, err := cli.ImagePull(ctx, config.Image, types.ImagePullOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error pulling image: %v\n", err)
		os.Exit(1)
	}
	return resp, nil
}

func createContainer(ctx context.Context, cli *client.Client, name string) (container.CreateResponse, error) {
	//options := types.ContainerListOptions{}
	config := &container.Config{
		Image: name,
		Cmd:   []string{"echo", "Hello, World!"},
	}
	hostConfig := &container.HostConfig{}
	resp, err := cli.ContainerCreate(ctx, config, hostConfig, nil, nil, name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating container: %v\n", err)
		os.Exit(1)
	}
	return resp, nil
}
