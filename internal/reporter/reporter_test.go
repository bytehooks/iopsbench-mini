package reporter

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestPrintDiskResult(t *testing.T) {
	var buf bytes.Buffer
	rep := New(&buf)
	rep.PrintDiskResult(DiskResults{
		Name:       "Test",
		ReadOps:    100,
		WriteOps:   200,
		ReadBytes:  1024 * 1024,
		WriteBytes: 2 * 1024 * 1024,
		Duration:   time.Second,
		AvgLatency: 100 * time.Microsecond,
	})
	out := buf.String()
	if !strings.Contains(out, "Test") {
		t.Error("missing benchmark name in output")
	}
	if !strings.Contains(out, "300.0") {
		t.Error("missing expected total IOPS")
	}
}

func TestPrintGPUResult(t *testing.T) {
	var buf bytes.Buffer
	rep := New(&buf)
	rep.PrintGPUResult(GPUResults{
		BufferSize:      1024,
		H2DLatencyUS:    10.0,
		D2HLatencyUS:    12.0,
		H2DBandwidthGBs: 0.1,
		D2HBandwidthGBs: 0.08,
	})
	out := buf.String()
	if !strings.Contains(out, "1.0 KB") {
		t.Error("missing buffer size in output")
	}
}
