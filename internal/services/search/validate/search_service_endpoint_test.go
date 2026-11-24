// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestSearchServiceEndpoint(t *testing.T) {
	cases := []struct {
		Input any
		Valid bool
	}{
		{
			Input: "https://example.search.windows.net",
			Valid: true,
		},
		{
			Input: "https://example.search.azure.us",
			Valid: true,
		},
		{
			Input: "https://example.search.azure.cn",
			Valid: true,
		},
		{
			Input: "https://example.search.windows.net/",
			Valid: false,
		},
		{
			Input: "https://example.search.windows.net//",
			Valid: false,
		},
		{
			Input: "http://example.search.windows.net",
			Valid: false,
		},
		{
			Input: "https://",
			Valid: false,
		},
		{
			Input: "not-a-url",
			Valid: false,
		},
		{
			Input: "",
			Valid: false,
		},
		{
			Input: 123,
			Valid: false,
		},
	}

	for _, tc := range cases {
		_, errors := SearchServiceEndpoint(tc.Input, "search_service_endpoint")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("expected %t but got %t for input %v", tc.Valid, valid, tc.Input)
		}
	}
}
