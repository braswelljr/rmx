// Package core implements the removal engine behind rmx, a drop-in replacement
// for the GNU/UNIX rm command. It mirrors rm's flags, prompting behavior and
// diagnostics while adding concurrent removal on the non-interactive path.
package core

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"unicode"

	"github.com/mattn/go-isatty"

	"github.com/braswelljr/rmx/internal/common"
)

// ErrFailed reports that one or more paths could not be removed. It is returned
// by Run after all diagnostics have already been written to Err, so callers
// should exit non-zero without printing it again.
var ErrFailed = errors.New("some files could not be removed")

// New builds a removal engine. If in is an *os.File attached to a terminal,
// write-protected prompting is enabled in PromptDefault mode.
func New(opts Options, in io.Reader, out, errw io.Writer) *Rm {
	r := &Rm{
		Options: opts,
		In:      in,
		Out:     out,
		Err:     errw,
		reader:  bufio.NewReader(in),
	}
	if f, ok := in.(*os.File); ok {
		r.isTTY = isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd())
	}
	return r
}

// Run removes each of paths according to the configured options. Diagnostics
// are written to Err as they occur; the returned error is ErrFailed if any
// removal failed, otherwise nil.
func (r *Rm) Run(paths []string) error {
	if len(paths) == 0 {
		if r.Force {
			return nil
		}
		r.errf("missing operand")
		r.warnf("Try '%s --help' for more information.", common.AppName)
		return r.result()
	}

	// -I: a single up-front prompt gates a bulk or recursive removal.
	if r.Interactive == PromptOnce {
		if n := len(paths); r.Recursive || n > 3 {
			q := fmt.Sprintf("remove %d argument%s?", n, plural(n))
			if r.Recursive {
				q = fmt.Sprintf("remove %d argument%s recursively?", n, plural(n))
			}
			if !r.ask(q) {
				return nil
			}
		}
		r.Interactive = PromptNever
	}

	// Prompting reads a shared stdin, and verbose output must stay in argument
	// order; both force sequential processing. Everything else fans out.
	if r.mayPrompt() || r.Verbose {
		for _, p := range paths {
			r.removeArg(p)
		}
	} else {
		r.runConcurrent(paths)
	}
	return r.result()
}

// runConcurrent removes independent top-level arguments in parallel, bounded to
// one worker per CPU. Only reached when no prompting can occur.
func (r *Rm) runConcurrent(paths []string) {
	if len(paths) == 1 {
		r.removeArg(paths[0])
		return
	}

	// len(paths) >= 2 here (the single-path case returned above), and
	// runtime.NumCPU() >= 1, so workers is always >= 1.
	workers := min(runtime.NumCPU(), len(paths))

	jobs := make(chan string)
	var wg sync.WaitGroup
	wg.Add(workers)
	for range workers {
		go func() {
			defer wg.Done()
			for p := range jobs {
				r.removeArg(p)
			}
		}()
	}
	for _, p := range paths {
		jobs <- p
	}
	close(jobs)
	wg.Wait()
}

// removeArg handles a single top-level operand.
func (r *Rm) removeArg(path string) {
	if base := filepath.Base(filepath.Clean(path)); base == "." || base == ".." {
		r.errf("refusing to remove '.' or '..' directory: skipping '%s'", path)
		return
	}

	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			if r.Force {
				return
			}
			r.errf("cannot remove '%s': No such file or directory", path)
			return
		}
		r.errf("cannot remove '%s': %s", path, errText(err))
		return
	}

	if info.IsDir() {
		if r.Recursive && r.PreserveRoot && isRoot(path) {
			r.errf("it is dangerous to operate recursively on '%s'", path)
			r.warnf("use --no-preserve-root to override this failsafe")
			return
		}
		r.removeDirArg(path, info)
		return
	}

	r.removeFileEntry(path, info)
}

// removeDirArg dispatches directory removal based on -r and -d.
func (r *Rm) removeDirArg(path string, info os.FileInfo) {
	switch {
	case r.Recursive:
		if r.fastPathOK() {
			if err := os.RemoveAll(path); err != nil {
				r.errf("cannot remove '%s': %s", path, errText(err))
			}
			return
		}
		dev, haveDev := uint64(0), false
		if r.OneFileSystem {
			dev, haveDev = deviceID(info)
		}
		r.removeTree(path, dev, haveDev)

	case r.Dir:
		// Report a non-empty directory with rm's canonical message rather than
		// the OS-native errno text (which differs across platforms).
		entries, err := os.ReadDir(path)
		if err != nil {
			r.errf("cannot remove '%s': %s", path, errText(err))
			return
		}
		if len(entries) > 0 {
			r.errf("cannot remove '%s': Directory not empty", path)
			return
		}
		if !r.confirmDir(path) {
			return
		}
		if err := r.del(path); err != nil {
			r.errf("cannot remove '%s': %s", path, errText(err))
			return
		}
		if r.Verbose {
			r.outf("removed directory '%s'", path)
		}

	default:
		r.errf("cannot remove '%s': Is a directory", path)
	}
}

// removeTree removes a directory and its contents post-order, honoring
// prompting, --one-file-system and --verbose. It returns true only if the whole
// subtree (including the directory itself) was removed.
func (r *Rm) removeTree(path string, rootDev uint64, haveDev bool) bool {
	entries, err := os.ReadDir(path)
	if err != nil {
		r.errf("cannot remove '%s': %s", path, errText(err))
		return false
	}

	if r.Interactive == PromptAlways && len(entries) > 0 {
		if !r.ask(fmt.Sprintf("descend into directory '%s'?", path)) {
			return false
		}
	}

	childSkipped := false
	for _, e := range entries {
		child := filepath.Join(path, e.Name())
		ci, err := os.Lstat(child)
		if err != nil {
			r.errf("cannot remove '%s': %s", child, errText(err))
			childSkipped = true
			continue
		}

		if ci.IsDir() {
			if r.OneFileSystem && haveDev {
				if cdev, ok := deviceID(ci); ok && cdev != rootDev {
					r.warnf("skipping '%s', since it's on a different device", child)
					childSkipped = true
					continue
				}
			}
			if !r.removeTree(child, rootDev, haveDev) {
				childSkipped = true
			}
		} else if !r.removeFileEntry(child, ci) {
			childSkipped = true
		}
	}

	if childSkipped {
		return false
	}

	if !r.confirmDir(path) {
		return false
	}
	if err := r.del(path); err != nil {
		r.errf("cannot remove '%s': %s", path, errText(err))
		return false
	}
	if r.Verbose {
		r.outf("removed directory '%s'", path)
	}
	return true
}

// removeFileEntry removes a single non-directory entry, prompting as required.
// It returns true if the entry was removed.
func (r *Rm) removeFileEntry(path string, info os.FileInfo) bool {
	switch {
	case r.Interactive == PromptAlways:
		if !r.ask(fmt.Sprintf("remove %s '%s'?", describe(info), path)) {
			return false
		}
	case r.Interactive == PromptDefault && r.isTTY && !writable(path):
		if !r.ask(fmt.Sprintf("remove write-protected %s '%s'?", describe(info), path)) {
			return false
		}
	}

	if err := r.del(path); err != nil {
		if r.Force && os.IsNotExist(err) {
			return false
		}
		r.errf("cannot remove '%s': %s", path, errText(err))
		return false
	}
	if r.Verbose {
		r.outf("removed '%s'", path)
	}
	return true
}

// confirmDir asks for confirmation before removing an (empty) directory,
// returning true when removal should proceed.
func (r *Rm) confirmDir(path string) bool {
	switch {
	case r.Interactive == PromptAlways:
		return r.ask(fmt.Sprintf("remove directory '%s'?", path))
	case r.Interactive == PromptDefault && r.isTTY && !writable(path):
		return r.ask(fmt.Sprintf("remove write-protected directory '%s'?", path))
	default:
		return true
	}
}

// mayPrompt reports whether any removal could trigger an interactive prompt.
func (r *Rm) mayPrompt() bool {
	return r.Interactive == PromptAlways || (r.Interactive == PromptDefault && r.isTTY)
}

// fastPathOK reports whether a directory can be removed via a single
// os.RemoveAll rather than a manual, decision-making walk.
func (r *Rm) fastPathOK() bool {
	return r.RemoveFn == nil && !r.mayPrompt() && !r.Verbose && !r.OneFileSystem
}

// del performs the actual unlink, honoring an injected RemoveFn.
func (r *Rm) del(path string) error {
	if r.RemoveFn != nil {
		return r.RemoveFn(path)
	}
	return os.Remove(path)
}

// ask poses a yes/no question, returning true on an affirmative answer.
func (r *Rm) ask(question string) bool {
	if r.PromptFn != nil {
		return r.PromptFn(question)
	}

	r.mu.Lock()
	fmt.Fprintf(r.Err, "%s: %s ", common.AppName, question)
	r.mu.Unlock()

	line, err := r.reader.ReadString('\n')
	if err != nil && line == "" {
		return false
	}
	line = strings.TrimSpace(line)
	return line != "" && (line[0] == 'y' || line[0] == 'Y')
}

func (r *Rm) errf(format string, a ...any) {
	r.mu.Lock()
	defer r.mu.Unlock()
	fmt.Fprintf(r.Err, common.AppName+": "+format+"\n", a...)
	r.errCount++
}

func (r *Rm) warnf(format string, a ...any) {
	r.mu.Lock()
	defer r.mu.Unlock()
	fmt.Fprintf(r.Err, common.AppName+": "+format+"\n", a...)
}

func (r *Rm) outf(format string, a ...any) {
	r.mu.Lock()
	defer r.mu.Unlock()
	fmt.Fprintf(r.Out, format+"\n", a...)
}

func (r *Rm) result() error {
	if r.errCount > 0 {
		return ErrFailed
	}
	return nil
}

// describe returns rm's noun for a filesystem entry (e.g. "regular empty file").
func describe(info os.FileInfo) string {
	m := info.Mode()
	switch {
	case m&os.ModeSymlink != 0:
		return "symbolic link"
	case m.IsDir():
		return "directory"
	case m&os.ModeNamedPipe != 0:
		return "fifo"
	case m&os.ModeSocket != 0:
		return "socket"
	case m&os.ModeDevice != 0:
		if m&os.ModeCharDevice != 0 {
			return "character special file"
		}
		return "block special file"
	case m.IsRegular():
		if info.Size() == 0 {
			return "regular empty file"
		}
		return "regular file"
	default:
		return "file"
	}
}

// errText renders an error the way rm does: the bare, capitalised syscall
// message (e.g. "Permission denied", "Directory not empty").
func errText(err error) string {
	var perr *os.PathError
	if errors.As(err, &perr) {
		err = perr.Err
	}
	var lerr *os.LinkError
	if errors.As(err, &lerr) {
		err = lerr.Err
	}
	msg := err.Error()
	if msg == "" {
		return msg
	}
	runes := []rune(msg)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// isRoot reports whether path resolves to a filesystem root (e.g. "/" or "C:\").
func isRoot(path string) bool {
	abs, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	abs = filepath.Clean(abs)
	return filepath.Dir(abs) == abs
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
