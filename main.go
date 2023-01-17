package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/braswelljr/rmx/cmd"
	"github.com/braswelljr/rmx/rm"
)

var (
	Rmx = rm.Rm{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
	}
)

func main() {

	// get root command
	command := &cobra.Command{
		Use:   "rmx " + color.MagentaString("[flags]...") + " " + color.YellowString("[directory / file]..."),
		Short: "A cross-platform replacement for UNIX " + color.YellowString("rm") + " command.",
		Long:  "A cross-platform replacement for UNIX " + color.YellowString("rm") + " command.",
		Args:  cmd.ArgumentValidator(&Rmx),
		RunE: func(command *cobra.Command, args []string) error {
			Rmx.Directory = "."
			if len(args) > 0 {
				Rmx.Directory = args[0]

				// run command and return any error
				return cmd.Run(&Rmx, command, args)
			}

			return nil
		},
	}

	// add flags
	command.PersistentFlags().BoolP("help", "h", false, "help for this command")
	command.Flags().BoolVarP(&Rmx.Flags.F, "force", "f", false, "ignore nonexistent files and arguments, never prompt")
	command.Flags().BoolVarP(&Rmx.Flags.Ii, "interactive", "i", false, "prompt before every removal")
	command.Flags().BoolVarP(&Rmx.Flags.II, "INTERACTIVE", "I", false, "prompt once before removing more than three files, or when removing recursively; less intrusive than -i, while still giving protection against most mistakes")
	command.Flags().BoolVarP(&Rmx.Flags.Rr, "recursive", "r", false, "remove directories and their contents recursively")
	command.Flags().BoolVarP(&Rmx.Flags.RR, "RECURSIVE", "R", false, "remove directories and their contents recursively")
	command.Flags().BoolVarP(&Rmx.Flags.D, "dir", "d", false, "remove empty directories")
	command.PersistentFlags().BoolVarP(&Rmx.Flags.V, "verbose", "v", false, "explain what is being done")
	command.Flags().Float32Var(&Rmx.Flags.Version, "version", 0.01, "output version information and exit")

	// clean up and exit on error
	if err := command.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err) // print error to stderr
		os.Exit(1)                   // exit on error
	}
}
