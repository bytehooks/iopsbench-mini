// Package config holds application configuration and CLI flag parsing.
package config

import (
	"flag"
	"time"
)

// Config holds all benchmark settings.
type Config struct {
	FileSize  int64
	BlockSize int
	Duration  time.Duration
	TestDir   string
	RandomIO  bool
	ReadRatio float64
	GPU       bool // run GPU benchmarks
}

// ParseFlags parses command-line flags and returns a populated Config.
func ParseFlags() Config {
	var cfg Config
	flag.Int64Var(&cfg.FileSize, "filesize", 500<<20, "Test file size in bytes")
	flag.IntVar(&cfg.BlockSize, "blocksize", 4096, "I/O block size in bytes")
	flag.DurationVar(&cfg.Duration, "duration", 10*time.Second, "Duration of each benchmark")
	flag.StringVar(&cfg.TestDir, "dir", ".", "Directory to create test files")
	flag.BoolVar(&cfg.RandomIO, "random", true, "Run random I/O benchmarks")
	flag.Float64Var(&cfg.ReadRatio, "readratio", 0.5, "Read ratio for mixed benchmarks (0-1)")
	flag.BoolVar(&cfg.GPU, "gpu", true, "Run GPU benchmarks")
	flag.Parse()
	return cfg
}
