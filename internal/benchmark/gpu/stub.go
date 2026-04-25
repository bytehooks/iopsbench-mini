//go:build stub

package gpu

// New returns a no-op benchmarker when built with the stub tag.
func New() Benchmarker {
	return &noopBenchmarker{}
}
