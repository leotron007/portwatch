package summary_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/summary"
)

func tmpFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "history.json")
}

func TestNew_NilHistoryReturnsError(t *testing.T) {
	_, err := summary.New(nil, nil, time.Hour)
	if err == nil {
		t.Fatal("expected error for nil history")
	}
}

func TestNew_NilWriterDefaultsToStdout(t *testing.T) {
	h, _ := history.New(tmpFile(t))
	r, err := summary.New(h, nil, time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}

func TestWrite_CountsOpenedAndClosed(t *testing.T) {
	h, _ := history.New(tmpFile(t))
	_ = h.Record("opened", []int{8080, 9090})
	_ = h.Record("closed", []int{3000})

	var buf bytes.Buffer
	r, _ := summary.New(h, &buf, time.Hour)
	if err := r.Write(); err != nil {
		t.Fatalf("Write() error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "opened=1") {
		t.Errorf("expected opened=1 in output, got: %s", out)
	}
	if !strings.Contains(out, "closed=1") {
		t.Errorf("expected closed=1 in output, got: %s", out)
	}
	if !strings.Contains(out, "events=2") {
		t.Errorf("expected events=2 in output, got: %s", out)
	}
}

func TestWrite_EmptyHistory(t *testing.T) {
	h, _ := history.New(tmpFile(t))
	var buf bytes.Buffer
	r, _ := summary.New(h, &buf, time.Hour)
	_ = r.Write()

	out := buf.String()
	if !strings.Contains(out, "events=0") {
		t.Errorf("expected events=0, got: %s", out)
	}
}

func TestNew_ZeroWindowDefaultsTo24h(t *testing.T) {
	h, _ := history.New(tmpFile(t))
	var buf bytes.Buffer
	r, err := summary.New(h, &buf, 0)
	if err != nil || r == nil {
		t.Fatalf("unexpected error or nil reporter: %v", err)
	}
	_ = os.Remove(tmpFile(t))
}
