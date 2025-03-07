// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestSystemCenterVirtualMachineManagerVirtualMachineInstanceStorageDiskName(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected bool
	}{
		{
			Input:    "",
			Expected: false,
		},
		{
			Input:    "test-test_",
			Expected: false,
		},
		{
			Input:    "s",
			Expected: true,
		},
		{
			Input:    "test-_test",
			Expected: true,
		},
		{
			Input:    "test",
			Expected: true,
		},
		{
			Input:    "test_1",
			Expected: true,
		},
		{
			Input:    strings.Repeat("s", 79),
			Expected: true,
		},
		{
			Input:    strings.Repeat("s", 80),
			Expected: true,
		},
		{
			Input:    strings.Repeat("s", 81),
			Expected: false,
		},
	}

	for _, v := range testCases {
		_, errors := SystemCenterVirtualMachineManagerVirtualMachineInstanceStorageDiskName(v.Input, "test")
		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %t but got %t (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
