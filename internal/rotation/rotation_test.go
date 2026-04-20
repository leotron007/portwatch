package rotation

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func tmpDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "rotation-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(dir) })
	return dir
}

func TestNew_InvalidMaxBytes(t *testing.T) {
	_, err := New("/tmp/x.log", 0, 3)
	if err == nil {
		t.Fatal("expected error for maxBytes=0")
	}
}

func TestNew_InvalidMaxFiles(t *testing.T) {
	_, err := New("/tmp/x.log", 1024, 0)
	if err == nil {
		t.Fatal("expected error for maxFiles=0")
	}
}

func TestWrite_CreatesFile(t *testing.T) {
	dir := tmpDir(t)
	path := filepath.Join(dir, "out.log")

	r, err := New(path, 1024, 3)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	if _, err := r.Write([]byte("hello\n")); err != nil {
		t.Fatalf("Write: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(data) != "hello\n" {
		t.Errorf("unexpected content: %q", data)
	}
}

func TestWrite_RotatesWhenExceedsMaxBytes(t *testing.T) {
	dir := tmpDir(t)
	path := filepath.Join(dir, "out.log")

	// maxBytes=10 so a 12-byte write triggers rotation.
	r, err := New(path, 10, 3)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	first := []byte("0123456789") // exactly 10 bytes — fills the file
	if _, err := r.Write(first); err != nil {
		t.Fatalf("first Write: %v", err)
	}

	second := []byte("ABCDEF") // triggers rotation
	if _, err := r.Write(second); err != nil {
		t.Fatalf("second Write: %v", err)
	}

	backup := path + ".1"
	data, err := os.ReadFile(backup)
	if err != nil {
		t.Fatalf("backup not found: %v", err)
	}
	if string(data) != string(first) {
		t.Errorf("backup content = %q, want %q", data, first)
	}

	current, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("current file missing: %v", err)
	}
	if string(current) != string(second) {
		t.Errorf("current content = %q, want %q", current, second)
	}
}

func TestWrite_MultipleRotations(t *testing.T) {
	dir := tmpDir(t)
	path := filepath.Join(dir, "out.log")

	r, err := New(path, 5, 2)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	payloads := []string{"AAAAA", "BBBBB", "CCCCC"}
	for _, p := range payloads {
		if _, err := r.Write([]byte(p)); err != nil {
			t.Fatalf("Write %q: %v", p, err)
		}
	}

	// After 3 writes with maxFiles=2, .1 and .2 should exist; .3 should not.
	if _, err := os.Stat(path + ".1"); err != nil {
		t.Errorf("expected backup .1 to exist: %v", err)
	}
	if _, err := os.Stat(path + ".2"); err != nil {
		t.Errorf("expected backup .2 to exist: %v", err)
	}
}

func TestClose_Idempotent(t *testing.T) {
	dir := tmpDir(t)
	path := filepath.Join(dir, "out.log")

	r, err := New(path, 1024, 3)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := r.Close(); err != nil {
		t.Fatalf("first Close: %v", err)
	}
	if err := r.Close(); err != nil {
		t.Fatalf("second Close: %v", err)
	}
}

func TestNew_BadPath(t *testing.T) {
	_, err := New(filepath.Join("nonexistent", "dir", "out.log"), 1024, 3)
	if err == nil {
		t.Fatal("expected error for bad path")
	}
	if !strings.Contains(err.Error(), "rotation:") {
		t.Errorf("error should be prefixed with 'rotation:': %v", err)
	}
}
