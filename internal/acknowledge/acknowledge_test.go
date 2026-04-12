package acknowledge_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/acknowledge"
)

func tmpFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "ack.json")
}

func TestNew_MissingFileIsOK(t *testing.T) {
	s, err := acknowledge.New(tmpFile(t))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(s.All()) != 0 {
		t.Fatalf("expected empty store")
	}
}

func TestNew_CorruptFileReturnsError(t *testing.T) {
	p := tmpFile(t)
	_ = os.WriteFile(p, []byte("not json"), 0o644)
	_, err := acknowledge.New(p)
	if err == nil {
		t.Fatal("expected error for corrupt file")
	}
}

func TestAck_PersistsAndReloads(t *testing.T) {
	p := tmpFile(t)
	s, _ := acknowledge.New(p)
	now := time.Now().UTC().Truncate(time.Second)

	if err := s.Ack(8080, "opened", now); err != nil {
		t.Fatalf("Ack: %v", err)
	}

	s2, err := acknowledge.New(p)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if !s2.IsAcked(8080, "opened") {
		t.Fatal("expected port 8080/opened to be acked after reload")
	}
}

func TestIsAcked_ReturnsFalseForUnknown(t *testing.T) {
	s, _ := acknowledge.New(tmpFile(t))
	if s.IsAcked(9999, "closed") {
		t.Fatal("expected not acked")
	}
}

func TestClear_RemovesEntry(t *testing.T) {
	p := tmpFile(t)
	s, _ := acknowledge.New(p)
	_ = s.Ack(443, "closed", time.Now())

	if !s.IsAcked(443, "closed") {
		t.Fatal("expected acked before clear")
	}
	if err := s.Clear(443, "closed"); err != nil {
		t.Fatalf("Clear: %v", err)
	}
	if s.IsAcked(443, "closed") {
		t.Fatal("expected not acked after clear")
	}

	// Verify persisted.
	s2, _ := acknowledge.New(p)
	if s2.IsAcked(443, "closed") {
		t.Fatal("expected not acked after reload")
	}
}

func TestAll_ReturnsAllEntries(t *testing.T) {
	s, _ := acknowledge.New(tmpFile(t))
	_ = s.Ack(80, "opened", time.Now())
	_ = s.Ack(443, "opened", time.Now())

	if got := len(s.All()); got != 2 {
		t.Fatalf("expected 2 entries, got %d", got)
	}
}
