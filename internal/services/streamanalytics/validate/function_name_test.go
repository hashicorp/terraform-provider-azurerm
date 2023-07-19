// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestFunctionName(t *testing.T) {
	testCases := []struct {
		input       string
		shouldError bool
	}{
		{"", true},
		{"ABC", false},
		{"abc", false},
		{"a-b", false},
		{"ab-", false},
		{"ab-1", false},
		{"ab#-1", false},
		{strings.Repeat("s", 3), false},
		{strings.Repeat("s", 63), false},
		{strings.Repeat("s", 64), true},
	}

	for _, test := range testCases {
		_, es := FunctionName(test.input, "name")

		if test.shouldError && len(es) == 0 {
			t.Fatalf("Expected validating name %q to fail", test.input)
		}
	}
}
