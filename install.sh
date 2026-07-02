#!/usr/bin/env bash
#
# Install rmx, the cross-platform drop-in replacement for rm.
#
# Usage:
#   ./install.sh                # install the latest version with `go install`
#   VERSION=v1.2.3 ./install.sh # install a specific tag
#   ./install.sh --from-source  # build from the current checkout instead
#
set -euo pipefail

MODULE="github.com/braswelljr/rmx"
VERSION="${VERSION:-latest}"
FROM_SOURCE=0
[ "${1:-}" = "--from-source" ] && FROM_SOURCE=1

if ! command -v go >/dev/null 2>&1; then
  echo "error: the Go toolchain is required. Install it from https://go.dev/dl/ and re-run." >&2
  exit 1
fi

if [ "${FROM_SOURCE}" -eq 1 ]; then
  version="$(git describe --tags --always --dirty 2>/dev/null || echo dev)"
  echo "Building rmx ${version} from source..."
  go install -trimpath \
    -ldflags "-s -w -X ${MODULE}/internal/common.Version=${version}" \
    "${MODULE}"
else
  echo "Installing ${MODULE}@${VERSION}..."
  go install "${MODULE}@${VERSION}"
fi

bindir="$(go env GOBIN)"
[ -z "${bindir}" ] && bindir="$(go env GOPATH)/bin"
echo "Installed rmx to ${bindir}/rmx"
echo "Ensure ${bindir} is on your PATH, then run: rmx --version"
