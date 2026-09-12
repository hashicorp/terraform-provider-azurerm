// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestStorageDiscoveryWorkspaceName(t *testing.T) {
	testCases := []struct {
		input       string
		shouldError bool
	}{
		{"", true},
		{"abc", true},
		{"abcd", false},
		{"Workspace123", false},
		{"workspace-name", false},
		{"1workspace", true},
		{"-workspace", true},
		{"workspace-", true},
		{"work--space", true},
		{"work_space", true},
		{"work space", true},
		{"workspace\n", true},
		{strings.Repeat("a", 64), false},
		{strings.Repeat("a", 65), true},
	}

	for _, tc := range testCases {
		_, errs := StorageDiscoveryWorkspaceName()(tc.input, "name")
		if tc.shouldError != (len(errs) != 0) {
			t.Errorf("unexpected validation result for %q: %v", tc.input, errs)
		}
	}
}
