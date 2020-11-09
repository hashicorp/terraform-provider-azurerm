package validate

import (
	"testing"
)

func TestDigitaltwinsName(t *testing.T) {
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
			name:  "Too short",
			input: "a",
			valid: false,
		},
		{
			name:  "Invalid character",
			input: "digital_twins",
			valid: false,
		},
		{
			name:  "Valid Name",
			input: "Digital-12-Twins",
			valid: true,
		},
		{
			name:  "End with `-`",
			input: "Digital-12-",
			valid: false,
		},
		{
			name:  "Start with `-`",
			input: "-Digital-12",
			valid: false,
		},
		{
			name:  "Invalid character",
			input: "digital.twins",
			valid: false,
		},
	}
	var validationFunction = DigitaltwinsName()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validationFunction(tt.input, "")
			valid := err == nil
			if valid != tt.valid {
				t.Errorf("Expected valid status %t but got %t for input %s", tt.valid, valid, tt.input)
			}
		})
	}
}
