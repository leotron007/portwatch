package trendlog_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/trendlog"
)

var base = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func makeEntries() []trendlog.Entry {
	return []trendlog.Entry{
		{Timestamp: base, OpenCount: 5, Delta: 5},
		{Timestamp: base.Add(10 * time.Minute), OpenCount: 8, Delta: 3},
		{Timestamp: base.Add(20 * time.Minute), OpenCount: 6, Delta: -2},
	}
}

func TestAnalyze_RisingTrend(t *testing.T) {
	entries := makeEntries()[:2] // delta sum = +8, net positive
	now := base.Add(25 * time.Minute)
	s := trendlog.Analyze(entries, 30*time.Minute, now)
	if s.Direction != trendlog.Rising {
		t.Errorf("direction: want Rising, got %s", s.Direction)
	}
	if s.NetDelta != 8 {
		t.Errorf("net delta: want 8, got %d", s.NetDelta)
	}
}

func TestAnalyze_FallingTrend(t *testing.T) {
	entries := []trendlog.Entry{
		{Timestamp: base, OpenCount: 10, Delta: 10},
		{Timestamp: base.Add(5 * time.Minute), OpenCount: 4, Delta: -6},
		{Timestamp: base.Add(10 * time.Minute), OpenCount: 2, Delta: -2},
	}
	now := base.Add(15 * time.Minute)
	s := trendlog.Analyze(entries, 30*time.Minute, now)
	if s.Direction != trendlog.Falling {
		t.Errorf("direction: want Falling, got %s", s.Direction)
	}
}

func TestAnalyze_StableWhenNoDelta(t *testing.T) {
	entries := []trendlog.Entry{
		{Timestamp: base, OpenCount: 5, Delta: 0},
	}
	s := trendlog.Analyze(entries, time.Hour, base.Add(time.Minute))
	if s.Direction != trendlog.Stable {
		t.Errorf("direction: want Stable, got %s", s.Direction)
	}
}

func TestAnalyze_EmptyWindowIsStable(t *testing.T) {
	s := trendlog.Analyze(nil, time.Hour, base)
	if s.Direction != trendlog.Stable {
		t.Errorf("direction: want Stable, got %s", s.Direction)
	}
	if s.Entries != 0 {
		t.Errorf("entries: want 0, got %d", s.Entries)
	}
}

func TestAnalyze_PeakOpenIsTracked(t *testing.T) {
	entries := makeEntries()
	now := base.Add(25 * time.Minute)
	s := trendlog.Analyze(entries, 30*time.Minute, now)
	if s.PeakOpen != 8 {
		t.Errorf("peak: want 8, got %d", s.PeakOpen)
	}
}

func TestAnalyze_WindowFiltersOldEntries(t *testing.T) {
	entries := makeEntries()
	// Use a narrow window that only captures the last entry.
	now := base.Add(25 * time.Minute)
	s := trendlog.Analyze(entries, 10*time.Minute, now)
	if s.Entries != 1 {
		t.Errorf("entries in window: want 1, got %d", s.Entries)
	}
}
