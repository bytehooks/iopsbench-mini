// Package platform provides host-platform detection utilities.
package platform

import "runtime"

// Info holds basic platform metadata.
type Info struct {
	OS     string
	Arch   string
	NumCPU int
}

// Detect returns the current platform information.
func Detect() Info {
	return Info{
		OS:     runtime.GOOS,
		Arch:   runtime.GOARCH,
		NumCPU: runtime.NumCPU(),
	}
}

// IOMode returns "Random" or "Sequential" based on the boolean.
func IOMode(random bool) string {
	if random {
		return "Random"
	}
	return "Sequential"
}
