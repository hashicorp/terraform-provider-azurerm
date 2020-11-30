package validate

import "testing"

func TestShareName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "invalid character",
			input: "9()",
			valid: false,
		},
		{
			name:  "less character",
			input: "a",
			valid: false,
		},
		{
			name:  "invalid character2",
			input: "adgeFG-98",
			valid: false,
		},
		{
			name:  "valid",
			input: "dfakF88u7_",
			valid: true,
		},
	}
	validationFunction := DatashareName()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validationFunction(tt.input, "")
			valid := err == nil
			if valid != tt.valid {
				t.Errorf("expected valid status %t but got %t for input %s", tt.valid, valid, tt.input)
			}
		})
	}
}
