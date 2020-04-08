package validate

import (
	"testing"
)

func TestAdvisorSuppressionName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "Valid one",
			input: "da _-.",
			valid: true,
		},
		{
			name:  "Valid two",
			input: "ew2 d-.dE",
			valid: true,
		},
		{
			name:  "Special Character",
			input: "dfs^&",
			valid: false,
		},
		{
			name:  "Empty",
			input: "",
			valid: false,
		},
		{
			name:  "Starts with Character",
			input: " .ds",
			valid: true,
		},
	}
	var validationFunction = AdvisorSuppressionName()
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
