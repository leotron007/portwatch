package digest

import "sync"

// Tracker maintains the last known digest and reports whether the port
// set has changed between consecutive observations.
type Tracker struct {
	mu   sync.Mutex
	last *Digest
}

// NewTracker returns an initialised Tracker with no prior digest.
func NewTracker() *Tracker {
	return &Tracker{}
}

// Update records a new observation and returns (changed, previous, current).
// On the very first call changed is always false because there is no prior
// state to compare against.
func (t *Tracker) Update(ports []int) (changed bool, prev *Digest, curr *Digest) {
	curr = New(ports)

	t.mu.Lock()
	defer t.mu.Unlock()

	prev = t.last
	if prev != nil && !prev.Equal(curr) {
		changed = true
	}
	t.last = curr
	return changed, prev, curr
}

// Last returns the most recently recorded digest, or nil if Update has never
// been called.
func (t *Tracker) Last() *Digest {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.last
}

// Reset clears the stored digest so the next Update is treated as a first
// observation.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.last = nil
}
