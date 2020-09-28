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
	var validationFunction = FlexibleServerName()
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

func TestFlexibleServerSkuName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "Standard_E64s_v3",
			input: "Standard_E64s_v3",
			valid: true,
		},
		{
			name:  "Standard",
			input: "Standard",
			valid: false,
		},
		{
			name:  "Empty",
			input: "",
			valid: false,
		},
		{
			name:  "Standard_E32s_v3",
			input: "Standard_E32s_v3",
			valid: true,
		},
		{
			name:  "Standard_E30s_v3",
			input: "Standard_E30s_v3",
			valid: false,
		},
		{
			name:  "Standard_E16s",
			input: "Standard_E16s",
			valid: false,
		},
		{
			name:  "Standard_E2s_v3",
			input: "Standard_E2s_v3",
			valid: true,
		},
		{
			name:  "Standard_B1ms",
			input: "Standard_B1ms",
			valid: true,
		},
		{
			name:  "Standard_B1",
			input: "Standard_B1",
			valid: false,
		},
		{
			name:  "Standard_D2s_v3",
			input: "Standard_D2s_v3",
			valid: true,
		},
	}
	var validationFunction = FlexibleServerSkuName()
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
