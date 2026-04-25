//go:build darwin && !cgo

package gpu

// New returns a no-op benchmarker when CGO is disabled on macOS.
func New() Benchmarker {
	return &noopBenchmarker{}
}
