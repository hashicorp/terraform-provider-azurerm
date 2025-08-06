// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestDevCenterNetworkConnectionDomainName(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected bool
	}{
		{
			Input:    "",
			Expected: false,
		},
		{
			Input:    "a",
			Expected: true,
		},
		{
			Input:    "aa",
			Expected: true,
		},
		{
			Input:    "Aa-",
			Expected: true,
		},
		{
			Input:    "a.a",
			Expected: true,
		},
		{
			Input:    "aa.",
			Expected: false,
		},
		{
			Input:    ".aa",
			Expected: false,
		},
		{
			Input:    strings.Repeat("s", 254),
			Expected: true,
		},
		{
			Input:    strings.Repeat("s", 255),
			Expected: true,
		},
		{
			Input:    strings.Repeat("s", 256),
			Expected: false,
		},
	}

	for _, v := range testCases {
		_, errors := DevCenterNetworkConnectionDomainName(v.Input, "domain_name")
		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %t but got %t (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
