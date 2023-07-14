// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestDatadogMonitorsPhoneNumber(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected bool
	}{
		{
			Input:    "",
			Expected: true,
		},
		{
			Input:    "1234567890",
			Expected: true,
		},
		{
			Input:    "12345678901234567890123456789012345678901234567890",
			Expected: false,
		},
	}

	for _, v := range testCases {
		_, errors := DatadogMonitorsPhoneNumber(v.Input, "phone_number")
		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %t but got %t (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
