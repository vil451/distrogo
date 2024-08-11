package cmd

import "github.com/spf13/cobra"

func stopContainer() *cobra.Command {
	var containerName string

	command := &cobra.Command{
		Use:     "stop",
		Short:   "Stop a container",
		Aliases: []string{"s"},
		Run: func(cmd *cobra.Command, args []string) {
			stop(containerName)
		},
	}

	command.Flags().StringVarP(&containerName, "name", "n", "", "container name")
	return command
}

func stop(containerName string) {

}
