// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestTriggerTimespan(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			// empty
			Input: "",
			Valid: false,
		},
		{
			// invalid timespan
			Input: "12:30",
			Valid: false,
		},
		{
			// invalid character
			Input: "hh:mm:ss",
			Valid: false,
		},
		{
			// invalid second
			Input: "12:34:61",
			Valid: false,
		},
		{
			// invalid minute
			Input: "24:61:23",
			Valid: false,
		},
		{
			// invalid hour
			Input: "24:34:61",
			Valid: false,
		},
		{
			Input: "12:34:56",
			Valid: true,
		},
		{
			Input: "1.12:34:56",
			Valid: true,
		},
		{
			Input: "-1.12:34:56",
			Valid: true,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := TriggerTimespan(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
