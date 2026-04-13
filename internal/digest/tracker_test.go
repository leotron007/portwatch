package digest_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/digest"
)

func TestTracker_FirstUpdateNotChanged(t *testing.T) {
	tr := digest.NewTracker()
	changed, prev, curr := tr.Update([]int{80, 443})

	if changed {
		t.Error("first update should not report a change")
	}
	if prev != nil {
		t.Error("prev should be nil on first update")
	}
	if curr == nil {
		t.Error("curr should not be nil")
	}
}

func TestTracker_SamePortsNoChange(t *testing.T) {
	tr := digest.NewTracker()
	tr.Update([]int{80, 443})
	changed, _, _ := tr.Update([]int{443, 80})

	if changed {
		t.Error("same port set should not report a change")
	}
}

func TestTracker_DifferentPortsReportsChange(t *testing.T) {
	tr := digest.NewTracker()
	tr.Update([]int{80})
	changed, prev, curr := tr.Update([]int{80, 443})

	if !changed {
		t.Error("different port set should report a change")
	}
	if prev == nil || curr == nil {
		t.Fatal("prev and curr must be non-nil when changed")
	}
	if prev.Equal(curr) {
		t.Error("prev and curr should differ")
	}
}

func TestTracker_Last_NilBeforeFirstUpdate(t *testing.T) {
	tr := digest.NewTracker()
	if tr.Last() != nil {
		t.Error("Last() should be nil before any update")
	}
}

func TestTracker_Last_ReturnsCurrentDigest(t *testing.T) {
	tr := digest.NewTracker()
	_, _, curr := tr.Update([]int{22, 80})

	if !tr.Last().Equal(curr) {
		t.Error("Last() should match the returned by Update")
	}
}

func TestTracker_Reset_ClearsState(t *testing.T) {
	tr := digest.NewTracker()
	tr.Update([]int{80})
	tr.Reset()

	if tr.Last() != nil {
		t.Error("Last() should be nil after Reset")
	}

	// After reset the next update should behave like the first.
	changed, prev, _ := tr.Update([]int{80})
	if changed {
		t.Error("update after reset should not report a change")
	}
	if prev != nil {
		t.Error("prev should be nil after reset")
	}
}
