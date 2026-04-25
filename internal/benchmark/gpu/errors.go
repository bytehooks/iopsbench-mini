package gpu

import "errors"

// Common errors returned by GPU benchmarkers.
var (
	ErrNoDevice        = errors.New("no GPU device available")
	ErrBenchmarkFailed = errors.New("GPU benchmark failed")
)
