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
	var status string
	command := &cobra.Command{
		Use:     "list",
		Short:   "List containers",
		Aliases: []string{"ps", "ls"},
		Run: func(cmd *cobra.Command, args []string) {
			if all {
				listContainers(true, "", "")
			} else {
				listContainers(false, containerName, status)
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

	// флаг для фильтрации по статусу
	command.Flags().StringVarP(
		&status,
		"status",
		"s",
		"",
		"Status of the containers to list",
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
	tableOut.AppendHeader(table.Row{"ID", "Name", "ImageNAme", "Status", "Status code"})

	for _, cont := range containers {
		statusEmoji := getStatusEmoji(cont.State)
		tableOut.AppendRows([]table.Row{
			{cont.ID, cont.Names[0], cont.Image, statusEmoji, cont.Status},
		})
	}
	tableOut.Render()
}

// функция listContainers выводит список контейнеров с учетом фильтров
func listContainers(all bool, containerName string, status string) {
	ctx := context.Background()
	cli, err := dockerclient.InitDockerClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating Docker client: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := dockerclient.CloseDockerClient(cli); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing Docker client: %v\n", err)
		}
	}()

	containers, err := getContainers(ctx, cli, all)
	containers = filterContainersByLabel(containers, "manager", "distrogo")
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing containers: %v\n", err)
		os.Exit(1)
	}

	// Если all == false, фильтруем контейнеры по имени и метке "manager=distrogo"
	if !all {
		containers = filterContainersByName(containers, containerName)

	}

	// Если статус указан, фильтруем контейнеры по статусу
	if status != "" {
		containers = filterContainersByStatus(containers, status)
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

func filterContainersByLabel(containers []types.Container, labelKey string, labelValue string) []types.Container {
	var filtered []types.Container
	for _, container := range containers {
		if val, ok := container.Labels[labelKey]; ok && val == labelValue {
			filtered = append(filtered, container)
		}
	}
	return filtered
}

// функция для фильтрации контейнеров по статусу
func filterContainersByStatus(containers []types.Container, status string) []types.Container {
	var filteredContainers []types.Container
	for _, cont := range containers {
		if cont.State == status {
			filteredContainers = append(filteredContainers, cont)
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
