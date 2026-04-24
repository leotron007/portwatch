package topology

import (
	"testing"
)

func buildGraph(t *testing.T, data map[string][]int) *Graph {
	t.Helper()
	g := New()
	for host, ports := range data {
		g.Set(host, ports)
	}
	return g
}

func TestDiff_NoChanges(t *testing.T) {
	prev := buildGraph(t, map[string][]int{"host-a": {22, 80}})
	next := buildGraph(t, map[string][]int{"host-a": {22, 80}})

	changes := Diff(prev, next)
	if len(changes) != 0 {
		t.Fatalf("expected no changes, got %v", changes)
	}
}

func TestDiff_OpenedPort(t *testing.T) {
	prev := buildGraph(t, map[string][]int{"host-a": {22}})
	next := buildGraph(t, map[string][]int{"host-a": {22, 443}})

	changes := Diff(prev, next)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if len(changes[0].Opened) != 1 || changes[0].Opened[0] != 443 {
		t.Fatalf("expected opened=[443], got %v", changes[0].Opened)
	}
}

func TestDiff_ClosedPort(t *testing.T) {
	prev := buildGraph(t, map[string][]int{"host-a": {22, 80}})
	next := buildGraph(t, map[string][]int{"host-a": {22}})

	changes := Diff(prev, next)
	if len(changes) != 1 || len(changes[0].Closed) != 1 {
		t.Fatalf("expected 1 closed change, got %v", changes)
	}
	if changes[0].Closed[0] != 80 {
		t.Fatalf("expected closed=[80], got %v", changes[0].Closed)
	}
}

func TestDiff_HostOnlyInPrev_TreatedAsClosed(t *testing.T) {
	prev := buildGraph(t, map[string][]int{"host-gone": {8080}})
	next := New()

	changes := Diff(prev, next)
	if len(changes) != 1 || changes[0].Host != "host-gone" {
		t.Fatalf("unexpected changes: %v", changes)
	}
	if len(changes[0].Closed) != 1 || changes[0].Closed[0] != 8080 {
		t.Fatalf("expected closed=[8080], got %v", changes[0].Closed)
	}
}

func TestDiff_NewHostAllOpened(t *testing.T) {
	prev := New()
	next := buildGraph(t, map[string][]int{"new-host": {22, 80}})

	changes := Diff(prev, next)
	if len(changes) != 1 || changes[0].Host != "new-host" {
		t.Fatalf("unexpected changes: %v", changes)
	}
	if len(changes[0].Opened) != 2 {
		t.Fatalf("expected 2 opened ports, got %v", changes[0].Opened)
	}
}
