// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestDatabaseAutoPauseDelay(t *testing.T) {
	testCases := []struct {
		input       string
		shouldError bool
	}{
		{"-1", false},
		{"-2", true},
		{"10", true},
		{"42", false},
		{"60", false},
		{"65", false},
		{"360", false},
		{"19900", true},
	}

	for _, test := range testCases {
		_, es := DatabaseAutoPauseDelay(test.input, "name")

		if test.shouldError && len(es) == 0 {
			t.Fatalf("Expected validating name %q to fail", test.input)
		}
	}
}
