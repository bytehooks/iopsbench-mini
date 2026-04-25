// Package gpu abstracts GPU benchmarking across different platforms.
package gpu

// Result holds the outcome of a GPU benchmark.
type Result struct {
	BufferSize      int64
	H2DLatencyUS    float64
	D2HLatencyUS    float64
	H2DBandwidthGBs float64
	D2HBandwidthGBs float64
	KernelLatencyUS float64
}

// Benchmarker defines the contract for GPU benchmarks.
type Benchmarker interface {
	Init() error
	Shutdown()
	BackendName() string
	Benchmark(bufferSize, iterations int) (Result, error)
}
