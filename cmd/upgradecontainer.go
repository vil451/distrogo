package cmd

import "github.com/spf13/cobra"

func upgradeContainer() *cobra.Command {
	var containerName string
	command := &cobra.Command{
		Use:     "upgrade",
		Short:   "Upgrade a container",
		Aliases: []string{"up"},
		Run: func(cmd *cobra.Command, args []string) {
			upgrade(containerName)
		},
	}

	command.Flags().StringVarP(&containerName, "container", "c", "", "Container name")
	return command
}

func upgrade(containerName string) {

}
