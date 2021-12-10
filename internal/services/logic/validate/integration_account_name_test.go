package validate

import "testing"

func TestIntegrationAccountName(t *testing.T) {
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
			name:  "1",
			input: "1",
			valid: true,
		},
		{
			name:  "abs_def",
			input: "abs_def",
			valid: true,
		},
		{
			name:  "abs.def",
			input: "abs.def",
			valid: true,
		},
		{
			name:  "abs def",
			input: "abs def",
			valid: false,
		},
		{
			name:  "abs-def",
			input: "abs-def",
			valid: true,
		},
		{
			name:  "AA-bb-",
			input: "AA-bb-",
			valid: true,
		},
		{
			name:  "-1-A-b",
			input: "-1-A-b",
			valid: true,
		},
	}
	validationFunction := IntegrationAccountName()
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
