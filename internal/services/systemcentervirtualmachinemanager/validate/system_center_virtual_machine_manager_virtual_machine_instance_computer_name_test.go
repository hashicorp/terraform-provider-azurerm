// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestSystemCenterVirtualMachineManagerVirtualMachineInstanceComputerName(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected bool
	}{
		{
			Input:    "",
			Expected: false,
		},
		{
			Input:    "test-test",
			Expected: false,
		},
		{
			Input:    "s",
			Expected: true,
		},
		{
			Input:    "test123",
			Expected: true,
		},
	}

	for _, v := range testCases {
		_, errors := SystemCenterVirtualMachineManagerVirtualMachineInstanceComputerName(v.Input, "test")
		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %t but got %t (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
