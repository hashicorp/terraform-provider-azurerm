// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
)

func TestFunctionAppFunctionName(t *testing.T) {
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
			Input: "-notValid",
			Valid: false,
		},
		{
			Input: "ThisIsALongAndValidNameThatWillWork",
			Valid: true,
		},
		{
			Input: "ThisNameIsTooLongThisNameIsTooLongThisNameIsTooLongThisNameIsTooLongThisNameIsTooLongThisNameIsTooLongThisNameIsTooLongThisNameIsTooLong",
			Valid: false,
		},
		{
			Input: "EndsInWrongChar-",
			Valid: true,
		},
	}

	for _, tc := range cases {
		_, errs := validate.FunctionAppFunctionName(tc.Input, "test")
		valid := len(errs) == 0

		if valid != tc.Valid {
			t.Fatalf("expected %s to be %t, got %t", tc.Input, tc.Valid, valid)
		}
	}
}
