package core

import (
	"bufio"
	"io"
	"sync"
)

// Mode controls when rmx prompts before removing a path.
type Mode int

const (
	// PromptDefault prompts only for write-protected files, and only when
	// standard input is an interactive terminal (mirrors rm with no -f/-i/-I).
	PromptDefault Mode = iota
	// PromptNever never prompts (rm -f / --interactive=never).
	PromptNever
	// PromptOnce prompts a single time before a bulk or recursive removal
	// (rm -I / --interactive=once).
	PromptOnce
	// PromptAlways prompts before every removal (rm -i / --interactive=always).
	PromptAlways
)

// Options mirrors the GNU rm flag surface.
type Options struct {
	Force         bool // -f: ignore nonexistent files, never prompt, suppress errors
	Recursive     bool // -r/-R: remove directories and their contents recursively
	Dir           bool // -d: remove empty directories
	Verbose       bool // -v: explain what is being done
	Interactive   Mode // -i/-I/--interactive
	OneFileSystem bool // --one-file-system: don't cross device boundaries
	PreserveRoot  bool // refuse to recurse on '/' (default true)
}

// Rm is the removal engine. Construct it with New.
type Rm struct {
	Options

	In  io.Reader
	Out io.Writer
	Err io.Writer

	// PromptFn, when non-nil, is consulted instead of reading from In. It is
	// the injection point used by tests to answer confirmation prompts.
	PromptFn func(question string) bool

	// RemoveFn, when non-nil, replaces os.Remove for the actual unlink. Tests
	// use it to record removals without touching the filesystem; setting it
	// also disables the os.RemoveAll fast path so every entry flows through it.
	RemoveFn func(path string) error

	isTTY    bool
	reader   *bufio.Reader
	mu       sync.Mutex // guards Out/Err writes and errCount during concurrent removal
	errCount int
}
