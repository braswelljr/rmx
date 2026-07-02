package cmd

import (
	"io"

	"github.com/spf13/cobra"

	"github.com/braswelljr/rmx/internal/common"
)

// NewRoot builds the fully-configured rmx root command with every flag
// registered and its removal engine wired to the given streams. main and the
// docs generator both build the command through here so they stay in sync.
func NewRoot(in io.Reader, out, errw io.Writer) *cobra.Command {
	root := &cobra.Command{
		Use:           common.AppName + " [OPTION]... [FILE]...",
		Short:         "Remove (unlink) files and directories — a cross-platform drop-in for rm.",
		Long:          "rmx removes each FILE. By default it does not remove directories; use -r to\nremove a directory and everything under it, or -d to remove an empty directory.",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(command *cobra.Command, args []string) error {
			return Execute(command, args, in, out, errw)
		},
	}

	// Each flag is owned by its own package under cmd/; Register installs them all.
	Register(root.Flags())

	return root
}
