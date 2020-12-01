package validate

import "testing"

func TestDatabaseSkuName(t *testing.T) {
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
	validationFunction := DatabaseSkuName()
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
