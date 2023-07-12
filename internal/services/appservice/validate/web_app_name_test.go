// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
)

func TestWebAppName(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			Input: "",
			Valid: false,
		},
		{
			Input: "a",
			Valid: true,
		},
		{
			Input: "-valid",
			Valid: true,
		},
		{
			// len is 35
			Input: "ThisIsALongAndValidNameThatWillWork",
			Valid: true,
		},
		{
			Input: "ThisNameIsTooLongThisNameIsTooLongThisNameIsTooLongThisNameIsTooLongThisNameIsTooLongThisNameIsTooLongThisNameIsTooLongThisNameIsTooLong",
			Valid: false,
		},
		{
			// len is 60 and should show the warning message
			Input: "012345678901234567890123456789012345678901234567890123456789",
			Valid: true,
		},
	}

	for _, tc := range cases {
		_, errs := validate.WebAppName(tc.Input, "test")
		valid := len(errs) == 0

		if valid != tc.Valid {
			t.Fatalf("expected %s to be %t, got %t", tc.Input, tc.Valid, valid)
		}
	}
}
