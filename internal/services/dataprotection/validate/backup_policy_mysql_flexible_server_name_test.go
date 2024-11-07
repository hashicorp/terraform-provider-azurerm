// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestBackupPolicyMySQLFlexibleServerName(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected bool
	}{
		{
			Input:    "",
			Expected: false,
		},
		{
			Input:    "aaaa",
			Expected: true,
		},
		{
			Input:    "a8a",
			Expected: true,
		},
		{
			Input:    "a-8.a",
			Expected: false,
		},
		{
			Input:    "a-a",
			Expected: true,
		},
		{
			Input:    "a.a",
			Expected: false,
		},
		{
			Input:    strings.Repeat("s", 149),
			Expected: true,
		},
		{
			Input:    strings.Repeat("s", 150),
			Expected: true,
		},
		{
			Input:    strings.Repeat("s", 151),
			Expected: false,
		},
	}

	for _, v := range testCases {
		_, errors := BackupPolicyMySQLFlexibleServerName(v.Input, "name")
		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %t but got %t (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
