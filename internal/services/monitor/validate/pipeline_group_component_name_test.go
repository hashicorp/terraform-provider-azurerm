// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestPipelineGroupComponentName(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{Input: "", Valid: false},
		{Input: "ab", Valid: false},
		{Input: "abc", Valid: false},
		{Input: "abcd", Valid: true},
		{Input: "custom-table-exporter", Valid: true},
		{Input: "-leading-hyphen", Valid: false},
		{Input: "trailing-hyphen-", Valid: false},
		{Input: "has_underscore", Valid: false},
		{Input: strings.Repeat("a", 33), Valid: true},
		{Input: strings.Repeat("a", 34), Valid: false},
	}

	for _, tc := range cases {
		t.Run(tc.Input, func(t *testing.T) {
			_, errors := PipelineGroupComponentName(tc.Input, "name")
			valid := len(errors) == 0
			if valid != tc.Valid {
				t.Fatalf("expected %q to be valid=%t, got valid=%t", tc.Input, tc.Valid, valid)
			}
		})
	}
}
