package cmd

import (
	"distrogo/internal/config"
	"distrogo/internal/config/data"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
		RunE:  run,
	}
)

func run(cmd *cobra.Command, args []string) error {

	if err := config.InitLogLocs(); err != nil {
		return err
	}
	file, err := os.OpenFile(
		*cmdFlags.LogFile,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		data.DefaultFileMod,
	)

	if err != nil {
		return fmt.Errorf("failed %q", *cmdFlags.LogFile, err)
	}

	if err != nil {
		if file != nil {
			_ = file.Close()
		}
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: file})

	return nil
}

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
			panic(err)
		}
	}
}
