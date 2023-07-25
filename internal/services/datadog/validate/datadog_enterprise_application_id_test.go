// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestDatadogEnterpriseApplicationID(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected bool
	}{
		{
			Input:    "",
			Expected: false,
		},
		{
			Input:    "1234567890-8665-8687",
			Expected: true,
		},
		{
			Input:    "12345678901234567890123456789012345678901234567890",
			Expected: false,
		},
	}

	for _, v := range testCases {
		_, errors := DatadogEnterpriseApplicationID(v.Input, "enterprise_application_id")
		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %t but got %t (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
