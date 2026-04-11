// Package filter provides port-level filtering for portwatch scan results.
//
// A Filter can restrict which ports are considered during a scan cycle by
// specifying an explicit allow-list, an ignore-list, and a min/max port range.
// Filters are composable with the scanner and daemon packages — pass the
// filtered port list to state.Compare to limit noise from well-known or
// intentionally open ports.
//
// # Quick start
//
//	f := filter.New()
//	f.IgnoredPorts = []int{22, 80, 443}
//	f.MinPort = 1024
//	open := f.Apply(scannedPorts)
//
// Filters can also be loaded from YAML configuration via filter.FromBytes or
// filter.FromConfig for integration with the rules.yaml config file.
package filter
