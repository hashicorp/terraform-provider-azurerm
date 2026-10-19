// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validation

import (
	"testing"
)

func TestIsCIDRIPv4(t *testing.T) {
	cases := []struct {
		CIDR   string
		Errors int
	}{
		{
			CIDR:   "",
			Errors: 1,
		},
		{
			CIDR:   "0.0.0.0",
			Errors: 0,
		},
		{
			CIDR:   "127.0.0.1/8",
			Errors: 0,
		},
		{
			CIDR:   "127.0.0.1/33",
			Errors: 1,
		},
		{
			CIDR:   "127.0.0.1/-1",
			Errors: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.CIDR, func(t *testing.T) {
			_, errors := IsCIDRIPv4(tc.CIDR, "test")

			if len(errors) != tc.Errors {
				t.Fatalf("Expected CIDR to return %d error(s) not %d", tc.Errors, len(errors))
			}
		})
	}
}
