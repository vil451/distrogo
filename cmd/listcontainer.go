package cmd

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"os"
)

func listContainer() *cobra.Command {

	command := cobra.Command{
		Use: "list",
		Run: func(cmd *cobra.Command, args []string) {
			listContainers()
		},
	}
	return &command
}

func listContainers() {
	ctx := context.Background()
	const format = "%s\n"
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer func(cli *client.Client) {
		err := cli.Close()
		if err != nil {

		}
	}(cli)

	containers, err := cli.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, cont := range containers {
		fmt.Fprintf(os.Stderr, format, cont.ID)
	}
}
