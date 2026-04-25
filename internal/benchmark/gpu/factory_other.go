//go:build !darwin && !linux

package gpu

// New returns a no-op benchmarker on unsupported platforms.
func New() Benchmarker {
	return &noopBenchmarker{}
}
