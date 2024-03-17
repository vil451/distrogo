package config

import "github.com/spf13/cobra"

type Flags struct {
	LogFile    *string
	SubCommand *string
}

type SubCommands struct {
	GetContainers *cobra.Command
}

var (
	AppLogFile    string
	GetSubCommand cobra.Command
)

func NewFlags() *Flags {
	return &Flags{
		LogFile: strPtr(AppLogFile),
	}
}

func NewSubCommands() *SubCommands {
	return &SubCommands{
		GetContainers: subCmdPtr(GetSubCommand),
	}
}

func strPtr(s string) *string {
	return &s
}

func subCmdPtr(c cobra.Command) *cobra.Command {
	return &c
}
