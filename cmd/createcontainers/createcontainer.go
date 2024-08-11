package createcontainers

import "github.com/spf13/cobra"

func CreateContainer() *cobra.Command {
	var containerName string
	var imageName string
	var pullImage bool
	command := &cobra.Command{
		Use:     "create",
		Short:   "Create a container",
		Aliases: []string{"c"},
		Run: func(cmd *cobra.Command, args []string) {
			create(containerName)
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
		"",
		"container name",
	)

	command.Flags().BoolVarP(
		&pullImage,
		"pull",
		"p",
		true,
		"pull image",
	)
	return command
}

func create(containerName string) {

}
