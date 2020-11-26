package validate

import "testing"

func TestSubscriptionID(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			Input: "",
			Valid: false,
		},
		{
			Input: "NotASubscriptionID",
			Valid: false,
		},
		{
			Input: "/subscriptions/",
			Valid: false,
		},
		{
			Input: "/subscriptions/StillNotAValidSubscription",
			Valid: false,
		},
		{
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000",
			Valid: true,
		},
		{
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups",
			Valid: false,
		},
		{
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup1",
			Valid: false,
		},
	}

	for _, tc := range cases {
		_, errors := SubscriptionID(tc.Input, "")

		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t for %q", tc.Valid, valid, tc.Input)
		}
	}
}
