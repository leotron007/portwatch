package digest_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/digest"
)

func TestNew_EmptyPorts(t *testing.T) {
	d := digest.New([]int{})
	if d.Hash() == "" {
		t.Fatal("expected non-empty hash for empty port list")
	}
}

func TestNew_DeterministicHash(t *testing.T) {
	d1 := digest.New([]int{80, 443, 8080})
	d2 := digest.New([]int{8080, 80, 443})

	if !d1.Equal(d2) {
		t.Errorf("expected equal digests for same ports in different order, got %s vs %s", d1.Hash(), d2.Hash())
	}
}

func TestNew_DifferentPortsDifferentHash(t *testing.T) {
	d1 := digest.New([]int{80, 443})
	d2 := digest.New([]int{80, 444})

	if d1.Equal(d2) {
		t.Error("expected different digests for different port sets")
	}
}

func TestPorts_ReturnsSortedCopy(t *testing.T) {
	input := []int{9000, 22, 80}
	d := digest.New(input)
	ports := d.Ports()

	expected := []int{22, 80, 9000}
	for i, p := range expected {
		if ports[i] != p {
			t.Errorf("index %d: expected %d, got %d", i, p, ports[i])
		}
	}

	// Mutating returned slice must not affect digest.
	ports[0] = 9999
	if d.Ports()[0] == 9999 {
		t.Error("Ports() returned a reference to internal slice")
	}
}

func TestEqual_NilOther(t *testing.T) {
	d := digest.New([]int{80})
	if d.Equal(nil) {
		t.Error("Equal(nil) should return false")
	}
}

func TestString_ShortPrefix(t *testing.T) {
	d := digest.New([]int{80, 443})
	s := d.String()
	if len(s) != 12 {
		t.Errorf("expected 12-char prefix, got %d chars: %s", len(s), s)
	}
}

func TestNew_SinglePort(t *testing.T) {
	d1 := digest.New([]int{22})
	d2 := digest.New([]int{22})
	if !d1.Equal(d2) {
		t.Error("identical single-port digests should be equal")
	}
}
