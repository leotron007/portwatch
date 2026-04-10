// Package scanner provides functionality to detect open TCP ports on a host.
//
// Usage:
//
//	s := scanner.New("127.0.0.1", 500*time.Millisecond)
//	ports, err := s.Scan(1, 1024)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, p := range ports {
//		fmt.Printf("open: %s/%d\n", p.Protocol, p.Number)
//	}
//
// The scanner performs sequential TCP dial attempts within the specified port
// range. Each dial respects the configured Timeout so that the overall scan
// duration remains predictable.
package scanner
