// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

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
			name:  "Valid Fsv2",
			input: "GP_Fsv2_10",
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
			name:  "Valid Free",
			input: "Free",
			valid: true,
		},
		{
			name:  "Valid Basic",
			input: "Basic",
			valid: true,
		},
		{
			name:  "Valid ElasticPool",
			input: "ElasticPool",
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
		{
			name:  "Valid DC",
			input: "GP_DC_8",
			valid: true,
		},
		{
			name:  "Valid DC",
			input: "BC_DC_2",
			valid: true,
		},
		{
			name:  "Valid DC",
			input: "HS_DC_6",
			valid: true,
		},
		{
			name:  "Invalid DC",
			input: "HS_DC_80",
			valid: false,
		},
		{
			name:  "Valid 8IM",
			input: "HS_8IM_24",
			valid: true,
		},
		{
			name:  "Valid Serverless 8IM",
			input: "HS_S_8IM_24",
			valid: true,
		},
		{
			name:  "Invalid Serverless 8IM",
			input: "HS_S_8IM_48",
			valid: false,
		},
		{
			name:  "Valid Premium 8IM",
			input: "HS_PRMS_8",
			valid: true,
		},
		{
			name:  "Invalid Premium 8IM",
			input: "HS_PRMS_28",
			valid: false,
		},
		{
			name:  "Valid Premium 8IH",
			input: "HS_MOPRMS_8",
			valid: true,
		},
		{
			name:  "Invalid Premium 8IH",
			input: "HS_MOPRMS_128",
			valid: false,
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

func TestDatabaseSkuNameWithoutFree(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "Free",
			input: "Free",
			valid: false,
		},
		{
			name:  "Free lower case",
			input: "free",
			valid: false,
		},
		{
			name:  "Valid Serverless GP",
			input: "GP_S_Gen5_2",
			valid: true,
		},
		{
			name:  "Valid Basic",
			input: "Basic",
			valid: true,
		},
		{
			name:  "Valid ElasticPool",
			input: "ElasticPool",
			valid: true,
		},
		{
			name:  "empty",
			input: "",
			valid: false,
		},
		{
			name:  "Wrong Family",
			input: "BC_Inv_2",
			valid: false,
		},
	}
	validationFunction := DatabaseSkuNameWithoutFree()
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

func TestDatabaseSkuName_FreeAliasBoundary(t *testing.T) {
	// The legacy `Free` alias must stay valid for the pre-5.0 validator but be rejected by the
	// 5.0-mode validator, which replaces it with the `free_limit_enabled`/`free_limit_exhaustion_behavior` properties.
	withFree := DatabaseSkuName()
	withoutFree := DatabaseSkuNameWithoutFree()

	for _, input := range []string{"Free", "free"} {
		if _, err := withFree(input, ""); err != nil {
			t.Errorf("expected DatabaseSkuName to accept %q, got error: %+v", input, err)
		}
		if _, err := withoutFree(input, ""); err == nil {
			t.Errorf("expected DatabaseSkuNameWithoutFree to reject %q, but it was accepted", input)
		}
	}
}
