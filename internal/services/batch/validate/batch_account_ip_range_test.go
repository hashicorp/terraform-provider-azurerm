// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestBatchAccountIpRange(t *testing.T) {
	cases := []struct {
		IPRange string
		Errors  int
	}{
		{
			IPRange: "0.0.0.0",
			Errors:  0,
		},
		{
			IPRange: "23.45.0.1/30",
			Errors:  0,
		},
		{
			IPRange: "",
			Errors:  1,
		},
		{
			IPRange: "23.45.0.1/31",
			Errors:  1,
		},
		{
			IPRange: "23.45.0.1/32",
			Errors:  1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.IPRange, func(t *testing.T) {
			_, errors := BatchAccountIpRange(tc.IPRange, "ip_range")

			if len(errors) != tc.Errors {
				t.Fatalf("Expected BatchAccountIpRange to return %d error(s) not %d", tc.Errors, len(errors))
			}
		})
	}
}
