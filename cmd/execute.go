// Package cmd wires the cobra CLI to the rm removal engine. Each flag lives in
// its own sub-package (cmd/force, cmd/interactive, …); this file registers them
// and folds their values into a single core.Options for one engine run.
package cmd

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/braswelljr/rmx/cmd/directory"
	"github.com/braswelljr/rmx/cmd/force"
	"github.com/braswelljr/rmx/cmd/interactive"
	"github.com/braswelljr/rmx/cmd/onefilesystem"
	"github.com/braswelljr/rmx/cmd/preserveroot"
	"github.com/braswelljr/rmx/cmd/recursive"
	"github.com/braswelljr/rmx/cmd/verbose"
	"github.com/braswelljr/rmx/cmd/version"
	"github.com/braswelljr/rmx/core"
)

// Register installs every flag onto fs. main wires this onto the root command.
func Register(fs *pflag.FlagSet) {
	force.Register(fs)
	interactive.Register(fs)
	recursive.Register(fs)
	directory.Register(fs)
	verbose.Register(fs)
	onefilesystem.Register(fs)
	preserveroot.Register(fs)
	version.Register(fs)
}

// Execute builds the removal options from the parsed flags and runs the engine
// against args. It returns core.ErrFailed when a removal failed (diagnostics are
// already on errw) or a usage error otherwise.
func Execute(command *cobra.Command, args []string, in io.Reader, out, errw io.Writer) error {
	fs := command.Flags()

	if version.Requested(fs) {
		version.Print(out)
		return nil
	}

	opts, err := buildOptions(fs)
	if err != nil {
		return err
	}

	return core.New(opts, in, out, errw).Run(args)
}

// buildOptions folds every flag package's value into a single core.Options. Order
// matters: force sets the baseline prompt mode, then interactive may raise it
// (see the package docs for the precedence rules).
func buildOptions(fs *pflag.FlagSet) (core.Options, error) {
	var o core.Options

	force.Apply(fs, &o)
	if err := interactive.Apply(fs, &o); err != nil {
		return core.Options{}, err
	}
	recursive.Apply(fs, &o)
	directory.Apply(fs, &o)
	verbose.Apply(fs, &o)
	onefilesystem.Apply(fs, &o)
	preserveroot.Apply(fs, &o)

	return o, nil
}
