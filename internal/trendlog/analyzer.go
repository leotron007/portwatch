package trendlog

import "time"

// Direction indicates whether the port count trend is rising, falling, or stable.
type Direction string

const (
	Rising  Direction = "rising"
	Falling Direction = "falling"
	Stable  Direction = "stable"
)

// Summary holds the result of a trend analysis over a window of entries.
type Summary struct {
	Window    time.Duration
	Entries   int
	NetDelta  int
	PeakOpen  int
	Direction Direction
}

// Analyze computes a trend summary for entries that fall within the given
// duration window ending at now. If fewer than two entries exist in the
// window, Direction is Stable and NetDelta is 0.
func Analyze(entries []Entry, window time.Duration, now time.Time) Summary {
	cutoff := now.Add(-window)
	var windowed []Entry
	for _, e := range entries {
		if !e.Timestamp.Before(cutoff) {
			windowed = append(windowed, e)
		}
	}

	s := Summary{Window: window, Entries: len(windowed)}
	if len(windowed) == 0 {
		s.Direction = Stable
		return s
	}

	for _, e := range windowed {
		if e.OpenCount > s.PeakOpen {
			s.PeakOpen = e.OpenCount
		}
		s.NetDelta += e.Delta
	}

	switch {
	case s.NetDelta > 0:
		s.Direction = Rising
	case s.NetDelta < 0:
		s.Direction = Falling
	default:
		s.Direction = Stable
	}
	return s
}
