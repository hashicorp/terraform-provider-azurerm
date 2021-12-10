package validate

import (
	"strings"
	"testing"
)

func TestExpressRoutePortName(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{

		{
			// empty
			Input: "",
			Valid: false,
		},

		{
			// starts with word (length: 1)
			Input: "a",
			Valid: true,
		},

		{
			// starts with number (length: 1)
			Input: "1",
			Valid: true,
		},

		{
			// starts with non alphnum  (length: 1)
			Input: "_",
			Valid: false,
		},

		{
			// ends with word
			Input: "aa",
			Valid: true,
		},

		{
			// ends with number
			Input: "a1",
			Valid: true,
		},

		{
			// ends with non underline
			Input: "a_",
			Valid: true,
		},

		{
			// ends with non alphnum or underline
			Input: "a.",
			Valid: false,
		},

		{
			// max length (80)
			Input: strings.Repeat("a", 80),
			Valid: true,
		},

		{
			// exceed max length (80)
			Input: strings.Repeat("a", 81),
			Valid: false,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := ExpressRoutePortName(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
