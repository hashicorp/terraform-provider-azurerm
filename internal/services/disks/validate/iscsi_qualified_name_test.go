// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestStorageDisksPoolIscsiIqnName(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			Input: "iqn.2021-11.com.microsoft",
			Valid: true,
		},
		{
			Input: "iqn.2021-11.com.microsoft:iscsi",
			Valid: true,
		},
		{
			Input: "iqn.2021-11.com.microsoft:",
			Valid: true,
		},
		{
			Input: "iqn.2021-11.com.microsoft:-",
			Valid: true,
		},
		{
			Input: "iqn.2021-11.com..microsoft",
			Valid: true,
		},

		{
			Input: "iqn.2021-11.com.microsoft:_",
			Valid: false,
		},
		{
			Input: "2021-11.com.microsoft",
			Valid: false,
		},
		{
			Input: "iqn.2021-11",
			Valid: false,
		},
		{
			Input: "iqn.a021-11.com.microsoft",
			Valid: false,
		},
		{
			Input: "iqn.2021-l1.com.microsoft",
			Valid: false,
		},
		{
			Input: "iqn.2021-11.com.microsoft:@",
			Valid: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.Input, func(t *testing.T) {
			_, errors := IQN(tc.Input, "test")
			valid := len(errors) == 0

			if tc.Valid != valid {
				t.Fatalf("Case %s Expected %t but got %t", tc.Input, tc.Valid, valid)
			}
		})
	}
}
