package cmd

import (
	"distrogo/internal/initContainer"
	"github.com/spf13/cobra"
)

var (
	initCommandFlags *initContainer.Flags
)

func init() {

}
func initContainerFlags() *cobra.Command {
	initCommandFlags = initContainer.NewFlags()
	command := cobra.Command{
		Use:     "init",
		Short:   "Init Distrogo",
		Aliases: []string{"i"},
		Run: func(cmd *cobra.Command, args []string) {
			initContainerArgs(cmd, args)
		},
	}
	return &command
}

func initContainerArgs(command *cobra.Command, args []string) {
	command.Flags().StringVarP(
		initCommandFlags.Name,
		"name", "n",
		"",
		"User name",
	)
	command.Flags()
}
