//go:build darwin && cgo

package gpu

/*
#cgo LDFLAGS: -framework Metal -framework Foundation
#include "metal_backend.h"
*/
import "C"

// New returns the Metal-backed GPU benchmarker for macOS.
func New() Benchmarker {
	return &metalBenchmarker{}
}

func newMetalBenchmarker() Benchmarker {
	return &metalBenchmarker{}
}

type metalBenchmarker struct{}

func (m *metalBenchmarker) Init() error {
	if C.metalInit() != 0 {
		return ErrNoDevice
	}
	return nil
}

func (m *metalBenchmarker) Shutdown() {
	C.metalShutdown()
}

func (m *metalBenchmarker) BackendName() string {
	return "Metal"
}

func (m *metalBenchmarker) Benchmark(bufferSize, iterations int) (Result, error) {
	res := C.metalBenchmark(C.int(bufferSize), C.int(iterations))
	if res.ok == 0 {
		return Result{}, ErrBenchmarkFailed
	}
	return Result{
		BufferSize:      int64(bufferSize),
		H2DLatencyUS:    float64(res.h2d_latency_us),
		D2HLatencyUS:    float64(res.d2h_latency_us),
		H2DBandwidthGBs: float64(res.h2d_bw_gbps),
		D2HBandwidthGBs: float64(res.d2h_bw_gbps),
		KernelLatencyUS: float64(res.kernel_latency_us),
	}, nil
}
