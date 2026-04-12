package notify

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/alert"
)

// Dispatcher fans out an event to multiple Channel implementations.
type Dispatcher struct {
	mu       sync.RWMutex
	channels []Channel
}

// NewDispatcher returns an empty Dispatcher.
func NewDispatcher() *Dispatcher {
	return &Dispatcher{}
}

// Register adds a channel to the dispatcher.
func (d *Dispatcher) Register(ch Channel) {
	if ch == nil {
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	d.channels = append(d.channels, ch)
}

// Dispatch sends the event to every registered channel, collecting errors.
func (d *Dispatcher) Dispatch(event alert.Event) error {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var errs []error
	for _, ch := range d.channels {
		if err := ch.Send(event); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("notify: %d channel(s) failed: %v", len(errs), errs)
}

// Len returns the number of registered channels.
func (d *Dispatcher) Len() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return len(d.channels)
}
