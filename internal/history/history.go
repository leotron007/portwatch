package history

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Entry records a single port-change event for historical reference.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Port      int       `json:"port"`
	Event     string    `json:"event"` // "opened" or "closed"
	Protocol  string    `json:"protocol"`
}

// History maintains an in-memory log of port-change events and can persist
// them to a newline-delimited JSON file.
type History struct {
	mu      sync.Mutex
	entries []Entry
	path    string
}

// New creates a History that persists to path. Existing entries are loaded
// from the file if it already exists.
func New(path string) (*History, error) {
	h := &History{path: path}
	if err := h.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return h, nil
}

// Record appends an entry and flushes to disk.
func (h *History) Record(port int, event, protocol string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	e := Entry{
		Timestamp: time.Now().UTC(),
		Port:      port,
		Event:     event,
		Protocol:  protocol,
	}
	h.entries = append(h.entries, e)
	return h.flush()
}

// Entries returns a copy of all recorded entries.
func (h *History) Entries() []Entry {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make([]Entry, len(h.entries))
	copy(out, h.entries)
	return out
}

func (h *History) load() error {
	data, err := os.ReadFile(h.path)
	if err != nil {
		return err
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return err
	}
	h.entries = entries
	return nil
}

func (h *History) flush() error {
	data, err := json.Marshal(h.entries)
	if err != nil {
		return err
	}
	return os.WriteFile(h.path, data, 0o644)
}
