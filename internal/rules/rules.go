// Package rules provides configuration and evaluation of port monitoring rules.
package rules

import (
	"fmt"
	"strings"
)

// Action defines what to do when a rule matches.
type Action string

const (
	ActionAlert  Action = "alert"
	ActionIgnore Action = "ignore"
)

// Rule defines a single port monitoring rule.
type Rule struct {
	Name    string  `yaml:"name"`
	Ports   []int   `yaml:"ports"`
	Action  Action  `yaml:"action"`
	Comment string  `yaml:"comment,omitempty"`
}

// Set holds a collection of rules.
type Set struct {
	Rules []Rule `yaml:"rules"`
}

// Evaluate checks whether the given port matches any rule and returns the
// matching rule along with a boolean indicating a match was found.
func (s *Set) Evaluate(port int) (*Rule, bool) {
	for i := range s.Rules {
		for _, p := range s.Rules[i].Ports {
			if p == port {
				return &s.Rules[i], true
			}
		}
	}
	return nil, false
}

// Validate checks that all rules in the set are well-formed.
func (s *Set) Validate() error {
	names := make(map[string]struct{})
	for _, r := range s.Rules {
		if strings.TrimSpace(r.Name) == "" {
			return fmt.Errorf("rule missing name")
		}
		if _, dup := names[r.Name]; dup {
			return fmt.Errorf("duplicate rule name: %q", r.Name)
		}
		names[r.Name] = struct{}{}
		if r.Action != ActionAlert && r.Action != ActionIgnore {
			return fmt.Errorf("rule %q has invalid action: %q", r.Name, r.Action)
		}
		if len(r.Ports) == 0 {
			return fmt.Errorf("rule %q has no ports defined", r.Name)
		}
	}
	return nil
}
