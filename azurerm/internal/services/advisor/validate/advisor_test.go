package validate

import (
	"testing"
)

func TestAdvisorSuppresionTTL(t *testing.T) {
	testCases := []struct {
		input       string
		shouldError bool
	}{
		{"7", false},
		{"7.12", false},
		{"7.0:30", false},
		{"7.0:0:21", false},
		{"0:30", false},
		{"0:30:21", false},
		{"0", true},
		{"0:0", true},
		{"0:0:0", true},
		{"", true},
		{"7.24:0:0", true},
		{"1000", false},
		{"-1", false},
		{"-10", true},
	}

	for _, test := range testCases {
		_, es := AdvisorSuppresionTTL(test.input, "name")

		if test.shouldError && len(es) == 0 {
			t.Fatalf("Expected validating name %q to fail", test.input)
		}
	}
}

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
