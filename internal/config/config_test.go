package config

import (
	"testing"
	"time"
)

func TestParseFlagsDefaults(t *testing.T) {
	cfg := ParseFlags()
	if cfg.FileSize != 500<<20 {
		t.Errorf("FileSize default wrong: got %d, want %d", cfg.FileSize, 500<<20)
	}
	if cfg.BlockSize != 4096 {
		t.Errorf("BlockSize default wrong: got %d, want %d", cfg.BlockSize, 4096)
	}
	if cfg.Duration != 10*time.Second {
		t.Errorf("Duration default wrong: got %v, want %v", cfg.Duration, 10*time.Second)
	}
	if cfg.TestDir != "." {
		t.Errorf("TestDir default wrong: got %s, want .", cfg.TestDir)
	}
	if !cfg.RandomIO {
		t.Error("RandomIO default should be true")
	}
	if cfg.ReadRatio != 0.5 {
		t.Errorf("ReadRatio default wrong: got %f, want 0.5", cfg.ReadRatio)
	}
	if !cfg.GPU {
		t.Error("GPU default should be true")
	}
}
