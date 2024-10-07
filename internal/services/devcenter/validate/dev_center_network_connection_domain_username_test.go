// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestDevCenterNetworkConnectionDomainUsername(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected bool
	}{
		{
			Input:    "a",
			Expected: false,
		},
		{
			Input:    "abc",
			Expected: false,
		},
		{
			Input:    "123",
			Expected: false,
		},
		{
			Input:    "test.com",
			Expected: false,
		},
		{
			Input:    "test@.com",
			Expected: false,
		},
		{
			Input:    "test.com",
			Expected: false,
		},
		{
			Input:    "tfuser@test.com",
			Expected: true,
		},
	}

	for _, v := range testCases {
		_, errors := DevCenterNetworkConnectionDomainUsername(v.Input, "domain_username")
		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %t but got %t (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
