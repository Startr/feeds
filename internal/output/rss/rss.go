// Package rss implements the RSS2PodcastRenderer — a file output sink that
// writes rewritten feed bytes atomically to disk.
//
// Atomic-ness comes from the classic write-to-temp + rename pattern: readers
// never see a partial file. The temp file is created in the same directory
// as the target so the rename is a simple inode swap, not a cross-device
// copy.
package rss

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

// Renderer writes XML bytes to a target path atomically.
type Renderer struct {
	Path string
}

// NewRenderer returns a Renderer that writes to path.
func NewRenderer(path string) *Renderer {
	return &Renderer{Path: path}
}

// Write reads the existing file at r.Path (if any) and compares to data.
// If they match byte-for-byte, Write is a no-op and returns false. If they
// differ (or the target doesn't exist yet), the bytes are written atomically
// and Write returns true. The unchanged=false return tells callers that a
// downstream CDN cache invalidation might be needed.
func (r *Renderer) Write(data []byte) (changed bool, err error) {
	// If the target exists and matches, skip the write (idempotent).
	if existing, readErr := os.ReadFile(r.Path); readErr == nil {
		if bytes.Equal(existing, data) {
			return false, nil
		}
	}

	dir := filepath.Dir(r.Path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return false, fmt.Errorf("create output dir: %w", err)
	}

	// Write to temp file in the same dir so the rename is atomic on POSIX.
	tmp, err := os.CreateTemp(dir, ".feeds-*.xml.tmp")
	if err != nil {
		return false, fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmp.Name()

	// Clean up the temp file on any failure path.
	defer func() {
		if err != nil {
			_ = os.Remove(tmpPath)
		}
	}()

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return false, fmt.Errorf("write temp file: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return false, fmt.Errorf("fsync temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return false, fmt.Errorf("close temp file: %w", err)
	}

	if err := os.Rename(tmpPath, r.Path); err != nil {
		return false, fmt.Errorf("atomic rename: %w", err)
	}

	return true, nil
}
