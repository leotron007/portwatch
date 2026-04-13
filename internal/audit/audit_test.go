package audit_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/audit"
)

func TestLog_WritesValidJSON(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	if err := l.Log("port_opened", 8080, "tcp", "test"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var e audit.Entry
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &e); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if e.Event != "port_opened" {
		t.Errorf("event = %q, want port_opened", e.Event)
	}
	if e.Port != 8080 {
		t.Errorf("port = %d, want 8080", e.Port)
	}
	if e.Proto != "tcp" {
		t.Errorf("proto = %q, want tcp", e.Proto)
	}
}

func TestLog_ContainsNewline(t *testing.T) {
	var buf bytes.Buffer
	audit.New(&buf).Log("x", 1, "tcp", "") //nolint:errcheck

	if !strings.HasSuffix(buf.String(), "\n") {
		t.Error("expected trailing newline")
	}
}

func TestNew_NilWriterUsesStderr(t *testing.T) {
	// Simply ensure no panic when writer is nil.
	l := audit.New(nil)
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestFileWriter_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "audit.log")

	w, err := audit.FileWriter(path)
	if err != nil {
		t.Fatalf("FileWriter: %v", err)
	}
	defer w.Close()

	if _, err := os.Stat(path); err != nil {
		t.Errorf("file not created: %v", err)
	}
}

func TestLog_EmptyReason_OmittedFromJSON(t *testing.T) {
	var buf bytes.Buffer
	audit.New(&buf).Log("port_closed", 443, "tcp", "") //nolint:errcheck

	if strings.Contains(buf.String(), "reason") {
		t.Error("reason field should be omitted when empty")
	}
}
