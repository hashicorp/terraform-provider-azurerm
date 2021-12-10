package validate

import (
	"testing"
)

func TestValidateCognitiveServicesAccountName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "empty name",
			input: "",
			valid: false,
		},
		{
			name:  "Valid short name",
			input: "abc",
			valid: true,
		},
		{
			name:  "Invalid short name",
			input: "a",
			valid: false,
		},
		{
			name:  "Valid short name",
			input: "ab",
			valid: true,
		},
		{
			name:  "Valid long name",
			input: "abc_-.123",
			valid: true,
		},
		{
			name:  "Valid with a digit at the end",
			input: "hello1",
			valid: true,
		},
		{
			name:  "Valid with a digit in the middle",
			input: "hello1",
			valid: true,
		},
		{
			name:  "Invalid with a digit at the start",
			input: "1hello",
			valid: true,
		},
		{
			name:  "Invalid with a period at the start",
			input: ".heyo",
			valid: false,
		},
		{
			name:  "Valid name with period in the middle",
			input: "a.bc",
			valid: true,
		},
		{
			name:  "Valid name with period at end",
			input: "a.",
			valid: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CognitiveServicesAccountName()(tt.input, "")
			valid := err == nil
			if valid != tt.valid {
				t.Errorf("Expected valid status %t but got %t for input %s", tt.valid, valid, tt.input)
			}
		})
	}
}
