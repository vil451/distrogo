package cmd

import "github.com/spf13/cobra"

func listContainer() *cobra.Command {

	command := cobra.Command{
		Use:   "list",
		Short: "List containers",
		Run:   list(),
	}
	return &command
}

func list() func(cmd *cobra.Command, args []string) {
	return nil
}
