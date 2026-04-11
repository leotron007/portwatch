package history_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/portwatch/internal/history"
)

func tmpFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "history.json")
}

func TestNew_MissingFileIsOK(t *testing.T) {
	_, err := history.New(tmpFile(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRecord_PersistsEntry(t *testing.T) {
	path := tmpFile(t)
	h, _ := history.New(path)

	if err := h.Record(8080, "opened", "tcp"); err != nil {
		t.Fatalf("Record: %v", err)
	}

	entries := h.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.Port != 8080 || e.Event != "opened" || e.Protocol != "tcp" {
		t.Errorf("unexpected entry: %+v", e)
	}
	if e.Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}
}

func TestNew_LoadsExistingEntries(t *testing.T) {
	path := tmpFile(t)
	h1, _ := history.New(path)
	_ = h1.Record(443, "opened", "tcp")
	_ = h1.Record(443, "closed", "tcp")

	h2, err := history.New(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if len(h2.Entries()) != 2 {
		t.Errorf("expected 2 entries after reload, got %d", len(h2.Entries()))
	}
}

func TestNew_CorruptFileReturnsError(t *testing.T) {
	path := tmpFile(t)
	_ = os.WriteFile(path, []byte("not json{"), 0o644)
	_, err := history.New(path)
	if err == nil {
		t.Fatal("expected error for corrupt file")
	}
}

func TestEntries_ReturnsCopy(t *testing.T) {
	h, _ := history.New(tmpFile(t))
	_ = h.Record(22, "opened", "tcp")

	a := h.Entries()
	a[0].Port = 9999

	b := h.Entries()
	if b[0].Port == 9999 {
		t.Error("Entries should return an independent copy")
	}
}
