package cmd

import "github.com/spf13/cobra"

func createContainer() *cobra.Command {
	var containerName string
	command := &cobra.Command{
		Use:     "create",
		Short:   "Create a container",
		Aliases: []string{"c"},
		Run: func(cmd *cobra.Command, args []string) {
			create(containerName)
		},
	}
	command.Flags().StringVarP(&containerName, "name", "n", "", "container name")
	return command
}

func create(containerName string) {

}
