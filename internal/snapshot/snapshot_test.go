package snapshot_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

var now = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

func TestNew_ValidPorts(t *testing.T) {
	s, err := snapshot.New([]int{80, 443, 8080}, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Ports) != 3 {
		t.Errorf("expected 3 ports, got %d", len(s.Ports))
	}
	if s.CapturedAt != now {
		t.Errorf("CapturedAt mismatch")
	}
}

func TestNew_ZeroTime(t *testing.T) {
	_, err := snapshot.New([]int{80}, time.Time{})
	if err == nil {
		t.Fatal("expected error for zero time")
	}
}

func TestEqual_SamePorts(t *testing.T) {
	a, _ := snapshot.New([]int{80, 443}, now)
	b, _ := snapshot.New([]int{443, 80}, now)
	if !a.Equal(b) {
		t.Error("expected snapshots with same ports to be equal")
	}
}

func TestEqual_DifferentPorts(t *testing.T) {
	a, _ := snapshot.New([]int{80}, now)
	b, _ := snapshot.New([]int{443}, now)
	if a.Equal(b) {
		t.Error("expected snapshots with different ports to be unequal")
	}
}

func TestEqual_NilOther(t *testing.T) {
	a, _ := snapshot.New([]int{80}, now)
	if a.Equal(nil) {
		t.Error("expected Equal(nil) to return false")
	}
}

func TestAdded_ReturnsNewPorts(t *testing.T) {
	old, _ := snapshot.New([]int{80, 443}, now)
	new_, _ := snapshot.New([]int{80, 443, 8080}, now)
	added := old.Added(new_)
	if len(added) != 1 || added[0] != 8080 {
		t.Errorf("expected [8080], got %v", added)
	}
}

func TestRemoved_ReturnsMissingPorts(t *testing.T) {
	old, _ := snapshot.New([]int{80, 443, 8080}, now)
	new_, _ := snapshot.New([]int{80, 443}, now)
	removed := old.Removed(new_)
	if len(removed) != 1 || removed[0] != 8080 {
		t.Errorf("expected [8080], got %v", removed)
	}
}

func TestAdded_NilOther(t *testing.T) {
	a, _ := snapshot.New([]int{80}, now)
	if got := a.Added(nil); got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}
