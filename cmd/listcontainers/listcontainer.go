package listcontainers

import (
	"os"

	"distrogo/internal/application"
	"distrogo/internal/logger"
	"distrogo/internal/services/container/labels"
	"github.com/spf13/cobra"
)

func ListContainers() *cobra.Command {
	var containerName string
	var all bool
	var state string
	command := &cobra.Command{
		Use:     "list",
		Short:   "List containers",
		Aliases: []string{"ps", "ls"},
		Run: func(cmd *cobra.Command, args []string) {
			app, err := application.New()
			if err != nil {
				os.Exit(1)
			}
			defer app.Terminate()

			containerSvc, err := app.GetContainerService()
			if err != nil {
				os.Exit(1)
			}

			containers, err := containerSvc.List(all)
			if err != nil {
				logger.Error(err)
			}
			containers = containerSvc.FilterByLabelValue(containers, labels.LabelManager, labels.LabelValueDistrogo)

			if containerName != "" {
				containers = containerSvc.FilterByName(containers, containerName)
			}

			if state != "" {
				containers = containerSvc.FilterByState(containers, state)
			}

			containerSvc.RenderTable(containers)
		},
	}

	command.Flags().StringVarP(
		&containerName,
		"name",
		"n",
		"",
		"Name of the container to list",
	)

	command.Flags().BoolVarP(
		&all,
		"all",
		"a",
		false,
		"List all containers",
	)

	// флаг для фильтрации по статусу
	command.Flags().StringVarP(
		&state,
		"state",
		"s",
		"",
		"State of the containers to list",
	)
	return command
}
