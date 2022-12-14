package main

import (
	"os"

	"github.com/braswelljr/rmx/rm"

	"github.com/spf13/cobra"
)

func main() {
	// get the root command
	rmc := rm.RM{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
	}

	// initialize cobra
	command := &cobra.Command{
		Use:   "rmx [OPTION]... [FILE]...",
		Short: "A cross-platform replacement for UNIX rm command",
		Long:  "A cross-platform replacement for UNIX rm command",
		Args:  rm.ArgumentValidator(&rmc),
		RunE: func(command *cobra.Command, args []string) error {
			rmc.Directory = "."
			if len(args) > 0 {
				rmc.Directory = args[0]

				// run command and return any error
				return rmc.Run(command, args)
			}

			return nil
		},
	}

	// add flags
	command.PersistentFlags().BoolP("help", "h", false, "help for this command")
	command.Flags().BoolVarP(&rmc.Flags.F, "force", "f", false, "ignore nonexistent files and arguments, never prompt")
	command.Flags().BoolVarP(&rmc.Flags.Ii, "interactive", "i", false, "prompt before every removal")
	//command.Flags().BoolVarP(&rmc.Flags.II, "interactive", "I", false, "prompt once before removing more than three files, or when removing recursively; less intrusive than -i, while still giving protection against most mistakes")
	command.Flags().BoolVarP(&rmc.Flags.Rr, "recursive", "r", false, "remove directories and their contents recursively")
	//command.Flags().BoolVarP(&rmc.Flags.RR, "recursive", "R", false, "remove directories and their contents recursively")
	command.Flags().BoolVarP(&rmc.Flags.D, "dir", "d", false, "remove empty directories")
	command.Flags().BoolVarP(&rmc.Flags.V, "verbose", "v", false, "explain what is being done")
	command.Flags().Float32Var(&rmc.Flags.Version, "version", 0.01, "output version information and exit")

	// clean up and exit on error
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
