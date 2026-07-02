// Package docmeta defines the documentation metadata that each cmd/<flag>
// package exposes via its Meta() function. It has no dependencies, so it is
// safe to import from the flag packages that link into the rmx binary — the
// doc generator (cmd/docs) turns this metadata into per-flag reference pages.
package docmeta

// Meta describes a single flag topic for documentation purposes. The flag
// details themselves (short/long name, default, usage) are read from the
// package's Register function; Meta supplies the prose, permissions notes,
// use cases and worked examples around them.
type Meta struct {
	Name        string    // topic/file base name, e.g. "force"
	Use         string    // synopsis line for the topic, e.g. "force [FILE]..."
	Short       string    // one-line summary
	Long        string    // full description
	Permissions string    // how the flag interacts with filesystem permissions
	UseCases    []string  // when to reach for this flag
	Examples    []Example // worked examples
}

// Example is one worked example: a description plus the equivalent command in
// each shell. PowerShell may be empty when it adds nothing over Bash.
type Example struct {
	Description string
	Bash        string
	PowerShell  string
}
