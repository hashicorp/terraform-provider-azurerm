// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestStorageDiscoveryScopeDisplayName(t *testing.T) {
	testCases := []struct {
		input       string
		shouldError bool
	}{
		{"", true},
		{"ab", true},
		{"abc", true},
		{"abcd", false},
		{"TestScopeOne", false},
		{"TestScopeTwo", false},
		{" Production Storage ", true},
		{"Production Storage", false},
		{"a1b2", true},
		{"1abc", true},
		{"abc1", true},
		{"ab  cd", true},
		{"ab--cd", true},
		{"-abc", true},
		{"abc-", true},
		{" abc", true},
		{"abc ", true},
		{"A-B C", false},
		{strings.Repeat("a", 65), true},
	}
	for _, tc := range testCases {
		_, errs := StorageDiscoveryScopeDisplayName(tc.input, "display_name")
		if tc.shouldError && len(errs) == 0 {
			t.Errorf("Expected %q to fail validation", tc.input)
		}
		if !tc.shouldError && len(errs) > 0 {
			t.Errorf("Expected %q to pass validation, got: %v", tc.input, errs)
		}
	}
}
