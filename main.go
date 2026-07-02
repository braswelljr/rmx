package main

//go:generate go run ./tools/gendocs docs

import (
	"errors"
	"fmt"
	"os"

	"github.com/braswelljr/rmx/cmd"
	"github.com/braswelljr/rmx/core"
	"github.com/braswelljr/rmx/internal/common"
)

func main() {
	root := cmd.NewRoot(os.Stdin, os.Stdout, os.Stderr)

	if err := root.Execute(); err != nil {
		// core.ErrFailed means diagnostics were already printed; anything else is
		// a usage/parse error we surface here.
		if !errors.Is(err, core.ErrFailed) {
			fmt.Fprintf(os.Stderr, "%s: %v\n", common.AppName, err)
		}
		os.Exit(1)
	}
}
