package rules_test

import (
	"testing"

	"github.com/user/portwatch/internal/rules"
)

func TestEvaluate_MatchFound(t *testing.T) {
	s := &rules.Set{
		Rules: []rules.Rule{
			{Name: "web", Ports: []int{80, 443}, Action: rules.ActionIgnore},
		},
	}
	rule, ok := s.Evaluate(80)
	if !ok {
		t.Fatal("expected match for port 80")
	}
	if rule.Name != "web" {
		t.Errorf("expected rule name 'web', got %q", rule.Name)
	}
}

func TestEvaluate_NoMatch(t *testing.T) {
	s := &rules.Set{
		Rules: []rules.Rule{
			{Name: "web", Ports: []int{80}, Action: rules.ActionIgnore},
		},
	}
	_, ok := s.Evaluate(9999)
	if ok {
		t.Fatal("expected no match for port 9999")
	}
}

func TestValidate_DuplicateName(t *testing.T) {
	s := &rules.Set{
		Rules: []rules.Rule{
			{Name: "dup", Ports: []int{80}, Action: rules.ActionAlert},
			{Name: "dup", Ports: []int{443}, Action: rules.ActionAlert},
		},
	}
	if err := s.Validate(); err == nil {
		t.Fatal("expected error for duplicate rule names")
	}
}

func TestValidate_InvalidAction(t *testing.T) {
	s := &rules.Set{
		Rules: []rules.Rule{
			{Name: "bad", Ports: []int{80}, Action: "unknown"},
		},
	}
	if err := s.Validate(); err == nil {
		t.Fatal("expected error for invalid action")
	}
}

func TestLoadFromBytes_Valid(t *testing.T) {
	yaml := []byte(`
rules:
  - name: ssh
    ports: [22]
    action: ignore
`)
	s, err := rules.LoadFromBytes(yaml)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Rules) != 1 {
		t.Errorf("expected 1 rule, got %d", len(s.Rules))
	}
}

func TestDefault_NotEmpty(t *testing.T) {
	s := rules.Default()
	if len(s.Rules) == 0 {
		t.Fatal("expected default rules to be non-empty")
	}
	if err := s.Validate(); err != nil {
		t.Fatalf("default rules failed validation: %v", err)
	}
}
