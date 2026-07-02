# syntax=docker/dockerfile:1

# ---- build stage ---------------------------------------------------------------
FROM golang:1.25-alpine AS build

WORKDIR /src

# Cache module downloads separately from the source.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Version is injected at link time; override with --build-arg VERSION=<tag>.
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -trimpath \
    -ldflags "-s -w -X github.com/braswelljr/rmx/internal/common.Version=${VERSION}" \
    -o /out/rmx .

# ---- runtime stage -------------------------------------------------------------
FROM alpine:3.20

# Run as an unprivileged user; mount the directory to operate on at /work.
RUN adduser -D -u 10001 rmx
COPY --from=build /out/rmx /usr/local/bin/rmx

USER rmx
WORKDIR /work

ENTRYPOINT ["rmx"]
CMD ["--help"]
