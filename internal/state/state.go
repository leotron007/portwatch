package state

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Snapshot holds the last known set of open ports at a point in time.
type Snapshot struct {
	Timestamp time.Time `json:"timestamp"`
	Ports     []int     `json:"ports"`
}

// Store persists and retrieves port snapshots.
type Store struct {
	mu   sync.RWMutex
	path string
	current Snapshot
}

// New creates a Store backed by the given file path.
// If the file exists its contents are loaded immediately.
func New(path string) (*Store, error) {
	s := &Store{path: path}
	if err := s.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return s, nil
}

// Current returns the most recently stored snapshot.
func (s *Store) Current() Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.current
}

// Save replaces the current snapshot and flushes it to disk.
func (s *Store) Save(ports []int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.current = Snapshot{Timestamp: time.Now().UTC(), Ports: ports}
	return s.flush()
}

func (s *Store) load() error {
	f, err := os.Open(s.path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(&s.current)
}

func (s *Store) flush() error {
	f, err := os.CreateTemp("", "portwatch-state-*")
	if err != nil {
		return err
	}
	if err := json.NewEncoder(f).Encode(s.current); err != nil {
		f.Close()
		return err
	}
	f.Close()
	return os.Rename(f.Name(), s.path)
}
