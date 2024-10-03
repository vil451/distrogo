package image

import (
	"context"
	"distrogo/cmd/dockerclient"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"os"
	"path/filepath"
)

func (image *Image) Create(imageName *Image) error {
	ctx := context.Background()
	cli, err := dockerclient.InitDockerClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating Docker client: %v\n", err)
		os.Exit(1)
	}

	defer func() {
		if err := dockerclient.CloseDockerClient(cli); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close Docker client: %v\n", err)
		}
	}()

	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		homeDir = filepath.Join("/home", os.Getenv("USER"))
	}

	labels := map[string]string{
		"manager": "distrogo",
	}

	volumes := map[string]struct{}{
		"/tmp":  {},
		homeDir: {},
	}

	config := &container.Config{
		Image:     image.imageName,
		Labels:    labels,
		Cmd:       []string{"/bin/sh"},
		Tty:       true,
		OpenStdin: true,
		Volumes:   volumes,
	}

	hostConfig := &container.HostConfig{
		Binds: []string{
			"/tmp:/tmp:rslave",
			homeDir + ":" + homeDir + ":rslave",
		},
	}

	_, err = cli.ContainerCreate(ctx, config, hostConfig, nil, nil, image.containerName)
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Error creating container: %v\n", err)
		if err != nil {
			return err
		}
		os.Exit(1)
	}

	return nil
}
