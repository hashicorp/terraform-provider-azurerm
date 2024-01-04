package validate

import (
	"strings"
	"testing"
)

func TestAdminUsername(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			Input: "",
			Valid: false,
		},
		{
			Input: "testAdmin",
			Valid: true,
		},
		{
			Input: strings.Repeat("s", 63),
			Valid: true,
		},
		{
			Input: strings.Repeat("s", 64),
			Valid: true,
		},
		{
			Input: strings.Repeat("s", 65),
			Valid: false,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := AdminUsername(tc.Input, "admin_username")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
