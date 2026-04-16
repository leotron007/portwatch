// Package probe provides a lightweight TCP port prober that measures
// reachability and round-trip latency for individual ports or batches.
//
// Basic usage:
//
//	p, err := probe.New("127.0.0.1", 2*time.Second)
//	if err != nil {
//		log.Fatal(err)
//	}
//	result := p.Probe(8080)
//	if result.Open {
//		fmt.Printf("port open, latency=%s\n", result.Latency)
//	}
//
// Configuration can also be loaded from YAML via FromBytes or FromConfig.
package probe
