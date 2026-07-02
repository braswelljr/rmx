# AGENTS.md

This file provides guidance to agents when working with code in this repository.

## Project Overview

**rmx** (`github.com/braswelljr/rmx`) is a cross-platform, drop-in replacement for the
GNU/UNIX `rm` command. It mirrors `rm`'s flag surface and semantics — prompting, recursion,
write-protected handling, `--preserve-root`, faithful diagnostics — while adding concurrent
removal on the non-interactive path.

It is a **single top-level command** invoked as `rmx [OPTION]... [FILE]...`, exactly like `rm`.

## Architecture at a glance

Two layers:

- **`core/`** — one correct removal engine. Owns the tree walk, prompting, permission checks,
  diagnostics and concurrency. Everything routes through it.
- **`cmd/`** — one package per flag (`cmd/force`, `cmd/interactive`, …). Each package is a thin
  **adapter**: it registers its flag(s) and folds the parsed value into a single `core.Options`.
  `cmd/execute.go` composes them and runs the engine once.

This split keeps flag composition (`rmx -rvf`, `-ri`) correct: there is exactly one walk, driven
by one `core.Options`, rather than each flag re-implementing traversal.

## Project Structure

```bash
rmx/
  main.go                  # Thin entry point; //go:generate hook for docs
  cmd/                     # One package per flag (adapters over the rm engine)
    root.go                # NewRoot(in,out,err) — builds the configured cobra command
    execute.go             # package cmd — Register(fs) + Execute() + buildOptions()
    docs/                  # Reference generator (cobra/doc: md/man/yaml/rest); tool-only dep
    force/                 # -f / --force
    interactive/           # -i, -I, --interactive[=WHEN]
    recursive/             # -r, -R, --recursive
    directory/             # -d / --dir
    verbose/               # -v / --verbose
    onefilesystem/         # --one-file-system
    preserveroot/          # --preserve-root / --no-preserve-root
    version/               # --version
  core/                    # The removal engine
    model.go               # Options, Mode, Rm struct, injection points
    core.go                # New, Run, the walk, prompting, diagnostics, concurrency
    permission_unix.go     # writable()/deviceID() via syscall (build tag !windows)
    permission_windows.go  # writable()/deviceID() via os.FileInfo (build tag windows)
  internal/
    common/                # AppName + link-time Version
    docmeta/               # zero-dep doc metadata type each cmd/<flag> exposes via Meta()
  tools/gendocs/           # `go run ./tools/gendocs [dir] [format]` — reference generator
  docs/                    # Generated reference (committed, deterministic), grouped by format:
                           #   <format>/{index.*, guide.*, flags/<flag>.*}
  .github/workflows/       # go.yml (build+test matrix), lint.yml (gofmt/tidy + golangci-lint)
```

## Development Commands

```bash
make build       # go build -o ./bin/ ./.
make install     # go install -v ./...
make test        # go test -race ./...
make lint        # golangci-lint run
make fix         # gofmt -s -w . && goimports
make docs        # regenerate docs/ in all formats (also: go generate ./...)

# Version is injected at link time:
go build -ldflags "-X github.com/braswelljr/rmx/internal/common.Version=$(git describe --tags --always)" -o bin/rmx .
```

## Tech Stack

- **Go 1.25** — module path `github.com/braswelljr/rmx`
- **Cobra / pflag** — CLI framework (single root command); `cobra/doc` for reference generation
- **mattn/go-isatty** — terminal detection (gates write-protected prompting)
- **syscall** (`access(2)` / `Stat_t.Dev`) — write-permission and device checks; build-tagged
- Tests use the standard library only (`testing`, `t.TempDir()`, injected funcs)

## CLI Surface

| Flag      | Long                                     | Behaviour                                                     |
|-----------|------------------------------------------|---------------------------------------------------------------|
| `-f`      | `--force`                                | Ignore nonexistent files, never prompt, suppress errors       |
| `-i`      | (prompt-each)                            | Prompt before **every** removal                               |
| `-I`      | (prompt-once)                            | Prompt **once** before removing >3 files or recursively       |
|           | `--interactive[=WHEN]`                   | `never` / `once` / `always` (bare `--interactive` = `always`) |
| `-r`,`-R` | `--recursive`                            | Remove directories and their contents recursively             |
| `-d`      | `--dir`                                  | Remove empty directories                                      |
| `-v`      | `--verbose`                              | Report each removal                                           |
|           | `--one-file-system`                      | When recursive, skip subdirs on a different device            |
|           | `--preserve-root` / `--no-preserve-root` | Refuse (default) / allow recursive `/`                        |
|           | `--version`                              | Print version and exit                                        |
| `-h`      | `--help`                                 | Usage                                                         |

## The flag-package contract (`cmd/*`)

Every flag package exposes the same two functions so `cmd/execute.go` can treat them uniformly:

```go
func Register(fs *pflag.FlagSet)             // define the flag(s)
func Apply(fs *pflag.FlagSet, o *core.Options) // fold parsed values into Options
```

`interactive.Apply` additionally returns an `error` for an invalid `--interactive` value.
`version` exposes `Requested(fs) bool` and `Print(w)` instead of `Apply`, since it short-circuits
before any removal.

Each flag package also exposes `Meta() docmeta.Meta` — the prose, permissions note, use cases and
per-shell examples the doc generator renders into that flag's reference page.

To **add a flag**: create `cmd/<flag>/`, implement `Register`/`Apply`/`Meta`, add the field to
`core.Options`, wire the calls into `cmd/execute.go` (`Register` and `buildOptions`), and add the
package to `flagTopics()` in `cmd/docs/flags.go`. Do not add traversal logic in the flag package —
extend the engine instead.

### Interactive-mode precedence

`buildOptions` applies `force` first (baseline `PromptNever`), then `interactive`. Conflicts
resolve toward the **most cautious** option: `-i` (always) > `--interactive=WHEN` > `-I` (once) >
`-f` (never). `-f` still suppresses nonexistent-file errors regardless of the final prompt mode.
(GNU's true last-wins-by-position is approximated because cobra does not expose flag order.)

## The engine (`core`)

- `New(opts, in, out, err) *Rm` — builds the engine; detects a terminal on `in` to gate
  write-protected prompting.
- `Run(paths) error` — entry point. Returns `core.ErrFailed` if any removal failed (diagnostics
  already written); `main` exits non-zero without reprinting it.
- Injection points for tests: `PromptFn func(string) bool` and `RemoveFn func(string) error`
  (setting `RemoveFn` also disables the `os.RemoveAll` fast path so every entry flows through it).

### Removal semantics (match GNU rm)

- Refuses `.` / `..`; refuses recursive `/` unless `--no-preserve-root`.
- `os.Lstat` throughout — symlinks are removed, never followed.
- Non-recursive directory → `Is a directory`; `-d` on a non-empty dir → `Directory not empty`.
- Recursive removal is **post-order**: with `-i`, prompts `descend into directory?`, then each
  child, then `remove directory?`.
- Diagnostics use rm's noun set (`describe`) and capitalised syscall strings (`errText`).

### Permissions / write-protected files

In the default mode (no `-f`/`-i`) on a terminal, a non-writable entry triggers
`remove write-protected <type> '<path>'?`. Writability is `access(2)` on Unix and the owner
write-bit on Windows (`core/permission_*.go`). `-f` skips these prompts; `-i` supersedes them with
the ordinary prompt.

### Concurrency model

- **Non-interactive, non-verbose** removal fans out one worker per CPU across the top-level
  arguments (`runConcurrent`); `errCount` and writer access are mutex-guarded.
- **Prompting or `-v`** runs sequentially — prompts share one stdin, and `-v` output must stay in
  argument order to match `rm`.
- Single-directory recursion uses `os.RemoveAll` (fast path) whenever no per-entry decision is
  needed (`fastPathOK`): no prompting, no `-v`, no `--one-file-system`, no injected `RemoveFn`.

## Documentation generation

`tools/gendocs` (also `make docs` and `go generate ./...`) builds the reference into `docs/` via
`cmd/docs`. It produces, in **markdown, man, yaml and reST**:

Each format directory is self-contained (`docs/markdown/` holds all markdown, and so on):

- `<format>/index.*` — the whole-command reference (the entry page).
- `<format>/guide.*` — the installation + usage guide (synthetic command; prose in `guide.go`).
- `<format>/flags/<flag>.*` — one rich page per flag: description, permissions, use cases, the
  flag listing, and per-shell (`bash` / `powershell`) examples.

Per-flag markdown is hand-rendered (`renderFlagMarkdown`) to get multi-section, multi-language
code blocks; man/yaml/reST reuse `cobra/doc` from a synthetic command that folds permissions and
use cases into its description. Output is **deterministic** — the cobra auto-gen tag is disabled
and the man date is pinned (`docs.ManDate`) — so committed docs don't churn. Fenced code blocks
are language-tagged (`bash` for commands, `text` for flag listings). `cobra/doc` pulls in
`go-md2man`, but only `tools/gendocs` imports `cmd/docs`, so it never links into the `rmx` binary.

## Code Style

- Doc comments on exported symbols; comment the *why*, not the obvious.
- US spelling (misspell runs with `locale: US`).
- Wrap errors with `fmt.Errorf("…: %w", err)`.
- Filenames are snake_case (`permission_unix.go`).
- Platform variance goes in build-tagged files (`*_unix.go` / `*_windows.go`), never runtime
  `runtime.GOOS` branching.

## Linting

`.golangci.yml` (golangci-lint **v2** schema): default set + `misspell`, `unconvert`, `revive`,
`gocritic`. Complexity linters (`funlen`, `cyclop`, `maintidx`, `gocognit`) are disabled — the
engine's dispatch uses intentional `switch` chains. `errcheck` ignores `*.Close()`, `os.Remove`
cleanup, and best-effort `fmt.Fprint*` diagnostic writes. `_test.go` files are exempt from
`gocritic`/`revive`.

## Testing

- Race detector always on: `go test -race ./...`.
- `core/core_test.go` is the primary suite (black-box `core_test`): covers file/dir/recursive
  removal, force, `-d` empty vs non-empty, interactive accept/decline, `-I` decline, verbose
  order, concurrent bulk removal, `.`-refusal, and post-order via injected `RemoveFn`.
- `cmd/execute_test.go` covers `buildOptions` — every flag mapping, the interactive-mode
  precedence, and the `--version` short-circuit.
- `cmd/docs/docs_test.go` covers `wrapCode`/`fenceLang` and an end-to-end `GenerateAll` (files
  present, markdown wraps references, man stays plain, output is deterministic).
- Tests use `t.TempDir()`; inject `PromptFn`/`RemoveFn` rather than driving a real terminal.

## CI

| File | What it does |
| --- | --- |
| `go.yml` | `go fmt` + `go build` + `go test -race` on ubuntu/macos/windows |
| `lint.yml` | `gofmt -s` / `go mod tidy` drift check + `golangci-lint-action@v7` |

Both read the Go toolchain from `go.mod` via `go-version-file: go.mod` — keep it there rather
than hardcoding a version (the engine uses `min` and range-over-int, which need Go ≥ 1.22).

## Possible future work

- `-i` write-protected prompt fidelity for symlink targets and special files.
- Read paths from stdin; `--` end-of-options handling.

## Important Files

- `main.go` — thin entry point; `//go:generate` hook, error/exit handling.
- `cmd/root.go` — `NewRoot` builds the configured cobra command.
- `cmd/execute.go` — flag registration, option composition, engine invocation.
- `cmd/docs/` — reference generator (`docs.go`, `flags.go`, `guide.go`).
- `core/core.go` — the engine: walk, prompting, diagnostics, concurrency.
- `core/model.go` — `Options`, `Mode`, `Rm`, injection points.
- `core/permission_unix.go` / `permission_windows.go` — writability + device id.
- `internal/docmeta/` — doc metadata type exposed by each flag package's `Meta()`.
- `.golangci.yml` — linter configuration (v2 schema).
