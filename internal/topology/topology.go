package topology

import (
	"fmt"
	"sort"
	"sync"
)

// Node represents a single host and its observed open ports.
type Node struct {
	Host  string
	Ports []int
}

// Graph holds the current topology: a map of host -> open ports.
type Graph struct {
	mu    sync.RWMutex
	nodes map[string][]int
}

// New returns an empty topology Graph.
func New() *Graph {
	return &Graph{nodes: make(map[string][]int)}
}

// Set records the current open ports for a host, replacing any prior value.
func (g *Graph) Set(host string, ports []int) {
	if host == "" {
		return
	}
	cp := make([]int, len(ports))
	copy(cp, ports)
	sort.Ints(cp)

	g.mu.Lock()
	defer g.mu.Unlock()
	g.nodes[host] = cp
}

// Get returns the recorded ports for a host and whether the host exists.
func (g *Graph) Get(host string) ([]int, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	p, ok := g.nodes[host]
	if !ok {
		return nil, false
	}
	cp := make([]int, len(p))
	copy(cp, p)
	return cp, true
}

// Nodes returns all nodes in deterministic (sorted) order.
func (g *Graph) Nodes() []Node {
	g.mu.RLock()
	defer g.mu.RUnlock()

	hosts := make([]string, 0, len(g.nodes))
	for h := range g.nodes {
		hosts = append(hosts, h)
	}
	sort.Strings(hosts)

	out := make([]Node, 0, len(hosts))
	for _, h := range hosts {
		cp := make([]int, len(g.nodes[h]))
		copy(cp, g.nodes[h])
		out = append(out, Node{Host: h, Ports: cp})
	}
	return out
}

// Remove deletes a host from the graph.
func (g *Graph) Remove(host string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.nodes, host)
}

// String returns a compact human-readable summary.
func (g *Graph) String() string {
	nodes := g.Nodes()
	s := fmt.Sprintf("topology(%d hosts)", len(nodes))
	return s
}
