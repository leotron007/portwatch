// Package rules defines the rule engine for portwatch.
//
// Rules are expressed in YAML and control how the daemon responds when a
// particular port is found to be open. Each rule associates a list of port
// numbers with an action:
//
//   - alert  – emit an alert notification (default for unknown ports)
//   - ignore – suppress notifications for expected ports
//
// Example rules file:
//
//	rules:
//	  - name: web-servers
//	    ports: [80, 443]
//	    action: ignore
//	    comment: "Standard HTTP/HTTPS ports"
//	  - name: backdoor-check
//	    ports: [4444, 1337]
//	    action: alert
//	    comment: "Flag suspicious listener ports"
//
// Use LoadFromFile or LoadFromBytes to parse a rules file, and Set.Evaluate
// to check a specific port against the loaded rule set.
package rules
