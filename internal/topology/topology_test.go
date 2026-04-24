package topology

import (
	"testing"
)

func TestSet_And_Get(t *testing.T) {
	g := New()
	g.Set("host-a", []int{80, 443, 22})

	ports, ok := g.Get("host-a")
	if !ok {
		t.Fatal("expected host-a to exist")
	}
	if len(ports) != 3 || ports[0] != 22 {
		t.Fatalf("expected sorted ports [22 80 443], got %v", ports)
	}
}

func TestGet_UnknownHost(t *testing.T) {
	g := New()
	_, ok := g.Get("ghost")
	if ok {
		t.Fatal("expected false for unknown host")
	}
}

func TestSet_EmptyHostIgnored(t *testing.T) {
	g := New()
	g.Set("", []int{80})
	if len(g.Nodes()) != 0 {
		t.Fatal("empty host should not be stored")
	}
}

func TestRemove_DeletesHost(t *testing.T) {
	g := New()
	g.Set("host-a", []int{80})
	g.Remove("host-a")
	_, ok := g.Get("host-a")
	if ok {
		t.Fatal("expected host-a to be removed")
	}
}

func TestNodes_ReturnsSortedHosts(t *testing.T) {
	g := New()
	g.Set("z-host", []int{9000})
	g.Set("a-host", []int{22})

	nodes := g.Nodes()
	if len(nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(nodes))
	}
	if nodes[0].Host != "a-host" || nodes[1].Host != "z-host" {
		t.Fatalf("unexpected order: %v", nodes)
	}
}
