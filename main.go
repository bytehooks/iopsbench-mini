package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"iopsbench-mini/internal/benchmark/disk"
	"iopsbench-mini/internal/benchmark/gpu"
	"iopsbench-mini/internal/config"
	"iopsbench-mini/internal/platform"
	"iopsbench-mini/internal/reporter"
	"iopsbench-mini/internal/version"
	"iopsbench-mini/pkg/utils"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	showVersion := flag.Bool("version", false, "Print version and exit")
	cfg := config.ParseFlags()

	if *showVersion {
		fmt.Println(version.String())
		os.Exit(0)
	}
	rep := reporter.New(os.Stdout)

	// Disk benchmarks
	runDiskBenchmarks(cfg, rep)

	// GPU benchmarks
	if cfg.GPU {
		runGPUBenchmarks(rep)
	}
}

func runDiskBenchmarks(cfg config.Config, rep *reporter.Reporter) {
	rep.PrintHeader("Disk IOPS Benchmark Tool")
	info := platform.Detect()
	rep.PrintPlatform(info.OS, info.Arch, info.NumCPU)
	rep.PrintDiskConfig(cfg.FileSize, cfg.BlockSize, cfg.Duration, cfg.TestDir, cfg.RandomIO, cfg.ReadRatio)

	rep.Println("Running benchmarks...")
	rep.Println()

	benchmarks := []struct {
		name      string
		random    bool
		readRatio float64
	}{
		{"Sequential Write", false, 0.0},
		{"Sequential Read", false, 1.0},
	}

	if cfg.RandomIO {
		benchmarks = append(benchmarks,
			struct {
				name      string
				random    bool
				readRatio float64
			}{"Random Mixed", true, cfg.ReadRatio},
			struct {
				name      string
				random    bool
				readRatio float64
			}{"Random Read", true, 1.0},
			struct {
				name      string
				random    bool
				readRatio float64
			}{"Random Write", true, 0.0},
		)
	}

	d := disk.Benchmarker{FS: disk.NewOSFS()}
	dcfg := disk.Config{
		FileSize:  cfg.FileSize,
		BlockSize: cfg.BlockSize,
		Duration:  cfg.Duration,
		TestDir:   cfg.TestDir,
	}

	for _, bm := range benchmarks {
		res, err := d.Run(dcfg, bm.name, bm.random, bm.readRatio)
		if err != nil {
			rep.Printf("  %s: error: %v\n", bm.name, err)
			continue
		}
		rep.PrintDiskResult(res)
	}

	_ = d.Cleanup(cfg.TestDir)
}

func runGPUBenchmarks(rep *reporter.Reporter) {
	rep.Println()
	rep.PrintHeader("GPU / CPU Latency Benchmark")
	info := platform.Detect()
	rep.PrintPlatform(info.OS, info.Arch, info.NumCPU)
	rep.Println()

	bench := gpu.New()
	rep.Printf("GPU Backend: %s\n", bench.BackendName())
	rep.Println()

	if err := bench.Init(); err != nil {
		rep.Printf("GPU init error: %v\n", err)
		return
	}
	defer bench.Shutdown()

	rep.PrintGPUHeader()

	sizes := []int{1024, 64 * 1024, 1024 * 1024, 16 * 1024 * 1024}
	iterations := 100

	for _, size := range sizes {
		res, err := bench.Benchmark(size, iterations)
		if err != nil {
			rep.Printf("  %7s: error: %v\n", utils.HumanSize(int64(size)), err)
			continue
		}
		rep.PrintGPUResult(reporter.GPUResults{
			BufferSize:      res.BufferSize,
			H2DLatencyUS:    res.H2DLatencyUS,
			D2HLatencyUS:    res.D2HLatencyUS,
			H2DBandwidthGBs: res.H2DBandwidthGBs,
			D2HBandwidthGBs: res.D2HBandwidthGBs,
		})
	}

	res, _ := bench.Benchmark(4, 1000)
	if res.KernelLatencyUS > 0 {
		rep.PrintKernelLatency(res.KernelLatencyUS)
	}

	rep.Println()
	rep.Println("Note: On Apple Silicon, unified memory may make private/shared")
	rep.Println("      copies behave differently than on discrete GPUs.")
}
