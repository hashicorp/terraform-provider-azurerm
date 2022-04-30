package validate_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
)

func TestWebAppName(t *testing.T) {
	cases := []struct {
		Input       string
		Valid       bool
		ShowWarning bool
	}{
		{
			Input:       "",
			Valid:       false,
			ShowWarning: false,
		},
		{
			Input:       "a",
			Valid:       true,
			ShowWarning: false,
		},
		{
			Input:       "-valid",
			Valid:       true,
			ShowWarning: false,
		},
		{
			// len is 35
			Input:       "ThisIsALongAndValidNameThatWillWork",
			Valid:       true,
			ShowWarning: true,
		},
		{
			Input:       "ThisNameIsTooLongThisNameIsTooLongThisNameIsTooLongThisNameIsTooLongThisNameIsTooLongThisNameIsTooLongThisNameIsTooLongThisNameIsTooLong",
			Valid:       false,
			ShowWarning: true,
		},
		{
			// len is 60 and should show the warning message
			Input:       "012345678901234567890123456789012345678901234567890123456789",
			Valid:       true,
			ShowWarning: true,
		},
	}

	for _, tc := range cases {
		warnings, errs := validate.WebAppName(tc.Input, "test")
		valid := len(errs) == 0
		showWarning := len(warnings) > 0

		if valid != tc.Valid && showWarning != tc.ShowWarning {
			t.Fatalf("expected %s to be %t, got %t", tc.Input, tc.Valid, valid)
		}
	}
}
