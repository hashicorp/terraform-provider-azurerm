// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"testing"
)

func TestFlexibleServerSkuName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "GP_Standard_E64s_v3",
			input: "GP_Standard_E64s_v3",
			valid: false,
		},
		{
			name:  "GP_Standard_D64s_v3",
			input: "GP_Standard_D64s_v3",
			valid: true,
		},
		{
			name:  "GP_Standard_D16ds_v4",
			input: "GP_Standard_D16ds_v4",
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
			valid: false,
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
			name:  "MO_Standard_E64ds_v3",
			input: "MO_Standard_E64ds_v3",
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
			name:  "MO_Standard_D2s_v3",
			input: "MO_Standard_D2s_v3",
			valid: false,
		},
		{
			name:  "MO_Standard_E16ds_v4",
			input: "MO_Standard_E16ds_v4",
			valid: true,
		},
		{
			name:  "B_Standard_B20ms",
			input: "B_Standard_B20ms",
			valid: true,
		},
		{
			name:  "GP_Standard_D16ds_v5",
			input: "GP_Standard_D16ds_v5",
			valid: true,
		},
		{
			name:  "GP_Standard_D16ads_v5",
			input: "GP_Standard_D16ads_v5",
			valid: true,
		},
		{
			name:  "MO_Standard_E16ds_v5",
			input: "MO_Standard_E16ds_v5",
			valid: true,
		},
		{
			name:  "MO_Standard_E16ads_v5",
			input: "MO_Standard_E16ads_v5",
			valid: true,
		},
		{
			name:  "MO_Standard_E96ds_v5",
			input: "MO_Standard_E96ds_v5",
			valid: true,
		},
		{
			name:  "GP_Standard_DC2ads_v5",
			input: "GP_Standard_DC2ads_v5",
			valid: true,
		},
		{
			name:  "GP_Standard_DC64ads_v5",
			input: "GP_Standard_DC64ads_v5",
			valid: true,
		},
		{
			name:  "GP_Standard_DC96ads_v5",
			input: "GP_Standard_DC96ads_v5",
			valid: true,
		},
		{
			name:  "MO_Standard_EC2ads_v5",
			input: "MO_Standard_EC2ads_v5",
			valid: true,
		},
		{
			name:  "MO_Standard_EC20ads_v5",
			input: "MO_Standard_EC20ads_v5",
			valid: true,
		},
		{
			name:  "MO_Standard_EC96ads_v5",
			input: "MO_Standard_EC96ads_v5",
			valid: true,
		},
		{
			name:  "MO_Standard_EC96as_v5",
			input: "MO_Standard_EC96as_v5",
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

func TestFlexibleServerSkuNameChange(t *testing.T) {
	tests := []struct {
		name  string
		sku1  string
		sku2  string
		valid bool
	}{
		{
			name:  "Non-confidential identical",
			sku1:  "GP_Standard_D16ds_v4",
			sku2:  "GP_Standard_D16ds_v4",
			valid: true,
		},
		{
			name:  "Non-confidential update",
			sku1:  "GP_Standard_D16ds_v4",
			sku2:  "MO_Standard_E16ds_v5",
			valid: true,
		},
		{
			name:  "Confidential identical",
			sku1:  "GP_Standard_DC96ads_v5",
			sku2:  "GP_Standard_DC96ads_v5",
			valid: true,
		},
		{
			name:  "Confidential update",
			sku1:  "GP_Standard_DC64ads_v5",
			sku2:  "MO_Standard_EC96ads_v5",
			valid: true,
		},
		{
			name:  "Non-confidential to confidential",
			sku1:  "GP_Standard_D16ds_v4",
			sku2:  "MO_Standard_EC2ads_v5",
			valid: false,
		},
		{
			name:  "Confidential to non-confidential",
			sku1:  "MO_Standard_EC20ads_v5",
			sku2:  "MO_Standard_E16ads_v5",
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := FlexibleServerSkuNameChange(tt.sku1, tt.sku2)
			valid := err == nil
			if valid != tt.valid {
				t.Errorf("Expected valid status %t but got %t for input %s", tt.valid, valid, fmt.Sprintf("%s <-> %s", tt.sku1, tt.sku2))
			}
		})
	}
}
