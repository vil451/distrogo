package cmd

import "github.com/spf13/cobra"

func enterContainer() *cobra.Command {
	var containerName string
	command := &cobra.Command{
		Use:     "enter",
		Short:   "Enter a container",
		Aliases: []string{"e", "ent"},
		Run: func(cmd *cobra.Command, args []string) {
			enter(containerName)
		},
	}
	command.Flags().StringVarP(&containerName, "name", "n", "", "container name")
	return command
}

func enter(name string) {

}
