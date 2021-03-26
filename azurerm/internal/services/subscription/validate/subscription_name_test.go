package validate

import "testing"

func TestSubscriptionName(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			Input: "",
			Valid: false,
		},
		{
			Input: "A Name With Invalid Char >",
			Valid: false,
		},
		{
			Input: "A Valid Name",
			Valid: true,
		},
		{
			Input: "This is a valid but max length name of a subscription so there!!",
			Valid: true,
		},
	}

	for _, tc := range cases {
		_, errors := SubscriptionName(tc.Input, "")

		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t for %q", tc.Valid, valid, tc.Input)
		}
	}
}
