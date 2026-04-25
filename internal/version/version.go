// Package version holds build-time version info.
package version

import "fmt"

// Version is set at build time via -ldflags.
var Version = "dev"

// String returns the formatted version.
func String() string {
	return fmt.Sprintf("iopsbench-mini %s", Version)
}
