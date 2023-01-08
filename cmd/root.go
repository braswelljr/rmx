package commands

import (
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/braswelljr/rmx/rm"
)

var Rmc = rm.RM{
	Stdin:  os.Stdin,
	Stdout: os.Stdout,
}

func RootCommand() *cobra.Command {
	// initialize cobra
	command := &cobra.Command{
		Use:   "rmx " + color.MagentaString("[flags]...") + " " + color.YellowString("[directory / file]..."),
		Short: "A cross-platform replacement for UNIX " + color.YellowString("rm") + " command.",
		Long:  "A cross-platform replacement for UNIX " + color.YellowString("rm") + " command.",
		Args:  ArgumentValidator(&Rmc),
		RunE: func(command *cobra.Command, args []string) error {
			Rmc.Directory = "."
			if len(args) > 0 {
				Rmc.Directory = args[0]

				// run command and return any error
				return Run(&Rmc, command, args)
			}

			return nil
		},
	}

	return command
}
