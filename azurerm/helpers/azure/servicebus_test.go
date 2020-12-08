package azure

import (
	"strings"
	"testing"
)

func TestValidateServiceBusTopicName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "Empty value",
			input: "",
			valid: false,
		},
		{
			name:  "Invalid name with only 1 letter",
			input: "a",
			valid: true,
		},
		{
			name:  "Invalid name starts with underscore",
			input: "_a",
			valid: false,
		},
		{
			name:  "Invalid name ends with period",
			input: "a.",
			valid: false,
		},
		{
			name:  "Valid name with numbers",
			input: "12345",
			valid: true,
		},
		{
			name:  "Valid name with only 1 number",
			input: "1",
			valid: true,
		},
		{
			name:  "Valid name with hyphens",
			input: "malcolm-in-the-middle",
			valid: true,
		},
		{
			name:  "Valid name with 259 characters",
			input: strings.Repeat("w", 259),
			valid: true,
		},
		{
			name:  "Valid name with 260 characters",
			input: strings.Repeat("w", 260),
			valid: true,
		},
		{
			name:  "Invalid name with 261 characters",
			input: strings.Repeat("w", 261),
			valid: false,
		},
	}

	validationFunction := ValidateServiceBusTopicName()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validationFunction(tt.input, "name")
			valid := err == nil
			if valid != tt.valid {
				t.Errorf("Expected valid status %t but got %t for input %s", tt.valid, valid, tt.input)
			}
		})
	}
}
