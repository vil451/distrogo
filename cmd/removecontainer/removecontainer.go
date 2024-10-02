package removecontainer

import (
	"context"
	"distrogo/cmd/dockerclient"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"os"
)

func RemoveContainer() *cobra.Command {
	var containerName string

	command := &cobra.Command{
		Use:   "remove [container name]",
		Short: "Remove a Docker container",
		Long:  `This command removes a Docker container`,
		Args:  cobra.MaximumNArgs(1), // Ожидаем максимум 1 позиционный аргумент
		Run: func(cmd *cobra.Command, args []string) {
			// Если имя контейнера передано как позиционный аргумент, использовать его
			if len(args) > 0 {
				containerName = args[0]
			}

			if containerName == "" {
				fmt.Println("Container name is required")
				return
			}

			remove(containerName)
		},
	}

	command.Flags().StringVarP(
		&containerName,
		"container",
		"c",
		"",
		"container name",
	)

	return command
}

func remove(containerName string) {
	if containerName == "" {
		fmt.Println("Container name is required")
		os.Exit(1)
	}

	ctx := context.Background()
	cli, err := dockerclient.InitDockerClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating Docker client: %v\n", err)
		os.Exit(1)
	}

	err = removeContainer(ctx, cli, containerName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error removing container: %v\n", err)
		os.Exit(1)
	}

	defer func() {
		if err := dockerclient.CloseDockerClient(cli); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing Docker client: %v\n", err)
		}
	}()
}

func removeContainer(ctx context.Context, cli *client.Client, containerName string) error {
	removeOptions := types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}

	if err := cli.ContainerRemove(ctx, containerName, removeOptions); err != nil {
		return fmt.Errorf("error removing container: %v", err)
	}

	fmt.Printf("Container %s is removed successfully\n", containerName)
	return nil
}
