package cmd

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/enescakir/emoji"
	"github.com/fatih/color"
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
	const format = "%s\t%s\t%s\n"
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
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
		if cont.State != "running" {
			color.Red(format, cont.ID, cont.Names[0], statusEmoji)
		} else {
			fmt.Printf(format, cont.ID, cont.Names[0], statusEmoji)
		}
	}
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
