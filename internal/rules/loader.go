package rules

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadFromFile reads a YAML rules file from the given path and returns a
// validated Set, or an error if the file cannot be parsed or is invalid.
func LoadFromFile(path string) (*Set, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading rules file %q: %w", path, err)
	}
	return LoadFromBytes(data)
}

// LoadFromBytes parses YAML bytes into a Set and validates it.
func LoadFromBytes(data []byte) (*Set, error) {
	var s Set
	if err := yaml.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parsing rules YAML: %w", err)
	}
	if err := s.Validate(); err != nil {
		return nil, fmt.Errorf("invalid rules: %w", err)
	}
	return &s, nil
}

// Default returns a Set with sensible built-in rules for common well-known ports.
func Default() *Set {
	return &Set{
		Rules: []Rule{
			{
				Name:    "common-services",
				Ports:   []int{22, 80, 443, 3306, 5432, 6379, 27017},
				Action:  ActionIgnore,
				Comment: "Well-known service ports, typically expected",
			},
			{
				Name:    "suspicious-high-ports",
				Ports:   []int{4444, 1337, 31337},
				Action:  ActionAlert,
				Comment: "Ports commonly associated with backdoors or test tools",
			},
		},
	}
}
