package entercontainer

import (
	"bufio"
	"context"
	"distrogo/cmd/listcontainers"
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
				containerName = args[0]
			}

			if containerName == "" {
				fmt.Println("Container name is required")
				return
			}

			enter(containerName)
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

func enter(containerName string) {
	if containerName == "" {
		fmt.Println("Container name is required")
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating Docker client: %v\n", err)
		os.Exit(1)
	}
	defer cli.Close()

	err = runContainer(ctx, cli, containerName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running container: %v\n", err)
		os.Exit(1)
	}

	execConfig := types.ExecConfig{
		Cmd:          strslice.StrSlice([]string{"/bin/sh"}),
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

	attachResp, err := cli.ContainerExecAttach(ctx, execIDResp.ID, types.ExecStartCheck{Tty: true})
	if err != nil {
		fmt.Println("Error attaching to exec session:", err)
		os.Exit(1)
	}
	defer attachResp.Close()

	// Канал для завершения горутин
	done := make(chan struct{})
	var once sync.Once
	var wg sync.WaitGroup

	// Канал для обработки системных сигналов
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Горутина для передачи вывода контейнера в stdout
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := io.Copy(os.Stdout, attachResp.Reader)
		if err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "Error reading output: %v\n", err)
		}
		// Завершение работы
		once.Do(func() {
			close(done)
		})
	}()

	// Обработка пользовательского ввода
	reader := bufio.NewReader(os.Stdin)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			default:
				input, err := reader.ReadString('\n')
				if err == io.EOF {
					fmt.Println("\nExiting container session (Ctrl+D)...")
					once.Do(func() {
						close(done)
					})
					return
				}
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
					continue
				}

				trimmedInput := strings.TrimSpace(input)
				if trimmedInput == "exit" {
					fmt.Println("Exiting container session (exit)...")
					once.Do(func() {
						close(done)
					})
					return
				}

				// Обработка записи в контейнер
				_, err = io.WriteString(attachResp.Conn, input)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error writing input: %v\n", err)
					// Завершение по ошибке записи
					once.Do(func() {
						close(done)
					})
					return
				}
			}
		}
	}()
	// Ожидание завершения программы по системному сигналу или пользовательскому вводу
	select {
	case <-done:
		cancel()
	case sig := <-sigs:
		fmt.Printf("Received signal: %v. Exiting...\n", sig)
		once.Do(func() {
			close(done)
		})
		cancel()
	}

	// Ожидание завершения всех горутин
	wg.Wait()
	fmt.Println("Session terminated.")
}

func runContainer(ctx context.Context, cli *client.Client, containerName string) error {
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
