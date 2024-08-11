package cmd

import "github.com/spf13/cobra"

func assembleConatiner() *cobra.Command {
	var containerName string

	command := &cobra.Command{
		Use:     "assemble",
		Short:   "Assemble a container",
		Aliases: []string{"a"},
		Run: func(cmd *cobra.Command, args []string) {
			assemble(containerName)
		},
	}

	command.Flags().StringVarP(&containerName, "container", "c", "", "Container name")
	return command

}

func assemble(containerName string) {

}
