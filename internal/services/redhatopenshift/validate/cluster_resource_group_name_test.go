// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestClusterResourceGroupName(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			Input: "",
			Valid: false,
		},

		{
			Input: "2",
			Valid: true,
		},

		{
			Input: "a.1",
			Valid: true,
		},

		{
			Input: "a-a",
			Valid: true,
		},

		{
			Input: "a.",
			Valid: false,
		},

		{
			Input: "2a-9",
			Valid: true,
		},

		{
			Input: "2a-9_",
			Valid: true,
		},

		// upper case
		{
			Input: "2A",
			Valid: false,
		},

		// 90 chars
		{
			Input: "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijkl",
			Valid: true,
		},

		// 91 chars
		{
			Input: "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklm",
			Valid: false,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := ClusterResourceGroupName(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
