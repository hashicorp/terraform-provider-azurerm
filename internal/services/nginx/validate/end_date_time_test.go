// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate_test

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/nginx/validate"
)

func TestValidateNetAppBackupPolicyCombinedRetention(t *testing.T) {
	cases := []struct {
		input interface{}
		valid bool
	}{
		{
			// invalid - empty
			input: "",
			valid: false,
		},
		{
			// invalid - wrong input type
			input: 123,
			valid: false,
		},
		{
			// invalid - not a valid RFC3339 format
			input: "Jan. 1 2025",
			valid: false,
		},
		{
			// invalid - 3 months in the past
			input: time.Now().AddDate(0, -3, 0).Format(time.RFC3339),
			valid: false,
		},
		{
			// invalid - 5 years into the future
			input: time.Now().AddDate(5, 0, 0).Format(time.RFC3339),
			valid: false,
		},
		{
			// valid - 3 months into the future
			input: time.Now().AddDate(0, 3, 0).Format(time.RFC3339),
			valid: true,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value: %+v", tc.input)
		_, errors := validate.EndDateTime(tc.input, "end_date_time")
		valid := len(errors) == 0

		if tc.valid != valid {
			t.Fatalf("Expected %t but got %t", tc.valid, valid)
		}
	}
}
