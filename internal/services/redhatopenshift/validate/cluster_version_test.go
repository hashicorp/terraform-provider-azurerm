// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestClusterVersion(t *testing.T) {
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
			Valid: false,
		},

		{
			Input: "1.1",
			Valid: false,
		},

		{
			Input: "3..",
			Valid: false,
		},

		{
			Input: "2.9.0.1",
			Valid: false,
		},

		{
			Input: "2.9.0",
			Valid: true,
		},

		{
			Input: "2.0.77",
			Valid: true,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := ClusterVersion(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
