// Package reporter handles formatting and printing benchmark results.
package reporter

import (
	"fmt"
	"io"
	"time"

	"iopsbench-mini/pkg/utils"
)

// Reporter formats and outputs benchmark results.
type Reporter struct {
	out io.Writer
}

// New creates a Reporter that writes to w.
func New(w io.Writer) *Reporter {
	return &Reporter{out: w}
}

// Printf writes formatted output.
func (r *Reporter) Printf(format string, args ...interface{}) {
	fmt.Fprintf(r.out, format, args...)
}

// Println writes a line of output.
func (r *Reporter) Println(args ...interface{}) {
	fmt.Fprintln(r.out, args...)
}

// DiskResults holds the outcome of a disk benchmark run.
type DiskResults struct {
	Name       string
	ReadOps    uint64
	WriteOps   uint64
	ReadBytes  uint64
	WriteBytes uint64
	Duration   time.Duration
	AvgLatency time.Duration
}

// GPUResults holds the outcome of a GPU benchmark run.
type GPUResults struct {
	BufferSize      int64
	H2DLatencyUS    float64
	D2HLatencyUS    float64
	H2DBandwidthGBs float64
	D2HBandwidthGBs float64
	KernelLatencyUS float64
}

// PrintHeader prints the top-level benchmark header.
func (r *Reporter) PrintHeader(title string) {
	r.Println("╔══════════════════════════════════════════╗")
	r.Printf("║ %-40s ║\n", title)
	r.Println("╚══════════════════════════════════════════╝")
}

// PrintPlatform prints OS/arch/CPU info.
func (r *Reporter) PrintPlatform(os, arch string, cpus int) {
	r.Printf("OS:           %s\n", os)
	r.Printf("Arch:         %s\n", arch)
	r.Printf("CPUs:         %d\n", cpus)
}

// PrintDiskConfig prints disk benchmark configuration.
func (r *Reporter) PrintDiskConfig(fileSize int64, blockSize int, duration time.Duration, dir string, random bool, readRatio float64) {
	r.Printf("Test Dir:     %s\n", dir)
	r.Printf("File Size:    %s\n", utils.HumanSize(fileSize))
	r.Printf("Block Size:   %s\n", utils.HumanSize(int64(blockSize)))
	r.Printf("Duration:     %s\n", duration)
	r.Printf("I/O Mode:     %s\n", ioMode(random))
	r.Printf("Read Ratio:   %.0f%%\n", readRatio*100)
	r.Println()
}

// PrintDiskResult prints a single disk benchmark result line.
func (r *Reporter) PrintDiskResult(res DiskResults) {
	d := res.Duration.Seconds()
	if d == 0 {
		d = 1e-9
	}
	totalOps := res.ReadOps + res.WriteOps
	totalBytes := res.ReadBytes + res.WriteBytes
	totalIOPS := float64(totalOps) / d
	totalMBps := float64(totalBytes) / d / 1024 / 1024

	r.Printf("  %s  IOPS: %8.1f  |  MB/s: %7.2f  |  Latency: %s\n",
		res.Name, totalIOPS, totalMBps, res.AvgLatency)

	if res.ReadOps > 0 && res.WriteOps > 0 {
		r.Printf("    Reads : %8.1f IOPS (%7.2f MB/s)\n", float64(res.ReadOps)/d, float64(res.ReadBytes)/d/1024/1024)
		r.Printf("    Writes: %8.1f IOPS (%7.2f MB/s)\n", float64(res.WriteOps)/d, float64(res.WriteBytes)/d/1024/1024)
	}
	r.Println()
}

// PrintGPUHeader prints the GPU section header.
func (r *Reporter) PrintGPUHeader() {
	r.Println("Transfer Benchmarks (Private ↔ Shared buffer copy):")
	r.Println()
}

// PrintGPUResult prints a single GPU benchmark result line.
func (r *Reporter) PrintGPUResult(res GPUResults) {
	r.Printf("  Buffer: %7s  |  H→D: %8.2f µs  %6.2f GB/s  |  D→H: %8.2f µs  %6.2f GB/s\n",
		utils.HumanSize(res.BufferSize),
		res.H2DLatencyUS, res.H2DBandwidthGBs,
		res.D2HLatencyUS, res.D2HBandwidthGBs)
}

// PrintKernelLatency prints the kernel latency.
func (r *Reporter) PrintKernelLatency(latencyUS float64) {
	r.Println()
	r.Printf("Kernel round-trip latency: %8.3f µs\n", latencyUS)
}

func ioMode(random bool) string {
	if random {
		return "Random"
	}
	return "Sequential"
}
