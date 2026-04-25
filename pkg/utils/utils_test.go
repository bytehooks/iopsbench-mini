package utils

import "testing"

func TestHumanSize(t *testing.T) {
	tests := []struct {
		bytes int64
		want  string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1024 * 1024, "1.0 MB"},
		{1024 * 1024 * 1024, "1.0 GB"},
	}

	for _, tt := range tests {
		if got := HumanSize(tt.bytes); got != tt.want {
			t.Errorf("HumanSize(%d) = %q, want %q", tt.bytes, got, tt.want)
		}
	}
}
