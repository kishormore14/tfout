package cli

import (
	"fmt"
	"io"
	"runtime"
)

var (
	// Version is injected at build time via -ldflags.
	Version = "v0.1.0"
	// Commit is injected at build time via -ldflags.
	Commit = "dev"
	// Date is injected at build time via -ldflags.
	Date = "unknown"
)

// RunVersion executes the `tfout version` command.
func RunVersion(stdout io.Writer) int {
	fmt.Fprintf(stdout, "tfout version %s (%s) built on %s with %s/%s (%s)\n",
		Version, Commit, Date, runtime.GOOS, runtime.GOARCH, runtime.Version())
	return 0
}
