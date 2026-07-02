// Package common holds app-wide constants shared across packages.
package common

// AppName is the program name used in diagnostics and prompts.
const AppName = "rmx"

// Version is the build version. It is overridden at link time via
//
//	-ldflags "-X github.com/braswelljr/rmx/internal/common.Version=<tag>"
var Version = "dev"
