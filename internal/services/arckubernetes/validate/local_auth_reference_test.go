// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestLocalAuthReference(t *testing.T) {
	cases := []struct {
		LocalAuthReference string
		Errors             int
	}{
		{
			LocalAuthReference: "",
			Errors:             1,
		},
		{
			LocalAuthReference: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			Errors:             1,
		},
		{
			LocalAuthReference: "a",
			Errors:             0,
		},
		{
			LocalAuthReference: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			Errors:             0,
		},
		{
			LocalAuthReference: "A",
			Errors:             1,
		},
		{
			LocalAuthReference: "-",
			Errors:             1,
		},
		{
			LocalAuthReference: "a-1",
			Errors:             0,
		},
		{
			LocalAuthReference: "1-a",
			Errors:             0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.LocalAuthReference, func(t *testing.T) {
			_, errors := LocalAuthReference(tc.LocalAuthReference, "test")

			if len(errors) != tc.Errors {
				t.Fatalf("Expected Local Auth Reference to return %d error(s) not %d", len(errors), tc.Errors)
			}
		})
	}
}
