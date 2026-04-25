// Package disk implements disk I/O benchmarking.
package disk

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"iopsbench-mini/internal/reporter"
)

// FS abstracts file-system operations for testability.
type FS interface {
	OpenFile(name string, flag int, perm os.FileMode) (File, error)
	Remove(name string) error
}

// File abstracts file handle operations.
type File interface {
	io.ReadWriteSeeker
	Sync() error
	Close() error
	Stat() (os.FileInfo, error)
}

// osFS is the production filesystem implementation.
type osFS struct{}

func (osFS) OpenFile(name string, flag int, perm os.FileMode) (File, error) {
	return os.OpenFile(name, flag, perm)
}

func (osFS) Remove(name string) error {
	return os.Remove(name)
}

// NewOSFS returns a real OS-backed FS.
func NewOSFS() FS {
	return osFS{}
}

// Benchmarker runs disk benchmarks.
type Benchmarker struct {
	FS FS
}

// Run executes a single benchmark and returns results.
func (b *Benchmarker) Run(cfg Config, name string, random bool, readRatio float64) (reporter.DiskResults, error) {
	if b.FS == nil {
		b.FS = NewOSFS()
	}

	filePath := filepath.Join(cfg.TestDir, "iopstest.tmp")
	f, err := b.FS.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return reporter.DiskResults{}, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	if fi, _ := f.Stat(); fi.Size() < cfg.FileSize {
		data := make([]byte, cfg.BlockSize)
		for i := range data {
			data[i] = byte(rand.Intn(256))
		}
		written := int64(0)
		for written < cfg.FileSize {
			n, err := f.Write(data)
			if err != nil {
				return reporter.DiskResults{}, fmt.Errorf("write: %w", err)
			}
			written += int64(n)
		}
		if err := f.Sync(); err != nil {
			return reporter.DiskResults{}, fmt.Errorf("sync: %w", err)
		}
	}

	if !random {
		if _, err := f.Seek(0, io.SeekStart); err != nil {
			return reporter.DiskResults{}, fmt.Errorf("seek: %w", err)
		}
	}

	blockCount := int(cfg.FileSize / int64(cfg.BlockSize))
	if blockCount == 0 {
		blockCount = 1
	}
	buf := make([]byte, cfg.BlockSize)
	for i := range buf {
		buf[i] = byte(rand.Intn(256))
	}

	var readOps, writeOps uint64
	var readBytes, writeBytes uint64
	start := time.Now()
	deadline := start.Add(cfg.Duration)
	latencySum := time.Duration(0)
	opsCount := uint64(0)

	for time.Now().Before(deadline) {
		opStart := time.Now()
		isRead := rand.Float64() < readRatio

		var offset int64
		if random {
			offset = int64(rand.Intn(blockCount)) * int64(cfg.BlockSize)
		} else {
			offset, err = f.Seek(0, io.SeekCurrent)
			if err != nil || offset >= cfg.FileSize {
				_, _ = f.Seek(0, io.SeekStart)
				offset = 0
			}
		}

		if _, err := f.Seek(offset, io.SeekStart); err != nil {
			break
		}

		if isRead {
			n, err := f.Read(buf)
			if err != nil && err != io.EOF {
				break
			}
			readOps++
			readBytes += uint64(n)
		} else {
			n, err := f.Write(buf)
			if err != nil {
				break
			}
			writeOps++
			writeBytes += uint64(n)
		}

		latencySum += time.Since(opStart)
		opsCount++
	}

	elapsed := time.Since(start)
	_ = f.Sync()

	avgLatency := time.Duration(0)
	if opsCount > 0 {
		avgLatency = latencySum / time.Duration(opsCount)
	}

	return reporter.DiskResults{
		Name:       name,
		ReadOps:    readOps,
		WriteOps:   writeOps,
		ReadBytes:  readBytes,
		WriteBytes: writeBytes,
		Duration:   elapsed,
		AvgLatency: avgLatency,
	}, nil
}

// Cleanup removes the test file.
func (b *Benchmarker) Cleanup(dir string) error {
	if b.FS == nil {
		b.FS = NewOSFS()
	}
	return b.FS.Remove(filepath.Join(dir, "iopstest.tmp"))
}

// Config is the local disk benchmark configuration.
type Config struct {
	FileSize  int64
	BlockSize int
	Duration  time.Duration
	TestDir   string
}
