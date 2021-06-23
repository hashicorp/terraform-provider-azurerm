package validate

import (
	"strings"
	"testing"
)

func TestFlexibleServerFirewallRuleName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "Empty",
			input: "",
			valid: false,
		},
		{
			name:  "Invalid Characters",
			input: "flexible%",
			valid: false,
		},
		{
			name:  "One character",
			input: "a",
			valid: true,
		},
		{
			name:  "End with `_`",
			input: "test_",
			valid: true,
		},
		{
			name:  "Start with `-`",
			input: "_test",
			valid: true,
		},
		{
			name:  "Valid",
			input: "flexible-6-test",
			valid: true,
		},
		{
			name:  "Valid2",
			input: "flex6ible6-6-te6st",
			valid: true,
		},
		{
			name:  "too long",
			input: strings.Repeat("a", 129),
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := FlexibleServerFirewallRuleName(tt.input, "name")
			valid := err == nil
			if valid != tt.valid {
				t.Errorf("Expected valid status %t but got %t for input %s", tt.valid, valid, tt.input)
			}
		})
	}
}
