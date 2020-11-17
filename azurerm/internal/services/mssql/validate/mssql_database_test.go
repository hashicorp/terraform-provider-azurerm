package validate

import (
	"testing"
)

func TestMsSqlDatabaseAutoPauseDelay(t *testing.T) {
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
		_, es := MsSqlDatabaseAutoPauseDelay(test.input, "name")

		if test.shouldError && len(es) == 0 {
			t.Fatalf("Expected validating name %q to fail", test.input)
		}
	}
}

func TestMsSqlDBSkuName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "DataWarehouse",
			input: "DW100c",
			valid: true,
		},
		{
			name:  "DataWarehouse",
			input: "DW102c",
			valid: false,
		},
		{
			name:  "Stretch",
			input: "DS100",
			valid: true,
		},
		{
			name:  "Stretch",
			input: "DS1001",
			valid: false,
		},
		{
			name:  "Valid GP",
			input: "GP_Gen4_3",
			valid: true,
		},
		{
			name:  "Valid Serverless GP",
			input: "GP_S_Gen5_2",
			valid: true,
		},
		{
			name:  "Valid HS",
			input: "HS_Gen5_2",
			valid: true,
		},
		{
			name:  "Valid BC",
			input: "BC_Gen4_5",
			valid: true,
		},
		{
			name:  "Valid BC",
			input: "BC_M_12",
			valid: true,
		},
		{
			name:  "Valid BC",
			input: "BC_Gen5_14",
			valid: true,
		},
		{
			name:  "Valid Standard",
			input: "S3",
			valid: true,
		},
		{
			name:  "Valid Basic",
			input: "Basic",
			valid: true,
		},
		{
			name:  "Valid Premium",
			input: "P15",
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
			input: "GP_S_Gen4_2",
			valid: false,
		},
		{
			name:  "Wrong Serverless",
			input: "BC_S_Gen5_2",
			valid: false,
		},
		{
			name:  "Lower case",
			input: "bc_gen5_2",
			valid: true,
		},
	}
	validationFunction := MsSqlDBSkuName()
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

func TestMsSqlDBCollation(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "SQL Collation",
			input: "SQL_Latin1_General_CP1_CI_AS",
			valid: true,
		},
		{
			name:  "Windows Collation",
			input: "Latin1_General_100_CI_AS_SC",
			valid: true,
		},
		{
			name:  "SQL Collation",
			input: "SQL_AltDiction_CP850_CI_AI",
			valid: true,
		},
		{
			name:  "SQL Collation",
			input: "SQL_Croatian_CP1250_CI_AS",
			valid: true,
		},
		{
			name:  "Windows Collation",
			input: "Chinese_Hong_Kong_Stroke_90_CI_AI",
			valid: true,
		},
		{
			name:  "Windows Collation",
			input: "Japanese_BIN",
			valid: true,
		},
		{
			name:  "lowercase",
			input: "sql_croatian_cp1250_ci_as",
			valid: false,
		},
		{
			name:  "extra dot",
			input: "SQL_Croatian_CP1250.",
			valid: false,
		},
		{
			name:  "Invalid collation",
			input: "CDD",
			valid: false,
		},
		{
			name:  "Double definition",
			input: "Latin1_General_100_CI_CS",
			valid: false,
		},
	}
	validationFunction := MsSqlDBCollation()
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
