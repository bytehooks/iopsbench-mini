package platform

import (
	"runtime"
	"testing"
)

func TestDetect(t *testing.T) {
	info := Detect()
	if info.OS != runtime.GOOS {
		t.Errorf("OS = %s, want %s", info.OS, runtime.GOOS)
	}
	if info.Arch != runtime.GOARCH {
		t.Errorf("Arch = %s, want %s", info.Arch, runtime.GOARCH)
	}
	if info.NumCPU != runtime.NumCPU() {
		t.Errorf("NumCPU = %d, want %d", info.NumCPU, runtime.NumCPU())
	}
}

func TestIOMode(t *testing.T) {
	if got := IOMode(true); got != "Random" {
		t.Errorf("IOMode(true) = %s, want Random", got)
	}
	if got := IOMode(false); got != "Sequential" {
		t.Errorf("IOMode(false) = %s, want Sequential", got)
	}
}
