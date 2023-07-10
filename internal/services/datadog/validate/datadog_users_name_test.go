// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestDatadogUsersName(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected bool
	}{
		{
			Input:    "",
			Expected: false,
		},
		{
			Input:    "Test",
			Expected: true,
		},
		{
			Input:    "qwertyuiopasdfghjklzxcvbnmqwertyuiopasdfghjklzxcvbnm",
			Expected: false,
		},
	}

	for _, v := range testCases {
		_, errors := DatadogUsersName(v.Input, "user_name")
		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %t but got %t (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
