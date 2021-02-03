package validate

import "testing"

func TestApimSkuName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "Consumption_0",
			input: "Consumption_0",
			valid: true,
		},
		{
			name:  "Consumption_1",
			input: "Consumption_1",
			valid: false,
		},
		{
			name:  "Basic_3",
			input: "Basic_3",
			valid: false,
		},
		{
			name:  "Basic_1",
			input: "Basic_1",
			valid: true,
		},
		{
			name:  "Developer_1",
			input: "Developer_1",
			valid: true,
		},
		{
			name:  "Premium_0",
			input: "Premium_0",
			valid: false,
		},
		{
			name:  "Premium_11",
			input: "Premium_11",
			valid: false,
		},
		{
			name:  "Premium_7",
			input: "Premium_7",
			valid: true,
		},
		{
			name:  "Standard_7",
			input: "Standard_7",
			valid: false,
		},
		{
			name:  "standard_2",
			input: "standard_2",
			valid: false,
		},
		{
			name:  "PREMIUM_7",
			input: "PREMIUM_7",
			valid: false,
		},
	}
	var validationFunction = ApimSkuName()
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
