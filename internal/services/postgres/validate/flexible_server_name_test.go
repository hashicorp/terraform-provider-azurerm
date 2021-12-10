package validate

import "testing"

func TestFlexibleServerName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "Upper Case",
			input: "Flexible",
			valid: false,
		},
		{
			name:  "Invalid Characters",
			input: "flexible_test",
			valid: false,
		},
		{
			name:  "Upper Case 2",
			input: "flexible_Test",
			valid: false,
		},
		{
			name:  "One character",
			input: "a",
			valid: true,
		},
		{
			name:  "Empty",
			input: "",
			valid: false,
		},
		{
			name:  "End with `_`",
			input: "test_",
			valid: false,
		},
		{
			name:  "Start with `-`",
			input: "_test",
			valid: false,
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := FlexibleServerName(tt.input, "name")
			valid := err == nil
			if valid != tt.valid {
				t.Errorf("Expected valid status %t but got %t for input %s", tt.valid, valid, tt.input)
			}
		})
	}
}
