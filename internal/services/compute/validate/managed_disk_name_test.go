// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestManagedDiskName(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{Input: "", Valid: false},
		{Input: "a", Valid: true},
		{Input: "1", Valid: true},
		{Input: "-abc", Valid: false},
		{Input: ".abc", Valid: false},
		{Input: "abc-", Valid: false},
		{Input: "abc.", Valid: false},
		{Input: "abc_", Valid: true},
		{Input: "my-disk_1.0", Valid: true},
		{Input: "my disk", Valid: false},
		{Input: "disk!", Valid: false},
		{Input: strings.Repeat("a", 80), Valid: true},
		{Input: strings.Repeat("a", 81), Valid: false},
	}

	for _, tc := range cases {
		_, errors := ManagedDiskName(tc.Input, "name")
		valid := len(errors) == 0
		if valid != tc.Valid {
			t.Fatalf("expected %q to have validity %t but got %t", tc.Input, tc.Valid, valid)
		}
	}
}
