// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestBatchMaxWaitTime(t *testing.T) {
	cases := map[string]bool{
		"":         false,
		"NotValid": false,
		"10:00":    false,
		"00:02:00": true,
		"00:00:00": true,
		"99:99:99": true,
		"2":        false,
	}
	for i, shouldBeValid := range cases {
		_, errors := BatchMaxWaitTime(i, "batch_max_wait_time")

		isValid := len(errors) == 0
		if shouldBeValid != isValid {
			t.Fatalf("Expected %s to be %t but got %t", i, shouldBeValid, isValid)
		}
	}
}
