package docs

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/pflag"

	"github.com/braswelljr/rmx/cmd/directory"
	"github.com/braswelljr/rmx/cmd/force"
	"github.com/braswelljr/rmx/cmd/interactive"
	"github.com/braswelljr/rmx/cmd/onefilesystem"
	"github.com/braswelljr/rmx/cmd/preserveroot"
	"github.com/braswelljr/rmx/cmd/recursive"
	"github.com/braswelljr/rmx/cmd/verbose"
	"github.com/braswelljr/rmx/cmd/version"
	"github.com/braswelljr/rmx/internal/common"
	"github.com/braswelljr/rmx/internal/docmeta"
)

// flagTopic pairs a flag package's documentation metadata with its flag
// registrar, so the generator can render each flag as a standalone doc.
type flagTopic struct {
	meta     docmeta.Meta
	register func(*pflag.FlagSet)
}

// flagTopics is the registry of documentable flag packages. Add new flag
// packages here to include them in the per-flag reference.
func flagTopics() []flagTopic {
	return []flagTopic{
		{force.Meta(), force.Register},
		{interactive.Meta(), interactive.Register},
		{recursive.Meta(), recursive.Register},
		{directory.Meta(), directory.Register},
		{verbose.Meta(), verbose.Register},
		{onefilesystem.Meta(), onefilesystem.Register},
		{preserveroot.Meta(), preserveroot.Register},
		{version.Meta(), version.Register},
	}
}

// flagSet builds a fresh flag set holding only this topic's flags, used to
// render the options listing without cobra's injected help flag.
func (t flagTopic) flagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet(t.meta.Name, pflag.ContinueOnError)
	t.register(fs)
	return fs
}

// command builds a standalone cobra command for the topic so cobra/doc can
// render it in man, yaml and reST. Permissions and use cases are folded into
// the description, and the bash examples into the example block, so those
// formats carry the same information as the markdown page.
func (t flagTopic) command() *cobra.Command {
	long := t.meta.Long
	if t.meta.Permissions != "" {
		long += "\n\nPERMISSIONS. " + t.meta.Permissions
	}
	if len(t.meta.UseCases) > 0 {
		long += "\n\nUSE CASES.\n"
		for _, u := range t.meta.UseCases {
			long += "  - " + u + "\n"
		}
	}

	c := &cobra.Command{
		Use:     t.meta.Use,
		Short:   t.meta.Short,
		Long:    strings.TrimRight(long, "\n"),
		Example: exampleBlock(t.meta.Examples),
		// A no-op Run marks the command runnable so the synopsis is emitted.
		Run: func(*cobra.Command, []string) {},
	}
	c.DisableAutoGenTag = true
	t.register(c.Flags())
	return c
}

// exampleBlock renders the examples as a plain shell block for the cobra-driven
// formats (man/yaml/rest), which support only a single example section.
func exampleBlock(examples []docmeta.Example) string {
	var b strings.Builder
	for i, e := range examples {
		if i > 0 {
			b.WriteString("\n")
		}
		if e.Description != "" {
			fmt.Fprintf(&b, "  # %s\n", strings.TrimRight(e.Description, ":"))
		}
		if e.Bash != "" {
			fmt.Fprintf(&b, "  %s\n", e.Bash)
		}
	}
	return strings.TrimRight(b.String(), "\n")
}

// GenerateFlags writes a per-flag reference for every flag topic, nested under
// each format directory: base/<format>/flags/<flag>.<ext>. This keeps each
// format self-contained (e.g. base/markdown holds all markdown, flags included).
func GenerateFlags(base string) error {
	for _, t := range flagTopics() {
		name := t.meta.Name

		if err := writeFile(base, filepath.Join("markdown", "flags"), name+".md", func(w io.Writer) error {
			_, err := w.Write(wrapCode([]byte(renderFlagMarkdown(t.meta, t.flagSet()))))
			return err
		}); err != nil {
			return err
		}

		c := t.command()
		if err := writeFile(base, filepath.Join("rest", "flags"), name+".rst", func(w io.Writer) error {
			return doc.GenReST(c, w)
		}); err != nil {
			return err
		}
		if err := writeFile(base, filepath.Join("yaml", "flags"), name+".yaml", func(w io.Writer) error {
			return doc.GenYaml(c, w)
		}); err != nil {
			return err
		}
		header := manHeader(common.AppName + "-" + name)
		if err := writeFile(base, filepath.Join("man", "flags"), name+".1", func(w io.Writer) error {
			return doc.GenMan(c, header, w)
		}); err != nil {
			return err
		}
	}
	return nil
}

// renderFlagMarkdown builds a rich, multi-section markdown page for one flag:
// description, permissions, use cases, the flag listing, and per-shell examples
// in language-tagged code blocks.
func renderFlagMarkdown(m docmeta.Meta, flags *pflag.FlagSet) string {
	var b strings.Builder

	fmt.Fprintf(&b, "## %s\n\n%s\n\n", m.Name, m.Short)
	fmt.Fprintf(&b, "### Description\n\n%s\n\n", m.Long)

	if m.Permissions != "" {
		fmt.Fprintf(&b, "### Permissions\n\n%s\n\n", m.Permissions)
	}

	if len(m.UseCases) > 0 {
		b.WriteString("### Use cases\n\n")
		for _, u := range m.UseCases {
			fmt.Fprintf(&b, "- %s\n", u)
		}
		b.WriteString("\n")
	}

	if usage := strings.TrimRight(flags.FlagUsages(), "\n"); usage != "" {
		fmt.Fprintf(&b, "### Flags\n\n```text\n%s\n```\n\n", usage)
	}

	if len(m.Examples) > 0 {
		b.WriteString("### Examples\n\n")
		for _, e := range m.Examples {
			if e.Description != "" {
				fmt.Fprintf(&b, "%s\n\n", e.Description)
			}
			if e.Bash != "" {
				fmt.Fprintf(&b, "```bash\n%s\n```\n\n", e.Bash)
			}
			if e.PowerShell != "" {
				fmt.Fprintf(&b, "```powershell\n%s\n```\n\n", e.PowerShell)
			}
		}
	}

	return strings.TrimRight(b.String(), "\n") + "\n"
}

// manHeader builds a section-1 man header with the pinned ManDate.
func manHeader(title string) *doc.GenManHeader {
	return &doc.GenManHeader{
		Title:   strings.ToUpper(title),
		Section: "1",
		Source:  common.AppName,
		Manual:  common.AppName + " flag reference",
		Date:    &ManDate,
	}
}

// writeFile renders gen into base/sub/filename, creating directories as needed.
func writeFile(base, sub, filename string, gen func(io.Writer) error) error {
	dir := filepath.Join(base, sub)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create docs dir %q: %w", dir, err)
	}
	buf := new(bytes.Buffer)
	if err := gen(buf); err != nil {
		return fmt.Errorf("render %q: %w", filename, err)
	}
	path := filepath.Join(dir, filename)
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("write %q: %w", path, err)
	}
	return nil
}

// genMarkdownTagged renders cmd as markdown, labels each fenced code block with
// a language, and wraps flag/command references in inline code.
func genMarkdownTagged(cmd *cobra.Command, w io.Writer) error {
	buf := new(bytes.Buffer)
	if err := doc.GenMarkdown(cmd, buf); err != nil {
		return err
	}
	_, err := w.Write(wrapCode(tagMarkdownFences(buf.Bytes())))
	return err
}

// Patterns for wrapping flag and command references in prose. Applied to
// markdown only — the source strings stay plain so terminal help and the
// man/yaml/reST outputs are unaffected.
var (
	// flags: -d, -rf, --dir, --no-preserve-root (preceded by a non-word/dash char)
	flagRefRe = regexp.MustCompile(`(^|[^\w-])(--?[A-Za-z][\w-]*)`)
	// multi-word tool commands: "go install", "make build", "docker run", ...
	cmdPhraseRe = regexp.MustCompile(`\b((?:go|make|docker) (?:install|build|generate|test|run|compose|mod|docs|lint|cover|tidy))\b`)
	// standalone command names specific to this project
	cmdWordRe = regexp.MustCompile(`\b(rmdir|rmx|rm)\b`)
)

// wrapCode wraps flag and command references in prose with backticks so they
// read as code and stand out from regular words. It skips fenced code blocks,
// inline code spans and headings, so existing code is never touched.
func wrapCode(md []byte) []byte {
	lines := strings.Split(string(md), "\n")
	inFence := false
	for i, line := range lines {
		if strings.HasPrefix(line, "```") {
			inFence = !inFence
			continue
		}
		if inFence || strings.HasPrefix(line, "#") {
			continue
		}
		lines[i] = wrapCodeInLine(line)
	}
	return []byte(strings.Join(lines, "\n"))
}

// wrapCodeInLine wraps references in a single prose line, leaving anything
// already inside inline-code spans (backtick pairs) untouched.
func wrapCodeInLine(line string) string {
	parts := strings.Split(line, "`")
	// Even indexes are outside inline code; odd indexes are inside it.
	for i := 0; i < len(parts); i += 2 {
		s := parts[i]
		s = cmdPhraseRe.ReplaceAllString(s, "`$1`")
		s = cmdWordRe.ReplaceAllString(s, "`$1`")
		s = flagRefRe.ReplaceAllString(s, "$1`$2`")
		parts[i] = s
	}
	return strings.Join(parts, "`")
}

// tagMarkdownFences rewrites each opening ``` fence to ```<lang>, choosing the
// language from the nearest preceding heading. Closing fences are left bare.
func tagMarkdownFences(md []byte) []byte {
	lines := strings.Split(string(md), "\n")
	section := ""
	inBlock := false
	for i, line := range lines {
		if strings.HasPrefix(line, "#") {
			section = strings.ToLower(line)
			continue
		}
		if strings.TrimSpace(line) == "```" {
			if inBlock {
				inBlock = false
				continue
			}
			inBlock = true
			lines[i] = "```" + fenceLang(section)
		}
	}
	return []byte(strings.Join(lines, "\n"))
}

// fenceLang maps a markdown section heading to a code-block language. The
// options listing is a plain table (text); synopsis and examples are shell.
func fenceLang(section string) string {
	if strings.Contains(section, "option") {
		return "text"
	}
	return "bash"
}
