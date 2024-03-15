package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func versionCmd() *cobra.Command {
	var short bool
	command := cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			printVersion(short)
		},
	}
	command.PersistentFlags().BoolVarP(&short, "short", "s", false, "Print version")

	return &command
}

func printVersion(short bool) {
	const format = "%-20s %s\n"
	fmt.Fprintf(os.Stderr, format, "Version", "dev")
}
