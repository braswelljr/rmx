package docs

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/braswelljr/rmx/cmd"
)

func TestWrapCode(t *testing.T) {
	tests := []struct {
		name, in, want string
	}{
		{"flag pair", "the -d/--dir flag", "the `-d`/`--dir` flag"},
		{"command words", "rmx wraps rm and rmdir", "`rmx` wraps `rm` and `rmdir`"},
		{"command phrase", "run go install then make build", "run `go install` then `make build`"},
		{"hyphenated words untouched", "a write-protected read-only file", "a write-protected read-only file"},
		{"heading untouched", "## rm and -d", "## rm and -d"},
		{"inline code untouched", "already `-d` here", "already `-d` here"},
		{"long flag with dashes", "the --no-preserve-root failsafe", "the `--no-preserve-root` failsafe"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := string(wrapCode([]byte(tt.in))); got != tt.want {
				t.Fatalf("wrapCode(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestWrapCodeSkipsFencedBlocks(t *testing.T) {
	in := "prose -d\n```bash\nrmx -d dir\n```\nmore -r"
	got := string(wrapCode([]byte(in)))

	if !strings.Contains(got, "prose `-d`") {
		t.Errorf("prose flag should be wrapped: %q", got)
	}
	if !strings.Contains(got, "more `-r`") {
		t.Errorf("trailing prose flag should be wrapped: %q", got)
	}
	if !strings.Contains(got, "\nrmx -d dir\n") {
		t.Errorf("fenced content must stay untouched: %q", got)
	}
}

func TestFenceLang(t *testing.T) {
	if got := fenceLang("### options"); got != "text" {
		t.Errorf("options => %q, want text", got)
	}
	if got := fenceLang("### examples"); got != "bash" {
		t.Errorf("examples => %q, want bash", got)
	}
}

func TestGenerateAll(t *testing.T) {
	dir := t.TempDir()
	root := cmd.NewRoot(strings.NewReader(""), io.Discard, io.Discard)
	if err := GenerateAll(root, dir); err != nil {
		t.Fatalf("GenerateAll: %v", err)
	}

	// A representative file exists and is non-empty in every layout position.
	for _, p := range []string{
		"markdown/index.md", "markdown/guide.md", "markdown/flags/force.md",
		"man/index.1", "man/flags/recursive.1",
		"yaml/index.yaml", "rest/flags/directory.rst",
	} {
		fi, err := os.Stat(filepath.Join(dir, p))
		if err != nil || fi.Size() == 0 {
			t.Errorf("expected non-empty %s (err=%v)", p, err)
		}
	}

	// Markdown wraps flag references; man stays plain.
	md, _ := os.ReadFile(filepath.Join(dir, "markdown/flags/directory.md"))
	if !strings.Contains(string(md), "`--dir`") {
		t.Errorf("markdown per-flag page should wrap flags in code")
	}
	man, _ := os.ReadFile(filepath.Join(dir, "man/flags/directory.1"))
	if strings.Contains(string(man), "`") {
		t.Errorf("man page should contain no backticks")
	}
}

func TestGenerateAllDeterministic(t *testing.T) {
	a, b := t.TempDir(), t.TempDir()
	for _, dir := range []string{a, b} {
		root := cmd.NewRoot(strings.NewReader(""), io.Discard, io.Discard)
		if err := GenerateAll(root, dir); err != nil {
			t.Fatalf("GenerateAll(%s): %v", dir, err)
		}
	}
	if diff := treeDiff(t, a, b); diff != "" {
		t.Fatalf("regeneration is not deterministic: %s", diff)
	}
}

// treeDiff returns a description of the first difference between two directory
// trees, or "" if their file sets and contents are identical.
func treeDiff(t *testing.T, a, b string) string {
	t.Helper()
	fa, fb := readTree(t, a), readTree(t, b)
	if len(fa) != len(fb) {
		return "different file counts"
	}
	for name, content := range fa {
		if fb[name] != content {
			return "differs: " + name
		}
	}
	return ""
}

func readTree(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, path)
		out[rel] = string(b)
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", root, err)
	}
	return out
}
