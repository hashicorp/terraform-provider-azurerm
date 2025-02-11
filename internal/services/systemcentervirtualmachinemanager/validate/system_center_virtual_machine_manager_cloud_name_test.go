// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestSystemCenterVirtualMachineManagerCloudName(t *testing.T) {
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
			Input:    "a8a",
			Expected: true,
		},
		{
			Input:    "a-8.a",
			Expected: true,
		},
		{
			Input:    "a-",
			Expected: false,
		},
		{
			Input:    "a.",
			Expected: false,
		},
		{
			Input:    strings.Repeat("s", 53),
			Expected: true,
		},
		{
			Input:    strings.Repeat("s", 54),
			Expected: true,
		},
		{
			Input:    strings.Repeat("s", 55),
			Expected: false,
		},
	}

	for _, v := range testCases {
		_, errors := SystemCenterVirtualMachineManagerCloudName(v.Input, "name")
		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %t but got %t (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
