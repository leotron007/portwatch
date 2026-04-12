// Package acknowledge provides a mechanism for tracking acknowledged port
// change events so that repeated alerts for the same condition can be
// suppressed until the operator explicitly clears the acknowledgement.
package acknowledge

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// Entry records when a specific port event was acknowledged.
type Entry struct {
	Port      int       `json:"port"`
	Event     string    `json:"event"` // "opened" or "closed"
	AckedAt   time.Time `json:"acked_at"`
}

// Store persists acknowledged events to a JSON file.
type Store struct {
	mu      sync.RWMutex
	path    string
	entries map[string]Entry
}

func key(port int, event string) string {
	return fmt.Sprintf("%d:%s", port, event)
}

// New loads an existing acknowledgement store from path, or returns an empty
// store if the file does not yet exist.
func New(path string) (*Store, error) {
	s := &Store{path: path, entries: make(map[string]Entry)}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return s, nil
	}
	if err != nil {
		return nil, fmt.Errorf("acknowledge: read %s: %w", path, err)
	}
	var list []Entry
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, fmt.Errorf("acknowledge: parse %s: %w", path, err)
	}
	for _, e := range list {
		s.entries[key(e.Port, e.Event)] = e
	}
	return s, nil
}

// Ack marks a port/event pair as acknowledged at the given time.
func (s *Store) Ack(port int, event string, at time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[key(port, event)] = Entry{Port: port, Event: event, AckedAt: at}
	return s.save()
}

// Clear removes an acknowledgement for the given port/event pair.
func (s *Store) Clear(port int, event string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, key(port, event))
	return s.save()
}

// IsAcked reports whether the given port/event pair is currently acknowledged.
func (s *Store) IsAcked(port int, event string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.entries[key(port, event)]
	return ok
}

// All returns a copy of all current acknowledgement entries.
func (s *Store) All() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		out = append(out, e)
	}
	return out
}

func (s *Store) save() error {
	list := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		list = append(list, e)
	}
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return fmt.Errorf("acknowledge: marshal: %w", err)
	}
	if err := os.WriteFile(s.path, data, 0o644); err != nil {
		return fmt.Errorf("acknowledge: write %s: %w", s.path, err)
	}
	return nil
}
