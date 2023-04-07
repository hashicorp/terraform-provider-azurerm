package validate

import (
	"strings"
	"testing"
)

func TestAvailabilitySetName(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			Input: "",
			Valid: false,
		},
		{
			Input: "testname",
			Valid: true,
		},
		{
			Input: strings.Repeat("s", 79),
			Valid: true,
		},
		{
			Input: strings.Repeat("s", 80),
			Valid: true,
		},
		{
			Input: strings.Repeat("s", 81),
			Valid: false,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := AvailabilitySetName(tc.Input, "availability_set_name")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
