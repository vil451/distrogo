package entercontainer

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
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
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating Docker client: %v\n", err)
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
