package validate

import "testing"

func TestFrontDoorRuleSetName(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			// Empty
			Input: "",
			Valid: false,
		},

		{
			// Starts invalid character
			Input: "-foo",
			Valid: false,
		},

		{
			// Ends with invalid character
			Input: "foo-",
			Valid: false,
		},

		{
			// Has embedded invalid character
			Input: "foo-bar",
			Valid: false,
		},

		{
			// Starts with number
			Input: "1foo",
			Valid: false,
		},

		{
			// Ends with number
			Input: "foo1",
			Valid: true,
		},

		{
			// Min Len
			Input: "f",
			Valid: true,
		},

		{
			// Max Len
			Input: "AAAAAAAAAAAAAHHHHHHHHHHHHHHHHIIIIIIIIIIIIIIIIEEEEEEEEEEEEEEE",
			Valid: true,
		},

		{
			// Too Long
			Input: "AAAAAAAAAAAAAHHHHHHHHHHHHHHHHIIIIIIIIIIIIIIIIEEEEEEEEEEEEEEEE",
			Valid: false,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := FrontDoorRuleSetName(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
