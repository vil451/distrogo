package createcontainers

import (
	iamgeService "distrogo/internal/services/image"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func CreateContainer() *cobra.Command {
	var containerName string
	var imageName string
	var pullImage bool
	command := &cobra.Command{
		Use:     "create [container name]",
		Short:   "Create a container",
		Aliases: []string{"c"},
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
			image, err := iamgeService.New(imageName, containerName)

			if err != nil {
				os.Exit(1)
			}
			if pullImage {
				err := image.PullImage(containerName)
				if err != nil {
					return
				}
			}
			err = image.Create(image)
			if err != nil {
				return
			}
		},
	}

	command.Flags().StringVarP(
		&imageName,
		"image",
		"i",
		"",
		"image name of a container",
	)

	command.Flags().StringVarP(
		&containerName,
		"name",
		"n",
		containerName,
		"container name",
	)

	command.Flags().BoolVarP(
		&pullImage,
		"pull",
		"p",
		false,
		"pull image",
	)
	return command
}
