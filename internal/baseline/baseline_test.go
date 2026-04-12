package baseline_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/baseline"
)

func tmpFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "baseline.json")
}

func TestNew_MissingFileIsOK(t *testing.T) {
	b, err := baseline.New(tmpFile(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := b.Entries(); len(got) != 0 {
		t.Fatalf("expected empty baseline, got %d entries", len(got))
	}
}

func TestNew_CorruptFileReturnsError(t *testing.T) {
	p := tmpFile(t)
	_ = os.WriteFile(p, []byte("not json{"), 0o600)
	_, err := baseline.New(p)
	if err == nil {
		t.Fatal("expected error for corrupt file")
	}
}

func TestAdd_PersistsEntry(t *testing.T) {
	p := tmpFile(t)
	b, _ := baseline.New(p)
	if err := b.Add(8080, "test"); err != nil {
		t.Fatalf("Add: %v", err)
	}
	if !b.Contains(8080) {
		t.Fatal("expected port 8080 to be in baseline")
	}
	// reload from disk
	b2, err := baseline.New(p)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if !b2.Contains(8080) {
		t.Fatal("reloaded baseline missing port 8080")
	}
}

func TestRemove_DeletesEntry(t *testing.T) {
	p := tmpFile(t)
	b, _ := baseline.New(p)
	_ = b.Add(9090, "test")
	if err := b.Remove(9090); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if b.Contains(9090) {
		t.Fatal("port 9090 should have been removed")
	}
}

func TestEntries_IncludesAddedBy(t *testing.T) {
	b, _ := baseline.New(tmpFile(t))
	_ = b.Add(443, "admin")
	entries := b.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].AddedBy != "admin" {
		t.Errorf("expected AddedBy=admin, got %q", entries[0].AddedBy)
	}
}

func TestNew_LoadsExistingEntries(t *testing.T) {
	p := tmpFile(t)
	list := []baseline.Entry{{Port: 22, AddedBy: "seed"}}
	data, _ := json.Marshal(list)
	_ = os.WriteFile(p, data, 0o600)
	b, err := baseline.New(p)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if !b.Contains(22) {
		t.Fatal("expected port 22 from seed file")
	}
}
