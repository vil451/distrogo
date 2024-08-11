package cmd

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/enescakir/emoji"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"os"
)

func listContainer() *cobra.Command {
	var containerName string

	command := &cobra.Command{
		Use:     "list",
		Short:   "List containers",
		Aliases: []string{"ps", "ls"},
		Run: func(cmd *cobra.Command, args []string) {
			listContainers(containerName)
		},
	}

	command.Flags().StringVarP(&containerName, "name", "n", "", "Name of the container to list")
	return command
}

func listContainers(name string) {
	ctx := context.Background()
	tableOut := table.NewWriter()
	tableOut.SetOutputMirror(os.Stdout)
	tableOut.SetStyle(table.StyleLight)
	tableOut.AppendHeader(table.Row{"ID", "Name", "Status"})
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Error creating Docker client: %v\n", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
	defer func(cli *client.Client) {
		err := cli.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error closing client: %v\n", err)
		}
	}(cli)

	containers, err := cli.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, cont := range containers {
		statusEmoji := getStatusEmoji(cont.State)
		tableOut.AppendRows([]table.Row{
			{cont.ID, cont.Names[0], statusEmoji},
		})
	}
	tableOut.Render()
}

func getStatusEmoji(state string) string {
	switch state {
	case "running":
		return emoji.OkHand.Tone(emoji.Light)
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
