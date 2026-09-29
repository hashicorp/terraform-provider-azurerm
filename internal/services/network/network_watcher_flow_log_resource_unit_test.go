// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network

import "testing"

func TestNetworkWatcherFlowLogVersionValidation(t *testing.T) {
	validateFunc := resourceNetworkWatcherFlowLog().Schema["version"].ValidateFunc

	testCases := []struct {
		input       int
		expectError bool
	}{
		{input: -1, expectError: true},
		{input: 0, expectError: false},
		{input: 1, expectError: false},
		{input: 2, expectError: false},
		{input: 3, expectError: false},
		{input: 100, expectError: false},
	}

	for _, tc := range testCases {
		_, errors := validateFunc(tc.input, "version")

		hasError := len(errors) > 0
		if hasError != tc.expectError {
			t.Fatalf("expected version %d to produce error=%t, got error=%t (errors: %v)", tc.input, tc.expectError, hasError, errors)
		}
	}
}
