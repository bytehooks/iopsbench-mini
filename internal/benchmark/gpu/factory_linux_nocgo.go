//go:build linux && !cgo

package gpu

// New returns a no-op benchmarker when CGO is disabled on Linux.
func New() Benchmarker {
	return &noopBenchmarker{}
}
