# rmx

[![Test](https://github.com/braswelljr/rmx/actions/workflows/test.yml/badge.svg)](https://github.com/braswelljr/rmx/actions/workflows/test.yml)
[![Build](https://github.com/braswelljr/rmx/actions/workflows/build.yml/badge.svg)](https://github.com/braswelljr/rmx/actions/workflows/build.yml)
[![Lint](https://github.com/braswelljr/rmx/actions/workflows/lint.yml/badge.svg)](https://github.com/braswelljr/rmx/actions/workflows/lint.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/braswelljr/rmx.svg)](https://pkg.go.dev/github.com/braswelljr/rmx)

**rmx** is a cross-platform, drop-in replacement for the GNU/UNIX `rm` command. It mirrors
`rm`'s flags and behavior — prompting, recursion, write-protected handling and the
`--preserve-root` failsafe — while removing independent arguments concurrently on the
non-interactive path.

```bash
rmx [OPTION]... [FILE]...
```

## Installation

### With `go install`

```bash
go install github.com/braswelljr/rmx@latest
```

### From source

```bash
git clone https://github.com/braswelljr/rmx.git
cd rmx
make install        # or: make build && ./bin/rmx --version
```

### Install script

```bash
./install.sh                 # latest via `go install`
./install.sh --from-source   # build from the current checkout
```

### Docker

```bash
docker build -t rmx .
docker run --rm -v "$(pwd):/work" rmx -rv build/
# or via compose:
docker compose run --rm rmx --version
```

## Usage

Invoke `rmx` exactly like `rm`. By default it removes files but refuses directories; use `-r`
to remove a directory tree or `-d` to remove an empty directory.

| Flag       | Long                   | Description                                                      |
|------------|------------------------|------------------------------------------------------------------|
| `-f`       | `--force`              | Ignore nonexistent files and arguments, never prompt             |
| `-i`       |                        | Prompt before **every** removal                                  |
| `-I`       |                        | Prompt **once** before removing >3 files or removing recursively |
|            | `--interactive[=WHEN]` | Prompt according to `WHEN`: `never`, `once`, or `always`         |
| `-r`, `-R` | `--recursive`          | Remove directories and their contents recursively                |
| `-d`       | `--dir`                | Remove empty directories                                         |
| `-v`       | `--verbose`            | Explain what is being removed                                    |
|            | `--one-file-system`    | When recursive, skip directories on a different file system      |
|            | `--preserve-root`      | Refuse to recurse on `/` (default)                               |
|            | `--no-preserve-root`   | Disable the `/` failsafe                                         |
| `-h`       | `--help`               | Show help                                                        |
|            | `--version`            | Print version and exit                                           |

### Examples

```bash
rmx notes.txt draft.md            # remove files
rmx -r build/                     # remove a directory and its contents
rmx -rf /tmp/cache                # unattended: no prompts, ignore missing paths
rmx -i important/*.conf           # confirm before each removal
rmx -I -r project/                # one confirmation, then recurse
rmx -d generated/                 # remove an empty directory only
rmx -rv dist/                     # recurse, printing each deletion
```

When several prompting flags are combined, `rmx` chooses the most cautious:
`-i` > `--interactive=WHEN` > `-I` > `-f`. `-f` still suppresses errors for missing files
regardless of the resulting prompt mode.

## Documentation

A full reference is generated into [`docs/`](docs/) in markdown, man, YAML and reST, grouped by
format so each directory is self-contained (`docs/markdown/` holds all markdown, etc.):

- `docs/<format>/index.*` — the whole-command reference (entry page)
- `docs/<format>/guide.*` — installation and usage guide
- `docs/<format>/flags/<flag>.*` — one page per flag (description, permissions, use cases, examples)

Regenerate it with `make docs` (or `go generate ./...`).

## Development

```bash
make build        # build ./bin/rmx (version injected via ldflags)
make test         # go test -race ./...
make cover        # tests + coverage.out summary
make lint         # golangci-lint
make fix          # gofmt -s + goimports
make docs         # regenerate docs/
make ci           # vet + lint + test + docs freshness check
make help         # list all targets
```

The project is a two-layer design: a single removal engine in [`core/`](core/) and one adapter
package per flag under [`cmd/`](cmd/). See [`AGENTS.md`](AGENTS.md) for the architecture,
conventions, and how to add a flag.

## License

See [LICENSE](LICENSE).
