// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestFrontDoorContentType(t *testing.T) {
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
			// missing subtype
			Input: "text/",
			Valid: false,
		},

		{
			// missing type
			Input: "/html",
			Valid: false,
		},

		{
			// no separator
			Input: "texthtml",
			Valid: false,
		},

		{
			// multiple separators
			Input: "text/html/xml",
			Valid: false,
		},

		{
			// whitespace
			Input: "text/ html",
			Valid: false,
		},

		{
			// parameters are not supported
			Input: "text/html; charset=utf-8",
			Valid: false,
		},

		{
			Input: "text/html",
			Valid: true,
		},

		{
			Input: "application/xml+rss",
			Valid: true,
		},

		{
			Input: "application/vnd.ms-fontobject",
			Valid: true,
		},

		{
			// previously blocked by the hardcoded allow list
			Input: "application/x-ns-proxy-autoconfig",
			Valid: true,
		},

		{
			Input: "application/wasm",
			Valid: true,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := FrontDoorContentType(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t for %q", tc.Valid, valid, tc.Input)
		}
	}
}
