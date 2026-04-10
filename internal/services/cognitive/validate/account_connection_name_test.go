package validate

import (
	"testing"
)

func TestAccountConnectionName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "empty",
			input: "",
			valid: false,
		},
		{
			name:  "too short - one char",
			input: "a",
			valid: false,
		},
		{
			name:  "too short - two chars",
			input: "ab",
			valid: false,
		},
		{
			name:  "minimum valid length",
			input: "abc",
			valid: true,
		},
		{
			name:  "maximum valid length - 33 chars",
			input: "abcdefghijklmnopqrstuvwxyz1234567",
			valid: true,
		},
		{
			name:  "too long - 34 chars",
			input: "abcdefghijklmnopqrstuvwxyz12345678",
			valid: false,
		},
		{
			name:  "valid with dashes",
			input: "abc-def",
			valid: true,
		},
		{
			name:  "valid with underscores",
			input: "abc_def",
			valid: true,
		},
		{
			name:  "valid with mixed characters",
			input: "Abc-123_def",
			valid: true,
		},
		{
			name:  "starts with digit",
			input: "1abc",
			valid: true,
		},
		{
			name:  "starts with dash",
			input: "-abc",
			valid: false,
		},
		{
			name:  "starts with underscore",
			input: "_abc",
			valid: false,
		},
		{
			name:  "contains period",
			input: "abc.def",
			valid: false,
		},
		{
			name:  "contains space",
			input: "abc def",
			valid: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := AccountConnectionName()(tt.input, "")
			valid := err == nil
			if valid != tt.valid {
				t.Errorf("Expected valid status %t but got %t for input %s", tt.valid, valid, tt.input)
			}
		})
	}
}
