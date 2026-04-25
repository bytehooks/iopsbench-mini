// Package utils provides shared utility functions used across the codebase.
package utils

import "fmt"

// HumanSize formats bytes as a human-readable string (e.g., 1.0 KB, 16.0 MB).
func HumanSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
