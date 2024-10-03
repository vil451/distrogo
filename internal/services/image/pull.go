package image

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"io"
	"os"
)

func (image *Image) PullImage(imageName string) error {
	config := &container.Config{
		Image: imageName,
	}
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	resp, err := cli.ImagePull(ctx, config.Image, types.ImagePullOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error pulling image: %v\n", err)
		return err
	}
	defer resp.Close()

	_, err = io.Copy(os.Stdout, resp)
	if err != nil {
		return err
	}

	return nil
}
