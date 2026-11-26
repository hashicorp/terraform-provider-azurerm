// Copyright IBM Corp. 2014, 2026
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
			name:  "GP_Standard_E64ads_v5",
			input: "GP_Standard_E64ads_v5",
			valid: false,
		},
		{
			name:  "GP_Standard_D64ads_v5",
			input: "GP_Standard_D64ads_v5",
			valid: true,
		},
		{
			name:  "GP_Standard_D16ds_v4",
			input: "GP_Standard_D16ds_v4",
			valid: true,
		},
		{
			name:  "GP_Standard_D4ds_v5",
			input: "GP_Standard_D4ds_v5",
			valid: true,
		},
		{
			name:  "GP_Standard_D4ds_v6",
			input: "GP_Standard_D4ds_v6",
			valid: true,
		},
		{
			name:  "GP_Standard_D16ads_v6",
			input: "GP_Standard_D16ads_v6",
			valid: true,
		},
		{
			name:  "GP_Standard_D96ads_v6",
			input: "GP_Standard_D96ads_v6",
			valid: true,
		},
		{
			name:  "GP_Standard_D128ds_v6",
			input: "GP_Standard_D128ds_v6",
			valid: true,
		},
		{
			name:  "GP_Standard_D192ds_v6",
			input: "GP_Standard_D192ds_v6",
			valid: true,
		},
		{
			name:  "GP_Standard_D256ds_v6",
			input: "GP_Standard_D256ds_v6",
			valid: false,
		},
		{
			name:  "MO_Standard_E20ds_v6",
			input: "MO_Standard_E20ds_v6",
			valid: true,
		},
		{
			name:  "MO_Standard_E128ds_v6",
			input: "MO_Standard_E128ds_v6",
			valid: true,
		},
		{
			name:  "GP_Standard_D12ads_v6",
			input: "GP_Standard_D12ads_v6",
			valid: false,
		},
		{
			name:  "MO_Standard_E20ads_v6",
			input: "MO_Standard_E20ads_v6",
			valid: true,
		},
		{
			name:  "MO_Standard_E96ads_v6",
			input: "MO_Standard_E96ads_v6",
			valid: true,
		},
		{
			name:  "MO_Standard_E12ads_v6",
			input: "MO_Standard_E12ads_v6",
			valid: false,
		},
		{
			name:  "GP_Standard_D2s_v5",
			input: "GP_Standard_D2s_v5",
			valid: true,
		},
		{
			name:  "GP_Standard_D96as_v5",
			input: "GP_Standard_D96as_v5",
			valid: true,
		},
		{
			name:  "GP_Standard_D4s_v6",
			input: "GP_Standard_D4s_v6",
			valid: true,
		},
		{
			name:  "GP_Standard_D192s_v6",
			input: "GP_Standard_D192s_v6",
			valid: true,
		},
		{
			name:  "GP_Standard_D2as_v6",
			input: "GP_Standard_D2as_v6",
			valid: true,
		},
		{
			name:  "GP_Standard_D128as_v6",
			input: "GP_Standard_D128as_v6",
			valid: false,
		},
		{
			name:  "GP_Standard_D20s_v6",
			input: "GP_Standard_D20s_v6",
			valid: false,
		},
		{
			name:  "MO_Standard_E20s_v5",
			input: "MO_Standard_E20s_v5",
			valid: true,
		},
		{
			name:  "MO_Standard_E20ads_v5",
			input: "MO_Standard_E20ads_v5",
			valid: true,
		},
		{
			name:  "MO_Standard_E96as_v5",
			input: "MO_Standard_E96as_v5",
			valid: true,
		},
		{
			name:  "MO_Standard_E104is_v5",
			input: "MO_Standard_E104is_v5",
			valid: true,
		},
		{
			name:  "MO_Standard_E112ias_v5",
			input: "MO_Standard_E112ias_v5",
			valid: true,
		},
		{
			name:  "MO_Standard_E128s_v6",
			input: "MO_Standard_E128s_v6",
			valid: true,
		},
		{
			name:  "MO_Standard_E20as_v6",
			input: "MO_Standard_E20as_v6",
			valid: true,
		},
		{
			name:  "MO_Standard_E192is_v6",
			input: "MO_Standard_E192is_v6",
			valid: true,
		},
		{
			name:  "MO_Standard_E192ids_v6",
			input: "MO_Standard_E192ids_v6",
			valid: true,
		},
		{
			name:  "MO_Standard_E192s_v6",
			input: "MO_Standard_E192s_v6",
			valid: false,
		},
		{
			name:  "MO_Standard_E104is_v6",
			input: "MO_Standard_E104is_v6",
			valid: false,
		},
		{
			name:  "GP_Standard_E4s_v6",
			input: "GP_Standard_E4s_v6",
			valid: false,
		},
		{
			name:  "MO_Standard_D4s_v6",
			input: "MO_Standard_D4s_v6",
			valid: false,
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
			name:  "B_Standard_E32ads_v5",
			input: "B_Standard_E32ads_v5",
			valid: false,
		},
		{
			name:  "B_Standard_E30ads_v5",
			input: "B_Standard_E30ads_v5",
			valid: false,
		},
		{
			name:  "MO_Standard_E16s",
			input: "MO_Standard_E16s",
			valid: false,
		},
		{
			name:  "MO_Standard_E2ads_v5",
			input: "MO_Standard_E2ads_v5",
			valid: true,
		},
		{
			name:  "MO_Standard_E64dads_v5",
			input: "MO_Standard_E64dads_v5",
			valid: false,
		},
		{
			name:  "B_Standard_B1ms",
			input: "B_Standard_B1ms",
			valid: true,
		},
		{
			name:  "B_Standard_B2ats_v2",
			input: "B_Standard_B2ats_v2",
			valid: true,
		},
		{
			name:  "B_Standard_B2afs_v2",
			input: "B_Standard_B2afs_v2",
			valid: true,
		},
		{
			name:  "B_Standard_B2als_v2",
			input: "B_Standard_B2als_v2",
			valid: true,
		},
		{
			name:  "B_Standard_B2as_v2",
			input: "B_Standard_B2as_v2",
			valid: true,
		},
		{
			name:  "B_Standard_B4as_v2",
			input: "B_Standard_B4as_v2",
			valid: true,
		},
		{
			name:  "B_Standard_B8as_v2",
			input: "B_Standard_B8as_v2",
			valid: true,
		},
		{
			name:  "B_Standard_B16as_v2",
			input: "B_Standard_B16as_v2",
			valid: true,
		},
		{
			name:  "B_Standard_B32as_v2",
			input: "B_Standard_B32as_v2",
			valid: true,
		},
		{
			name:  "B_Standard_B1as_v2",
			input: "B_Standard_B1as_v2",
			valid: false,
		},
		{
			name:  "B_Standard_B2ats",
			input: "B_Standard_B2ats",
			valid: false,
		},
		{
			name:  "B_Standard_B64as_v2",
			input: "B_Standard_B64as_v2",
			valid: false,
		},
		{
			name:  "B_Standard_B2ats_v3",
			input: "B_Standard_B2ats_v3",
			valid: false,
		},
		{
			name:  "B_Standard_B1",
			input: "B_Standard_B1",
			valid: false,
		},
		{
			name:  "MO_Standard_D2ads_v5",
			input: "MO_Standard_D2ads_v5",
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
			name:  "MO_Standard_E96ds_v6",
			input: "MO_Standard_E96ds_v6",
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
