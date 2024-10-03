package container

import (
	"context"
	"os"
	"strings"

	"distrogo/internal/helpers"
	"distrogo/internal/services/container/state"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/jedib0t/go-pretty/v6/table"
)

type CliService interface {
	GetCLI() (*client.Client, error)
}

type Tty interface {
	SendSignal(signal os.Signal)
}

type Service struct {
	ctx context.Context

	cliService CliService
	attachTty  Tty
}

func New(ctx context.Context, cliService CliService) (*Service, error) {
	return &Service{
		ctx:        ctx,
		cliService: cliService,
	}, nil
}

func (s *Service) IsAttached() bool {
	return s.attachTty != nil
}

func (s *Service) SendSignalToAttach(signal os.Signal) {
	if s.IsAttached() {
		s.attachTty.SendSignal(signal)
	}
}

func (s *Service) RenderTable(containers []types.Container) {
	tableOut := table.NewWriter()
	tableOut.SetOutputMirror(os.Stdout)
	tableOut.SetStyle(table.StyleLight)
	tableOut.AppendHeader(table.Row{"ID", "Name", "ImageNAme", "Status", "Status code"})

	for _, cont := range containers {
		tableOut.AppendRows([]table.Row{
			{cont.ID, cont.Names[0][1:], cont.Image, state.GetEmoji(&cont), cont.Status},
		})
	}
	tableOut.Render()
}

func (s *Service) FilterByLabelKeyExist(containers []types.Container, labelKey string) []types.Container {
	return helpers.SliceFilter(containers, func(container types.Container) bool {
		if _, ok := container.Labels[labelKey]; ok {
			return true
		}
		return false
	})
}

func (s *Service) FilterByLabelValue(containers []types.Container, labelKey, labelValue string) []types.Container {
	return helpers.SliceFilter(containers, func(container types.Container) bool {
		if val, ok := container.Labels[labelKey]; ok && val == labelValue {
			return true
		}
		return false
	})
}

func (s *Service) FilterByName(containers []types.Container, name string) []types.Container {
	return helpers.SliceFilter(containers, func(container types.Container) bool {
		for _, cname := range container.Names {
			if strings.TrimPrefix(cname, "/") == name {
				return true
			}
		}
		return false
	})
}

func (s *Service) FilterByState(containers []types.Container, state string) []types.Container {
	return helpers.SliceFilter(containers, func(container types.Container) bool {
		if container.State == state {
			return true
		}
		return false
	})
}
