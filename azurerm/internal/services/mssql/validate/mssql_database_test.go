package validate

import (
	"testing"
)

func TestValidateMsSqlDatabaseAutoPauseDelay(t *testing.T) {
	testCases := []struct {
		input       string
		shouldError bool
	}{
		{"-1", false},
		{"-2", true},
		{"30", true},
		{"60", false},
		{"65", true},
		{"360", false},
		{"19900", true},
	}

	for _, test := range testCases {
		_, es := ValidateMsSqlDatabaseAutoPauseDelay(test.input, "name")

		if test.shouldError && len(es) == 0 {
			t.Fatalf("Expected validating name %q to fail", test.input)
		}
	}
}

func TestValidateMsSqlDBMinCapacity(t *testing.T) {
	testCases := []struct {
		input       string
		shouldError bool
	}{
		{"-1", true},
		{"0.25", true},
		{"0.5", false},
		{"1.25", false},
		{"2", false},
		{"2.25", true},
	}

	for _, test := range testCases {
		_, es := ValidateMsSqlDBMinCapacity(test.input, "name")

		if test.shouldError && len(es) == 0 {
			t.Fatalf("Expected validating name %q to fail", test.input)
		}
	}
}

func TestValidateMsSqlDBSkuName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "Valid GP",
			input: "GP_Gen5_2",
			valid: true,
		},
		{
			name:  "Valid Serverless GP",
			input: "GP_S_Gen5_2",
			valid: true,
		},
		{
			name:  "Valid HS",
			input: "HS_Gen4_1",
			valid: true,
		},
		{
			name:  "Valid BC",
			input: "BC_Gen5_4",
			valid: true,
		},
		{
			name:  "Valid Standard",
			input: "Standard",
			valid: true,
		},
		{
			name:  "Valid Basic",
			input: "Basic",
			valid: true,
		},
		{
			name:  "Valid Premium",
			input: "Premium",
			valid: true,
		},
		{
			name:  "empty",
			input: "",
			valid: false,
		},
		{
			name:  "Extra dot",
			input: "BC_Gen5_3.",
			valid: false,
		},
		{
			name:  "Wrong capacity",
			input: "BC_Gen5_3",
			valid: false,
		},
		{
			name:  "Wrong Family",
			input: "BC_Inv_2",
			valid: false,
		},
		{
			name:  "Wrong Serverless",
			input: "BC_S_Gen4_2",
			valid: false,
		},
		{
			name:  "Wrong Serverless",
			input: "BC_S_Gen4_2",
			valid: false,
		},
		{
			name:  "Lower case",
			input: "bc_gen5_2",
			valid: true,
		},
	}
	var validationFunction = ValidateMsSqlDBSkuName()
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
