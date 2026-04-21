package trendlog

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Entry records a single port-count observation used for trend analysis.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	OpenCount int       `json:"open_count"`
	Delta     int       `json:"delta"` // change relative to previous entry
}

// TrendLog accumulates periodic open-port counts and persists them to a
// newline-delimited JSON file so that long-running trend analysis is possible.
type TrendLog struct {
	mu      sync.Mutex
	entries []Entry
	path    string
	w       io.Writer
}

// New opens (or creates) the trend log at path. Existing entries are loaded
// into memory. Pass a nil writer to default to os.Stderr for error output.
func New(path string, w io.Writer) (*TrendLog, error) {
	if w == nil {
		w = os.Stderr
	}
	tl := &TrendLog{path: path, w: w}
	if err := tl.load(); err != nil {
		return nil, err
	}
	return tl, nil
}

// Record appends a new observation. The delta is computed automatically.
func (tl *TrendLog) Record(openCount int, at time.Time) error {
	tl.mu.Lock()
	defer tl.mu.Unlock()

	delta := openCount
	if len(tl.entries) > 0 {
		delta = openCount - tl.entries[len(tl.entries)-1].OpenCount
	}
	e := Entry{Timestamp: at, OpenCount: openCount, Delta: delta}
	tl.entries = append(tl.entries, e)
	return tl.append(e)
}

// Entries returns a snapshot of all recorded entries.
func (tl *TrendLog) Entries() []Entry {
	tl.mu.Lock()
	defer tl.mu.Unlock()
	out := make([]Entry, len(tl.entries))
	copy(out, tl.entries)
	return out
}

// load reads existing newline-delimited JSON entries from disk.
func (tl *TrendLog) load() error {
	f, err := os.Open(tl.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("trendlog: open %s: %w", tl.path, err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	for dec.More() {
		var e Entry
		if err := dec.Decode(&e); err != nil {
			return fmt.Errorf("trendlog: decode %s: %w", tl.path, err)
		}
		tl.entries = append(tl.entries, e)
	}
	return nil
}

// append writes a single entry to the end of the file.
func (tl *TrendLog) append(e Entry) error {
	f, err := os.OpenFile(tl.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("trendlog: open for append %s: %w", tl.path, err)
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(e)
}
