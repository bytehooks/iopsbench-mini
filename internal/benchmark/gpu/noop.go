package gpu

import "errors"

type noopBenchmarker struct{}

func (n *noopBenchmarker) Init() error {
	return errors.New("GPU benchmarking not supported on this platform")
}

func (n *noopBenchmarker) Shutdown() {}

func (n *noopBenchmarker) BackendName() string {
	return "noop"
}

func (n *noopBenchmarker) Benchmark(_, _ int) (Result, error) {
	return Result{}, errors.New("GPU benchmarking not supported on this platform")
}
