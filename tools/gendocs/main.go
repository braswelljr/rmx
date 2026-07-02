// Command gendocs generates the rmx CLI reference in every cobra format.
//
// Usage:
//
//	go run ./tools/gendocs [dir] [format]
//
//	dir      output base directory (default "docs")
//	format   one of: all (default), markdown, man, yaml, rest, guide, flags
//
// It is also wired as a `go generate` step (see main.go) and `make docs`.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/braswelljr/rmx/cmd"
	"github.com/braswelljr/rmx/cmd/docs"
)

func main() {
	dir := "docs"
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	format := "all"
	if len(os.Args) > 2 {
		format = os.Args[2]
	}

	root := cmd.NewRoot(os.Stdin, os.Stdout, os.Stderr)

	if err := generate(root, dir, format); err != nil {
		log.Fatalf("gendocs: %v", err)
	}

	log.Printf("gendocs: wrote %s reference to %s/", format, dir)
}

func generate(root *cobra.Command, dir, format string) error {
	switch format {
	case "all":
		return docs.GenerateAll(root, dir)
	case "markdown", "md":
		return docs.GenerateMarkdown(root, dir)
	case "man":
		return docs.GenerateMan(root, dir)
	case "yaml", "yml":
		return docs.GenerateYaml(root, dir)
	case "rest", "rst":
		return docs.GenerateReST(root, dir)
	case "guide":
		return docs.GenerateGuide(dir)
	case "flags":
		return docs.GenerateFlags(dir)
	default:
		return fmt.Errorf("unknown format %q (want: all, markdown, man, yaml, rest, guide, flags)", format)
	}
}
