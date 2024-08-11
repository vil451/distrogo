package listcontainers

import (
	"context"
	"distrogo/cmd/dockerclient"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/enescakir/emoji"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"os"
)

func ListContainers() *cobra.Command {
	var containerName string
	var all bool
	command := &cobra.Command{
		Use:     "list",
		Short:   "List containers",
		Aliases: []string{"ps", "ls"},
		Run: func(cmd *cobra.Command, args []string) {
			if all {
				listContainers(true, "")
			} else {
				listContainers(false, containerName)
			}
		},
	}

	command.Flags().StringVarP(
		&containerName,
		"name",
		"n",
		"",
		"Name of the container to list",
	)

	command.Flags().BoolVarP(
		&all,
		"all",
		"a",
		false,
		"List all containers",
	)
	return command
}

func getContainers(ctx context.Context, cli *client.Client, all bool) ([]types.Container, error) {
	options := types.ContainerListOptions{}
	if all {
		options.All = true
	}
	containers, err := cli.ContainerList(ctx, options)
	if err != nil {
		return nil, err
	}
	return containers, nil
}

func renderTable(containers []types.Container) {
	tableOut := table.NewWriter()
	tableOut.SetOutputMirror(os.Stdout)
	tableOut.SetStyle(table.StyleLight)
	tableOut.AppendHeader(table.Row{"ID", "Name", "Image", "Status", "Status code"})

	for _, cont := range containers {
		statusEmoji := getStatusEmoji(cont.State)
		tableOut.AppendRows([]table.Row{
			{cont.ID, cont.Names[0], cont.ImageID, statusEmoji, cont.Status},
		})
	}
	tableOut.Render()
}

func listContainers(all bool, containerName string) {
	ctx := context.Background()
	cli, err := dockerclient.InitDockerClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating Docker client: %v\n", err)
		os.Exit(1)
	}
	defer dockerclient.CloseDockerClient(cli)

	containers, err := getContainers(ctx, cli, all)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing containers: %v\n", err)
		os.Exit(1)
	}

	if !all {
		containers = filterContainersByName(containers, containerName)
	}

	renderTable(containers)
}

func filterContainersByName(containers []types.Container, name string) []types.Container {
	if name == "" {
		return containers
	}

	var filteredContainers []types.Container
	for _, cont := range containers {
		for _, cname := range cont.Names {
			if cname == name {
				filteredContainers = append(filteredContainers, cont)
				break
			}
		}
	}
	return filteredContainers
}

func getStatusEmoji(state string) string {
	switch state {
	case "running":
		return emoji.OkHand.String()
	case "exited":
		return emoji.Door.String()
	case "paused":
		return emoji.PauseButton.String()
	case "restarting":
		return emoji.ClockwiseVerticalArrows.String()
	default:
		return emoji.QuestionMark.String()
	}
}
