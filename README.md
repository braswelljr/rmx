# rmx

A cross-platform replacement for the UNIX `rm` command.

## Installation

```bash
go get github.com/braswelljr/rmx
```

## Usage

```bash
rmx [flags] [path ...]
```

### Flags

`-f`, `--force` - Force removal of files and directories.

`-i` - Prompt before removal.

`-I` - Prompt before removal of directories.

`-r`, `-R` - Recursively remove directories.

`-v` - Verbose output.

`-h` - Print help.

`--version` - Print version.
