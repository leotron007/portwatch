// Package snapshot captures and compares periodic port scan results,
// producing a stable digest-keyed record suitable for diffing over time.
package snapshot

import (
	"fmt"
	"time"

	"github.com/user/portwatch/internal/digest"
)

// Snapshot holds a single point-in-time capture of open ports.
type Snapshot struct {
	CapturedAt time.Time
	Ports      []int
	Digest     *digest.Digest
}

// New creates a Snapshot from the given list of open ports.
// Ports are deduplicated and sorted inside the digest.
func New(ports []int, capturedAt time.Time) (*Snapshot, error) {
	if capturedAt.IsZero() {
		return nil, fmt.Errorf("snapshot: capturedAt must not be zero")
	}
	d, err := digest.New(ports)
	if err != nil {
		return nil, fmt.Errorf("snapshot: %w", err)
	}
	return &Snapshot{
		CapturedAt: capturedAt,
		Ports:      d.Ports(),
		Digest:     d,
	}, nil
}

// Equal returns true when both snapshots carry the same port set.
func (s *Snapshot) Equal(other *Snapshot) bool {
	if other == nil {
		return false
	}
	return s.Digest.Equal(other.Digest)
}

// Added returns ports present in other but not in s.
func (s *Snapshot) Added(other *Snapshot) []int {
	if other == nil {
		return nil
	}
	current := toSet(s.Ports)
	var added []int
	for _, p := range other.Ports {
		if !current[p] {
			added = append(added, p)
		}
	}
	return added
}

// Removed returns ports present in s but not in other.
func (s *Snapshot) Removed(other *Snapshot) []int {
	if other == nil {
		return nil
	}
	next := toSet(other.Ports)
	var removed []int
	for _, p := range s.Ports {
		if !next[p] {
			removed = append(removed, p)
		}
	}
	return removed
}

func toSet(ports []int) map[int]bool {
	m := make(map[int]bool, len(ports))
	for _, p := range ports {
		m[p] = true
	}
	return m
}
