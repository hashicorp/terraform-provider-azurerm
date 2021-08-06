package validate

import "testing"

func TestFlexibleServerSkuName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "GP_Standard_E64s_v3",
			input: "GP_Standard_E64s_v3",
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
			name:  "B_Standard_E32s_v3",
			input: "B_Standard_E32s_v3",
			valid: true,
		},
		{
			name:  "B_Standard_E30s_v3",
			input: "B_Standard_E30s_v3",
			valid: false,
		},
		{
			name:  "MO_Standard_E16s",
			input: "MO_Standard_E16s",
			valid: false,
		},
		{
			name:  "MO_Standard_E2s_v3",
			input: "MO_Standard_E2s_v3",
			valid: true,
		},
		{
			name:  "B_Standard_B1ms",
			input: "B_Standard_B1ms",
			valid: true,
		},
		{
			name:  "B_Standard_B1",
			input: "B_Standard_B1",
			valid: false,
		},
		{
			name:  "MO_Standard_D2s_v3",
			input: "MO_Standard_D2s_v3",
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
