// Package docs generates reference documentation for the rmx CLI in every
// format cobra supports (markdown, man, yaml, ReST). It is used only by the
// tools/gendocs generator, so the go-md2man dependency it pulls in never links
// into the rmx binary itself.
package docs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"

	"github.com/braswelljr/rmx/internal/common"
)

// ManDate is stamped into generated man pages. It is fixed (rather than
// time.Now) so committed documentation is reproducible; release tooling can
// override it before calling GenerateMan or GenerateAll.
var ManDate = time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)

// prepare creates dir and disables cobra's auto-gen timestamp so the markdown,
// yaml and ReST output is deterministic and safe to commit.
func prepare(root *cobra.Command, dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create docs dir %q: %w", dir, err)
	}
	root.DisableAutoGenTag = true
	return nil
}

// GenerateMarkdown writes a markdown reference tree into dir, with each fenced
// code block labeled with a language for syntax highlighting.
func GenerateMarkdown(root *cobra.Command, dir string) error {
	if err := prepare(root, dir); err != nil {
		return err
	}
	return writeFile(dir, "", "index.md", func(w io.Writer) error {
		return genMarkdownTagged(root, w)
	})
}

// GenerateReST writes the reStructuredText reference into dir as index.rst.
func GenerateReST(root *cobra.Command, dir string) error {
	if err := prepare(root, dir); err != nil {
		return err
	}
	return writeFile(dir, "", "index.rst", func(w io.Writer) error {
		return doc.GenReST(root, w)
	})
}

// GenerateYaml writes the YAML reference into dir as index.yaml.
func GenerateYaml(root *cobra.Command, dir string) error {
	if err := prepare(root, dir); err != nil {
		return err
	}
	return writeFile(dir, "", "index.yaml", func(w io.Writer) error {
		return doc.GenYaml(root, w)
	})
}

// GenerateMan writes the section-1 man page into dir as index.1.
func GenerateMan(root *cobra.Command, dir string) error {
	if err := prepare(root, dir); err != nil {
		return err
	}
	header := &doc.GenManHeader{
		Title:   strings.ToUpper(common.AppName),
		Section: "1",
		Source:  common.AppName,
		Manual:  common.AppName + " manual",
		Date:    &ManDate,
	}
	return writeFile(dir, "", "index.1", func(w io.Writer) error {
		return doc.GenMan(root, header, w)
	})
}

// GenerateAll writes the complete documentation set under base, grouped by
// format so each format directory is self-contained:
//
//	base/<format>/index.*        whole-command reference (entry page)
//	base/<format>/guide.*        installation/usage guide
//	base/<format>/flags/<flag>.* per-flag reference
func GenerateAll(root *cobra.Command, base string) error {
	formats := []struct {
		name string
		gen  func(*cobra.Command, string) error
	}{
		{"markdown", GenerateMarkdown},
		{"man", GenerateMan},
		{"yaml", GenerateYaml},
		{"rest", GenerateReST},
	}
	for _, f := range formats {
		if err := f.gen(root, filepath.Join(base, f.name)); err != nil {
			return err
		}
	}
	if err := GenerateGuide(base); err != nil {
		return err
	}
	return GenerateFlags(base)
}
