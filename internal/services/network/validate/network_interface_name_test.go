// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestNetworkInterfaceName(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			Input: "",
			Valid: false,
		},
		{
			Input: "-a",
			Valid: false,
		},
		{
			Input: "1_n.i-c",
			Valid: true,
		},
		{
			Input: "01nic",
			Valid: true,
		},
		{
			Input: "01nic01",
			Valid: true,
		},
		{
			Input: strings.Repeat("s", 63),
			Valid: true,
		},
		{
			Input: strings.Repeat("s", 64),
			Valid: true,
		},
		{
			Input: strings.Repeat("s", 65),
			Valid: false,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := NetworkInterfaceName(tc.Input, "name")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
