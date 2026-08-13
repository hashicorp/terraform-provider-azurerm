// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestFrontDoorRuleRequestPathMatchValue(t *testing.T) {
	testCases := []struct {
		name  string
		input interface{}
		valid bool
	}{
		{
			name:  "empty",
			input: "",
			valid: false,
		},
		{
			name:  "root path",
			input: "/",
			valid: true,
		},
		{
			name:  "path without leading slash",
			input: "legacy-login",
			valid: true,
		},
		{
			name:  "path with leading slash",
			input: "/legacy-login",
			valid: false,
		},
		{
			name:  "non-string value",
			input: 1,
			valid: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, errors := FrontDoorRuleRequestPathMatchValue(testCase.input, "test")
			valid := len(errors) == 0

			if testCase.valid != valid {
				t.Fatalf("expected %t but got %t", testCase.valid, valid)
			}
		})
	}
}
