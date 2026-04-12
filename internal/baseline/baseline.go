// Package baseline manages the trusted set of ports that are expected to be
// open on the host. It persists the baseline to disk so that portwatch can
// distinguish genuinely new ports from ports that were already known.
package baseline

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

// Entry represents a single port that has been accepted into the baseline.
type Entry struct {
	Port      int       `json:"port"`
	AddedAt   time.Time `json:"added_at"`
	AddedBy   string    `json:"added_by"`
}

// Baseline holds the current set of trusted ports.
type Baseline struct {
	mu      sync.RWMutex
	path    string
	entries map[int]Entry
}

// New loads an existing baseline file from path, or returns an empty Baseline
// if the file does not yet exist. A corrupt file returns an error.
func New(path string) (*Baseline, error) {
	b := &Baseline{
		path:    path,
		entries: make(map[int]Entry),
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return b, nil
	}
	if err != nil {
		return nil, err
	}
	var list []Entry
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	for _, e := range list {
		b.entries[e.Port] = e
	}
	return b, nil
}

// Contains reports whether port is part of the trusted baseline.
func (b *Baseline) Contains(port int) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	_, ok := b.entries[port]
	return ok
}

// Add adds port to the baseline and persists the updated set to disk.
func (b *Baseline) Add(port int, addedBy string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.entries[port] = Entry{Port: port, AddedAt: time.Now().UTC(), AddedBy: addedBy}
	return b.save()
}

// Remove removes port from the baseline and persists the change.
func (b *Baseline) Remove(port int) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.entries, port)
	return b.save()
}

// Entries returns a snapshot of all baseline entries.
func (b *Baseline) Entries() []Entry {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := make([]Entry, 0, len(b.entries))
	for _, e := range b.entries {
		out = append(out, e)
	}
	return out
}

// save serialises the current entries to disk. Caller must hold b.mu.
func (b *Baseline) save() error {
	list := make([]Entry, 0, len(b.entries))
	for _, e := range b.entries {
		list = append(list, e)
	}
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(b.path, data, 0o600)
}
