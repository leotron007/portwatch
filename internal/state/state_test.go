package state

import (
	"os"
	"path/filepath"
	"testing"
)

func tmpPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "state.json")
}

func TestNew_MissingFileIsOK(t *testing.T) {
	_, err := New(tmpPath(t))
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
}

func TestSave_PersistsAndReloads(t *testing.T) {
	path := tmpPath(t)
	s, _ := New(path)
	ports := []int{80, 443, 8080}
	if err := s.Save(ports); err != nil {
		t.Fatalf("Save: %v", err)
	}

	s2, err := New(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	snap := s2.Current()
	if len(snap.Ports) != len(ports) {
		t.Fatalf("expected %d ports, got %d", len(ports), len(snap.Ports))
	}
}

func TestNew_CorruptFile(t *testing.T) {
	path := tmpPath(t)
	os.WriteFile(path, []byte("not-json{"), 0o644)
	_, err := New(path)
	if err == nil {
		t.Fatal("expected error for corrupt state file")
	}
}

func TestCompare_OpenedAndClosed(t *testing.T) {
	prev := []int{80, 443, 22}
	curr := []int{80, 8080}
	d := Compare(prev, curr)
	if len(d.Opened) != 1 || d.Opened[0] != 8080 {
		t.Errorf("Opened: expected [8080], got %v", d.Opened)
	}
	if len(d.Closed) != 2 {
		t.Errorf("Closed: expected 2 ports, got %v", d.Closed)
	}
}

func TestCompare_NoDiff(t *testing.T) {
	ports := []int{80, 443}
	d := Compare(ports, ports)
	if !d.IsEmpty() {
		t.Errorf("expected empty diff, got %+v", d)
	}
}

func TestCompare_EmptyPrevious(t *testing.T) {
	d := Compare(nil, []int{22, 80})
	if len(d.Opened) != 2 {
		t.Errorf("expected 2 opened ports, got %v", d.Opened)
	}
	if len(d.Closed) != 0 {
		t.Errorf("expected 0 closed ports, got %v", d.Closed)
	}
}
