package entercontainer

import (
	"bufio"
	"context"
	"distrogo/cmd/listcontainers"
	"distrogo/internal/cancelable_reader"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"io"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

func EnterContainer() *cobra.Command {
	var containerName string

	command := &cobra.Command{
		Use:     "enter [container name]",
		Short:   "Enter container",
		Aliases: []string{"e"},
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				if containerName != "" {
					_, err := fmt.Fprintf(os.Stderr, "Error: container name provided in both argument and --name flag\n")
					if err != nil {
						return
					}
					return
				}
				containerName = args[0]
			}

			if containerName == "" {
				_, err := fmt.Fprintf(os.Stderr, "Container name is required\n")
				if err != nil {
					return
				}
				return
			}

			enter(containerName)
		},
	}

	command.Flags().StringVarP(
		&containerName,
		"name",
		"n",
		containerName,
		"container name",
	)

	return command
}

func enter(containerName string) {
	if containerName == "" {
		fmt.Println("Container name is required")
		os.Exit(1)
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Printf("Error creating Docker client: %v\n", err)
		os.Exit(1)
	}
	defer cli.Close()

	err = runContainer(cli, containerName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running container: %v\n", err)
		return
	}

	attachResp, ctx, ctxCancel, err := attachToContainer(cli, containerName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error attaching to container: %v\n", err)
		return
	}

	ttyReader := cancelable_reader.New(ctx, attachResp.Reader)

	detach := func() {
		ctxCancel()
		attachResp.Close()
	}

	var wg sync.WaitGroup

	// Горутина для передачи вывода контейнера в stdout
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				data := make([]byte, 1024)
				_, errRead := ttyReader.Read(data)
				if errRead != nil {
					detach()
				}
				fmt.Printf("%s", string(data))
			}
		}
	}()

	// Обработка пользовательского ввода
	reader := bufio.NewReader(os.Stdin)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				input, errRead := reader.ReadString('\n')
				if errRead != nil {
					if errors.Is(errRead, io.EOF) {
						fmt.Println("\nExiting container session (Ctrl+D)...")
					} else {
						fmt.Fprintf(os.Stderr, "Error reading input: %v\n", errRead)
					}
					detach()
				}

				trimmedInput := strings.TrimSpace(input)
				if trimmedInput == "detach" {
					fmt.Println("Exiting container session (detach)...")
					detach()
				}

				// Обработка записи в контейнер
				_, err = io.WriteString(attachResp.Conn, input)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error writing input: %v\n", err)
					detach()
				}
			}
		}
	}()

	// Канал для обработки системных сигналов
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			return
		case sig := <-sigs:
			fmt.Printf("Received signal: %v. Exiting...\n", sig)
			detach()
		}
	}()

	wg.Wait()
	fmt.Println("Session terminated.")
}

func runContainer(cli *client.Client, containerName string) error {
	ctx := context.Background()
	containers, err := listcontainers.GetContainers(ctx, cli, true)
	if err != nil {
		return fmt.Errorf("Error listing containers: %v", err)
	}

	containers = listcontainers.FilterContainersByLabel(containers, "manager", "distrogo")
	var resultContainerID, state string
	for _, container := range containers {
		if container.Names[0][1:] == containerName {
			resultContainerID = container.ID
			state = container.State
		}
	}
	if state == "running" {
		return nil
	}
	if resultContainerID == "" {
		return fmt.Errorf("container not found: %s", containerName)
	}

	startOptions := container.StartOptions{}
	if err := cli.ContainerStart(ctx, resultContainerID, startOptions); err != nil {
		return fmt.Errorf("Error starting container: %v", err)
	}

	fmt.Printf("Container %s is started with ID: %s\n", containerName, resultContainerID)
	return nil
}

func attachToContainer(cli *client.Client, containerName string) (*types.HijackedResponse, context.Context, context.CancelFunc, error) {
	ctx, cancel := context.WithCancel(context.Background())

	execConfig := types.ExecConfig{
		Cmd:          strslice.StrSlice([]string{"/bin/sh"}),
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
	}

	execIDResp, err := cli.ContainerExecCreate(ctx, containerName, execConfig)
	if err != nil {
		cancel()
		return nil, nil, nil, fmt.Errorf("Error creating exec instance: %v\n", err)
	}

	attachResp, err := cli.ContainerExecAttach(ctx, execIDResp.ID, types.ExecStartCheck{Tty: true})
	if err != nil {
		cancel()
		return nil, nil, nil, fmt.Errorf("Error attaching to exec session:%v\n", err)
	}

	return &attachResp, ctx, cancel, nil
}
