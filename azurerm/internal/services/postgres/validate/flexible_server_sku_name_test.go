package validate

import "testing"

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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := FlexibleServerSkuName(tt.input, "name")
			valid := err == nil
			if valid != tt.valid {
				t.Errorf("Expected valid status %t but got %t for input %s", tt.valid, valid, tt.input)
			}
		})
	}
}
