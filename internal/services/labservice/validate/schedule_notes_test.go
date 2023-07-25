// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestValidateScheduleNotes(t *testing.T) {
	cases := []struct {
		Input       string
		ExpectError bool
	}{
		{
			Input:       "",
			ExpectError: true,
		},
		{
			Input:       "hello",
			ExpectError: false,
		},
		{
			Input:       strings.Repeat("s", 999),
			ExpectError: false,
		},
		{
			Input:       strings.Repeat("s", 1000),
			ExpectError: false,
		},
		{
			Input:       strings.Repeat("s", 1001),
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		_, errors := ScheduleNotes(tc.Input, "notes")

		hasError := len(errors) > 0
		if tc.ExpectError && !hasError {
			t.Fatalf("Expected the Lab Plan Name to trigger a validation error for '%s'", tc.Input)
		}
	}
}
