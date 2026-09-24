// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package helper

import (
	"testing"
)

func TestCompareDatabaseSkuScaleUp(t *testing.T) {
	cases := []struct {
		Name     string
		Sku1     string
		Sku2     string
		Expected bool
	}{
		{
			Name:     "higher tier DTU",
			Sku1:     "P1",
			Sku2:     "S3",
			Expected: true,
		},
		{
			Name:     "lower tier DTU",
			Sku1:     "S3",
			Sku2:     "P1",
			Expected: false,
		},
		{
			Name:     "higher tier vCore over DTU",
			Sku1:     "GP_Gen5_2",
			Sku2:     "S3",
			Expected: true,
		},
		{
			Name:     "business critical higher than premium",
			Sku1:     "BC_Gen5_2",
			Sku2:     "P1",
			Expected: true,
		},
		{
			Name:     "same tier higher capacity DTU",
			Sku1:     "S3",
			Sku2:     "S1",
			Expected: true,
		},
		{
			// Regression: databaseSkuCapacity("S0") is 0, so the previous capacity2 > 0 guard
			// incorrectly reported that S1 was not a scale up from S0.
			Name:     "scale up from S0",
			Sku1:     "S1",
			Sku2:     "S0",
			Expected: true,
		},
		{
			Name:     "scale up from S0 to highest standard",
			Sku1:     "S12",
			Sku2:     "S0",
			Expected: true,
		},
		{
			Name:     "scale down to S0",
			Sku1:     "S0",
			Sku2:     "S1",
			Expected: false,
		},
		{
			Name:     "same DTU sku",
			Sku1:     "S1",
			Sku2:     "S1",
			Expected: false,
		},
		{
			Name:     "same S0 sku",
			Sku1:     "S0",
			Sku2:     "S0",
			Expected: false,
		},
		{
			Name:     "same tier higher capacity vCore",
			Sku1:     "GP_Gen5_8",
			Sku2:     "GP_Gen5_2",
			Expected: true,
		},
		{
			Name:     "same tier lower capacity vCore",
			Sku1:     "GP_Gen5_2",
			Sku2:     "GP_Gen5_8",
			Expected: false,
		},
		{
			Name:     "same vCore sku",
			Sku1:     "GP_Gen5_2",
			Sku2:     "GP_Gen5_2",
			Expected: false,
		},
		{
			Name:     "serverless general purpose higher capacity",
			Sku1:     "GP_S_Gen5_4",
			Sku2:     "GP_S_Gen5_2",
			Expected: true,
		},
		{
			Name:     "different hardware family same tier",
			Sku1:     "GP_Fsv2_8",
			Sku2:     "GP_Gen5_2",
			Expected: false,
		},
		{
			Name:     "same premium sku",
			Sku1:     "P1",
			Sku2:     "P1",
			Expected: false,
		},
		{
			Name:     "higher premium capacity",
			Sku1:     "P2",
			Sku2:     "P1",
			Expected: true,
		},
		{
			Name:     "basic same sku",
			Sku1:     "Basic",
			Sku2:     "Basic",
			Expected: false,
		},
		{
			Name:     "empty sku1",
			Sku1:     "",
			Sku2:     "S1",
			Expected: false,
		},
		{
			Name:     "empty sku2",
			Sku1:     "S1",
			Sku2:     "",
			Expected: false,
		},
		{
			Name:     "both empty",
			Sku1:     "",
			Sku2:     "",
			Expected: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			if actual := CompareDatabaseSkuScaleUp(tc.Sku1, tc.Sku2); actual != tc.Expected {
				t.Fatalf("expected CompareDatabaseSkuScaleUp(%q, %q) to be %t, got %t", tc.Sku1, tc.Sku2, tc.Expected, actual)
			}
		})
	}
}
