package cmd

import (
	"distrogo/internal/config"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

const (
	appName      = config.AppName
	shortAppDesc = "Distro Tools"
)

var (
	cmdFlags *config.Flags
	rootCmd  = &cobra.Command{
		Use:   appName,
		Short: shortAppDesc,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
			}
		},
	}
)

type flagError struct {
	err error
}

func (e flagError) Error() string {
	return e.err.Error()
}

func init() {
	if err := config.InitLogLocs(); err != nil {
		fmt.Printf("fail initContainer logs location %s\n", err)
	}
	rootCmd.SetFlagErrorFunc(func(command *cobra.Command, err error) error {
		return flagError{err: err}
	})
	rootCmd.AddCommand(versionCmd())
	rootCmd.AddCommand(listContainer())
	rootCmd.AddCommand(initContainerFlags())
	initFlags()
}

func initFlags() {
	cmdFlags = config.NewFlags()
	rootCmd.Flags().StringVarP(
		cmdFlags.LogFile,
		"logsFile", "l",
		config.AppLogFile,
		"Specify the log file",
	)
	rootCmd.Flags()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if !errors.As(err, &flagError{}) {
			fmt.Printf("Execution error: %s\n", err)
			os.Exit(1)
		}
	}
}
