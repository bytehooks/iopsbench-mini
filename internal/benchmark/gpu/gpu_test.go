package gpu

import "testing"

func TestMockBenchmarker(t *testing.T) {
	mock := &MockBenchmarker{
		BenchmarkRes: Result{
			BufferSize:      1024,
			H2DLatencyUS:    5.0,
			D2HLatencyUS:    6.0,
			H2DBandwidthGBs: 0.2,
			D2HBandwidthGBs: 0.15,
			KernelLatencyUS: 1.0,
		},
	}

	if err := mock.Init(); err != nil {
		t.Fatalf("Init error: %v", err)
	}
	mock.Shutdown()

	res, err := mock.Benchmark(1024, 10)
	if err != nil {
		t.Fatalf("Benchmark error: %v", err)
	}
	if res.H2DLatencyUS != 5.0 {
		t.Errorf("H2DLatencyUS = %f, want 5.0", res.H2DLatencyUS)
	}
}

func TestNoopBenchmarker(t *testing.T) {
	n := &noopBenchmarker{}
	if err := n.Init(); err == nil {
		t.Error("expected error from noop Init")
	}
	if _, err := n.Benchmark(0, 0); err == nil {
		t.Error("expected error from noop Benchmark")
	}
}
