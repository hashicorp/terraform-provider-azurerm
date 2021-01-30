package validate

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccountName(t *testing.T) {
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
			Valid: false,
		},
		{
			Input: "ab",
			Valid: false,
		},
		{
			Input: "abc-",
			Valid: false,
		},
		{
			Input: "-abc",
			Valid: false,
		},
		{
			Input: acctest.RandString(64),
			Valid: false,
		},
		{
			Input: "abc",
			Valid: true,
		},
		{
			Input: "ABC",
			Valid: true,
		},
		{
			Input: "1a3",
			Valid: true,
		},
		{
			Input: "account-purview-01",
			Valid: true,
		},
		{
			Input: acctest.RandString(63),
			Valid: true,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := AccountName()(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
