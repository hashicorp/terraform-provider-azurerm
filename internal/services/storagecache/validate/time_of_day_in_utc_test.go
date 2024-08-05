// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestTimeOfDayInUTC(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected bool
	}{
		{
			Input:    "",
			Expected: false,
		},
		{
			Input:    "t",
			Expected: false,
		},
		{
			Input:    "23",
			Expected: false,
		},
		{
			Input:    "00:00",
			Expected: true,
		},
		{
			Input:    "09:00",
			Expected: true,
		},
		{
			Input:    "10:10",
			Expected: true,
		},
		{
			Input:    "23:59",
			Expected: true,
		},
		{
			Input:    "24:01",
			Expected: false,
		},
		{
			Input:    "123:123",
			Expected: false,
		},
	}

	for _, v := range testCases {
		_, errors := TimeOfDayInUTC(v.Input, "time_of_day_in_utc")
		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %t but got %t (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
