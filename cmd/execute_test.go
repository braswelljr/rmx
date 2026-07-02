package cmd

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/spf13/pflag"

	"github.com/braswelljr/rmx/core"
)

// options registers every flag, parses args and returns the resolved Options.
func options(t *testing.T, args ...string) (core.Options, error) {
	t.Helper()
	fs := pflag.NewFlagSet("rmx", pflag.ContinueOnError)
	Register(fs)
	if err := fs.Parse(args); err != nil {
		t.Fatalf("parse %v: %v", args, err)
	}
	return buildOptions(fs)
}

func mustOptions(t *testing.T, args ...string) core.Options {
	t.Helper()
	o, err := options(t, args...)
	if err != nil {
		t.Fatalf("buildOptions %v: %v", args, err)
	}
	return o
}

func TestBuildOptions(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want core.Options
	}{
		{"defaults", nil, core.Options{PreserveRoot: true, Interactive: core.PromptDefault}},
		{"force", []string{"-f"}, core.Options{Force: true, PreserveRoot: true, Interactive: core.PromptNever}},
		{"interactive each", []string{"-i"}, core.Options{PreserveRoot: true, Interactive: core.PromptAlways}},
		{"interactive once", []string{"-I"}, core.Options{PreserveRoot: true, Interactive: core.PromptOnce}},
		{"recursive short", []string{"-r"}, core.Options{Recursive: true, PreserveRoot: true}},
		{"recursive upper", []string{"-R"}, core.Options{Recursive: true, PreserveRoot: true}},
		{"dir", []string{"-d"}, core.Options{Dir: true, PreserveRoot: true}},
		{"verbose", []string{"-v"}, core.Options{Verbose: true, PreserveRoot: true}},
		{"one file system", []string{"--one-file-system"}, core.Options{OneFileSystem: true, PreserveRoot: true}},
		{"no preserve root", []string{"--no-preserve-root", "-r"}, core.Options{Recursive: true, PreserveRoot: false}},
		{"interactive=never", []string{"--interactive=never"}, core.Options{PreserveRoot: true, Interactive: core.PromptNever}},
		{"interactive=once", []string{"--interactive=once"}, core.Options{PreserveRoot: true, Interactive: core.PromptOnce}},
		{"interactive bare", []string{"--interactive"}, core.Options{PreserveRoot: true, Interactive: core.PromptAlways}},
		{"combined -rvf", []string{"-rvf"}, core.Options{Force: true, Recursive: true, Verbose: true, PreserveRoot: true, Interactive: core.PromptNever}},

		// Precedence: -i (always) is the most cautious and wins over -f/-I.
		{"force then -i wins", []string{"-f", "-i"}, core.Options{Force: true, PreserveRoot: true, Interactive: core.PromptAlways}},
		{"force then -I", []string{"-fI"}, core.Options{Force: true, PreserveRoot: true, Interactive: core.PromptOnce}},
		{"-I then -i wins", []string{"-Ii"}, core.Options{PreserveRoot: true, Interactive: core.PromptAlways}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mustOptions(t, tt.args...)
			if got != tt.want {
				t.Fatalf("buildOptions(%v) = %+v, want %+v", tt.args, got, tt.want)
			}
		})
	}
}

func TestBuildOptionsInvalidInteractive(t *testing.T) {
	if _, err := options(t, "--interactive=bogus"); err == nil {
		t.Fatal("expected an error for --interactive=bogus")
	}
}

func TestExecuteVersion(t *testing.T) {
	var out bytes.Buffer
	root := NewRoot(strings.NewReader(""), &out, io.Discard)
	root.SetArgs([]string{"--version"})
	if err := root.Execute(); err != nil {
		t.Fatalf("execute --version: %v", err)
	}
	if !strings.Contains(out.String(), "rmx") {
		t.Fatalf("version output = %q, want it to mention rmx", out.String())
	}
}
