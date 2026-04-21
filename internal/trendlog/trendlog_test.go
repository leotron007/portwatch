package trendlog_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/trendlog"
)

func tmpFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "trend.jsonl")
}

func TestNew_MissingFileIsOK(t *testing.T) {
	tl, err := trendlog.New(tmpFile(t), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := len(tl.Entries()); got != 0 {
		t.Fatalf("expected 0 entries, got %d", got)
	}
}

func TestNew_CorruptFileReturnsError(t *testing.T) {
	p := tmpFile(t)
	_ = os.WriteFile(p, []byte("{bad json\n"), 0o644)
	_, err := trendlog.New(p, nil)
	if err == nil {
		t.Fatal("expected error for corrupt file")
	}
}

func TestRecord_PersistsEntry(t *testing.T) {
	p := tmpFile(t)
	tl, _ := trendlog.New(p, nil)
	now := time.Now().UTC().Truncate(time.Second)

	if err := tl.Record(5, now); err != nil {
		t.Fatalf("Record: %v", err)
	}

	entries := tl.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].OpenCount != 5 {
		t.Errorf("open_count: want 5, got %d", entries[0].OpenCount)
	}
	if entries[0].Delta != 5 {
		t.Errorf("delta on first entry: want 5, got %d", entries[0].Delta)
	}
}

func TestRecord_DeltaIsRelativeToPrevious(t *testing.T) {
	p := tmpFile(t)
	tl, _ := trendlog.New(p, nil)
	now := time.Now().UTC()

	_ = tl.Record(10, now)
	_ = tl.Record(7, now.Add(time.Minute))

	entries := tl.Entries()
	if entries[1].Delta != -3 {
		t.Errorf("delta: want -3, got %d", entries[1].Delta)
	}
}

func TestNew_LoadsExistingEntries(t *testing.T) {
	p := tmpFile(t)
	tl, _ := trendlog.New(p, nil)
	now := time.Now().UTC()
	_ = tl.Record(3, now)
	_ = tl.Record(6, now.Add(time.Minute))

	// Re-open and verify persistence.
	tl2, err := trendlog.New(p, nil)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if got := len(tl2.Entries()); got != 2 {
		t.Fatalf("expected 2 entries after reload, got %d", got)
	}
}
