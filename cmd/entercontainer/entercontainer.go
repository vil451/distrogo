package entercontainer

import (
	"context"
	"distrogo/cmd/dockerclient"
	"distrogo/cmd/listcontainers"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"os"
)

func EnterContainer() *cobra.Command {
	var containerName string

	command := &cobra.Command{
		Use:     "enter [container name]",
		Short:   "Enter container",
		Aliases: []string{"e"},
		Args:    cobra.MaximumNArgs(1), // Ожидаем максимум 1 позиционный аргумент
		Run: func(cmd *cobra.Command, args []string) {
			// Если имя контейнера передано как позиционный аргумент, использовать его
			if len(args) > 0 {
				containerName = args[0]
			}

			if containerName == "" {
				fmt.Println("Container name is required")
				return
			}

			enter(containerName)
		},
	}

	// Добавляем флаг для указания контейнера
	command.Flags().StringVarP(
		&containerName,
		"container",
		"c",
		"",
		"container name",
	)

	return command
}

func enter(containerName string) {
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

	err = runContainer(ctx, cli, containerName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error run container: %v\n", err)
		os.Exit(1)
	}

	defer cli.Close()

	execConfig := types.ExecConfig{
		User:         "root",
		Cmd:          strslice.StrSlice([]string{"/bin/bash"}),
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
	}

	execIDResp, err := cli.ContainerExecCreate(ctx, containerName, execConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating exec instance: %v\n", err)
		os.Exit(1)
	}

	// Запуск интерактивного сеанса внутри контейнера
	err = cli.ContainerExecStart(ctx, execIDResp.ID, types.ExecStartCheck{Tty: true})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting exec instance: %v\n", err)
		os.Exit(1)
	}
}

func runContainer(ctx context.Context, cli *client.Client, containerName string) error {
	containers, err := listcontainers.GetContainers(ctx, cli, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating exec instance: %v\n", err)
		os.Exit(1)
	}

	containers = listcontainers.FilterContainersByLabel(containers, "manager", "distrogo")
	result_container_id := ""
	state := ""
	for _, container := range containers {
		if container.Names[0][1:] == containerName {
			result_container_id = container.ID
			state = container.State
		}
	}
	if state == "running" {
		return nil
	}
	if result_container_id == "" {
		fmt.Fprintf(os.Stderr, "container not found: %s", containerName)
		os.Exit(1)
	}

	// Опции для запуска контейнера
	startOptions := container.StartOptions{}

	// Запуск контейнера
	if err := cli.ContainerStart(ctx, result_container_id, startOptions); err != nil {
		return fmt.Errorf("error starting container: %v", err)
	}

	fmt.Printf("Container %s is started with ID: %s\n", containerName, result_container_id)
	return nil
}
