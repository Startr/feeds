// Package cache persists conditional-GET state (ETag, Last-Modified) to a
// small JSON file so the next run can short-circuit on HTTP 304.
//
// The eng review target is ~98% short-circuit rate in steady state: for a
// feed that updates weekly, a 15-minute cron hits 304 on 671 of the 672
// weekly runs.
package cache

import (
	"encoding/json"
	"errors"
	"os"
)

// State is the JSON shape written to disk.
type State struct {
	ETag         string `json:"etag,omitempty"`
	LastModified string `json:"last_modified,omitempty"`
}

// Cache reads and writes a single state file.
type Cache struct {
	Path string
}

// New returns a Cache backed by path.
func New(path string) *Cache {
	return &Cache{Path: path}
}

// Load reads the state from disk. A missing file returns a zero State and
// nil error — that's the first-run case.
func (c *Cache) Load() (State, error) {
	raw, err := os.ReadFile(c.Path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return State{}, nil
		}
		return State{}, err
	}

	var s State
	if err := json.Unmarshal(raw, &s); err != nil {
		// Corrupt state file — treat as first run rather than failing loud.
		// The worst case is one wasted full-GET.
		return State{}, nil
	}
	return s, nil
}

// Save writes the state to disk. Not atomic because state file corruption
// is recoverable (worst case: one wasted full-GET on next run).
func (c *Cache) Save(s State) error {
	raw, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.Path, raw, 0o644)
}
