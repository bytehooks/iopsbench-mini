package disk

import (
	"errors"
	"io"
	"os"
	"testing"
	"time"
)

// mockFile is an in-memory file for testing.
type mockFile struct {
	data   []byte
	offset int64
	closed bool
}

func (m *mockFile) Read(p []byte) (int, error) {
	if m.closed {
		return 0, errors.New("file closed")
	}
	n := copy(p, m.data[m.offset:])
	m.offset += int64(n)
	if n == 0 {
		return 0, io.EOF
	}
	return n, nil
}

func (m *mockFile) Write(p []byte) (int, error) {
	if m.closed {
		return 0, errors.New("file closed")
	}
	need := int(m.offset) + len(p)
	if need > len(m.data) {
		newData := make([]byte, need)
		copy(newData, m.data)
		m.data = newData
	}
	n := copy(m.data[m.offset:], p)
	m.offset += int64(n)
	return n, nil
}

func (m *mockFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		m.offset = offset
	case io.SeekCurrent:
		m.offset += offset
	case io.SeekEnd:
		m.offset = int64(len(m.data)) + offset
	}
	return m.offset, nil
}

func (m *mockFile) Sync() error  { return nil }
func (m *mockFile) Close() error { m.closed = true; return nil }

func (m *mockFile) Stat() (os.FileInfo, error) {
	return &mockFileInfo{size: int64(len(m.data))}, nil
}

type mockFileInfo struct {
	size int64
}

func (fi *mockFileInfo) Name() string       { return "mock" }
func (fi *mockFileInfo) Size() int64        { return fi.size }
func (fi *mockFileInfo) Mode() os.FileMode  { return 0644 }
func (fi *mockFileInfo) ModTime() time.Time { return time.Now() }
func (fi *mockFileInfo) IsDir() bool        { return false }
func (fi *mockFileInfo) Sys() interface{}   { return nil }

// mockFS implements disk.FS using in-memory files.
type mockFS struct {
	files map[string]*mockFile
}

func newMockFS() *mockFS {
	return &mockFS{files: make(map[string]*mockFile)}
}

func (fs *mockFS) OpenFile(name string, _ int, _ os.FileMode) (File, error) {
	if f, ok := fs.files[name]; ok {
		return f, nil
	}
	f := &mockFile{data: make([]byte, 0)}
	fs.files[name] = f
	return f, nil
}

func (fs *mockFS) Remove(name string) error {
	delete(fs.files, name)
	return nil
}

func TestDiskBenchmarkMock(t *testing.T) {
	fs := newMockFS()
	b := &Benchmarker{FS: fs}

	cfg := Config{
		FileSize:  64 * 1024,
		BlockSize: 4096,
		Duration:  100 * time.Millisecond,
		TestDir:   ".",
	}

	res, err := b.Run(cfg, "MockTest", false, 1.0)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if res.Name != "MockTest" {
		t.Errorf("Name = %q, want MockTest", res.Name)
	}
	if res.ReadOps == 0 {
		t.Error("expected some read operations")
	}
	if err := b.Cleanup("."); err != nil {
		t.Fatalf("Cleanup failed: %v", err)
	}
}
