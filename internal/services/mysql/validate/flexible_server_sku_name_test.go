// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestFlexibleServerSkuName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "GP_Standard_D2ds_v4",
			input: "GP_Standard_D2ds_v4",
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
			name:  "MO_Standard_E2ds",
			input: "MO_Standard_E2ds",
			valid: false,
		},
		{
			name:  "MO_Standard_E16ds_v4",
			input: "MO_Standard_E16ds_v4",
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
			name:  "MO_Standard_E64ds_v4",
			input: "MO_Standard_E64ds_v4",
			valid: true,
		},
		{
			name:  "MO_Standard_E96ds_v5",
			input: "MO_Standard_E96ds_v5",
			valid: true,
		},
		{
			name:  "B_Standard_B20ms",
			input: "B_Standard_B20ms",
			valid: true,
		},
		{
			name:  "GP_Standard_D64ds_v4",
			input: "GP_Standard_D64ds_v4",
			valid: true,
		},
		{
			name:  "B_Standard_D64ds_v4",
			input: "B_Standard_D64ds_v4",
			valid: false,
		},
		{
			name:  "MO_Standard_B2s",
			input: "MO_Standard_B2s",
			valid: false,
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
