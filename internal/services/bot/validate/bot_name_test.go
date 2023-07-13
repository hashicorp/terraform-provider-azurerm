// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestBotName(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected bool
	}{
		{
			Input:    "Test123",
			Expected: true,
		},
		{
			Input:    "Test_123",
			Expected: true,
		},
		{
			Input:    "Test-123",
			Expected: true,
		},
		{
			Input:    strings.Repeat("s", 41),
			Expected: true,
		},
		{
			Input:    strings.Repeat("s", 42),
			Expected: true,
		},
		{
			Input:    strings.Repeat("s", 43),
			Expected: false,
		},
	}
	for _, v := range testCases {
		_, errors := BotName(v.Input, "bot_name")
		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %t but got %t (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
