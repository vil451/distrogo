package entercontainer

import (
	"fmt"
	"os"

	"distrogo/internal/application"
	"distrogo/internal/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const errEnter = "enter container"

func EnterContainer() *cobra.Command {
	logger.SetLogLevel(logger.LogLevelError)

	var containerName string

	command := &cobra.Command{
		Use:     "enter [container name]",
		Short:   "Enter container",
		Aliases: []string{"e"},
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				if containerName != "" {
					_, err := fmt.Fprintf(os.Stderr, "Error: container name provided in both argument and --name flag\n")
					if err != nil {
						return
					}
					return
				}
				containerName = args[0]
			}

			if containerName == "" {
				_, err := fmt.Fprintf(os.Stderr, "Container name is required\n")
				if err != nil {
					return
				}
				return
			}

			app, err := application.New()
			if err != nil {
				os.Exit(1)
			}
			defer app.Terminate()

			containerSvc, err := app.GetContainerService()
			if err != nil {
				os.Exit(1)
			}

			err = containerSvc.Enter(containerName)
			if err != nil {
				logger.Error(errors.Wrap(err, errEnter))
			}
		},
	}

	command.Flags().StringVarP(
		&containerName,
		"name",
		"n",
		containerName,
		"container name",
	)

	return command
}
