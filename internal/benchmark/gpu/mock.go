package gpu

// MockBenchmarker is a test-double that returns pre-canned values.
type MockBenchmarker struct {
	InitErr      error
	BenchmarkRes Result
	BenchmarkErr error
}

// Init implements Benchmarker.
func (m *MockBenchmarker) Init() error { return m.InitErr }

// Shutdown implements Benchmarker.
func (m *MockBenchmarker) Shutdown() {}

// BackendName implements Benchmarker.
func (m *MockBenchmarker) BackendName() string {
	return "mock"
}

// Benchmark implements Benchmarker.
func (m *MockBenchmarker) Benchmark(_, _ int) (Result, error) {
	return m.BenchmarkRes, m.BenchmarkErr
}
